package places

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"poroto.app/poroto/planner/internal/domain/utils"

	"googlemaps.github.io/maps"
)

type PlacesApi struct {
	apiKey     string
	mapsClient *maps.Client
	logger     *zap.Logger
}

func NewPlacesApi() (*PlacesApi, error) {
	apiKey := os.Getenv("GOOGLE_PLACES_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("env variable GOOGLE_PLACES_API_KEY is not set")
	}

	opt := maps.WithAPIKey(apiKey)
	c, err := maps.NewClient(opt)
	if err != nil {
		return nil, fmt.Errorf("error while initializing maps api client: %v", err)
	}

	logger, err := utils.NewLogger(utils.LoggerOption{
		Tag: "PlacesApi",
	})
	if err != nil {
		return nil, fmt.Errorf("error while initializing logger: %v", err)
	}

	return &PlacesApi{
		apiKey:     apiKey,
		mapsClient: c,
		logger:     logger,
	}, nil
}
