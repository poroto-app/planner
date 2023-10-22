package place

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/repository"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
	"poroto.app/poroto/planner/internal/infrastructure/firestore"
)

type Service struct {
	placesApi                      places.PlacesApi
	placeInPlanCandidateRepository repository.PlaceInPlanCandidateRepository
	placeSearchResultRepository    repository.GooglePlaceSearchResultRepository
	planCandidateRepository        repository.PlanCandidateRepository
}

func NewPlaceService(ctx context.Context) (*Service, error) {
	placesApi, err := places.NewPlacesApi()
	if err != nil {
		return nil, fmt.Errorf("error while initializing places api: %v", err)
	}

	placeSearchResultRepository, err := firestore.NewGooglePlaceSearchResultRepository(ctx)
	if err != nil {
		return nil, err
	}

	placeInPlanCandidateRepository, err := firestore.NewPlaceInPlanCandidateRepository(ctx)
	if err != nil {
		return nil, err
	}

	planCandidateRepository, err := firestore.NewPlanCandidateRepository(ctx)
	if err != nil {
		return nil, err
	}

	return &Service{
		placesApi:                      *placesApi,
		placeInPlanCandidateRepository: *placeInPlanCandidateRepository,
		placeSearchResultRepository:    placeSearchResultRepository,
		planCandidateRepository:        planCandidateRepository,
	}, nil
}
