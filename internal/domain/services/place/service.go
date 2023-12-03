package place

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/repository"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
	"poroto.app/poroto/planner/internal/infrastructure/firestore"
)

type Service struct {
	placesApi               places.PlacesApi
	placeRepository         repository.PlaceRepository
	planCandidateRepository repository.PlanCandidateRepository
}

func NewPlaceService(ctx context.Context) (*Service, error) {
	placesApi, err := places.NewPlacesApi()
	if err != nil {
		return nil, fmt.Errorf("error while initializing places api: %v", err)
	}

	placeRepository, err := firestore.NewPlaceRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while initializing place repository: %v", err)
	}

	planCandidateRepository, err := firestore.NewPlanCandidateRepository(ctx)
	if err != nil {
		return nil, err
	}

	return &Service{
		placesApi:               *placesApi,
		placeRepository:         *placeRepository,
		planCandidateRepository: planCandidateRepository,
	}, nil
}
