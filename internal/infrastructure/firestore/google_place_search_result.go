package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/firestore/entity"
)

const (
	subCollectionGooglePlaceSearchResults = "google_places_api_search_results"
	subCollectionReviews                  = "google_places_api_reviews"
	subCollectionPhotos                   = "google_places_api_photos"
)

type GooglePlaceSearchResultRepository struct {
	client *firestore.Client
}

func NewGooglePlaceSearchResultRepository(ctx context.Context) (*GooglePlaceSearchResultRepository, error) {
	var options []option.ClientOption
	if os.Getenv("GCP_CREDENTIAL_FILE_PATH") != "" {
		options = append(options, option.WithCredentialsFile(os.Getenv("GCP_CREDENTIAL_FILE_PATH")))
	}

	client, err := firestore.NewClient(ctx, os.Getenv("GCP_PROJECT_ID"), options...)
	if err != nil {
		return nil, fmt.Errorf("error while initializing firestore client: %v", err)
	}

	return &GooglePlaceSearchResultRepository{
		client: client,
	}, nil
}

func (p GooglePlaceSearchResultRepository) saveTx(tx *firestore.Transaction, planCandidateId string, googlePlace models.GooglePlace) error {
	doc := p.doc(planCandidateId, googlePlace.PlaceId)
	if err := tx.Set(doc, entity.GooglePlaceEntityFromGooglePlace(googlePlace)); err != nil {
		return fmt.Errorf("error while saving place search result: %v", err)
	}
	return nil
}

func (p GooglePlaceSearchResultRepository) find(ctx context.Context, planCandidateId string) ([]models.GooglePlace, error) {
	collection := p.subCollection(planCandidateId)

	snapshots, err := collection.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("error while getting place search results: %v", err)
	}

	// 写真を取得
	photoEntities, err := p.fetchPhotosByPlanCandidateId(ctx, planCandidateId)
	if err != nil {
		return nil, fmt.Errorf("error while fetching image entities: %v", err)
	}

	// レビューを取得
	reviewEntities, err := p.fetchReviewsByPlanCandidateId(ctx, planCandidateId)
	if err != nil {
		return nil, fmt.Errorf("error while fetching reviews: %v", err)
	}

	googlePlaceIdToPhotos := make(map[string][]entity.GooglePlacePhotoEntity)
	for _, photoEntity := range *photoEntities {
		googlePlaceIdToPhotos[photoEntity.GooglePlaceId] = append(googlePlaceIdToPhotos[photoEntity.GooglePlaceId], photoEntity)
	}

	googlePlaceIdToReviews := make(map[string][]entity.GooglePlaceReviewEntity)
	for _, reviewEntity := range *reviewEntities {
		googlePlaceIdToReviews[reviewEntity.GooglePlaceId] = append(googlePlaceIdToReviews[reviewEntity.GooglePlaceId], reviewEntity)
	}

	var places []models.GooglePlace
	for _, snapshot := range snapshots {
		var googlePlaceEntity entity.GooglePlaceEntity
		if err = snapshot.DataTo(&googlePlaceEntity); err != nil {
			return nil, fmt.Errorf("error while converting snapshot to place search result entity: %v", err)
		}

		places = append(places, googlePlaceEntity.ToGooglePlace(googlePlaceIdToPhotos[googlePlaceEntity.PlaceID], googlePlaceIdToReviews[googlePlaceEntity.PlaceID]))
	}

	return places, nil
}

func (p GooglePlaceSearchResultRepository) findGooglePlace(ctx context.Context, planCandidateId string, googlePlaceId string) (*models.GooglePlace, error) {
	doc := p.doc(planCandidateId, googlePlaceId)

	snapshotGooglePlaceEntity, err := doc.Get(ctx)
	if err != nil && status.Code(err) != codes.NotFound {
		return nil, fmt.Errorf("error while getting user: %v", err)
	}

	if !snapshotGooglePlaceEntity.Exists() {
		return nil, nil
	}

	var googlePlaceEntity entity.GooglePlaceEntity
	if err = snapshotGooglePlaceEntity.DataTo(&googlePlaceEntity); err != nil {
		return nil, fmt.Errorf("error while converting snapshotGooglePlaceEntity to place search result entity: %v", err)
	}

	// レビューを取得
	reviews, err := p.fetchReviewsByGooglePlaceId(ctx, planCandidateId, googlePlaceId)
	if err != nil {
		return nil, fmt.Errorf("error while fetching reviews: %v", err)
	}

	// 写真を取得
	photos, err := p.fetchPhotosByGooglePlaceId(ctx, planCandidateId, googlePlaceId)
	if err != nil {
		return nil, fmt.Errorf("error while fetching image entities: %v", err)
	}

	googlePlace := googlePlaceEntity.ToGooglePlace(*photos, *reviews)
	return &googlePlace, nil
}

// fetchPhotosByGooglePlaceId は、指定したGoogle Place IDに紐づく写真を取得する
// 一箇所だけ取得する場合はこの方法がクエリ回数が少ない
func (p GooglePlaceSearchResultRepository) fetchPhotosByGooglePlaceId(ctx context.Context, planCandidateId string, googlePlaceId string) (*[]entity.GooglePlacePhotoEntity, error) {
	snapshots, err := p.subCollectionPhotos(planCandidateId).Where("google_place_id", "==", googlePlaceId).Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("error while getting photos: %v", err)
	}

	var photos []entity.GooglePlacePhotoEntity
	for _, snapshot := range snapshots {
		var photoEntity entity.GooglePlacePhotoEntity
		if err = snapshot.DataTo(&photoEntity); err != nil {
			return nil, fmt.Errorf("error while converting snapshot to photo entity: %v", err)
		}
		photos = append(photos, photoEntity)
	}

	return &photos, nil
}

// fetchPhotosByPlanCandidateId は、指定したプラン候補で取得されたすべての写真を取得する
// すべてのプランを取得する場合はこの方法がクエリ回数が少ない
func (p GooglePlaceSearchResultRepository) fetchPhotosByPlanCandidateId(ctx context.Context, planCandidateId string) (*[]entity.GooglePlacePhotoEntity, error) {
	photosSnapshots, err := p.subCollectionPhotos(planCandidateId).Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("error while getting photos: %v", err)
	}

	var photos []entity.GooglePlacePhotoEntity
	for _, snapshot := range photosSnapshots {
		var photoEntity entity.GooglePlacePhotoEntity
		if err = snapshot.DataTo(&photoEntity); err != nil {
			return nil, fmt.Errorf("error while converting snapshot to photo entity: %v", err)
		}
		photos = append(photos, photoEntity)
	}

	return &photos, nil
}

// fetchReviewsByPlanCandidateId は、Google Place IDに紐づくレビューを取得する
// 一箇所だけ取得する場合はこの方法がクエリ回数が少ない
func (p GooglePlaceSearchResultRepository) fetchReviewsByGooglePlaceId(ctx context.Context, planCandidateId string, googlePlaceId string) (*[]entity.GooglePlaceReviewEntity, error) {
	snapshots, err := p.subCollectionReviews(planCandidateId).Where("google_place_id", "==", googlePlaceId).Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("error while getting reviews: %v", err)
	}

	var reviews []entity.GooglePlaceReviewEntity
	for _, snapshot := range snapshots {
		var reviewEntity entity.GooglePlaceReviewEntity
		if err = snapshot.DataTo(&reviewEntity); err != nil {
			return nil, fmt.Errorf("error while converting snapshot to review entity: %v", err)
		}
		reviews = append(reviews, reviewEntity)
	}

	return &reviews, nil
}

func (p GooglePlaceSearchResultRepository) fetchReviewsByPlanCandidateId(ctx context.Context, planCandidateId string) (*[]entity.GooglePlaceReviewEntity, error) {
	snapshots, err := p.subCollectionReviews(planCandidateId).Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("error while getting photos: %v", err)
	}

	var reviews []entity.GooglePlaceReviewEntity
	for _, snapshot := range snapshots {
		var reviewEntity entity.GooglePlaceReviewEntity
		if err = snapshot.DataTo(&reviewEntity); err != nil {
			return nil, fmt.Errorf("error while converting snapshot to review entity: %v", err)
		}
		reviews = append(reviews, reviewEntity)
	}

	return &reviews, nil
}

func (p GooglePlaceSearchResultRepository) updateOpeningHoursTx(tx *firestore.Transaction, planCandidateId string, googlePlaceId string, openingHours models.GooglePlaceOpeningHours) error {
	doc := p.doc(planCandidateId, googlePlaceId)

	openingHoursEntity := entity.GooglePlaceOpeningsEntityFromGooglePlaceOpeningHours(openingHours)
	if err := tx.Update(doc, []firestore.Update{
		{
			Path:  "opening_hours",
			Value: openingHoursEntity,
		},
		{
			Path:  "updated_at",
			Value: firestore.ServerTimestamp,
		},
	}); err != nil {
		return fmt.Errorf("error while updating opening hours: %v", err)
	}

	return nil
}

func (p GooglePlaceSearchResultRepository) saveGooglePlacePhotosTx(tx *firestore.Transaction, planCandidateId string, googlePlaceId string, photos []models.GooglePlacePhoto) error {
	for _, photo := range photos {
		photoEntity := entity.GooglePlacePhotoEntityFromGooglePlacePhoto(photo, googlePlaceId)

		var updates []firestore.Update
		if photoEntity.Small != nil {
			updates = append(updates, firestore.Update{
				Path:  "small",
				Value: photoEntity.Small,
			})
		}
		if photoEntity.Large != nil {
			updates = append(updates, firestore.Update{
				Path:  "large",
				Value: photoEntity.Large,
			})
		}

		// 画像が取得されていない場合は何もしない
		if len(updates) == 0 {
			continue
		}

		updates = append(updates, firestore.Update{
			Path:  "updated_at",
			Value: firestore.ServerTimestamp,
		})

		if err := tx.Update(p.subCollectionPhotos(planCandidateId).Doc(photoEntity.PhotoReference), updates); err != nil {
			return fmt.Errorf("error while saving photo: %v", err)
		}
	}
	return nil
}

func (p GooglePlaceSearchResultRepository) savePhotoReferencesTx(tx *firestore.Transaction, planCandidateId string, googlePlaceId string, photoReferences []models.GooglePlacePhotoReference) error {
	for _, photoReference := range photoReferences {
		doc := p.subCollectionPhotos(planCandidateId).Doc(photoReference.PhotoReference)
		if err := tx.Set(doc, entity.GooglePlacePhotoEntityFromGooglePhotoReference(photoReference, googlePlaceId)); err != nil {
			return fmt.Errorf("error while saving photo reference: %v", err)
		}
	}
	return nil
}

func (p GooglePlaceSearchResultRepository) reviewAlreadySavedTx(tx *firestore.Transaction, planCandidateId string, googlePlaceId string) (*bool, error) {
	query := p.subCollectionReviews(planCandidateId).Where("google_place_id", "==", googlePlaceId).Limit(1)
	snapshots, err := tx.Documents(query).GetAll()
	if err != nil {
		return nil, fmt.Errorf("error while getting reviews: %v", err)
	}

	isAlreadySaved := len(snapshots) > 0
	return &isAlreadySaved, nil
}

func (p GooglePlaceSearchResultRepository) saveReviewsTx(tx *firestore.Transaction, planCandidateId string, googlePlaceId string, reviews []models.GooglePlaceReview) error {
	for _, review := range reviews {
		if err := tx.Set(p.subCollectionReviews(planCandidateId).NewDoc(), entity.GooglePlaceReviewEntityFromGooglePlaceReview(review, googlePlaceId)); err != nil {
			return fmt.Errorf("error while saving review: %v", err)
		}
	}

	return nil
}

func (p GooglePlaceSearchResultRepository) deleteByPlanCandidateIdTx(tx *firestore.Transaction, planCandidateId string) error {
	collection := p.subCollection(planCandidateId)

	docIter := tx.DocumentRefs(collection)
	for {
		doc, err := docIter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return fmt.Errorf("error while iterating documents: %v", err)
		}
		if err := tx.Delete(doc); err != nil {
			return fmt.Errorf("error while deleting place search result: %v", err)
		}
	}

	return nil
}

func (p GooglePlaceSearchResultRepository) subCollection(planCandidateId string) *firestore.CollectionRef {
	return p.client.Collection(collectionPlanCandidates).Doc(planCandidateId).Collection(subCollectionGooglePlaceSearchResults)
}

func (p GooglePlaceSearchResultRepository) doc(planCandidateId string, googlePlaceId string) *firestore.DocumentRef {
	return p.subCollection(planCandidateId).Doc(googlePlaceId)
}

func (p GooglePlaceSearchResultRepository) subCollectionPhotos(planCandidateId string) *firestore.CollectionRef {
	return p.client.Collection(collectionPlanCandidates).Doc(planCandidateId).Collection(subCollectionPhotos)
}

func (p GooglePlaceSearchResultRepository) subCollectionReviews(planCandidateId string) *firestore.CollectionRef {
	return p.client.Collection(collectionPlanCandidates).Doc(planCandidateId).Collection(subCollectionReviews)
}
