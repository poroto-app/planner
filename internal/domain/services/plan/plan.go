package plan

import (
	"context"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/repository"
	"poroto.app/poroto/planner/internal/domain/services/user"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/infrastructure/rdb"
)

type Service struct {
	planRepository          repository.PlanRepository
	planCandidateRepository repository.PlanCandidateRepository
	placeRepository         repository.PlaceRepository
	userService             *user.Service
	logger                  *zap.Logger
}

func NewService(ctx context.Context, db *sql.DB) (*Service, error) {
	planRepository, err := rdb.NewPlanRepository(db)
	if err != nil {
		return nil, err
	}

	planCandidateRepository, err := rdb.NewPlanCandidateRepository(db)
	if err != nil {
		return nil, err
	}

	placeRepository, err := rdb.NewPlaceRepository(db)
	if err != nil {
		return nil, err
	}

	userService, err := user.NewService(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("error while initializing user service: %v", err)
	}

	logger, err := utils.NewLogger(utils.LoggerOption{
		Tag: "PlanService",
	})

	return &Service{
		planRepository:          planRepository,
		planCandidateRepository: planCandidateRepository,
		placeRepository:         placeRepository,
		userService:             userService,
		logger:                  logger,
	}, err
}
