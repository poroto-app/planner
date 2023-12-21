package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"os"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/env"
	"poroto.app/poroto/planner/internal/infrastructure/rdb"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/entities"
)

func init() {
	env.LoadEnv()
}

func main() {
	dns := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?parseTime=true&loc=%s&tls=%v&interpolateParams=%v",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
		"Asia%2FTokyo",
		os.Getenv("ENV") != "development",
		true,
	)

	db, err := sql.Open("mysql", dns)
	if err != nil {
		panic(err)
	}

	boil.SetDB(db)
	boil.DebugMode = true

	defer db.Close()

	repository, err := rdb.NewUserRepository(db)
	if err != nil {
		panic(err)
	}

	testUser := models.User{
		Id:          uuid.New().String(),
		FirebaseUID: uuid.New().String(),
		Name:        "テスト",
		Email:       utils.StrOmitEmpty("test@example.com"),
		PhotoUrl:    utils.StrOmitEmpty("https://example.com"),
	}

	ctx := context.Background()
	if err := repository.Create(ctx, testUser); err != nil {
		panic(err)
	}

	user, err := repository.Find(ctx, testUser.Id)
	if err != nil {
		panic(err)
	}

	if user == nil {
		panic("user not found")
	}

	if err := validateUser(user, testUser); err != nil {
		panic(err)
	}

	user, err = repository.FindByFirebaseUID(ctx, testUser.FirebaseUID)
	if err != nil {
		panic(err)
	}

	if user == nil {
		panic("user not found")
	}

	if err := validateUser(user, testUser); err != nil {
		panic(err)
	}

	if _, err := entities.Users().DeleteAll(ctx, db); err != nil {
		panic(err)
	}
}

func validateUser(found *models.User, expected models.User) error {
	if found.Id != expected.Id {
		return fmt.Errorf("expected id: %s, found id: %s", expected.Id, found.Id)
	}

	if found.FirebaseUID != expected.FirebaseUID {
		return fmt.Errorf("expected firebase uid: %s, found firebase uid: %s", expected.FirebaseUID, found.FirebaseUID)
	}

	if found.Name != expected.Name {
		return fmt.Errorf("expected name: %s, found name: %s", expected.Name, found.Name)
	}

	if utils.StrEmptyIfNil(found.Email) != utils.StrEmptyIfNil(expected.Email) {
		return fmt.Errorf("expected email: %s, found email: %s", utils.StrEmptyIfNil(expected.Email), utils.StrEmptyIfNil(found.Email))
	}

	if utils.StrEmptyIfNil(found.PhotoUrl) != utils.StrEmptyIfNil(expected.PhotoUrl) {
		return fmt.Errorf("expected photo url: %s, found photo url: %s", utils.StrEmptyIfNil(expected.PhotoUrl), utils.StrEmptyIfNil(found.PhotoUrl))
	}

	return nil
}
