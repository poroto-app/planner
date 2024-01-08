package place

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/repository"
	"poroto.app/poroto/planner/internal/domain/services/placesearch"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/infrastructure/firestore"
)

type Service struct {
	placeSearchService      placesearch.Service
	planCandidateRepository repository.PlanCandidateRepository
	logger                  zap.Logger
}

func NewService(ctx context.Context) (*Service, error) {
	planCandidateRepository, err := firestore.NewPlanCandidateRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while initializing plan candidate repository: %v", err)
	}

	placeSearchService, err := placesearch.NewPlaceSearchService(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while initializing place search service: %v", err)
	}

	logger, err := utils.NewLogger(utils.LoggerOption{
		Tag: "PlaceService",
	})
	if err != nil {
		return nil, fmt.Errorf("error while initializing logger: %v", err)
	}

	return &Service{
		placeSearchService:      *placeSearchService,
		planCandidateRepository: planCandidateRepository,
		logger:                  *logger,
	}, nil
}
