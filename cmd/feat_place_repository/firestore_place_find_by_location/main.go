package main

import (
	"context"
	"log"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/env"
	repo "poroto.app/poroto/planner/internal/infrastructure/firestore"
)

func main() {
	env.LoadEnv()
	ctx := context.Background()

	placeRepository, err := repo.NewPlaceRepository(ctx)
	if err != nil {
		log.Fatalf("failed to initialize place repository: %v", err)
	}

	location := models.GeoLocation{
		Latitude:  35.5684909,
		Longitude: 139.3952879,
	}

	places, err := placeRepository.FindByLocation(ctx, location)
	if err != nil {
		log.Fatalf("failed to find places by location: %v", err)
	}

	for _, place := range places {
		log.Printf("{ id: %s, name: %s, distance: %f }", place.Id, place.Name, place.Location.DistanceInMeter(location))
	}
}
