package place

import (
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/repository"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
	"poroto.app/poroto/planner/internal/infrastructure/rdb"
)

type Service struct {
	placesApi               places.PlacesApi
	placeRepository         repository.PlaceRepository
	planCandidateRepository repository.PlanCandidateRepository
	logger                  *zap.Logger
}

func NewPlaceService(db *sql.DB) (*Service, error) {
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
		return nil, err
	}

	logger, err := utils.NewLogger(utils.LoggerOption{
		Tag: "PlaceService",
	})
	if err != nil {
		return nil, fmt.Errorf("error while initializing logger: %v", err)
	}

	return &Service{
		placesApi:               *placesApi,
		placeRepository:         *placeRepository,
		planCandidateRepository: planCandidateRepository,
		logger:                  logger,
	}, nil
}
