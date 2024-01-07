package rdb

import (
	"context"
	"database/sql"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/mock"
	"testing"
	"time"
)

func TestPlaceRepository_SavePlacesFromGooglePlace(t *testing.T) {
	cases := []struct {
		name        string
		googlePlace models.GooglePlace
	}{
		{
			name:        "save places from google place with nearby search result",
			googlePlace: mock.NewMockGooglePlaceTokyo(false, false),
		},
		{
			name:        "save places from google place with place detail result",
			googlePlace: mock.NewMockGooglePlaceTokyo(true, true),
		},
	}

	placeRepository, err := NewPlaceRepository(testDB)
	if err != nil {
		t.Fatalf("error while initializing place repository: %v", err)
	}

	for _, c := range cases {
		testContext := context.Background()
		t.Run(c.name, func(t *testing.T) {
			defer func(ctx context.Context, db *sql.DB) {
				err := cleanup(ctx, db)
				if err != nil {
					t.Fatalf("error while cleaning up: %v", err)
				}
			}(testContext, testDB)

			actualFirstSave, err := placeRepository.SavePlacesFromGooglePlace(testContext, c.googlePlace)
			if err != nil {
				t.Fatalf("error while saving places: %v", err)
			}

			// GooglePlace が保存されているか確認
			isGooglePlaceSaved, err := generated.
				GooglePlaces(generated.GooglePlaceWhere.GooglePlaceID.EQ(c.googlePlace.PlaceId)).
				Exists(testContext, testDB)
			if err != nil {
				t.Fatalf("error while checking google place existence: %v", err)
			}
			if !isGooglePlaceSaved {
				t.Fatalf("google place is not saved")
			}

			// GooglePlaceType が保存されているか確認
			placeTypeCount, err := generated.
				GooglePlaceTypes(generated.GooglePlaceTypeWhere.GooglePlaceID.EQ(c.googlePlace.PlaceId)).
				Count(testContext, testDB)
			if err != nil {
				t.Fatalf("error while counting place types: %v", err)
			}

			if int(placeTypeCount) != len(c.googlePlace.Types) {
				t.Fatalf("place type expected: %d, actual: %d", len(c.googlePlace.Types), placeTypeCount)
			}

			// GooglePhotoReference が保存されているか確認
			for _, photoReference := range c.googlePlace.PhotoReferences {
				isPhotoReferenceSaved, err := generated.
					GooglePlacePhotoReferences(generated.GooglePlacePhotoReferenceWhere.PhotoReference.EQ(photoReference.PhotoReference)).
					Exists(testContext, testDB)
				if err != nil {
					t.Fatalf("error while checking photo reference existence: %v", err)
				}
				if !isPhotoReferenceSaved {
					t.Fatalf("photo is not saved")
				}
			}

			// HTMLAttributions が保存されているか確認
			for _, photoReference := range c.googlePlace.PhotoReferences {
				htmlAttributionCount, err := generated.
					GooglePlacePhotoAttributions(generated.GooglePlacePhotoAttributionWhere.PhotoReference.EQ(photoReference.PhotoReference)).
					Count(testContext, testDB)
				if err != nil {
					t.Fatalf("error while counting html attributions: %v", err)
				}

				if int(htmlAttributionCount) != len(photoReference.HTMLAttributions) {
					t.Fatalf("html attribution expected: %d, actual: %d", len(photoReference.HTMLAttributions), htmlAttributionCount)
				}
			}

			// Photo が保存されているか確認
			if c.googlePlace.Photos != nil {
				for _, photo := range *c.googlePlace.Photos {
					// 大・小の２つのバリエーションが保存されているか確認
					photoVariation := 0
					if photo.Small != nil {
						photoVariation++
					}
					if photo.Large != nil {
						photoVariation++
					}

					photoCount, err := generated.
						GooglePlacePhotos(generated.GooglePlacePhotoWhere.PhotoReference.EQ(photo.PhotoReference)).
						Count(testContext, testDB)
					if err != nil {
						t.Fatalf("error while counting google photos: %v", err)
					}
					if int(photoCount) != photoVariation {
						t.Fatalf("google photo expected: %d, actual: %d", photoVariation, photoCount)
					}
				}
			}

			if c.googlePlace.PlaceDetail != nil {
				// GooglePlaceOpeningPeriods が保存されているか確認
				if c.googlePlace.PlaceDetail.OpeningHours != nil {
					openingPeriodCount, err := generated.
						GooglePlaceOpeningPeriods(generated.GooglePlaceOpeningPeriodWhere.GooglePlaceID.EQ(c.googlePlace.PlaceId)).
						Count(testContext, testDB)
					if err != nil {
						t.Fatalf("error while counting opening periods: %v", err)
					}

					if int(openingPeriodCount) != len(c.googlePlace.PlaceDetail.OpeningHours.Periods) {
						t.Fatalf("opening period expected: %d, actual: %d", len(c.googlePlace.PlaceDetail.OpeningHours.Periods), openingPeriodCount)
					}
				}

				// GooglePlaceReviews が保存されているか確認
				reviewCount, err := generated.
					GooglePlaceReviews(generated.GooglePlaceReviewWhere.GooglePlaceID.EQ(c.googlePlace.PlaceId)).
					Count(testContext, testDB)
				if err != nil {
					t.Fatalf("error while counting reviews: %v", err)
				}

				if int(reviewCount) != len(c.googlePlace.PlaceDetail.Reviews) {
					t.Fatalf("review expected: %d, actual: %d", len(c.googlePlace.PlaceDetail.Reviews), reviewCount)
				}

				// GooglePhotoReference が保存されているか確認
				for _, photoReference := range c.googlePlace.PlaceDetail.PhotoReferences {
					isPhotoReferenceSaved, err := generated.
						GooglePlacePhotoReferences(generated.GooglePlacePhotoReferenceWhere.PhotoReference.EQ(photoReference.PhotoReference)).
						Exists(testContext, testDB)
					if err != nil {
						t.Fatalf("error while checking photo reference existence: %v", err)
					}
					if !isPhotoReferenceSaved {
						t.Fatalf("photo is not saved")
					}
				}

				// HTMLAttributions が保存されているか確認
				for _, photoReference := range c.googlePlace.PlaceDetail.PhotoReferences {
					htmlAttributionCount, err := generated.
						GooglePlacePhotoAttributions(generated.GooglePlacePhotoAttributionWhere.PhotoReference.EQ(photoReference.PhotoReference)).
						Count(testContext, testDB)
					if err != nil {
						t.Fatalf("error while counting html attributions: %v", err)
					}

					if int(htmlAttributionCount) != len(photoReference.HTMLAttributions) {
						t.Fatalf("html attribution expected: %d, actual: %d", len(photoReference.HTMLAttributions), htmlAttributionCount)
					}
				}
			}

			// 一度保存したあとは、すでに保存されたものが取得される
			actualSecondSave, err := placeRepository.SavePlacesFromGooglePlace(testContext, c.googlePlace)
			if err != nil {
				t.Fatalf("error while saving places second time: %v", err)
			}

			if actualFirstSave.Id != actualSecondSave.Id {
				t.Fatalf("place id expected: %s, actual: %s", actualFirstSave.Id, actualSecondSave.Id)
			}

			if len(actualFirstSave.Google.Types) != len(actualSecondSave.Google.Types) {
				t.Fatalf("place type expected: %d, actual: %d", len(actualFirstSave.Google.Types), len(actualSecondSave.Google.Types))
			}

			if len(actualFirstSave.Google.PhotoReferences) != len(actualSecondSave.Google.PhotoReferences) {
				t.Fatalf("photo reference expected: %d, actual: %d", len(actualFirstSave.Google.PhotoReferences), len(actualSecondSave.Google.PhotoReferences))
			}

			if c.googlePlace.Photos != nil {
				if len(*c.googlePlace.Photos) != len(*actualSecondSave.Google.Photos) {
					t.Fatalf("photo expected: %d, actual: %d", len(*c.googlePlace.Photos), len(*actualSecondSave.Google.Photos))
				}
			}

			if c.googlePlace.PlaceDetail != nil {
				if len(c.googlePlace.PlaceDetail.Reviews) != len(actualSecondSave.Google.PlaceDetail.Reviews) {
					t.Fatalf("review expected: %d, actual: %d", len(c.googlePlace.PlaceDetail.Reviews), len(actualSecondSave.Google.PlaceDetail.Reviews))
				}

				if c.googlePlace.PlaceDetail.OpeningHours != nil {
					if len(c.googlePlace.PlaceDetail.OpeningHours.Periods) != len(actualSecondSave.Google.PlaceDetail.OpeningHours.Periods) {
						t.Fatalf("opening period expected: %d, actual: %d", len(c.googlePlace.PlaceDetail.OpeningHours.Periods), len(actualSecondSave.Google.PlaceDetail.OpeningHours.Periods))
					}
				}
			}
		})
	}
}

func TestPlaceRepository_FindByGooglePlaceID(t *testing.T) {
	cases := []struct {
		name          string
		savedPlaces   []models.Place
		googlePlaceId string
		expectedPlace *models.Place
	}{
		{
			name: "find place by google place id",
			savedPlaces: []models.Place{
				{Id: "place_id_1", Google: models.GooglePlace{PlaceId: "google_place_id_1"}},
				{Id: "place_id_2", Google: models.GooglePlace{PlaceId: "google_place_id_2"}},
			},
			googlePlaceId: "google_place_id_1",
			expectedPlace: &models.Place{Id: "place_id_1", Google: models.GooglePlace{PlaceId: "google_place_id_1"}},
		},
	}

	placeRepository, err := NewPlaceRepository(testDB)
	if err != nil {
		t.Fatalf("error while initializing place repository: %v", err)
	}

	for _, c := range cases {
		testContext := context.Background()
		t.Run(c.name, func(t *testing.T) {
			t.Cleanup(func() {
				err := cleanup(testContext, testDB)
				if err != nil {
					t.Fatalf("error while cleaning up: %v", err)
				}
			})

			// 事前にPlaceを保存しておく
			if err := savePlaces(testContext, testDB, c.savedPlaces); err != nil {
				t.Fatalf("error while saving places: %v", err)
			}

			actualPlace, err := placeRepository.FindByGooglePlaceID(testContext, c.googlePlaceId)
			if err != nil {
				t.Fatalf("error while finding place: %v", err)
			}

			if diff := cmp.Diff(c.expectedPlace, actualPlace); diff != "" {
				t.Fatalf("(-want +got):\n%s", diff)
			}
		})
	}
}

func TestPlaceRepository_FindByGooglePlaceID_WithLikeCount(t *testing.T) {
	cases := []struct {
		name                            string
		savedPlaces                     []models.Place
		savedPlanCandidateSets          []models.PlanCandidate
		savedPlanCandidateSetLikePlaces []generated.PlanCandidateSetLikePlace
		googlePlaceId                   string
		expectedPlace                   *models.Place
	}{
		{
			name: "find place by google place id with like count",
			savedPlaces: []models.Place{
				{Id: "place_id_1", Google: models.GooglePlace{PlaceId: "google_place_id_1"}},
			},
			savedPlanCandidateSets: []models.PlanCandidate{
				{Id: "plan_candidate_set_id_1", ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local)},
			},
			savedPlanCandidateSetLikePlaces: []generated.PlanCandidateSetLikePlace{
				{
					ID:                 uuid.New().String(),
					PlanCandidateSetID: "plan_candidate_set_id_1",
					PlaceID:            "place_id_1",
				},
			},
			googlePlaceId: "google_place_id_1",
			expectedPlace: &models.Place{
				Id:        "place_id_1",
				Google:    models.GooglePlace{PlaceId: "google_place_id_1"},
				LikeCount: 1,
			},
		},
	}

	placeRepository, err := NewPlaceRepository(testDB)
	if err != nil {
		t.Fatalf("error while initializing place repository: %v", err)
	}

	for _, c := range cases {
		testContext := context.Background()
		t.Run(c.name, func(t *testing.T) {
			t.Cleanup(func() {
				err := cleanup(testContext, testDB)
				if err != nil {
					t.Fatalf("error while cleaning up: %v", err)
				}
			})

			// 事前にPlaceを保存しておく
			if err := savePlaces(testContext, testDB, c.savedPlaces); err != nil {
				t.Fatalf("error while saving places: %v", err)
			}

			// 事前にPlanCandidateSetを保存しておく
			for _, planCandidateSet := range c.savedPlanCandidateSets {
				if err := savePlanCandidate(testContext, testDB, planCandidateSet); err != nil {
					t.Fatalf("error while saving plan candidate set: %v", err)
				}
			}

			// 事前にPlanCandidateSetLikePlaceを保存しておく
			for _, planCandidateSetLikePlace := range c.savedPlanCandidateSetLikePlaces {
				if err := planCandidateSetLikePlace.Insert(testContext, testDB, boil.Infer()); err != nil {
					t.Fatalf("error while saving plan candidate set like place: %v", err)
				}
			}

			actualPlace, err := placeRepository.FindByGooglePlaceID(testContext, c.googlePlaceId)
			if err != nil {
				t.Fatalf("error while finding place: %v", err)
			}

			if diff := cmp.Diff(c.expectedPlace, actualPlace); diff != "" {
				t.Fatalf("(-want +got):\n%s", diff)
			}
		})
	}
}

func TestPlaceRepository_FindByPlanCandidateId(t *testing.T) {
	cases := []struct {
		name                               string
		planCandidateId                    string
		savedPlaces                        []models.Place
		savedPlanCandidateSet              models.PlanCandidate
		savedPlanCandidateSearchedPlaceIds []string
		expectedPlaces                     []models.Place
	}{
		{
			name:            "find places by plan candidate id",
			planCandidateId: "plan_candidate_id",
			savedPlaces: []models.Place{
				{Id: "place_id_1", Google: models.GooglePlace{PlaceId: "google_place_id_1"}},
				{Id: "place_id_2", Google: models.GooglePlace{PlaceId: "google_place_id_2"}},
				{Id: "place_id_3", Google: models.GooglePlace{PlaceId: "google_place_id_3"}},
			},
			savedPlanCandidateSet: models.PlanCandidate{
				Id:        "plan_candidate_id",
				ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
			},
			savedPlanCandidateSearchedPlaceIds: []string{
				"place_id_1",
				"place_id_2",
			},
			expectedPlaces: []models.Place{
				{Id: "place_id_1", Google: models.GooglePlace{PlaceId: "google_place_id_1"}},
				{Id: "place_id_2", Google: models.GooglePlace{PlaceId: "google_place_id_2"}},
			},
		},
	}

	placeRepository, err := NewPlaceRepository(testDB)
	if err != nil {
		t.Fatalf("error while initializing place repository: %v", err)
	}

	for _, c := range cases {
		testContext := context.Background()
		t.Run(c.name, func(t *testing.T) {
			defer func(ctx context.Context, db *sql.DB) {
				err := cleanup(ctx, db)
				if err != nil {
					t.Fatalf("error while cleaning up: %v", err)
				}
			}(testContext, testDB)

			// 事前にPlaceを保存しておく
			if err := savePlaces(testContext, testDB, c.savedPlaces); err != nil {
				t.Fatalf("error while saving places: %v", err)
			}

			// 事前にPlanCandidateSetを保存しておく
			if err := savePlanCandidate(testContext, testDB, c.savedPlanCandidateSet); err != nil {
				t.Fatalf("error while saving plan candidate set: %v", err)
			}

			// 事前にPlanCandidateSearchedPlaceを保存しておく
			for _, searchedPlaceId := range c.savedPlanCandidateSearchedPlaceIds {
				planCandidateSearchedPlaceEntity := generated.PlanCandidateSetSearchedPlace{
					ID:                 uuid.New().String(),
					PlanCandidateSetID: c.savedPlanCandidateSet.Id,
					PlaceID:            searchedPlaceId,
				}
				if err := planCandidateSearchedPlaceEntity.Insert(testContext, testDB, boil.Infer()); err != nil {
					t.Fatalf("error while saving plan candidate searched place: %v", err)
				}
			}

			actualPlaces, err := placeRepository.FindByPlanCandidateId(testContext, c.planCandidateId)
			if err != nil {
				t.Fatalf("error while finding places: %v", err)
			}

			if len(actualPlaces) != len(c.expectedPlaces) {
				t.Fatalf("place expected: %d, actual: %d", len(c.expectedPlaces), len(actualPlaces))
			}
		})
	}
}
