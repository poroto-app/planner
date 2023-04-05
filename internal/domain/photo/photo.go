package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"poroto.app/poroto/planner/internal/infrastructure/api/google"
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
	api := google.NewPlacesApi()
	places, _ := api.FindPlacesFromLocation(context.Background(), &google.FindPlacesFromLocationRequest{
		Location: google.Location{Latitude: 35.5689, Longitude: 139.3952},
		Radius:   1000,
	})

	placeDetails, err := api.GetPlaceDetailsFromPlaces(context.Background(), places)
	if err != nil {
		log.Fatal(err)
	}

	for _, placeDetail := range placeDetails {
		log.Printf("%s has %v photos", placeDetail.Place.Name, len(placeDetail.Photos))
	}

}
