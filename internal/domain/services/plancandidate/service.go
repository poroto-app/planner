package plancandidate

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/repository"
	"poroto.app/poroto/planner/internal/domain/services/placesearch"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
	"poroto.app/poroto/planner/internal/infrastructure/firestore"
)

type Service struct {
	placesApi               places.PlacesApi
	planCandidateRepository repository.PlanCandidateRepository
	placeSearchService      placesearch.Service
	logger                  *zap.Logger
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

	placeSearchService, err := placesearch.NewPlaceSearchService(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while initializing place search service: %v", err)
	}

	logger, err := utils.NewLogger(utils.LoggerOption{
		Tag: "PlanCandidateService",
	})
	if err != nil {
		return nil, fmt.Errorf("error while initializing logger: %v", err)
	}

	return &Service{
		placesApi:               *placesApi,
		planCandidateRepository: planCandidateRepository,
		placeSearchService:      *placeSearchService,
		logger:                  logger,
	}, nil
}
