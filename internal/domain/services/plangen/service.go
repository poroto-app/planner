package plangen

import (
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/repository"
	"poroto.app/poroto/planner/internal/domain/services/placesearch"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
	"poroto.app/poroto/planner/internal/infrastructure/api/openai"
	"poroto.app/poroto/planner/internal/infrastructure/rdb"
)

type Service struct {
	placesApi                  places.PlacesApi
	placeSearchService         placesearch.Service
	placeRepository            repository.PlaceRepository
	planCandidateRepository    repository.PlanCandidateRepository
	openaiChatCompletionClient openai.ChatCompletionClient
	logger                     *zap.Logger
}

func NewService(db *sql.DB) (*Service, error) {
	placesApi, err := places.NewPlacesApi()
	if err != nil {
		return nil, fmt.Errorf("error while initizalizing Places api: %v", err)
	}

	placeSearchService, err := placesearch.NewPlaceSearchService(db)
	if err != nil {
		return nil, fmt.Errorf("error while initializing place search service: %v", err)
	}

	placeRepository, err := rdb.NewPlaceRepository(db)
	if err != nil {
		return nil, fmt.Errorf("error while initializing place repository: %v", err)
	}

	planCandidateRepository, err := rdb.NewPlanCandidateRepository(db)
	if err != nil {
		return nil, fmt.Errorf("error while initializing plan candidate repository: %v", err)
	}

	openaiChatCompletionClient, err := openai.NewChatCompletionClient()
	if err != nil {
		return nil, fmt.Errorf("error while initializing openai chat completion client: %v", err)
	}

	logger, err := utils.NewLogger(utils.LoggerOption{
		Tag: "PlanGenService",
	})
	if err != nil {
		return nil, fmt.Errorf("error while initializing logger: %v", err)
	}

	return &Service{
		placesApi:                  *placesApi,
		placeSearchService:         *placeSearchService,
		placeRepository:            *placeRepository,
		planCandidateRepository:    planCandidateRepository,
		openaiChatCompletionClient: *openaiChatCompletionClient,
		logger:                     logger,
	}, nil
}
