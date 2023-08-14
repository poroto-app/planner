package plangenerator

import (
	"context"
	"fmt"

	"poroto.app/poroto/planner/internal/domain/repository"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
	"poroto.app/poroto/planner/internal/infrastructure/api/openai"
	"poroto.app/poroto/planner/internal/infrastructure/firestore"
)

type Service struct {
	placesApi                   places.PlacesApi
	openaiChatCompletionClient  openai.ChatCompletionClient
	placeSearchResultRepository repository.PlaceSearchResultRepository
}

func NewService(ctx context.Context) (*Service, error) {
	placesApi, err := places.NewPlacesApi()
	if err != nil {
		return nil, fmt.Errorf("error while initizalizing places api: %v", err)
	}

	openaiChatCompletionClient, err := openai.NewChatCompletionClient()
	if err != nil {
		return nil, fmt.Errorf("error while initializing openai chat completion client: %v", err)
	}

	planSearchResultRepository, err := firestore.NewPlaceSearchResultRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while initializing plan search result repository: %v", err)
	}

	return &Service{
		placesApi:                   *placesApi,
		openaiChatCompletionClient:  *openaiChatCompletionClient,
		placeSearchResultRepository: planSearchResultRepository,
	}, nil
}
