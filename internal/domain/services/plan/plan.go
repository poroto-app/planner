package plan

import (
	"context"
	"fmt"

	"poroto.app/poroto/planner/internal/domain/repository"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
	"poroto.app/poroto/planner/internal/infrastructure/firestore"
)

type PlanService struct {
	placesApi               places.PlacesApi
	planRepository          repository.PlanRepository
	planCandidateRepository repository.PlanCandidateRepository
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

	return &PlanService{
		placesApi:               *placesApi,
		planRepository:          planRepository,
		planCandidateRepository: planCandidateRepository,
	}, err
}
