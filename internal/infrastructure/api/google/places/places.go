package places

import (
	"context"
	"fmt"
	"os"

	"googlemaps.github.io/maps"
	"poroto.app/poroto/planner/internal/domain/array"
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
	PlaceID         string
	Name            string
	Types           []string
	Location        Location
	photoReferences []string
}

type Location struct {
	Latitude  float64
	Longitude float64
}

type FindPlacesFromLocationRequest struct {
	Location Location
	Radius   uint
}

func (r PlacesApi) FindPlacesFromLocation(ctx context.Context, req *FindPlacesFromLocationRequest) ([]Place, error) {
	res, err := r.mapsClient.NearbySearch(ctx, &maps.NearbySearchRequest{
		Location: &maps.LatLng{
			Lat: req.Location.Latitude,
			Lng: req.Location.Longitude,
		},
		Radius: req.Radius,
	})
	if err != nil {
		return nil, err
	}

	// Set objective place.Types
	var categories map[string][]string = make(map[string][]string)
	// Need initialization for ensure memory of map
	categories["amusements"] = []string{"amusement_park", "aquarium", "art_gallary", "museum"}
	categories["restaurants"] = []string{"bakery", "bar", "cafe", "food", "restaurant"}

	// Refactoring map to slice for hasIntersection
	var categoriesSlice []string
	for _, value := range categories {
		categoriesSlice = append(categoriesSlice, value...)
	}

	// Getting places nearby
	var places []Place
	for _, place := range res.Results {
		// To extract places
		// TODO: フィルタリングするカテゴリを `FindPlacesFromLocationRequest`で指定できるようにする
		if !array.HasIntersection(place.Types, categoriesSlice) {
			continue
		}

		// TODO: 現在時刻でフィルタリングするかを `FindPlacesFromLocationRequest`で指定できるようにする
		if place.OpeningHours == nil || place.OpeningHours.OpenNow == nil {
			continue
		}

		if !*place.OpeningHours.OpenNow {
			continue
		}

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
			photoReferences: photoReferences,
		})
	}

	return places, nil
}
