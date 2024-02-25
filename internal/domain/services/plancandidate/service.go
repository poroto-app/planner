package plancandidate

import (
	"context"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/repository"
	"poroto.app/poroto/planner/internal/domain/services/placesearch"
	"poroto.app/poroto/planner/internal/domain/services/user"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
	"poroto.app/poroto/planner/internal/infrastructure/rdb"
)

type Service struct {
	placesApi               places.PlacesApi
	placeRepository         repository.PlaceRepository
	planCandidateRepository repository.PlanCandidateRepository
	userService             user.Service
	placeSearchService      placesearch.Service
	logger                  *zap.Logger
}

func NewService(ctx context.Context, db *sql.DB) (*Service, error) {
	placesApi, err := places.NewPlacesApi()
	if err != nil {
		return nil, fmt.Errorf("error while initializing places api: %v", err)
	}

	placeRepository, err := rdb.NewPlaceRepository(db)
	if err != nil {
		return nil, fmt.Errorf("error while initializing place repository: %v", err)
	}

	planCandidateRepository, err := rdb.NewPlanCandidateRepository(db)
	if err != nil {
		return nil, fmt.Errorf("error while initializing plan candidate repository: %v", err)
	}

	userService, err := user.NewService(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("error while initializing user service: %v", err)
	}

	placeSearchService, err := placesearch.NewPlaceSearchService(db)
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
		placeRepository:         placeRepository,
		planCandidateRepository: planCandidateRepository,
		userService:             *userService,
		placeSearchService:      *placeSearchService,
		logger:                  logger,
	}, nil
}
