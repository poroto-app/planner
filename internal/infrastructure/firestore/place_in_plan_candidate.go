package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"log"
	"os"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/firestore/entity"
	"time"
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

func (p PlaceInPlanCandidateRepository) Save(ctx context.Context, planCandidateId string, place models.PlaceInPlanCandidate) error {
	if err := p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		doc := p.collectionPlaces(planCandidateId).Doc(place.Id)
		if err := tx.Set(doc, entity.ToPlaceInPlanCandidateEntity(place)); err != nil {
			return fmt.Errorf("error while saving place in plan candidate: %v", err)
		}

		// Google Places APIの検索結果を保存
		if err := p.googlePlaceSearchResultRepository.saveTx(tx, planCandidateId, place.Google); err != nil {
			return fmt.Errorf("error while saving google place: %v", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error while saving place in plan candidate: %v", err)
	}

	return nil
}

func (p PlaceInPlanCandidateRepository) SavePlaces(ctx context.Context, planCandidateId string, places []models.PlaceInPlanCandidate) error {
	if err := p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		for _, place := range places {
			doc := p.collectionPlaces(planCandidateId).Doc(place.Id)
			if err := tx.Set(doc, entity.ToPlaceInPlanCandidateEntity(place)); err != nil {
				return fmt.Errorf("error while saving place in plan candidate: %v", err)
			}

			// Google Places APIの検索結果を保存
			if err := p.googlePlaceSearchResultRepository.saveTx(tx, planCandidateId, place.Google); err != nil {
				return fmt.Errorf("error while saving google place: %v", err)
			}
		}
		return nil
	}, firestore.MaxAttempts(3)); err != nil {
		return fmt.Errorf("error while saving place in plan candidates: %v", err)
	}

	return nil
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

func (p PlaceInPlanCandidateRepository) SaveGoogleImages(ctx context.Context, planCandidateId string, googlePlaceId string, images []models.Image) error {
	if err := p.googlePlaceSearchResultRepository.saveImagesIfNotExist(ctx, planCandidateId, googlePlaceId, images); err != nil {
		return fmt.Errorf("error while saving google images: %v", err)
	}
	return nil
}

func (p PlaceInPlanCandidateRepository) SaveGoogleReviews(ctx context.Context, planCandidateId string, googlePlaceId string, reviews []models.GooglePlaceReview) error {
	if err := p.googlePlaceSearchResultRepository.saveReviewsIfNotExist(ctx, planCandidateId, googlePlaceId, reviews); err != nil {
		return fmt.Errorf("error while saving google reviews: %v", err)
	}
	return nil
}

func (p PlaceInPlanCandidateRepository) collectionPlaces(planCandidateId string) *firestore.CollectionRef {
	return p.client.Collection(collectionPlanCandidates).Doc(planCandidateId).Collection(collectionPlacesInPlanCandidate)
}
