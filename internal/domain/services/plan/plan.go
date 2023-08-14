package plan

import (
	"context"
	"fmt"

	"poroto.app/poroto/planner/internal/domain/repository"
	"poroto.app/poroto/planner/internal/domain/services/plangenerator"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
	"poroto.app/poroto/planner/internal/infrastructure/firestore"
)

type PlanService struct {
	placesApi                   places.PlacesApi
	planRepository              repository.PlanRepository
	planCandidateRepository     repository.PlanCandidateRepository
	placeSearchResultRepository repository.PlaceSearchResultRepository
	planGeneratorService        plangenerator.Service
}

func NewPlanService(ctx context.Context) (*PlanService, error) {
	placesApi, err := places.NewPlacesApi()
	if err != nil {
		return nil, fmt.Errorf("error while initizalizing places api: %v", err)
	}

	planRepository, err := firestore.NewPlanRepository(ctx)
	if err != nil {
		return nil, err
	}

	planCandidateRepository, err := firestore.NewPlanCandidateRepository(ctx)
	if err != nil {
		return nil, err
	}

	placeSearchResultRepository, err := firestore.NewPlaceSearchResultRepository(ctx)
	if err != nil {
		return nil, err
	}

	planGeneratorService, err := plangenerator.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while initializing plan generator service: %v", err)
	}

	return &PlanService{
		placesApi:                   *placesApi,
		planRepository:              planRepository,
		planCandidateRepository:     planCandidateRepository,
		placeSearchResultRepository: placeSearchResultRepository,
		planGeneratorService:        *planGeneratorService,
	}, err
}
