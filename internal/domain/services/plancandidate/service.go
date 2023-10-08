package plancandidate

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/services/place"
	"poroto.app/poroto/planner/internal/domain/services/plangen"

	"poroto.app/poroto/planner/internal/domain/repository"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
	"poroto.app/poroto/planner/internal/infrastructure/firestore"
)

type Service struct {
	placesApi                   places.PlacesApi
	planCandidateRepository     repository.PlanCandidateRepository
	placeSearchResultRepository repository.PlaceSearchResultRepository
	planGeneratorService        plangen.Service
	placeService                place.Service
}

func NewService(ctx context.Context) (*Service, error) {
	placesApi, err := places.NewPlacesApi()
	if err != nil {
		return nil, fmt.Errorf("error while initializing places api: %v", err)
	}

	planCandidateRepository, err := firestore.NewPlanCandidateRepository(ctx)
	if err != nil {
		return nil, err
	}

	placeSearchResultRepository, err := firestore.NewPlaceSearchResultRepository(ctx)
	if err != nil {
		return nil, err
	}

	placeService, err := place.NewPlaceService(ctx)
	if err != nil {
		return nil, err
	}

	planGeneratorService, err := plangen.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while initializing plan generator service: %v", err)
	}

	return &Service{
		placesApi:                   *placesApi,
		planCandidateRepository:     planCandidateRepository,
		placeSearchResultRepository: placeSearchResultRepository,
		planGeneratorService:        *planGeneratorService,
		placeService:                *placeService,
	}, nil
}
