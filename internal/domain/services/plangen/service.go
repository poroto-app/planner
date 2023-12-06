package plangen

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"os"
	"poroto.app/poroto/planner/internal/domain/repository"
	"poroto.app/poroto/planner/internal/domain/services/place"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
	"poroto.app/poroto/planner/internal/infrastructure/api/openai"
	"poroto.app/poroto/planner/internal/infrastructure/firestore"
)

type Service struct {
	placesApi                  places.PlacesApi
	placeService               place.Service
	planCandidateRepository    repository.PlanCandidateRepository
	openaiChatCompletionClient openai.ChatCompletionClient
	logger                     *zap.Logger
}

func NewService(ctx context.Context) (*Service, error) {
	placesApi, err := places.NewPlacesApi()
	if err != nil {
		return nil, fmt.Errorf("error while initizalizing places api: %v", err)
	}

	placeService, err := place.NewPlaceService(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while initializing place service: %v", err)
	}

	planCandidateRepository, err := firestore.NewPlanCandidateRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while initializing plan candidate repository: %v", err)
	}

	openaiChatCompletionClient, err := openai.NewChatCompletionClient()
	if err != nil {
		return nil, fmt.Errorf("error while initializing openai chat completion client: %v", err)
	}

	logger, err := utils.NewLogger(utils.LoggerOption{
		Tag:        "PlanGenService",
		Production: os.Getenv("ENV") != "development",
	})
	if err != nil {
		return nil, fmt.Errorf("error while initializing logger: %v", err)
	}

	return &Service{
		placesApi:                  *placesApi,
		placeService:               *placeService,
		planCandidateRepository:    planCandidateRepository,
		openaiChatCompletionClient: *openaiChatCompletionClient,
		logger:                     logger,
	}, nil
}
