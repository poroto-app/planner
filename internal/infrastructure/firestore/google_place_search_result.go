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
	googleplaces "poroto.app/poroto/planner/internal/infrastructure/api/google/places"
	"poroto.app/poroto/planner/internal/infrastructure/firestore/entity"
)

const (
	collectionPlaceSearchResults = "google_place_api_search_results"
	collectionReviews            = "reviews"
	subCollectionPhotos          = "photos"
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
	collection := p.collection(planCandidateId)

	snapshots, err := collection.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("error while getting place search results: %v", err)
	}

	var places []models.GooglePlace
	for _, snapshot := range snapshots {
		var placeEntity googleplaces.Place
		if err = snapshot.DataTo(&placeEntity); err != nil {
			return nil, fmt.Errorf("error while converting snapshot to place search result entity: %v", err)
		}

		photos, err := p.fetchPhotos(ctx, planCandidateId, placeEntity.PlaceID)
		if err != nil {
			return nil, fmt.Errorf("error while fetching photos: %v", err)
		}

		reviews, err := p.fetchReviews(ctx, planCandidateId, placeEntity.PlaceID)
		if err != nil {
			return nil, fmt.Errorf("error while fetching reviews: %v", err)
		}

		places = append(places, factory.GooglePlaceFromPlaceEntity(
			placeEntity,
			photos,
			reviews,
		))
	}

	return places, nil
}

// TODO: 個々の画像をIDで区別できるようにする
func (p GooglePlaceSearchResultRepository) saveImagesIfNotExist(ctx context.Context, planCandidateId string, googlePlaceId string, images []models.Image) error {
	subCollectionImages := p.subCollectionPhotos(planCandidateId, googlePlaceId)

	snapshots, err := subCollectionImages.Limit(1).Documents(ctx).GetAll()
	if err != nil {
		return fmt.Errorf("error while getting images: %v", err)
	}

	if len(snapshots) > 0 {
		// すでに画像が保存されている場合は何もしない
		return fmt.Errorf("images already exist")
	}

	for _, image := range images {
		if _, err := subCollectionImages.NewDoc().Set(ctx, entity.ToImageEntity(image)); err != nil {
			return fmt.Errorf("error while saving image: %v", err)
		}
	}

	return nil
}

func (p GooglePlaceSearchResultRepository) saveReviewsIfNotExist(ctx context.Context, planCandidateId string, googlePlaceId string, reviews []models.GooglePlaceReview) error {
	subCollectionReviews := p.subCollectionReviews(planCandidateId, googlePlaceId)

	snapshots, err := subCollectionReviews.Limit(1).Documents(ctx).GetAll()
	if err != nil {
		return fmt.Errorf("error while getting reviews: %v", err)
	}

	if len(snapshots) > 0 {
		// すでにレビューが保存されている場合は何もしない
		return fmt.Errorf("reviews already exist")
	}

	for _, review := range reviews {
		if _, err := subCollectionReviews.NewDoc().Set(ctx, entity.ToGooglePlaceReviewEntity(review)); err != nil {
			return fmt.Errorf("error while saving review: %v", err)
		}
	}

	return nil
}

func (p GooglePlaceSearchResultRepository) deleteByPlanCandidateIdTx(tx *firestore.Transaction, planCandidateId string) error {
	collection := p.collection(planCandidateId)

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

func (p GooglePlaceSearchResultRepository) fetchPhotos(ctx context.Context, planCandidateId string, googlePlaceId string) ([]entity.ImageEntity, error) {
	subCollectionPhotos := p.subCollectionPhotos(planCandidateId, googlePlaceId)
	photosSnapshots, err := subCollectionPhotos.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("error while getting photos: %v", err)
	}

	var photos []entity.ImageEntity
	for _, photoSnapshot := range photosSnapshots {
		var photoEntity entity.ImageEntity
		if err = photoSnapshot.DataTo(&photoEntity); err != nil {
			return nil, fmt.Errorf("error while converting snapshot to photo entity: %v", err)
		}
		photos = append(photos, photoEntity)
	}

	return photos, nil
}

func (p GooglePlaceSearchResultRepository) fetchReviews(ctx context.Context, planCandidateId string, googlePlaceId string) ([]entity.GooglePlaceReviewEntity, error) {
	subCollectionReviews := p.subCollectionReviews(planCandidateId, googlePlaceId)
	reviewsSnapshots, err := subCollectionReviews.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("error while getting reviews: %v", err)
	}

	var reviews []entity.GooglePlaceReviewEntity
	for _, reviewSnapshot := range reviewsSnapshots {
		var reviewEntity entity.GooglePlaceReviewEntity
		if err = reviewSnapshot.DataTo(&reviewEntity); err != nil {
			return nil, fmt.Errorf("error while converting snapshot to review entity: %v", err)
		}
		reviews = append(reviews, reviewEntity)
	}

	return reviews, nil
}

func (p GooglePlaceSearchResultRepository) collection(planCandidateId string) *firestore.CollectionRef {
	return p.client.Collection(collectionPlanCandidates).Doc(planCandidateId).Collection(collectionPlaceSearchResults)
}

func (p GooglePlaceSearchResultRepository) doc(planCandidateId string, googlePlaceId string) *firestore.DocumentRef {
	return p.collection(planCandidateId).Doc(googlePlaceId)
}

func (p GooglePlaceSearchResultRepository) subCollectionPhotos(planCandidateId string, googlePlaceId string) *firestore.CollectionRef {
	return p.doc(planCandidateId, googlePlaceId).Collection(subCollectionPhotos)
}

func (p GooglePlaceSearchResultRepository) subCollectionReviews(planCandidateId string, googlePlaceId string) *firestore.CollectionRef {
	return p.doc(planCandidateId, googlePlaceId).Collection(collectionReviews)
}
