package main

import (
	"context"
	"log"
	"os"

	"googlemaps.github.io/maps"
	"poroto.app/poroto/planner/internal/domain/array"
)

type GooglePlacesApi struct {
	apiKey string
}

func NewGooglePlacesApi() GooglePlacesApi {
	apiKey := os.Getenv("GOOGLE_PLACES_API_KEY")
	if apiKey == "" {
		log.Fatalln("env variable GOOGLE_PLACES_API_KEY is not set")
	}
	return GooglePlacesApi{
		apiKey: apiKey,
	}
}

func (r GooglePlacesApi) FindPlacesFromLocation(latitude float64, longitude float64) {
	googlePlacesApi := NewGooglePlacesApi()
	opt := maps.WithAPIKey(googlePlacesApi.apiKey)
	c, err := maps.NewClient(opt)
	if err != nil {
		log.Fatalln(err)
	}

	res, err := c.NearbySearch(context.Background(), &maps.NearbySearchRequest{
		Location: &maps.LatLng{
			Lat: latitude,
			Lng: longitude,
		},
		Radius: 1000,
	})
	if err != nil {
		log.Fatalln(err)
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
	for _, place := range res.Results {
		// To extract places
		if !array.HasIntersection(place.Types, categoriesSlice) {
			continue
		}

		if place.OpeningHours.OpenNow == nil {
			continue
		}

		if *place.OpeningHours.OpenNow {
			log.Println(place.Name, "[Open Now]")
		}
	}
}
