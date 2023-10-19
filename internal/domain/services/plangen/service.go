package plangen

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/repository"
	"poroto.app/poroto/planner/internal/domain/services/place"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
	"poroto.app/poroto/planner/internal/infrastructure/api/openai"
	"poroto.app/poroto/planner/internal/infrastructure/firestore"
)

type Service struct {
	placesApi                   places.PlacesApi
	placeService                place.Service
	planCandidateRepository     repository.PlanCandidateRepository
	placeSearchResultRepository repository.GooglePlaceSearchResultRepository
	openaiChatCompletionClient  openai.ChatCompletionClient
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

	placeSearchResultRepository, err := firestore.NewGooglePlaceSearchResultRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while initializing place search result repository: %v", err)
	}

	openaiChatCompletionClient, err := openai.NewChatCompletionClient()
	if err != nil {
		return nil, fmt.Errorf("error while initializing openai chat completion client: %v", err)
	}

	return &Service{
		placesApi:                   *placesApi,
		placeService:                *placeService,
		planCandidateRepository:     planCandidateRepository,
		placeSearchResultRepository: placeSearchResultRepository,
		openaiChatCompletionClient:  *openaiChatCompletionClient,
	}, nil
}
