package place

import (
	"database/sql"
	"fmt"

	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/repository"
	"poroto.app/poroto/planner/internal/domain/services/placesearch"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/infrastructure/rdb"
)

type Service struct {
	placeSearchService      placesearch.Service
	placeRepository         repository.PlaceRepository
	planCandidateRepository repository.PlanCandidateRepository
	logger                  zap.Logger
}

func NewService(db *sql.DB) (*Service, error) {
	planCandidateRepository, err := rdb.NewPlanCandidateRepository(db)
	if err != nil {
		return nil, fmt.Errorf("error while initializing plan candidate repository: %v", err)
	}

	placeRepository, err := rdb.NewPlaceRepository(db)
	if err != nil {
		return nil, fmt.Errorf("error while initializing place repository: %v", err)
	}

	placeSearchService, err := placesearch.NewPlaceSearchService(db)
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
		placeRepository:         placeRepository,
		planCandidateRepository: planCandidateRepository,
		logger:                  *logger,
	}, nil
}
