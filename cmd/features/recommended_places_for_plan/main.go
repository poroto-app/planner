package main

import (
	"context"
	"database/sql"
	"flag"
	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"log"
	"poroto.app/poroto/planner/internal/env"
	"poroto.app/poroto/planner/internal/infrastructure/rdb"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
)

func init() {
	env.LoadEnv()
}

func main() {
	placeId := flag.String("place", "", "おすすめのとして追加したい場所のID")
	registerPlace := flag.Bool("register", false, "おすすめの場所を登録する")
	deletePlaceRecommendation := flag.Bool("delete", false, "おすすめの場所を削除する")

	flag.Parse()

	if placeId == nil || *placeId == "" {
		flag.PrintDefaults()
		return
	}

	db, err := rdb.InitDB(false)
	if err != nil {
		log.Fatalf("error while initializing db: %v", err)
	}

	place, err := fetchPlace(context.Background(), db, *placeId)
	if err != nil {
		log.Fatalf("error while fetching place: %v", err)
	}

	// おすすめの場所として追加
	if registerPlace != nil && *registerPlace {
		googlePlace := place.R.GooglePlaces
		if len(googlePlace) == 0 {
			log.Fatalf("google place not found")
		}

		// 写真が登録されていない場合は、おすすめの場所として登録できないようにする
		googlePlacePhotos := place.R.GooglePlaces[0].R.GetGooglePlacePhotos()
		placePhotos := place.R.PlacePhotos
		if len(googlePlacePhotos) == 0 && len(placePhotos) == 0 {
			log.Fatalf("place photos not found")
		}

		placeRecommendation, err := addRecommendPlace(context.Background(), db, *placeId)
		if err != nil {
			log.Fatalf("error while adding recommend place: %v", err)
		}
		log.Printf("place recommendation added:\nid: %s\nname: %s", placeRecommendation.ID, place.Name)
		return
	}
	// おすすめの場所を削除
	if deletePlaceRecommendation != nil && *deletePlaceRecommendation {
		if err := deleteRecommendPlace(context.Background(), db, *placeId); err != nil {
			log.Fatalf("error while deleting recommend place: %v", err)
		}
		log.Printf("place recommendation deleted:\nid: %s\nname: %s", *placeId, place.Name)
		return
	}

	log.Printf("place name: %s", place.Name)
}

func fetchPlace(ctx context.Context, db *sql.DB, placeId string) (*generated.Place, error) {
	place, err := generated.Places(
		generated.PlaceWhere.ID.EQ(placeId),
		qm.Load(generated.PlaceRels.PlacePhotos),
		qm.Load(generated.PlaceRels.GooglePlaces+"."+generated.GooglePlaceRels.GooglePlacePhotos),
	).One(ctx, db)
	if err != nil {
		return nil, err
	}
	return place, nil
}

func addRecommendPlace(ctx context.Context, db *sql.DB, placeId string) (*generated.PlaceRecommendation, error) {
	placeRecommendation := generated.PlaceRecommendation{
		ID:      uuid.New().String(),
		PlaceID: placeId,
	}

	if err := placeRecommendation.Insert(ctx, db, boil.Infer()); err != nil {
		return nil, err
	}

	return &placeRecommendation, nil
}

func deleteRecommendPlace(ctx context.Context, db *sql.DB, placeId string) error {
	_, err := generated.PlaceRecommendations(generated.PlaceRecommendationWhere.PlaceID.EQ(placeId)).DeleteAll(ctx, db)
	return err
}
