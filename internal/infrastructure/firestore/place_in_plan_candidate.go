package firestore

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/firestore/entity"
)

const (
	collectionPlacesInPlanCandidate = "places"
)

// PlaceInPlanCandidateRepository models.PlaceInPlanCandidate を管理するリポジトリ
// 実際には models.PlaceInPlanCandidate#Google の内容を `google_place_search_result` コレクションの中に 保存し、
// `places` コレクションの中で検索結果との対応関係を取る
type PlaceInPlanCandidateRepository struct {
	client                            *firestore.Client
	googlePlaceSearchResultRepository *GooglePlaceSearchResultRepository
}

func NewPlaceInPlanCandidateRepository(ctx context.Context) (*PlaceInPlanCandidateRepository, error) {
	var options []option.ClientOption
	if os.Getenv("GCP_CREDENTIAL_FILE_PATH") != "" {
		options = append(options, option.WithCredentialsFile(os.Getenv("GCP_CREDENTIAL_FILE_PATH")))
	}

	client, err := firestore.NewClient(ctx, os.Getenv("GCP_PROJECT_ID"), options...)
	if err != nil {
		return nil, fmt.Errorf("error while initializing firestore client: %v", err)
	}

	googlePlaceSearchResultRepository, err := NewGooglePlaceSearchResultRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while initializing google place search result repository: %v", err)
	}

	return &PlaceInPlanCandidateRepository{
		client:                            client,
		googlePlaceSearchResultRepository: googlePlaceSearchResultRepository,
	}, nil
}

func (p PlaceInPlanCandidateRepository) SavePlaces(ctx context.Context, planCandidateId string, places []models.PlaceInPlanCandidate) error {
	for _, place := range places {
		if err := p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
			if err := p.saveTx(tx, planCandidateId, place); err != nil {
				return fmt.Errorf("error while saving place in plan candidate: %v", err)
			}

			return nil
		}, firestore.MaxAttempts(3)); err != nil {
			return fmt.Errorf("error while saving place in plan candidates: %v", err)
		}
	}
	return nil
}

func (p PlaceInPlanCandidateRepository) SaveGooglePlaceDetail(ctx context.Context, planCandidateId string, googlePlaceId string, googlePlaceDetail models.GooglePlaceDetail) error {
	if err := p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		isReviewAlreadySaved, err := p.googlePlaceSearchResultRepository.reviewAlreadySavedTx(tx, planCandidateId, googlePlaceId)
		if err != nil {
			return fmt.Errorf("error while checking review already saved: %v", err)
		}

		if err := p.savePlaceDetailTx(tx, planCandidateId, googlePlaceId, googlePlaceDetail, *isReviewAlreadySaved); err != nil {
			return fmt.Errorf("error while saving google place detail: %v", err)
		}

		return nil
	}, firestore.MaxAttempts(3)); err != nil {
		return fmt.Errorf("error while saving google place detail: %v", err)
	}

	return nil
}

func (p PlaceInPlanCandidateRepository) SaveGooglePlacePhotos(ctx context.Context, planCandidateId string, googlePlaceId string, photos []models.GooglePlacePhoto) error {
	if err := p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		return p.saveGooglePlacePhotosTx(tx, planCandidateId, googlePlaceId, photos)
	}, firestore.MaxAttempts(3)); err != nil {
		return fmt.Errorf("error while saving google place photos: %v", err)
	}
	return nil
}

func (p PlaceInPlanCandidateRepository) saveTx(tx *firestore.Transaction, planCandidateId string, place models.PlaceInPlanCandidate) error {
	// 事前にレビューが保存されているかを確認
	isReviewAlreadySaved, err := p.googlePlaceSearchResultRepository.reviewAlreadySavedTx(tx, planCandidateId, place.Google.PlaceId)
	if err != nil {
		return fmt.Errorf("error while checking review already saved: %v", err)
	}

	doc := p.collectionPlaces(planCandidateId).Doc(place.Id)
	if err := tx.Set(doc, entity.ToPlaceInPlanCandidateEntity(place)); err != nil {
		return fmt.Errorf("error while saving place in plan candidate: %v", err)
	}

	// Google Places APIの検索結果を保存
	if err := p.googlePlaceSearchResultRepository.saveTx(tx, planCandidateId, place.Google); err != nil {
		return fmt.Errorf("error while saving google place: %v", err)
	}

	if place.Google.PlaceDetail != nil {
		if err := p.savePlaceDetailTx(tx, planCandidateId, place.Google.PlaceId, *place.Google.PlaceDetail, *isReviewAlreadySaved); err != nil {
			return fmt.Errorf("error while saving google place detail: %v", err)
		}
	}

	// 写真を取得している場合は保存
	if place.Google.Photos != nil {
		if err := p.googlePlaceSearchResultRepository.saveTx(tx, planCandidateId, place.Google); err != nil {
			return err
		}
	}

	return nil
}

func (p PlaceInPlanCandidateRepository) savePlaceDetailTx(tx *firestore.Transaction, planCandidateId string, googlePlaceId string, placeDetail models.GooglePlaceDetail, isReviewAlreadySaved bool) error {
	// レビューを保存(レビューは一意性の判断ができないため、保存されていない場合のみ保存を行う)
	// TODO: コンテンツや投稿日などで一意性を判断する
	if !isReviewAlreadySaved {
		if err := p.googlePlaceSearchResultRepository.saveReviewsTx(tx, planCandidateId, googlePlaceId, placeDetail.Reviews); err != nil {
			return fmt.Errorf("error while saving google place detail: %v", err)
		}
	}

	// OpeningHoursを保存
	if placeDetail.OpeningHours != nil {
		if err := p.googlePlaceSearchResultRepository.updateOpeningHoursTx(tx, planCandidateId, googlePlaceId, *placeDetail.OpeningHours); err != nil {
			return fmt.Errorf("error while saving google place detail: %v", err)
		}
	}

	// PhotoReferenceを保存
	if err := p.googlePlaceSearchResultRepository.savePhotoReferencesTx(tx, planCandidateId, googlePlaceId, placeDetail.PhotoReferences); err != nil {
		return fmt.Errorf("error while saving google place detail: %v", err)
	}

	return nil
}

func (p PlaceInPlanCandidateRepository) saveGooglePlacePhotosTx(tx *firestore.Transaction, planCandidateId string, googlePlaceId string, photos []models.GooglePlacePhoto) error {
	return p.googlePlaceSearchResultRepository.saveGooglePlacePhotosTx(tx, planCandidateId, googlePlaceId, photos)
}

func (p PlaceInPlanCandidateRepository) FindByPlanCandidateId(ctx context.Context, planCandidateId string) (*[]models.PlaceInPlanCandidate, error) {
	collection := p.collectionPlaces(planCandidateId)

	snapshots, err := collection.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("error while getting place in plan candidates: %v", err)
	}

	// Google Places APIの検索結果を取得
	performanceTimer := time.Now()
	googlePlaces, err := p.googlePlaceSearchResultRepository.find(ctx, planCandidateId)
	if err != nil {
		return nil, fmt.Errorf("error while fetching google places: %v", err)
	}
	log.Printf("fetching google places took %v\n", time.Since(performanceTimer))

	var places []models.PlaceInPlanCandidate
	for _, snapshot := range snapshots {
		var place entity.PlaceInPlanCandidateEntity
		if err := snapshot.DataTo(&place); err != nil {
			return nil, fmt.Errorf("error while converting place in plan candidate: %v", err)
		}

		var googlePlace models.GooglePlace
		for _, gp := range googlePlaces {
			if gp.PlaceId == place.GooglePlaceId {
				googlePlace = gp
				break
			}
		}

		places = append(places, entity.FromPlaceInPlanCandidateEntity(place, googlePlace))
	}

	return &places, nil
}

func (p PlaceInPlanCandidateRepository) FindByGooglePlaceId(ctx context.Context, planCandidateId string, googlePlaceId string) (*models.PlaceInPlanCandidate, error) {
	snapshots, err := p.collectionPlaces(planCandidateId).Where("google_place_id", "==", googlePlaceId).Limit(1).Documents(ctx).GetAll()
	if err != nil && !errors.Is(err, iterator.Done) {
		return nil, fmt.Errorf("error while getting place in plan candidates: %v", err)
	}

	if len(snapshots) == 0 {
		return nil, nil
	}

	var place entity.PlaceInPlanCandidateEntity
	if err := snapshots[0].DataTo(&place); err != nil {
		return nil, fmt.Errorf("error while converting place in plan candidate: %v", err)
	}

	googlePlace, err := p.googlePlaceSearchResultRepository.findGooglePlace(ctx, planCandidateId, googlePlaceId)
	if err != nil {
		return nil, fmt.Errorf("error while fetching google place: %v", err)
	}

	placeInPlanCandidate := googlePlace.ToPlaceInPlanCandidate(place.Id)
	return &placeInPlanCandidate, nil
}

func (p PlaceInPlanCandidateRepository) DeleteByPlanCandidateId(ctx context.Context, planCandidateId string) error {
	if err := p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		if err := p.googlePlaceSearchResultRepository.deleteByPlanCandidateIdTx(tx, planCandidateId); err != nil {
			return fmt.Errorf("error while deleting google place search results: %v", err)
		}

		docIter := tx.DocumentRefs(p.collectionPlaces(planCandidateId))
		for {
			doc, err := docIter.Next()
			if errors.Is(err, iterator.Done) {
				break
			}
			if err != nil {
				return fmt.Errorf("error while iterating place in plan candidates: %v", err)
			}
			if err := tx.Delete(doc); err != nil {
				return fmt.Errorf("error while deleting place in plan candidates: %v", err)
			}
		}

		return nil
	}, firestore.MaxAttempts(3)); err != nil {
		return fmt.Errorf("error while deleting place in plan candidates: %v", err)
	}

	return nil
}

func (p PlaceInPlanCandidateRepository) collectionPlaces(planCandidateId string) *firestore.CollectionRef {
	return p.client.Collection(collectionPlanCandidates).Doc(planCandidateId).Collection(collectionPlacesInPlanCandidate)
}
