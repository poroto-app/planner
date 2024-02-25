package rdb

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
)

var (
	testDB *sql.DB
)

const (
	// project root ディレクトリからの深さ
	depth = 3
)

func TestMain(m *testing.M) {
	// .env.test を読み込む
	rootDir := strings.Repeat("../", depth)

	if err := godotenv.Load(fmt.Sprintf("%s.env.test", rootDir)); err != nil {
		log.Fatalln("failed to load .env.test")
	}

	// DB 接続
	dns := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?parseTime=true&loc=%s&tls=%v&interpolateParams=%v",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
		"Asia%2FTokyo",
		false,
		true,
	)

	db, err := sql.Open("mysql", dns)
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}

	testDB = db
	boil.SetDB(testDB)

	if err := cleanup(context.Background(), testDB); err != nil {
		log.Fatalf("failed to setup database: %v", err)
	}

	// テスト実行
	code := m.Run()

	if err := cleanup(context.Background(), testDB); err != nil {
		log.Fatalf("failed to cleanup database: %v", err)
	}

	os.Exit(code)
}

func cleanup(ctx context.Context, db *sql.DB) error {
	tables := []interface {
		DeleteAll(context.Context, boil.ContextExecutor) (int64, error)
	}{
		// PlanCandidate
		generated.PlanCandidateSetLikePlaces(),
		generated.PlanCandidatePlaces(),
		generated.PlanCandidateSetMetaDataCategories(),
		generated.PlanCandidateSetMetaData(),
		generated.PlanCandidateSetSearchedPlaces(),
		generated.PlanCandidates(),
		generated.PlanCandidateSets(),
		// Plan
		generated.PlanPlaces(),
		generated.Plans(),
		// GooglePlace
		generated.GooglePlaceOpeningPeriods(),
		generated.GooglePlaceReviews(),
		generated.GooglePlacePhotos(),
		generated.GooglePlacePhotoAttributions(),
		generated.GooglePlacePhotoReferences(),
		generated.GooglePlaceTypes(),
		generated.GooglePlaces(),
		// Place
		generated.UserLikePlaces(),
		generated.PlacePhotos(),
		generated.Places(),
		// User
		generated.Users(),
	}

	for _, table := range tables {
		if _, err := table.DeleteAll(ctx, db); err != nil {
			return fmt.Errorf("failed to delete table: %w", err)
		}
	}

	return nil
}
