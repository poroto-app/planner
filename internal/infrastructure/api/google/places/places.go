package places

import (
	"context"
	"fmt"
	"os"

	"googlemaps.github.io/maps"
	"poroto.app/poroto/planner/internal/domain/models"
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

type Place struct {
	PlaceID         string   `firestore:"place_id"`
	Name            string   `firestore:"name"`
	Types           []string `firestore:"types"`
	Location        Location `firestore:"location"`
	PhotoReferences []string `firestore:"photo_references"`
	OpenNow         bool     `firestore:"open_now"`
	Rating          float32  `firestore:"rating"`
}

type Location struct {
	Latitude  float64 `firestore:"latitude"`
	Longitude float64 `firestore:"longitude"`
}

func (r Location) ToGeoLocation() models.GeoLocation {
	return models.GeoLocation{
		Latitude:  r.Latitude,
		Longitude: r.Longitude,
	}
}

type FindPlacesFromLocationRequest struct {
	Location Location
	Radius   uint
	Language string
}

func (r PlacesApi) FindPlacesFromLocation(ctx context.Context, req *FindPlacesFromLocationRequest) ([]Place, error) {
	placeSearchResults, err := r.nearBySearch(ctx, &maps.NearbySearchRequest{
		Location: &maps.LatLng{
			Lat: req.Location.Latitude,
			Lng: req.Location.Longitude,
		},
		Radius:   req.Radius,
		Language: req.Language,
	})
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

		places = append(places, Place{
			PlaceID: place.PlaceID,
			Name:    place.Name,
			Types:   place.Types,
			Location: Location{
				Latitude:  place.Geometry.Location.Lat,
				Longitude: place.Geometry.Location.Lng,
			},
			OpenNow:         place.OpeningHours != nil && place.OpeningHours.OpenNow != nil && *place.OpeningHours.OpenNow,
			PhotoReferences: photoReferences,
			Rating:          place.Rating,
		})
	}

	return places, nil
}
