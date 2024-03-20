package place

import (
	"context"
	"database/sql"
	"fmt"

	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/repository"
	"poroto.app/poroto/planner/internal/domain/services/placesearch"
	"poroto.app/poroto/planner/internal/domain/services/user"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/infrastructure/rdb"
)

type Service struct {
	placeSearchService      placesearch.Service
	planCandidateRepository repository.PlanCandidateRepository
	planRepository          repository.PlanRepository
	placeRepository         repository.PlaceRepository
	userService             *user.Service
	logger                  zap.Logger
}

func NewService(ctx context.Context, db *sql.DB) (*Service, error) {
	placeSearchService, err := placesearch.NewPlaceSearchService(db)
	if err != nil {
		return nil, fmt.Errorf("error while initializing place search service: %v", err)
	}

	planCandidateRepository, err := rdb.NewPlanCandidateRepository(db)
	if err != nil {
		return nil, fmt.Errorf("error while initializing plan candidate repository: %v", err)
	}

	planRepository, err := rdb.NewPlanRepository(db)
	if err != nil {
		return nil, fmt.Errorf("error while initializing plan repository: %v", err)
	}

	placeRepository, err := rdb.NewPlaceRepository(db)
	if err != nil {
		return nil, fmt.Errorf("error while initializing place repository: %v", err)
	}

	userService, err := user.NewService(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("error while initializing user service: %v", err)
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
		planRepository:          planRepository,
		placeRepository:         placeRepository,
		userService:             userService,
		logger:                  *logger,
	}, nil
}
