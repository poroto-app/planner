package plancandidate

import (
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/repository"
	"poroto.app/poroto/planner/internal/domain/services/place"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
	"poroto.app/poroto/planner/internal/infrastructure/rdb"
)

type Service struct {
	placesApi               places.PlacesApi
	planCandidateRepository repository.PlanCandidateRepository
	placeService            place.Service
	logger                  *zap.Logger
}

func NewService(db *sql.DB) (*Service, error) {
	placesApi, err := places.NewPlacesApi()
	if err != nil {
		return nil, fmt.Errorf("error while initializing places api: %v", err)
	}

	planCandidateRepository, err := rdb.NewPlanCandidateRepository(db)
	if err != nil {
		return nil, fmt.Errorf("error while initializing plan candidate repository: %v", err)
	}

	placeService, err := place.NewPlaceService(db)
	if err != nil {
		return nil, fmt.Errorf("error while initializing place service: %v", err)
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
		placeService:            *placeService,
		logger:                  logger,
	}, nil
}
