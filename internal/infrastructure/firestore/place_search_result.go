package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"google.golang.org/api/option"
	"os"
	"poroto.app/poroto/planner/internal/domain/models"
	google_places "poroto.app/poroto/planner/internal/infrastructure/api/google/places"
	"poroto.app/poroto/planner/internal/infrastructure/firestore/entity"
)

const (
	collectionPlaceSearchResults = "google_place_api_search_results"
	subCollectionPhotos          = "photos"
)

type PlaceSearchResultRepository struct {
	client *firestore.Client
}

func NewPlaceSearchResultRepository(ctx context.Context) (*PlaceSearchResultRepository, error) {
	var options []option.ClientOption
	if os.Getenv("GCP_CREDENTIAL_FILE_PATH") != "" {
		options = append(options, option.WithCredentialsFile(os.Getenv("GCP_CREDENTIAL_FILE_PATH")))
	}

	client, err := firestore.NewClient(ctx, os.Getenv("GCP_PROJECT_ID"), options...)
	if err != nil {
		return nil, fmt.Errorf("error while initializing firestore client: %v", err)
	}

	return &PlaceSearchResultRepository{
		client: client,
	}, nil
}

func (p PlaceSearchResultRepository) Save(ctx context.Context, planCandidateId string, places []google_places.Place) error {
	if err := p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		for _, place := range places {
			doc := p.doc(planCandidateId, place.PlaceID)
			if _, err := doc.Set(ctx, place); err != nil {
				return fmt.Errorf("error while saving place search result: %v", err)
			}
		}
		return nil
	}, firestore.MaxAttempts(3)); err != nil {
		return fmt.Errorf("error while saving place search results: %v", err)
	}

	return nil
}

func (p PlaceSearchResultRepository) Find(ctx context.Context, planCandidateId string) ([]google_places.Place, error) {
	collection := p.collection(planCandidateId)

	snapshots, err := collection.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("error while getting place search results: %v", err)
	}

	var places []google_places.Place
	for _, snapshot := range snapshots {
		var placeEntity google_places.Place
		if err = snapshot.DataTo(&placeEntity); err != nil {
			return nil, fmt.Errorf("error while converting snapshot to place search result entity: %v", err)
		}

		places = append(places, placeEntity)
	}

	return places, nil
}

func (p PlaceSearchResultRepository) SaveImage(ctx context.Context, planCandidateId string, placeId string, image models.Image) error {
	subCollectionImages := p.subCollectionPhotos(planCandidateId, placeId)
	if _, err := subCollectionImages.NewDoc().Set(ctx, entity.ToImageEntity(image)); err != nil {
		return fmt.Errorf("error while saving image: %v", err)
	}

	return nil
}

func (p PlaceSearchResultRepository) DeleteAll(ctx context.Context, planCandidateIds []string) error {
	for _, planCandidateId := range planCandidateIds {
		if err := p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
			collection := p.collection(planCandidateId)

			// collection内のドキュメントをすべて削除
			snapshots, err := collection.Documents(ctx).GetAll()
			if err != nil {
				return fmt.Errorf("error while getting place search results: %v", err)
			}

			for _, snapshot := range snapshots {
				if _, err := snapshot.Ref.Delete(ctx); err != nil {
					return fmt.Errorf("error while deleting place search result: %v", err)
				}
			}

			return nil
		}, firestore.MaxAttempts(3)); err != nil {
			return fmt.Errorf("error while deleting place search results: %v", err)
		}
	}

	return nil
}

func (p PlaceSearchResultRepository) collection(planCandidateId string) *firestore.CollectionRef {
	return p.client.Collection(collectionPlanCandidates).Doc(planCandidateId).Collection(collectionPlaceSearchResults)
}

func (p PlaceSearchResultRepository) doc(planCandidateId string, placeId string) *firestore.DocumentRef {
	return p.collection(planCandidateId).Doc(placeId)
}

func (p PlaceSearchResultRepository) subCollectionPhotos(planCandidateId string, placeId string) *firestore.CollectionRef {
	return p.doc(planCandidateId, placeId).Collection(subCollectionPhotos)
}
