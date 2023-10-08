package firestore

import (
	"context"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
	"poroto.app/poroto/planner/internal/infrastructure/firestore/entity"
)

const (
	collectionPlaceSearchResults = "google_place_api_search_results"
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

func (p PlaceSearchResultRepository) Save(ctx context.Context, planCandidateId string, places []places.Place) error {
	placeSearchResultEntity := entity.PlaceSearchResultEntity{
		PlanCandidateId: planCandidateId,
		Places:          places,
		UpdatedAt:       time.Now(),
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

func (p PlaceSearchResultRepository) DeleteAll(ctx context.Context, planCandidateIds []string) error {
	if err := p.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		for _, planCandidateId := range planCandidateIds {
			doc := p.doc(planCandidateId)
			if err := tx.Delete(doc); err != nil {
				return fmt.Errorf("error while deleting place search result: %v", err)
			}
		}
		return nil
	}, firestore.MaxAttempts(3)); err != nil {
		return fmt.Errorf("error while deleting place search results: %v", err)
	}
	return nil
}

func (p PlaceSearchResultRepository) doc(planCandidateId string) *firestore.DocumentRef {
	return p.client.Collection(collectionPlanCandidates).Doc(planCandidateId).Collection(collectionPlaceSearchResults).Doc(planCandidateId)
}
