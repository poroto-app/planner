package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"log"
	"os"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/env"
	"poroto.app/poroto/planner/internal/infrastructure/rdb"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
	"time"
)

func init() {
	os.Setenv("ENV", "development")
	env.LoadEnv()
}

func main() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?parseTime=true&loc=%s&tls=%v&interpolateParams=%v",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
		"Asia%2FTokyo",
		os.Getenv("ENV") != "development",
		true,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	boil.SetDB(db)
	boil.DebugMode = true

	ctx := context.Background()

	cleanup(ctx, db)

	// 事前に Place のデータを登録しておく
	placeRepository, err := rdb.NewPlaceRepository(db)
	if err != nil {
		log.Fatalf("failed to create place repository: %v", err)
	}
	testPlace, err := placeRepository.SavePlacesFromGooglePlace(ctx, models.GooglePlace{
		PlaceId: "test-google-place-id",
		Name:    "test-place",
	})
	if err != nil {
		log.Fatalf("failed to save places from google place: %v", err)
	}

	// 終了時に削除する
	defer func() {
		cleanup(ctx, db)
		if _, err := generated.GooglePlaces(generated.GooglePlaceWhere.GooglePlaceID.EQ(testPlace.Google.PlaceId)).DeleteAll(ctx, db); err != nil {
			log.Fatalf("failed to delete google places: %v", err)
		}
		if _, err := generated.Places(generated.PlaceWhere.ID.EQ(testPlace.Id)).DeleteAll(ctx, db); err != nil {
			log.Fatalf("failed to delete places: %v", err)
		}
	}()

	planCandidateRepository, err := rdb.NewPlanCandidateRepository(db)
	if err != nil {
		log.Fatalf("failed to create plan candidate repository: %v", err)
	}

	// PlanCandidateSet の作成
	if err := planCandidateRepository.Create(ctx, "test-plan-candidate-set", time.Now().Add(time.Hour*24)); err != nil {
		log.Fatalf("failed to create plan candidate set: %v", err)
	}

	// Places API による検索結果の保存
	if err := planCandidateRepository.AddSearchedPlacesForPlanCandidate(ctx, "test-plan-candidate-set", []string{testPlace.Id}); err != nil {
		log.Fatalf("failed to add searched places for plan candidate: %v", err)
	}

	// メタデータを追加
	if err := planCandidateRepository.UpdatePlanCandidateMetaData(ctx, "test-plan-candidate-set", models.PlanCandidateMetaData{
		CreatedBasedOnCurrentLocation: true,
		CategoriesPreferred:           &[]models.LocationCategory{models.CategoryRestaurant},
		CategoriesRejected:            &[]models.LocationCategory{models.CategoryCafe},
		LocationStart:                 &models.GeoLocation{Latitude: 35.681236, Longitude: 139.767125},
	}); err != nil {
		log.Fatalf("failed to update plan candidate meta data: %v", err)
	}

	// プランを追加
	if err := planCandidateRepository.AddPlan(ctx, "test-plan-candidate-set", models.Plan{
		Id:       uuid.New().String(),
		Name:     "test-plan",
		Places:   []models.Place{*testPlace},
		AuthorId: nil,
	}); err != nil {
		log.Fatalf("failed to add plan: %v", err)
	}

	planCandidate, err := planCandidateRepository.Find(ctx, "test-plan-candidate-set", time.Now())
	if err != nil {
		log.Fatalf("failed to find plan candidate: %v", err)
	}

	log.Printf("plan candidate: %+v", planCandidate)
}

func cleanup(ctx context.Context, db *sql.DB) {
	type Deletable interface {
		DeleteAll(context.Context, boil.ContextExecutor) (int64, error)
	}

	tables := []Deletable{
		generated.PlanCandidateSetSearchedPlaces(),
		generated.PlanCandidatePlaces(),
		generated.PlanCandidateSetMetaData(),
		generated.PlanCandidateSetMetaDataCategories(),
		generated.PlanCandidates(),
		generated.PlanCandidateSets(),
	}

	for _, table := range tables {
		if _, err := table.DeleteAll(ctx, db); err != nil {
			panic(err)
		}
	}
}
