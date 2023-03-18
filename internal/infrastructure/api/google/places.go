package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"googlemaps.github.io/maps"
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

// Check whether including or not
func contains(target, words []string) bool {
	for _, element := range target {
		for _, word := range words {
			if element == word {
				return true
			}
		}
	}
	return false
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

	for _, place := range res.Results {
		log.Println(place.Name, place.Types)
	}
}
