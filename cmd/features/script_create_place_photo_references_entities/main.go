package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
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
	placePhotoSliceToUpdate, err := generated.PlacePhotos(
		qm.Where(fmt.Sprintf("%s is null", generated.PlacePhotoColumns.PlacePhotoReferenceID)),
	).All(ctx, db)
	if err != nil {
		log.Fatalf("error while fetching place photos: %v", err)
	}

	if len(placePhotoSliceToUpdate) == 0 {
		log.Println("no place photos to create place photo references")
		return
	}

	placePhotoReferenceSliceToSave := make(generated.PlacePhotoReferenceSlice, 0, len(placePhotoSliceToUpdate))

	for _, placePhoto := range placePhotoSliceToUpdate {

		placePhotoReferenceToSave := &generated.PlacePhotoReference{
			ID:      uuid.New().String(),
			PlaceID: placePhoto.PlaceID,
			UserID:  placePhoto.UserID,
		}

		placePhoto.PlacePhotoReferenceID = null.StringFrom(placePhotoReferenceToSave.ID)
		placePhotoReferenceSliceToSave = append(placePhotoReferenceSliceToSave, placePhotoReferenceToSave)
		log.Printf("place photo reference to save: %+v", placePhotoReferenceToSave)
	}

	if _, err := placePhotoReferenceSliceToSave.InsertAll(ctx, db, boil.Infer()); err != nil {
		log.Fatalf("error while inserting place photo references: %v", err)
	}

	for _, placePhoto := range placePhotoSliceToUpdate {
		if _, err := placePhoto.Update(ctx, db, boil.Infer()); err != nil {
			log.Fatalf("error while updating place photo: %v", err)
		}
	}
}
