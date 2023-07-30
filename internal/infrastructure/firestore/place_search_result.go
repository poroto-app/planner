package firestore

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
	"poroto.app/poroto/planner/internal/infrastructure/firestore/entity"
)

const (
	collectionPlaceSearchResults = "place_search_results"
)

type PlaceSearchResultRepository struct {
	client *firestore.Client
}

func NewPlaceSearchResultRepository(ctx context.Context) (*PlaceSearchResultRepository, error) {
	client, err := firestore.NewClient(ctx, os.Getenv("GCP_PROJECT_ID"), option.WithCredentialsFile("secrets/google-credential.json"))
	if err != nil {
		return nil, fmt.Errorf("error while initializing firestore client: %v", err)
	}

	return &PlaceSearchResultRepository{
		client: client,
	}, nil
}

func (p PlaceSearchResultRepository) Save(ctx context.Context, planCandidateId string, places []places.Place) error {
	placeSearchResultEntity := entity.PlaceSearchResultEntity{
		PlanCandidateId: planCandidateId,
		Places:          places,
	}

	doc := p.doc(planCandidateId)
	if _, err := doc.Set(ctx, placeSearchResultEntity); err != nil {
		return fmt.Errorf("error while saving place search result: %v", err)
	}

	return nil
}

func (p PlaceSearchResultRepository) Find(ctx context.Context, planCandidateId string) ([]places.Place, error) {
	doc := p.doc(planCandidateId)
	snapshot, err := doc.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("error while finding place search result: %v", err)
	}

	var placeSearchResultEntity entity.PlaceSearchResultEntity
	if err = snapshot.DataTo(&placeSearchResultEntity); err != nil {
		return nil, fmt.Errorf("error while converting snapshot to place search result entity: %v", err)
	}

	return placeSearchResultEntity.Places, nil
}

func (p PlaceSearchResultRepository) collection() *firestore.CollectionRef {
	return p.client.Collection(collectionPlaceSearchResults)
}

func (p PlaceSearchResultRepository) doc(planCandidateId string) *firestore.DocumentRef {
	return p.client.Collection(collectionPlaceSearchResults).Doc(planCandidateId)
}
