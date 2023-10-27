package places

import (
	"context"
	"fmt"
	"os"

	"googlemaps.github.io/maps"
)

type PlacesApi struct {
	apiKey     string
	mapsClient *maps.Client
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

	return &PlacesApi{
		apiKey:     apiKey,
		mapsClient: c,
	}, nil
}

type FindPlacesFromLocationRequest struct {
	Location    Location
	Radius      uint
	Language    string
	Type        *maps.PlaceType
	SearchCount int
}

func (r PlacesApi) FindPlacesFromLocation(ctx context.Context, req *FindPlacesFromLocationRequest) ([]Place, error) {
	var placeType maps.PlaceType
	if req.Type != nil {
		placeType = *req.Type
	}

	placeSearchResults, err := r.nearBySearch(ctx, &maps.NearbySearchRequest{
		Location: &maps.LatLng{
			Lat: req.Location.Latitude,
			Lng: req.Location.Longitude,
		},
		Radius:   req.Radius,
		Language: req.Language,
		Type:     placeType,
	}, req.SearchCount)
	if err != nil {
		return nil, err
	}

	// Getting places nearby
	var places []Place
	for _, place := range placeSearchResults {
		var photoReferences []string
		for _, photo := range place.Photos {
			photoReferences = append(photoReferences, photo.PhotoReference)
		}

		places = append(places, createPlace(
			place.PlaceID,
			place.Name,
			place.Types,
			place.Geometry,
			photoReferences,
			place.OpeningHours != nil && place.OpeningHours.OpenNow != nil && *place.OpeningHours.OpenNow,
			place.Rating,
			place.UserRatingsTotal,
			place.PriceLevel,
		))
	}

	return places, nil
}
