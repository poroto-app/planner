package plancandidate

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/repository"
	"poroto.app/poroto/planner/internal/domain/services/place"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
	"poroto.app/poroto/planner/internal/infrastructure/firestore"
)

type Service struct {
	placesApi                   places.PlacesApi
	planCandidateRepository     repository.PlanCandidateRepository
	placeSearchResultRepository repository.GooglePlaceSearchResultRepository
	placeService                place.Service
}

func NewService(ctx context.Context) (*Service, error) {
	placesApi, err := places.NewPlacesApi()
	if err != nil {
		return nil, fmt.Errorf("error while initializing places api: %v", err)
	}

	planCandidateRepository, err := firestore.NewPlanCandidateRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while initializing plan candidate repository: %v", err)
	}

	placeSearchResultRepository, err := firestore.NewGooglePlaceSearchResultRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while initializing place search result repository: %v", err)
	}

	placeService, err := place.NewPlaceService(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while initializing place service: %v", err)
	}

	return &Service{
		placesApi:                   *placesApi,
		planCandidateRepository:     planCandidateRepository,
		placeSearchResultRepository: placeSearchResultRepository,
		placeService:                *placeService,
	}, nil
}
