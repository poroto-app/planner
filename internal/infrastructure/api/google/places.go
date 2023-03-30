package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"googlemaps.github.io/maps"

	"poroto.app/poroto/planner/internal/domain/array"
)

func init() {
	env := os.Getenv("ENV")
	if "" == env {
		env = "development"
	}

	if err := godotenv.Load(".env.local"); err != nil {
		log.Fatalf("error while loading .env.local: %v", err)
	}

	if err := godotenv.Load(".env." + env); err != nil {
		log.Fatalf("error while loading .env.%s: %v", env, err)
	}
}

func main() {
	opt := maps.WithAPIKey(os.Getenv("GOOGLE_PLACES_API_KEY"))
	c, err := maps.NewClient(opt)
	if err != nil {
		log.Fatalln(err)
	}

	res, err := c.NearbySearch(context.Background(), &maps.NearbySearchRequest{
		Location: &maps.LatLng{
			Lat: 35.5689,
			Lng: 139.3952,
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
