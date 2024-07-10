package main

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"poroto.app/poroto/planner/internal/env"
	"poroto.app/poroto/planner/internal/infrastructure/rdb"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
)

func main() {
	env.LoadEnv()

	db, err := rdb.InitDB(false)
	if err != nil {
		log.Fatalf("error while initializing db: %v", err)
	}

	ctx := context.Background()
	placePhotoSlice, err := generated.PlacePhotos().All(ctx, db)
	if err != nil {
		log.Fatalf("error while fetching place photos: %v", err)
	}
	placePhotoReferenceSlice := make(generated.PlacePhotoReferenceSlice, 0, len(placePhotoSlice))

	for _, placePhoto := range placePhotoSlice {
		placePhotoReference := &generated.PlacePhotoReference{
			ID:      uuid.New().String(),
			PlaceID: placePhoto.PlaceID,
			UserID:  placePhoto.UserID,
		}
		placePhotoReferenceSlice = append(placePhotoReferenceSlice, placePhotoReference)
	}

	if _, err := placePhotoReferenceSlice.InsertAll(ctx, db, boil.Infer()); err != nil {
		log.Fatalf("error while inserting place photo references: %v", err)
	}
}
