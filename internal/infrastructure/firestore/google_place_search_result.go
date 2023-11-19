package firestore

import (
	"context"
	"errors"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"poroto.app/poroto/planner/internal/domain/factory"
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
	if err := tx.Set(doc, factory.PlaceEntityFromGooglePlace(googlePlace)); err != nil {
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

	photos := make(map[string][]models.GooglePlacePhoto)
	for _, photoEntity := range photoEntities {
		photos[photoEntity.GooglePlaceId] = append(photos[photoEntity.GooglePlaceId], photoEntity.ToGooglePlacePhoto())
	}

	// TODO: Place Detailを復元する

	var places []models.GooglePlace
	for _, snapshot := range snapshots {
		var googlePlaceEntity entity.GooglePlaceEntity
		if err = snapshot.DataTo(&googlePlaceEntity); err != nil {
			return nil, fmt.Errorf("error while converting snapshot to place search result entity: %v", err)
		}

		imagesOfPlace := photos[googlePlaceEntity.PlaceID]
		places = append(places, googlePlaceEntity.ToGooglePlace(&imagesOfPlace))
	}

	return places, nil
}

func (p GooglePlaceSearchResultRepository) updateOpeningHours(ctx context.Context, planCandidateId string, googlePlaceId string, openingHours []models.GooglePlaceOpeningPeriod) error {
	doc := p.doc(planCandidateId, googlePlaceId)

	if err := p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		openingHoursEntity := entity.GooglePlaceOpeningsEntityFromGooglePlaceOpeningPeriod(openingHours)
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
	}, firestore.MaxAttempts(3)); err != nil {
		return fmt.Errorf("error while updating opening hours: %v", err)
	}
	return nil
}

func (p GooglePlaceSearchResultRepository) saveImages(ctx context.Context, planCandidateId string, googlePlaceId string, photos []models.GooglePlacePhoto) error {
	subCollectionImages := p.subCollectionPhotos(planCandidateId)

	if err := p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		for _, photo := range photos {
			photoEntity := entity.GooglePlacePhotoEntityFromGooglePlacePhoto(photo, googlePlaceId)
			if err := tx.Set(subCollectionImages.Doc(photoEntity.PhotoReference), photoEntity); err != nil {
				return fmt.Errorf("error while saving photo: %v", err)
			}
		}
		return nil
	}, firestore.MaxAttempts(3)); err != nil {
		return fmt.Errorf("error while saving images: %v", err)
	}

	return nil
}

func (p GooglePlaceSearchResultRepository) saveReviewsIfNotExist(ctx context.Context, planCandidateId string, googlePlaceId string, reviews []models.GooglePlaceReview) error {
	subCollectionReviews := p.subCollectionReviews(planCandidateId)

	snapshots, err := subCollectionReviews.Where("google_place_id", "==", googlePlaceId).Limit(1).Documents(ctx).GetAll()
	if err != nil {
		return fmt.Errorf("error while getting reviews: %v", err)
	}

	if len(snapshots) > 0 {
		// すでにレビューが保存されている場合は何もしない
		return fmt.Errorf("reviews already exist")
	}

	for _, review := range reviews {
		if _, err := subCollectionReviews.NewDoc().Set(ctx, entity.GooglePlaceReviewEntityFromGooglePlaceReview(review, googlePlaceId)); err != nil {
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

func (p GooglePlaceSearchResultRepository) fetchPhotosByPlanCandidateId(ctx context.Context, planCandidateId string) ([]entity.GooglePlacePhotoEntity, error) {
	subCollectionPhotos := p.subCollectionPhotos(planCandidateId)
	photosSnapshots, err := subCollectionPhotos.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("error while getting photos: %v", err)
	}

	var photos []entity.GooglePlacePhotoEntity
	for _, photoSnapshot := range photosSnapshots {
		var photoEntity entity.GooglePlacePhotoEntity
		if err = photoSnapshot.DataTo(&photoEntity); err != nil {
			return nil, fmt.Errorf("error while converting snapshot to photo entity: %v", err)
		}
		photos = append(photos, photoEntity)
	}

	return photos, nil
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
