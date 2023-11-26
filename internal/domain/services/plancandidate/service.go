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
	placesApi                      places.PlacesApi
	planCandidateRepository        repository.PlanCandidateRepository
	placeInPlanCandidateRepository repository.PlaceInPlanCandidateRepository
	placeRepository                repository.PlaceRepository
	placeService                   place.Service
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

	placeInPlanCandidateRepository, err := firestore.NewPlaceInPlanCandidateRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while initializing place in plan candidate repository: %v", err)
	}

	placeRepository, err := firestore.NewPlaceRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while initializing place repository: %v", err)
	}

	placeService, err := place.NewPlaceService(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while initializing place service: %v", err)
	}

	return &Service{
		placesApi:                      *placesApi,
		planCandidateRepository:        planCandidateRepository,
		placeInPlanCandidateRepository: placeInPlanCandidateRepository,
		placeRepository:                placeRepository,
		placeService:                   *placeService,
	}, nil
}
