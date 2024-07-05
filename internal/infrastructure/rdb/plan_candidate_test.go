package rdb

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
)

func TestPlanCandidateRepository_Create(t *testing.T) {
	cases := []struct {
		name               string
		planCandidateSetId string
		expiresAt          time.Time
	}{
		{
			name:               "success",
			planCandidateSetId: uuid.New().String(),
			expiresAt:          time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	planCandidateRepository, err := NewPlanCandidateRepository(testDB)
	if err != nil {
		t.Fatalf("failed to create plan candidate repository: %v", err)
	}

	testContext := context.Background()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Cleanup(func() {
				err := cleanup(testContext, testDB)
				if err != nil {
					t.Fatalf("failed to cleanup: %v", err)
				}
			})

			if err := planCandidateRepository.Create(testContext, c.planCandidateSetId, c.expiresAt); err != nil {
				t.Fatalf("failed to create plan candidate: %v", err)
			}

			exists, err := generated.PlanCandidateSetExists(testContext, testDB, c.planCandidateSetId)
			if err != nil {
				t.Fatalf("failed to check plan candidate existence: %v", err)
			}

			if !exists {
				t.Fatalf("plan candidate does not exist")
			}

		})
	}
}

func TestPlanCandidateRepository_Find(t *testing.T) {
	cases := []struct {
		name                  string
		now                   time.Time
		savedPlans            generated.PlanSlice
		savedPlanCandidateSet models.PlanCandidateSet
		planCandidateSetId    string
		expected              *models.PlanCandidateSet
	}{
		{
			name: "plan candidate set with only id",
			now:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
			savedPlans: generated.PlanSlice{
				{ID: "plan-parent"},
			},
			savedPlanCandidateSet: models.PlanCandidateSet{
				Id:              "test",
				ExpiresAt:       time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
				IsPlaceSearched: true,
			},
			planCandidateSetId: "test",
			expected: &models.PlanCandidateSet{
				Id:              "test",
				ExpiresAt:       time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
				IsPlaceSearched: true,
				Plans:           []models.Plan{},
			},
		},
		{
			name: "plan candidate set with plan candidate",
			now:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
			savedPlans: generated.PlanSlice{
				{ID: "plan-parent"},
			},
			savedPlanCandidateSet: models.PlanCandidateSet{
				Id:              "test",
				ExpiresAt:       time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
				IsPlaceSearched: true,
				MetaData: models.PlanCandidateMetaData{
					CreatedBasedOnCurrentLocation: true,
					LocationStart:                 &models.GeoLocation{Latitude: 139.767125, Longitude: 35.681236},
					CategoriesPreferred:           &[]models.LocationCategory{models.CategoryRestaurant},
					CategoriesRejected:            &[]models.LocationCategory{models.CategoryCafe},
				},
				Plans: []models.Plan{
					{
						Id: "test-plan",
						Places: []models.Place{
							{
								Id:     "test-place",
								Google: models.GooglePlace{PlaceId: "test-google-place"},
							},
						},
						ParentPlanId: utils.ToPointer("plan-parent"),
					},
				},
			},
			planCandidateSetId: "test",
			expected: &models.PlanCandidateSet{
				Id:              "test",
				ExpiresAt:       time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
				IsPlaceSearched: true,
				MetaData: models.PlanCandidateMetaData{
					CreatedBasedOnCurrentLocation: true,
					LocationStart:                 &models.GeoLocation{Latitude: 139.767125, Longitude: 35.681236},
					CategoriesPreferred:           &[]models.LocationCategory{models.CategoryRestaurant},
					CategoriesRejected:            &[]models.LocationCategory{models.CategoryCafe},
				},
				Plans: []models.Plan{
					{
						Id: "test-plan",
						Places: []models.Place{
							{
								Id:     "test-place",
								Google: models.GooglePlace{PlaceId: "test-google-place"},
							},
						},
						ParentPlanId: utils.ToPointer("plan-parent"),
					},
				},
			},
		},
		{
			name:               "plan candidate set without PlanCandidateSetMetaData",
			now:                time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
			planCandidateSetId: "test",
			savedPlanCandidateSet: models.PlanCandidateSet{
				Id:              "test",
				ExpiresAt:       time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
				IsPlaceSearched: true,
				Plans: []models.Plan{
					{
						Id: "test-plan",
						Places: []models.Place{
							{Id: "test-place", Google: models.GooglePlace{PlaceId: "test-google-place"}},
						},
					},
				},
			},
			expected: &models.PlanCandidateSet{
				Id:              "test",
				ExpiresAt:       time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
				IsPlaceSearched: true,
				Plans: []models.Plan{
					{
						Id: "test-plan",
						Places: []models.Place{
							{Id: "test-place", Google: models.GooglePlace{PlaceId: "test-google-place"}},
						},
					},
				},
			},
		},
		{
			name:               "plan candidate set with IsPlaceSearched false",
			now:                time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
			planCandidateSetId: "test",
			savedPlanCandidateSet: models.PlanCandidateSet{
				Id:              "test",
				ExpiresAt:       time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
				IsPlaceSearched: false,
			},
			expected: &models.PlanCandidateSet{
				Id:              "test",
				ExpiresAt:       time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
				IsPlaceSearched: false,
				Plans:           []models.Plan{},
			},
		},
		{
			name:               "create by category",
			now:                time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
			planCandidateSetId: "test",
			savedPlanCandidateSet: models.PlanCandidateSet{
				Id:        "test",
				ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
				MetaData: models.PlanCandidateMetaData{
					CreatedBasedOnCurrentLocation: false,
					LocationStart:                 &models.GeoLocation{Latitude: 139.767125, Longitude: 35.681236},
					CreateByCategoryMetaData: &models.CreateByCategoryMetaData{
						Category:   models.LocationCategorySetCreatePlanAmusements.Categories[0],
						Location:   models.GeoLocation{Latitude: 139.767125, Longitude: 35.681236},
						RadiusInKm: 1.0,
					},
				},
			},
			expected: &models.PlanCandidateSet{
				Id:        "test",
				ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
				Plans:     []models.Plan{},
				MetaData: models.PlanCandidateMetaData{
					CreatedBasedOnCurrentLocation: false,
					LocationStart:                 &models.GeoLocation{Latitude: 139.767125, Longitude: 35.681236},
					CreateByCategoryMetaData: &models.CreateByCategoryMetaData{
						Category:   models.LocationCategorySetCreatePlanAmusements.Categories[0],
						Location:   models.GeoLocation{Latitude: 139.767125, Longitude: 35.681236},
						RadiusInKm: 1.0,
					},
				},
			},
		},
	}

	planCandidateRepository, err := NewPlanCandidateRepository(testDB)
	if err != nil {
		t.Fatalf("failed to create plan candidate repository: %v", err)
	}

	testContext := context.Background()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Cleanup(func() {
				err := cleanup(testContext, testDB)
				if err != nil {
					t.Fatalf("failed to cleanup: %v", err)
				}
			})

			// データ準備
			placesInPlanCandidates := array.Flatten(array.Map(c.savedPlanCandidateSet.Plans, func(plan models.Plan) []models.Place { return plan.Places }))
			if err := savePlaces(testContext, testDB, placesInPlanCandidates); err != nil {
				t.Fatalf("failed to save places: %v", err)
			}

			if _, err := c.savedPlans.InsertAll(testContext, testDB, boil.Infer()); err != nil {
				t.Fatalf("failed to save plans: %v", err)
			}

			if err := savePlanCandidateSet(testContext, testDB, c.savedPlanCandidateSet); err != nil {
				t.Fatalf("failed to save plan candidate: %v", err)
			}

			actual, err := planCandidateRepository.Find(testContext, c.planCandidateSetId, c.now)
			if err != nil {
				t.Fatalf("failed to find plan candidate: %v", err)
			}

			if diff := cmp.Diff(
				c.expected,
				actual,
				cmpopts.IgnoreUnexported(models.Plan{}),
			); diff != "" {
				t.Fatalf("wrong plan candidate (-expected, +actual): %v", diff)
			}
		})
	}
}

func TestPlanCandidateRepository_Find_ShouldReturnNil(t *testing.T) {
	cases := []struct {
		name                  string
		now                   time.Time
		savedPlanCandidateSet models.PlanCandidateSet
		planCandidateId       string
	}{
		{
			name: "expired plan candidate set will not be returned",
			now:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
			savedPlanCandidateSet: models.PlanCandidateSet{
				Id:        "test",
				ExpiresAt: time.Date(2019, 12, 1, 0, 0, 0, 0, time.Local),
			},
			planCandidateId: "test",
		},
	}

	planCandidateRepository, err := NewPlanCandidateRepository(testDB)
	if err != nil {
		t.Fatalf("failed to create plan candidate repository: %v", err)
	}

	testContext := context.Background()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Cleanup(func() {
				err := cleanup(testContext, testDB)
				if err != nil {
					t.Fatalf("failed to cleanup: %v", err)
				}
			})

			// 事前にPlaceを作成しておく
			placesInPlans := array.Map(c.savedPlanCandidateSet.Plans, func(plan models.Plan) []models.Place { return plan.Places })
			if err := savePlaces(testContext, testDB, array.Flatten(placesInPlans)); err != nil {
				t.Fatalf("failed to save places: %v", err)
			}

			// 事前にPlanCandidateSetを作成しておく
			if err := savePlanCandidateSet(testContext, testDB, c.savedPlanCandidateSet); err != nil {
				t.Fatalf("failed to save plan candidate: %v", err)
			}

			actual, err := planCandidateRepository.Find(testContext, c.planCandidateId, c.now)
			if err != nil {
				t.Fatalf("failed to find plan candidate: %v", err)
			}

			if actual != nil {
				t.Fatalf("plan candidate should not be found")
			}
		})
	}
}

func TestPlanCandidateRepository_Find_WithPlaceLikeCount(t *testing.T) {
	cases := []struct {
		name                                   string
		now                                    time.Time
		savedPlaces                            []models.Place
		savedUsers                             generated.UserSlice
		savedPlanCandidateSets                 []models.PlanCandidateSet
		savedPlanCandidateSetLikePlaceEntities []generated.PlanCandidateSetLikePlace
		savedUserLikePlaceEntities             generated.UserLikePlaceSlice
		planCandidateSetId                     string
		expected                               models.PlanCandidateSet
	}{
		{
			name: "plan candidate set with place like count",
			now:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
			savedPlaces: []models.Place{
				{Id: "test-place-1", Google: models.GooglePlace{PlaceId: "test-google-place-1"}},
				{Id: "test-place-2", Google: models.GooglePlace{PlaceId: "test-google-place-2"}},
			},
			savedUsers: generated.UserSlice{
				{ID: "test-user-1", FirebaseUID: uuid.New().String()},
				{ID: "test-user-2", FirebaseUID: uuid.New().String()},
			},
			savedPlanCandidateSets: []models.PlanCandidateSet{
				{
					Id:        "plan-candidate-set-1",
					ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
					Plans: []models.Plan{
						{
							Id: "plan-candidate-1",
							Places: []models.Place{
								{Id: "test-place-1", Google: models.GooglePlace{PlaceId: "test-google-place-1"}},
								{Id: "test-place-2", Google: models.GooglePlace{PlaceId: "test-google-place-2"}},
							},
						},
					},
				},
				{
					Id:        "plan-candidate-set-2",
					ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
				},
			},
			savedPlanCandidateSetLikePlaceEntities: []generated.PlanCandidateSetLikePlace{
				{ID: uuid.New().String(), PlanCandidateSetID: "plan-candidate-set-1", PlaceID: "test-place-1"},
				{ID: uuid.New().String(), PlanCandidateSetID: "plan-candidate-set-2", PlaceID: "test-place-1"},
				{ID: uuid.New().String(), PlanCandidateSetID: "plan-candidate-set-2", PlaceID: "test-place-2"},
			},
			savedUserLikePlaceEntities: generated.UserLikePlaceSlice{
				{ID: uuid.New().String(), UserID: "test-user-1", PlaceID: "test-place-1"},
				{ID: uuid.New().String(), UserID: "test-user-2", PlaceID: "test-place-2"},
			},
			planCandidateSetId: "plan-candidate-set-1",
			expected: models.PlanCandidateSet{
				Id:        "plan-candidate-set-1",
				ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
				Plans: []models.Plan{
					{
						Id: "plan-candidate-1",
						Places: []models.Place{
							{Id: "test-place-1", Google: models.GooglePlace{PlaceId: "test-google-place-1"}, LikeCount: 3},
							{Id: "test-place-2", Google: models.GooglePlace{PlaceId: "test-google-place-2"}, LikeCount: 2},
						},
					},
				},
				LikedPlaceIds: []string{"test-place-1"},
			},
		},
	}

	planCandidateRepository, err := NewPlanCandidateRepository(testDB)
	if err != nil {
		t.Fatalf("failed to create plan candidate repository: %v", err)
	}

	for _, c := range cases {
		textContext := context.Background()
		t.Run(c.name, func(t *testing.T) {
			t.Cleanup(func() {
				err := cleanup(textContext, testDB)
				if err != nil {
					t.Fatalf("failed to cleanup: %v", err)
				}
			})

			// 事前にPlace・User・PlanCandidateSetLikePlace・UserLikePlaceを作成しておく
			if err := savePlaces(textContext, testDB, c.savedPlaces); err != nil {
				t.Fatalf("failed to save places: %v", err)
			}

			if _, err := c.savedUsers.InsertAll(textContext, testDB, boil.Infer()); err != nil {
				t.Fatalf("failed to save users: %v", err)
			}

			for _, planCandidateSet := range c.savedPlanCandidateSets {
				if err := savePlanCandidateSet(textContext, testDB, planCandidateSet); err != nil {
					t.Fatalf("failed to save plan candidate: %v", err)
				}
			}

			for _, planCandidateSetLikePlaceEntity := range c.savedPlanCandidateSetLikePlaceEntities {
				if err := planCandidateSetLikePlaceEntity.Insert(textContext, testDB, boil.Infer()); err != nil {
					t.Fatalf("failed to save plan candidate set like place: %v", err)
				}
			}

			if _, err := c.savedUserLikePlaceEntities.InsertAll(textContext, testDB, boil.Infer()); err != nil {
				t.Fatalf("failed to save user like place: %v", err)
			}

			actual, err := planCandidateRepository.Find(textContext, c.planCandidateSetId, c.now)
			if err != nil {
				t.Fatalf("failed to find plan candidate: %v", err)
			}

			if actual == nil {
				t.Fatalf("plan candidate should be found")
			}

			if diff := cmp.Diff(c.expected, *actual); diff != "" {
				t.Fatalf("wrong plan candidate (-expected, +actual): %v", diff)
			}
		})
	}
}

func TestPlanCandidateRepository_FindPlan(t *testing.T) {
	cases := []struct {
		name                  string
		planCandidateSetId    string
		planCandidateId       string
		savedPlaces           []models.Place
		savedPlanCandidateSet models.PlanCandidateSet
		expected              models.Plan
	}{
		{
			name:               "success",
			planCandidateSetId: "test-plan-candidate-set",
			planCandidateId:    "test-plan-candidate",
			savedPlaces: []models.Place{
				{
					Id: "test-place",
					Google: models.GooglePlace{
						PlaceId: "test-google-place",
						Name:    "東京駅",
						Types:   []string{"train_station", "transit_station", "point_of_interest", "establishment"},
						PhotoReferences: []models.GooglePlacePhotoReference{
							{
								PhotoReference:   "photo-1-AWU5eFjiROQJEeMpt7Hh2Pv-fdsabvls-wKBKNsJwobLXjjnbzXSBxTTW3bOtTbsrxkaoE1xx8RU3XFzv64gtTL137nfZtz0YAwpRsWThU7FtEpuJ3xGYOEQ2BFIHKLF5OLpVoGUybE-NryBdtAF7MDlYwBS7XACG",
								Width:            4032,
								Height:           3024,
								HTMLAttributions: []string{"<a href=\"https://maps.google.com/maps/contrib/100969420913538879622\">A Google User</a>"},
							},
						},
						Photos: &[]models.GooglePlacePhoto{
							{
								PhotoReference:   "photo-1-AWU5eFjiROQJEeMpt7Hh2Pv-fdsabvls-wKBKNsJwobLXjjnbzXSBxTTW3bOtTbsrxkaoE1xx8RU3XFzv64gtTL137nfZtz0YAwpRsWThU7FtEpuJ3xGYOEQ2BFIHKLF5OLpVoGUybE-NryBdtAF7MDlYwBS7XACG",
								Width:            4032,
								Height:           3024,
								HTMLAttributions: []string{"<a href=\"https://maps.google.com/maps/contrib/100969420913538879622\">A Google User</a>"},
								Small: &models.Image{
									Width:  400,
									Height: 400,
									URL:    "https://lh3.googleusercontent.com/places/photo-1=s1600-w400-h400",
								},
								Large: &models.Image{
									Width:  4032,
									Height: 3024,
									URL:    "https://lh3.googleusercontent.com/places/photo-1=s1600-w4032-h3024",
								},
							},
						},
						PlaceDetail: &models.GooglePlaceDetail{
							Reviews: []models.GooglePlaceReview{
								{
									Rating:                4,
									Text:                  utils.StrPointer("とても大きな駅です。地下街も広く、お店もたくさんあります。駅員さんも多く、親切です。"),
									Time:                  1648126226,
									AuthorName:            "Alice Alicia",
									AuthorProfileImageUrl: utils.StrPointer("https://lh3.googleusercontent.com/a/ACg8ocKaPr9FWIiqs88c_Fugafugafugafugagfuagaufaugafufa=s128-c0x00000000-cc-rp-mo-ba5"),
									AuthorUrl:             utils.StrPointer("https://www.google.com/maps/contrib/117028493732372946396/reviews"),
								},
							},
						},
					},
				},
			},
			savedPlanCandidateSet: models.PlanCandidateSet{
				Id:        "test-plan-candidate-set",
				ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
				Plans: []models.Plan{
					{
						Id:     "test-plan-candidate",
						Places: []models.Place{{Id: "test-place"}},
					},
				},
			},
			expected: models.Plan{
				Id: "test-plan-candidate",
				Places: []models.Place{
					{
						Id: "test-place",
						Google: models.GooglePlace{
							PlaceId: "test-google-place",
							Name:    "東京駅",
							Types:   []string{"train_station", "transit_station", "point_of_interest", "establishment"},
							PhotoReferences: []models.GooglePlacePhotoReference{
								{
									PhotoReference:   "photo-1-AWU5eFjiROQJEeMpt7Hh2Pv-fdsabvls-wKBKNsJwobLXjjnbzXSBxTTW3bOtTbsrxkaoE1xx8RU3XFzv64gtTL137nfZtz0YAwpRsWThU7FtEpuJ3xGYOEQ2BFIHKLF5OLpVoGUybE-NryBdtAF7MDlYwBS7XACG",
									Width:            4032,
									Height:           3024,
									HTMLAttributions: []string{"<a href=\"https://maps.google.com/maps/contrib/100969420913538879622\">A Google User</a>"},
								},
							},
							Photos: &[]models.GooglePlacePhoto{
								{
									PhotoReference:   "photo-1-AWU5eFjiROQJEeMpt7Hh2Pv-fdsabvls-wKBKNsJwobLXjjnbzXSBxTTW3bOtTbsrxkaoE1xx8RU3XFzv64gtTL137nfZtz0YAwpRsWThU7FtEpuJ3xGYOEQ2BFIHKLF5OLpVoGUybE-NryBdtAF7MDlYwBS7XACG",
									Width:            4032,
									Height:           3024,
									HTMLAttributions: []string{"<a href=\"https://maps.google.com/maps/contrib/100969420913538879622\">A Google User</a>"},
									Small: &models.Image{
										Width:          400,
										Height:         400,
										URL:            "https://lh3.googleusercontent.com/places/photo-1=s1600-w400-h400",
										IsGooglePhotos: true,
									},
									Large: &models.Image{
										Width:          4032,
										Height:         3024,
										URL:            "https://lh3.googleusercontent.com/places/photo-1=s1600-w4032-h3024",
										IsGooglePhotos: true,
									},
								},
							},
							PlaceDetail: &models.GooglePlaceDetail{
								OpeningHours: &models.GooglePlaceOpeningHours{},
								Reviews: []models.GooglePlaceReview{
									{
										Rating:                4,
										Text:                  utils.StrPointer("とても大きな駅です。地下街も広く、お店もたくさんあります。駅員さんも多く、親切です。"),
										Time:                  1648126226,
										AuthorName:            "Alice Alicia",
										AuthorProfileImageUrl: utils.StrPointer("https://lh3.googleusercontent.com/a/ACg8ocKaPr9FWIiqs88c_Fugafugafugafugagfuagaufaugafufa=s128-c0x00000000-cc-rp-mo-ba5"),
										AuthorUrl:             utils.StrPointer("https://www.google.com/maps/contrib/117028493732372946396/reviews"),
									},
								},
								PhotoReferences: []models.GooglePlacePhotoReference{
									{
										PhotoReference:   "photo-1-AWU5eFjiROQJEeMpt7Hh2Pv-fdsabvls-wKBKNsJwobLXjjnbzXSBxTTW3bOtTbsrxkaoE1xx8RU3XFzv64gtTL137nfZtz0YAwpRsWThU7FtEpuJ3xGYOEQ2BFIHKLF5OLpVoGUybE-NryBdtAF7MDlYwBS7XACG",
										Width:            4032,
										Height:           3024,
										HTMLAttributions: []string{"<a href=\"https://maps.google.com/maps/contrib/100969420913538879622\">A Google User</a>"},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	planCandidateRepository, err := NewPlanCandidateRepository(testDB)
	if err != nil {
		t.Fatalf("failed to create plan candidate repository: %v", err)
	}

	for _, c := range cases {
		testContext := context.Background()
		t.Run(c.name, func(t *testing.T) {
			t.Cleanup(func() {
				err := cleanup(testContext, testDB)
				if err != nil {
					t.Fatalf("failed to cleanup: %v", err)
				}
			})

			// 事前にPlaceを作成しておく
			if err := savePlaces(testContext, testDB, c.savedPlaces); err != nil {
				t.Fatalf("failed to save places: %v", err)
			}

			// 事前にPlanCandidateSetを作成しておく
			if err := savePlanCandidateSet(testContext, testDB, c.savedPlanCandidateSet); err != nil {
				t.Fatalf("failed to save plan candidate: %v", err)
			}

			actual, err := planCandidateRepository.FindPlan(testContext, c.planCandidateSetId, c.planCandidateId)
			if err != nil {
				t.Fatalf("failed to find plan: %v", err)
			}

			if actual == nil {
				t.Fatalf("plan should be found")
			}

			// Id の値が一致する
			if actual.Id != c.expected.Id {
				t.Fatalf("wrong plan id expected: %v, actual: %v", c.expected.Id, actual.Id)
			}

			// Place が一致する
			if diff := cmp.Diff(
				c.expected.Places,
				actual.Places,
				cmpopts.SortSlices(func(a, b models.Place) bool { return a.Id < b.Id }),
			); diff != "" {
				t.Fatalf("wrong places (-expected, +actual): %v", diff)
			}

			if len(actual.Places) != len(c.expected.Places) {
				t.Fatalf("wrong number of places expected: %v, actual: %v", len(c.expected.Places), len(actual.Places))
			}

			// Place の順番が一致する
			for i, place := range actual.Places {
				if place.Id != c.expected.Places[i].Id {
					t.Fatalf("wrong place id expected: %v, actual: %v", c.expected.Places[i].Id, place.Id)
				}
			}
		})
	}
}

func TestPlanCandidateRepository_AddPlan(t *testing.T) {
	cases := []struct {
		name                   string
		planCandidateId        string
		savedPlanCandidateSet  generated.PlanCandidateSet
		savedPlanCandidates    generated.PlanCandidateSlice
		savedPlans             generated.PlanSlice
		plans                  []models.Plan
		expectedPlanCandidates generated.PlanCandidateSlice
	}{
		{
			name:            "add plan from empty",
			planCandidateId: "plan-candidate",
			savedPlanCandidateSet: generated.PlanCandidateSet{
				ID:        "plan-candidate",
				ExpiresAt: time.Now().Add(time.Hour),
			},
			savedPlans: generated.PlanSlice{
				// もとになったプラン
				{ID: "test-plan-parent"},
			},
			plans: []models.Plan{
				{
					Id: "plan-candidate-1",
					Places: []models.Place{
						{Id: "tokyo-station"},
						{Id: "shinagawa-station"},
					},
					ParentPlanId: utils.ToPointer("test-plan-parent"),
				},
				{
					Id: "plan-candidate-2",
					Places: []models.Place{
						{Id: "yokohama-station"},
						{Id: "shin-yokohama-station"},
					},
					ParentPlanId: nil,
				},
			},
			expectedPlanCandidates: generated.PlanCandidateSlice{
				{
					ID:                 "plan-candidate-1",
					PlanCandidateSetID: "plan-candidate",
					SortOrder:          0,
					ParentPlanID:       null.StringFrom("test-plan-parent"),
				},
				{
					ID:                 "plan-candidate-2",
					PlanCandidateSetID: "plan-candidate",
					SortOrder:          1,
					ParentPlanID:       null.String{},
				},
			},
		},
		{
			name:                  "add plan from existing",
			planCandidateId:       "plan-candidate",
			savedPlanCandidateSet: generated.PlanCandidateSet{ID: "plan-candidate", ExpiresAt: time.Now().Add(time.Hour)},
			savedPlanCandidates: generated.PlanCandidateSlice{
				// すでに作成されているプラン
				{ID: "plan-candidate-saved", PlanCandidateSetID: "plan-candidate", SortOrder: 0},
			},
			plans: []models.Plan{
				{
					Id: "plan-candidate-1",
					Places: []models.Place{
						{Id: "yokohama-station"},
						{Id: "shin-yokohama-station"},
					},
					ParentPlanId: nil,
				},
			},
			expectedPlanCandidates: generated.PlanCandidateSlice{
				{
					ID:                 "plan-candidate-saved",
					PlanCandidateSetID: "plan-candidate",
					SortOrder:          0,
					ParentPlanID:       null.String{},
				},
				{
					ID:                 "plan-candidate-1",
					PlanCandidateSetID: "plan-candidate",
					SortOrder:          1,
					ParentPlanID:       null.String{},
				},
			},
		},
	}

	planCandidateRepository, err := NewPlanCandidateRepository(testDB)
	if err != nil {
		t.Fatalf("failed to create plan candidate repository: %v", err)
	}

	testContext := context.Background()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Cleanup(func() {
				err := cleanup(testContext, testDB)
				if err != nil {
					t.Fatalf("failed to cleanup: %v", err)
				}
			})

			// データの準備
			placesInPlans := array.Map(c.plans, func(plan models.Plan) []models.Place { return plan.Places })
			if err := savePlaces(testContext, testDB, array.Flatten(placesInPlans)); err != nil {
				t.Fatalf("failed to save places: %v", err)
			}

			if err := c.savedPlanCandidateSet.Insert(testContext, testDB, boil.Infer()); err != nil {
				t.Fatalf("failed to save plan candidate set: %v", err)
			}

			if _, err := c.savedPlanCandidates.InsertAll(testContext, testDB, boil.Infer()); err != nil {
				t.Fatalf("failed to save plan candidates: %v", err)
			}

			if _, err := c.savedPlans.InsertAll(testContext, testDB, boil.Infer()); err != nil {
				t.Fatalf("failed to save plans: %v", err)
			}

			if err := planCandidateRepository.AddPlan(testContext, c.planCandidateId, c.plans...); err != nil {
				t.Fatalf("failed to add plan: %v", err)
			}

			// すべてのPlanCandidateが保存されている
			planCandidates, err := generated.
				PlanCandidates(generated.PlanCandidateWhere.PlanCandidateSetID.EQ(c.planCandidateId)).
				All(testContext, testDB)
			if err != nil {
				t.Fatalf("failed to get plan candidates: %v", err)
			}

			if diff := cmp.Diff(
				arrayFromPointerSlice(c.expectedPlanCandidates),
				arrayFromPointerSlice(planCandidates),
				cmpopts.IgnoreFields(generated.PlanCandidate{}, "CreatedAt", "UpdatedAt"),
				cmpopts.SortSlices(func(a, b generated.PlanCandidate) bool { return a.SortOrder < b.SortOrder }),
			); diff != "" {
				t.Fatalf("wrong plan candidates (-expected, +actual): %v", diff)
			}

			// すべてのPlanCandidateに対して、すべてのPlaceが保存されている
			for _, plan := range c.plans {
				numPlanCandidatePlaces, err := generated.
					PlanCandidatePlaces(generated.PlanCandidatePlaceWhere.PlanCandidateID.EQ(plan.Id)).
					Count(testContext, testDB)
				if err != nil {
					t.Fatalf("failed to get plan candidate places: %v", err)
				}
				if int(numPlanCandidatePlaces) != len(plan.Places) {
					t.Fatalf("wrong number of plan candidate places expected: %v, actual: %v", len(plan.Places), numPlanCandidatePlaces)
				}
			}
		})
	}
}

func TestPlanCandidateRepository_AddPlaceToPlan(t *testing.T) {
	cases := []struct {
		name                          string
		planCandidateSetId            string
		planCandidateId               string
		previousPlaceId               string
		savedPlanCandidatePlaces      []models.Place
		place                         models.Place
		expectedPlanCandidatePlaceIds []string
	}{
		{
			name:               "success",
			planCandidateSetId: uuid.New().String(),
			planCandidateId:    uuid.New().String(),
			previousPlaceId:    "first-place",
			savedPlanCandidatePlaces: []models.Place{
				{Id: "first-place"},
				{Id: "second-place"},
			},
			place:                         models.Place{Id: "third-place"},
			expectedPlanCandidatePlaceIds: []string{"first-place", "third-place", "second-place"},
		},
	}

	planCandidateRepository, err := NewPlanCandidateRepository(testDB)
	if err != nil {
		t.Fatalf("failed to create plan candidate repository: %v", err)
	}

	testContext := context.Background()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Cleanup(func() {
				err := cleanup(testContext, testDB)
				if err != nil {
					t.Fatalf("failed to cleanup: %v", err)
				}
			})

			// 事前にPlaceを作成しておく
			var placeEntitySlice generated.PlaceSlice
			placeEntitySlice = append(placeEntitySlice, &generated.Place{ID: c.place.Id})
			for _, place := range c.savedPlanCandidatePlaces {
				placeEntitySlice = append(placeEntitySlice, &generated.Place{ID: place.Id})
			}
			for _, placeEntity := range placeEntitySlice {
				if err := placeEntity.Insert(testContext, testDB, boil.Infer()); err != nil {
					t.Fatalf("failed to insert place: %v", err)
				}
			}

			// 事前にPlanCandidateSetを作成しておく
			if err := planCandidateRepository.Create(testContext, c.planCandidateSetId, time.Now().Add(time.Hour)); err != nil {
				t.Fatalf("failed to create plan candidate: %v", err)
			}

			// 事前にPlanCandidateを作成しておく
			if err := planCandidateRepository.AddPlan(testContext, c.planCandidateSetId, models.Plan{Id: c.planCandidateId, Places: c.savedPlanCandidatePlaces}); err != nil {
				t.Fatalf("failed to add plan: %v", err)
			}

			if err := planCandidateRepository.AddPlaceToPlan(testContext, c.planCandidateSetId, c.planCandidateId, c.previousPlaceId, c.place); err != nil {
				t.Fatalf("failed to add place to plan: %v", err)
			}

			savedPlanCandidatePlaceSlice, err := generated.
				PlanCandidatePlaces(
					generated.PlanCandidatePlaceWhere.PlanCandidateID.EQ(c.planCandidateId),
					qm.OrderBy(generated.PlanCandidatePlaceColumns.SortOrder),
				).All(testContext, testDB)
			if err != nil {
				t.Fatalf("failed to get plan candidate places: %v", err)
			}

			if len(savedPlanCandidatePlaceSlice) != len(c.expectedPlanCandidatePlaceIds) {
				t.Fatalf("wrong number of plan candidate places expected: %v, actual: %v", len(c.expectedPlanCandidatePlaceIds), len(savedPlanCandidatePlaceSlice))
			}

			for i, planCandidatePlaceEntity := range savedPlanCandidatePlaceSlice {
				if planCandidatePlaceEntity.PlaceID != c.expectedPlanCandidatePlaceIds[i] {
					t.Fatalf("wrong order of plan candidate places expected: %v, actual: %v", c.expectedPlanCandidatePlaceIds[i], planCandidatePlaceEntity.PlaceID)
				}
				if planCandidatePlaceEntity.SortOrder != i {
					t.Fatalf("wrong order of plan candidate places expected: %v, actual: %v", i, planCandidatePlaceEntity.SortOrder)
				}
			}
		})
	}
}

func TestPlanCandidateRepository_RemovePlaceFromPlan(t *testing.T) {
	cases := []struct {
		name                  string
		planCandidateSetId    string
		planCandidateId       string
		placeIdToDelete       string
		savedPlanCandidateSet models.PlanCandidateSet
	}{
		{
			name:               "success",
			planCandidateSetId: "test-plan-candidate-set",
			planCandidateId:    "test-plan-candidate",
			placeIdToDelete:    "second-place",
			savedPlanCandidateSet: models.PlanCandidateSet{
				Id:        "test-plan-candidate-set",
				ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
				Plans: []models.Plan{
					{
						Id: "test-plan-candidate",
						Places: []models.Place{
							{Id: "first-place"},
							{Id: "second-place"},
							{Id: "third-place"},
						},
					},
				},
			},
		},
		{
			name:               "delete not existing place",
			planCandidateSetId: "test-plan-candidate-set",
			planCandidateId:    "test-plan-candidate",
			placeIdToDelete:    "not-existing-place",
			savedPlanCandidateSet: models.PlanCandidateSet{
				Id:        "test-plan-candidate-set",
				ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
				Plans: []models.Plan{
					{
						Id: "test-plan-candidate",
						Places: []models.Place{
							{Id: "first-place"},
							{Id: "second-place"},
							{Id: "third-place"},
						},
					},
				},
			},
		},
	}

	planCandidateRepository, err := NewPlanCandidateRepository(testDB)
	if err != nil {
		t.Fatalf("failed to create plan candidate repository: %v", err)
	}

	for _, c := range cases {
		testContext := context.Background()
		t.Run(c.name, func(t *testing.T) {
			t.Cleanup(func() {
				err := cleanup(testContext, testDB)
				if err != nil {
					t.Fatalf("failed to cleanup: %v", err)
				}
			})

			// 事前にPlaceを作成しておく
			placesInPlanCandidates := array.Flatten(array.Map(c.savedPlanCandidateSet.Plans, func(plan models.Plan) []models.Place { return plan.Places }))
			if err := savePlaces(testContext, testDB, placesInPlanCandidates); err != nil {
				t.Fatalf("failed to save places: %v", err)
			}

			// 事前にPlanCandidateSetを作成しておく
			if err := savePlanCandidateSet(testContext, testDB, c.savedPlanCandidateSet); err != nil {
				t.Fatalf("failed to save plan candidate: %v", err)
			}

			if err := planCandidateRepository.RemovePlaceFromPlan(testContext, c.planCandidateSetId, c.planCandidateId, c.placeIdToDelete); err != nil {
				t.Fatalf("failed to remove place from plan: %v", err)
			}

			isExistPlanCandidatePlace, err := generated.PlanCandidatePlaces(
				generated.PlanCandidatePlaceWhere.PlanCandidateID.EQ(c.planCandidateId),
				generated.PlanCandidatePlaceWhere.PlaceID.EQ(c.placeIdToDelete),
			).Exists(testContext, testDB)
			if err != nil {
				t.Fatalf("failed to check existence of plan candidate place: %v", err)
			}

			if isExistPlanCandidatePlace {
				t.Fatalf("plan candidate place should be deleted")
			}
		})
	}
}

func TestPlanCandidateRepository_UpdatePlacesOrder(t *testing.T) {
	cases := []struct {
		name                  string
		planCandidateSetId    string
		planCandidateId       string
		placeIdsOrdered       []string
		savedPlanCandidateSet models.PlanCandidateSet
	}{
		{
			name:               "success",
			planCandidateSetId: "test-plan-candidate-set",
			planCandidateId:    "test-plan-candidate",
			placeIdsOrdered:    []string{"third-place", "first-place", "second-place"},
			savedPlanCandidateSet: models.PlanCandidateSet{
				Id:        "test-plan-candidate-set",
				ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
				Plans: []models.Plan{
					{
						Id: "test-plan-candidate",
						Places: []models.Place{
							{Id: "first-place"},
							{Id: "second-place"},
							{Id: "third-place"},
						},
					},
				},
			},
		},
	}

	planCandidateRepository, err := NewPlanCandidateRepository(testDB)
	if err != nil {
		t.Fatalf("failed to create plan candidate repository: %v", err)
	}

	for _, c := range cases {
		testContext := context.Background()
		t.Run(c.name, func(t *testing.T) {
			t.Cleanup(func() {
				err := cleanup(testContext, testDB)
				if err != nil {
					t.Fatalf("failed to cleanup: %v", err)
				}
			})

			// 事前にPlaceを作成しておく
			placesInPlanCandidates := array.Flatten(array.Map(c.savedPlanCandidateSet.Plans, func(plan models.Plan) []models.Place { return plan.Places }))
			if err := savePlaces(testContext, testDB, placesInPlanCandidates); err != nil {
				t.Fatalf("failed to save places: %v", err)
			}

			// 事前にPlanCandidateSetを作成しておく
			if err := savePlanCandidateSet(testContext, testDB, c.savedPlanCandidateSet); err != nil {
				t.Fatalf("failed to save plan candidate: %v", err)
			}

			err := planCandidateRepository.UpdatePlacesOrder(testContext, c.planCandidateId, c.planCandidateSetId, c.placeIdsOrdered)
			if err != nil {
				t.Fatalf("failed to update places order: %v", err)
			}

			for i, placeId := range c.placeIdsOrdered {
				planCandidatePlaceEntity, err := generated.PlanCandidatePlaces(
					generated.PlanCandidatePlaceWhere.PlanCandidateID.EQ(c.planCandidateId),
					generated.PlanCandidatePlaceWhere.PlaceID.EQ(placeId),
				).One(testContext, testDB)
				if err != nil {
					t.Fatalf("failed to get plan candidate place: %v", err)
				}

				if planCandidatePlaceEntity.SortOrder != i {
					t.Fatalf("wrong order of plan candidate place expected: %v, actual: %v", i, planCandidatePlaceEntity.SortOrder)
				}
			}
		})
	}
}

func TestPlanCandidateRepository_UpdatePlacesOrder_ShouldReturnError(t *testing.T) {
	cases := []struct {
		name                  string
		planCandidateSetId    string
		planCandidateId       string
		placeIdsOrdered       []string
		savedPlanCandidateSet models.PlanCandidateSet
	}{
		{
			name:               "reorder with not existing place",
			planCandidateSetId: "test-plan-candidate-set",
			planCandidateId:    "test-plan-candidate",
			placeIdsOrdered:    []string{"third-place", "first-place", "not-existing-place"},
			savedPlanCandidateSet: models.PlanCandidateSet{
				Id:        "test-plan-candidate-set",
				ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
				Plans: []models.Plan{
					{
						Id: "test-plan-candidate",
						Places: []models.Place{
							{Id: "first-place"},
							{Id: "second-place"},
							{Id: "third-place"},
						},
					},
				},
			},
		},
	}

	planCandidateRepository, err := NewPlanCandidateRepository(testDB)
	if err != nil {
		t.Fatalf("failed to create plan candidate repository: %v", err)
	}

	for _, c := range cases {
		testContext := context.Background()
		t.Run(c.name, func(t *testing.T) {
			t.Cleanup(func() {
				err := cleanup(testContext, testDB)
				if err != nil {
					t.Fatalf("failed to cleanup: %v", err)
				}
			})

			// 事前にPlaceを作成しておく
			placesInPlanCandidates := array.Flatten(array.Map(c.savedPlanCandidateSet.Plans, func(plan models.Plan) []models.Place { return plan.Places }))
			if err := savePlaces(testContext, testDB, placesInPlanCandidates); err != nil {
				t.Fatalf("failed to save places: %v", err)
			}

			// 事前にPlanCandidateSetを作成しておく
			if err := savePlanCandidateSet(testContext, testDB, c.savedPlanCandidateSet); err != nil {
				t.Fatalf("failed to save plan candidate: %v", err)
			}

			err := planCandidateRepository.UpdatePlacesOrder(testContext, c.planCandidateId, c.planCandidateSetId, c.placeIdsOrdered)
			if err == nil {
				t.Fatalf("error should be returned")
			}
		})
	}
}

func TestPlanCandidateRepository_UpdatePlanCandidateMetaData(t *testing.T) {
	cases := []struct {
		name                                             string
		planCandidateSetId                               string
		savedPlanCandidateSet                            generated.PlanCandidateSet
		savedPlanCandidateSetMetaData                    *generated.PlanCandidateSetMetaDatum
		savedPlanCandidateSetMetaDataCategorySlice       generated.PlanCandidateSetMetaDataCategorySlice
		metaData                                         models.PlanCandidateMetaData
		expectedPlanCandidateSetMetaData                 *generated.PlanCandidateSetMetaDatum
		expectedPlanCandidateSetMetaDataCategorySlice    generated.PlanCandidateSetMetaDataCategorySlice
		expectedPlanCandidateSetMetaDataCreateByCategory *generated.PlanCandidateSetMetaDataCreateByCategory
	}{
		{
			name:               "save plan candidate meta data",
			planCandidateSetId: "test-plan-candidate-set",
			savedPlanCandidateSet: generated.PlanCandidateSet{
				ID:        "test-plan-candidate-set",
				ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
			},
			metaData: models.PlanCandidateMetaData{
				CreatedBasedOnCurrentLocation: true,
				CategoriesPreferred:           &[]models.LocationCategory{models.CategoryRestaurant},
				CategoriesRejected:            &[]models.LocationCategory{models.CategorySpa},
				LocationStart:                 &models.GeoLocation{Latitude: 35.681236, Longitude: 139.767125},
				FreeTime:                      utils.ToPointer(60),
			},
			expectedPlanCandidateSetMetaData: &generated.PlanCandidateSetMetaDatum{
				PlanCandidateSetID:           "test-plan-candidate-set",
				LatitudeStart:                35.681236,
				LongitudeStart:               139.767125,
				IsCreatedFromCurrentLocation: true,
				PlanDurationMinutes:          null.IntFrom(60),
			},
			expectedPlanCandidateSetMetaDataCategorySlice: generated.PlanCandidateSetMetaDataCategorySlice{
				{
					PlanCandidateSetID: "test-plan-candidate-set",
					Category:           models.CategoryRestaurant.Name,
					IsSelected:         true,
				},
				{
					PlanCandidateSetID: "test-plan-candidate-set",
					Category:           models.CategorySpa.Name,
					IsSelected:         false,
				},
			},
		},
		{
			name:               "update plan candidate meta data",
			planCandidateSetId: "test-plan-candidate-set",
			savedPlanCandidateSet: generated.PlanCandidateSet{
				ID:        "test-plan-candidate-set",
				ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
			},
			savedPlanCandidateSetMetaData: &generated.PlanCandidateSetMetaDatum{
				PlanCandidateSetID:           "test-plan-candidate-set",
				LatitudeStart:                35.681236,
				LongitudeStart:               139.767125,
				IsCreatedFromCurrentLocation: false,
				PlanDurationMinutes:          null.IntFrom(60),
			},
			savedPlanCandidateSetMetaDataCategorySlice: generated.PlanCandidateSetMetaDataCategorySlice{
				{
					ID:                 uuid.New().String(),
					PlanCandidateSetID: "test-plan-candidate-set",
					Category:           models.CategoryRestaurant.Name,
					IsSelected:         true,
				},
				{
					ID:                 uuid.New().String(),
					PlanCandidateSetID: "test-plan-candidate-set",
					Category:           models.CategorySpa.Name,
					IsSelected:         false,
				},
			},
			metaData: models.PlanCandidateMetaData{
				CreatedBasedOnCurrentLocation: true,
				CategoriesPreferred:           &[]models.LocationCategory{models.CategoryRestaurant, models.CategoryBakery},
				CategoriesRejected:            &[]models.LocationCategory{models.CategoryShopping, models.CategoryAmusements},
				LocationStart:                 &models.GeoLocation{Latitude: 36.681236, Longitude: 140.767125},
				FreeTime:                      utils.ToPointer(120),
			},
			expectedPlanCandidateSetMetaData: &generated.PlanCandidateSetMetaDatum{
				PlanCandidateSetID:           "test-plan-candidate-set",
				LatitudeStart:                36.681236,
				LongitudeStart:               140.767125,
				IsCreatedFromCurrentLocation: true,
				PlanDurationMinutes:          null.IntFrom(120),
			},
			expectedPlanCandidateSetMetaDataCategorySlice: generated.PlanCandidateSetMetaDataCategorySlice{
				{
					PlanCandidateSetID: "test-plan-candidate-set",
					Category:           models.CategoryRestaurant.Name,
					IsSelected:         true,
				},
				{
					PlanCandidateSetID: "test-plan-candidate-set",
					Category:           models.CategoryBakery.Name,
					IsSelected:         true,
				},
				{
					PlanCandidateSetID: "test-plan-candidate-set",
					Category:           models.CategoryAmusements.Name,
					IsSelected:         false,
				},
				{
					PlanCandidateSetID: "test-plan-candidate-set",
					Category:           models.CategoryShopping.Name,
					IsSelected:         false,
				},
			},
		},
		{
			name:               "create plan by category",
			planCandidateSetId: "test-plan-candidate-set",
			savedPlanCandidateSet: generated.PlanCandidateSet{
				ID:        "test-plan-candidate-set",
				ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
			},
			metaData: models.PlanCandidateMetaData{
				CreatedBasedOnCurrentLocation: true,
				LocationStart:                 &models.GeoLocation{Latitude: 36.681236, Longitude: 140.767125},
				CreateByCategoryMetaData: &models.CreateByCategoryMetaData{
					Category:   models.LocationCategorySetCreatePlanAmusements.Categories[0],
					Location:   models.GeoLocation{Latitude: 36.681236, Longitude: 140.767125},
					RadiusInKm: 1,
				},
			},
			expectedPlanCandidateSetMetaData: &generated.PlanCandidateSetMetaDatum{
				PlanCandidateSetID:           "test-plan-candidate-set",
				LatitudeStart:                36.681236,
				LongitudeStart:               140.767125,
				IsCreatedFromCurrentLocation: true,
				PlanDurationMinutes:          null.IntFromPtr(nil),
			},
			expectedPlanCandidateSetMetaDataCreateByCategory: &generated.PlanCandidateSetMetaDataCreateByCategory{
				PlanCandidateSetID: "test-plan-candidate-set",
				CategoryID:         models.LocationCategorySetCreatePlanAmusements.Categories[0].Id,
				Latitude:           36.681236,
				Longitude:          140.767125,
				RangeInMeters:      1000,
			},
		},
	}

	planCandidateRepository, err := NewPlanCandidateRepository(testDB)
	if err != nil {
		t.Fatalf("failed to create plan candidate repository: %v", err)
	}

	for _, c := range cases {
		testContext := context.Background()
		t.Run(c.name, func(t *testing.T) {
			t.Cleanup(func() {
				err := cleanup(testContext, testDB)
				if err != nil {
					t.Fatalf("failed to cleanup: %v", err)
				}
			})

			// データ準備
			if err := c.savedPlanCandidateSet.Insert(testContext, testDB, boil.Infer()); err != nil {
				t.Fatalf("failed to save plan candidate: %v", err)
			}

			if _, err := c.savedPlanCandidateSetMetaDataCategorySlice.InsertAll(testContext, testDB, boil.Infer()); err != nil {
				t.Fatalf("failed to save plan candidate set meta data category: %v", err)
			}

			err := planCandidateRepository.UpdatePlanCandidateMetaData(testContext, c.planCandidateSetId, c.metaData)
			if err != nil {
				t.Fatalf("failed to update plan candidate meta data: %v", err)
			}

			planCandidateSetMetaDataEntity, err := generated.
				PlanCandidateSetMetaData(generated.PlanCandidateSetMetaDatumWhere.PlanCandidateSetID.EQ(c.planCandidateSetId)).
				One(testContext, testDB)
			if err != nil {
				t.Fatalf("failed to get plan candidate set meta data: %v", err)
			}

			if diff := cmp.Diff(
				c.expectedPlanCandidateSetMetaData,
				planCandidateSetMetaDataEntity,
				cmpopts.IgnoreFields(generated.PlanCandidateSetMetaDatum{}, "ID", "CreatedAt", "UpdatedAt"),
			); diff != "" {
				t.Fatalf("wrong plan candidate set meta data (-expected, +actual): %v", diff)
			}

			// Category が保存されている
			categories, err := generated.PlanCandidateSetMetaDataCategories(generated.PlanCandidateSetMetaDataCategoryWhere.PlanCandidateSetID.EQ(c.planCandidateSetId)).All(testContext, testDB)
			if err != nil {
				t.Fatalf("failed to get plan candidate set meta data categories: %v", err)
			}

			if diff := cmp.Diff(
				c.expectedPlanCandidateSetMetaDataCategorySlice,
				categories,
				cmpopts.SortSlices(func(a, b *generated.PlanCandidateSetMetaDataCategory) bool {
					return strings.Compare(a.Category, b.Category) < 0
				}),
				cmpopts.IgnoreFields(generated.PlanCandidateSetMetaDataCategory{}, "ID", "CreatedAt", "UpdatedAt"),
			); diff != "" {
				t.Fatalf("wrong plan candidate set meta data categories (-expected, +actual): %v", diff)
			}

			// CreateByCategory が保存されている
			createByCategoryMetaData, err := generated.PlanCandidateSetMetaDataCreateByCategories(generated.PlanCandidateSetMetaDataCreateByCategoryWhere.PlanCandidateSetID.EQ(c.planCandidateSetId)).One(testContext, testDB)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				t.Fatalf("failed to get plan candidate set meta data create by category: %v", err)
			}

			if diff := cmp.Diff(
				c.expectedPlanCandidateSetMetaDataCreateByCategory,
				createByCategoryMetaData,
				cmpopts.IgnoreFields(generated.PlanCandidateSetMetaDataCreateByCategory{}, "ID", "CreatedAt", "UpdatedAt"),
			); diff != "" {
				t.Fatalf("wrong plan candidate set meta data create by category (-expected, +actual): %v", diff)
			}
		})
	}
}

func TestPlanCandidateRepository_UpdateIsPlaceSearched(t *testing.T) {
	cases := []struct {
		name                  string
		savedPlanCandidateSet generated.PlanCandidateSetSlice
		planCandidateSetId    string
		isPlaceSearched       bool
		expected              bool
	}{
		{
			name: "success",
			savedPlanCandidateSet: generated.PlanCandidateSetSlice{
				{
					ID:              "test-plan-candidate-set",
					ExpiresAt:       time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
					IsPlaceSearched: false,
				},
			},
			planCandidateSetId: "test-plan-candidate-set",
			isPlaceSearched:    true,
			expected:           true,
		},
	}

	planCandidateRepository, err := NewPlanCandidateRepository(testDB)
	if err != nil {
		t.Fatalf("failed to create plan candidate repository: %v", err)
	}

	for _, c := range cases {
		testContext := context.Background()
		t.Run(c.name, func(t *testing.T) {
			t.Cleanup(func() {
				err := cleanup(testContext, testDB)
				if err != nil {
					t.Fatalf("failed to cleanup: %v", err)
				}
			})

			// データの準備
			if _, err := c.savedPlanCandidateSet.InsertAll(testContext, testDB, boil.Infer()); err != nil {
				t.Fatalf("failed to insert plan candidate set: %v", err)
			}

			err := planCandidateRepository.UpdateIsPlaceSearched(testContext, c.planCandidateSetId, c.isPlaceSearched)
			if err != nil {
				t.Fatalf("failed to update is place searched: %v", err)
			}

			planCandidateSetEntity, err := generated.PlanCandidateSets(
				generated.PlanCandidateSetWhere.ID.EQ(c.planCandidateSetId),
			).One(testContext, testDB)
			if err != nil {
				t.Fatalf("failed to get plan candidate set: %v", err)
			}

			if planCandidateSetEntity.IsPlaceSearched != c.expected {
				t.Fatalf("wrong is place searched expected: %v, actual: %v", c.expected, planCandidateSetEntity.IsPlaceSearched)
			}
		})
	}
}

func TestPlanCandidateRepository_ReplacePlace(t *testing.T) {
	cases := []struct {
		name                  string
		planCandidateSetId    string
		planCandidateId       string
		placeIdToReplace      string
		placeToReplace        models.Place
		savedPlanCandidateSet models.PlanCandidateSet
	}{
		{
			name:               "success",
			planCandidateSetId: "test-plan-candidate-set",
			planCandidateId:    "test-plan-candidate",
			placeIdToReplace:   "second-place",
			placeToReplace:     models.Place{Id: "replaced-place"},
			savedPlanCandidateSet: models.PlanCandidateSet{
				Id:        "test-plan-candidate-set",
				ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
				Plans: []models.Plan{
					{
						Id: "test-plan-candidate",
						Places: []models.Place{
							{Id: "first-place"},
							{Id: "second-place"},
							{Id: "third-place"},
						},
					},
				},
			},
		},
	}

	planCandidateRepository, err := NewPlanCandidateRepository(testDB)
	if err != nil {
		t.Fatalf("failed to create plan candidate repository: %v", err)
	}

	for _, c := range cases {
		testContext := context.Background()
		t.Run(c.name, func(t *testing.T) {
			t.Cleanup(func() {
				err := cleanup(testContext, testDB)
				if err != nil {
					t.Fatalf("failed to cleanup: %v", err)
				}
			})

			// 事前に Place を作成しておく
			placesInPlanCandidates := array.Map(c.savedPlanCandidateSet.Plans, func(plan models.Plan) []models.Place { return plan.Places })
			if err := savePlaces(testContext, testDB, array.Flatten(placesInPlanCandidates)); err != nil {
				t.Fatalf("failed to save places: %v", err)
			}
			if err := savePlaces(testContext, testDB, []models.Place{c.placeToReplace}); err != nil {
				t.Fatalf("failed to save places: %v", err)
			}

			// 事前にPlanCandidateSetを作成しておく
			if err := savePlanCandidateSet(testContext, testDB, c.savedPlanCandidateSet); err != nil {
				t.Fatalf("failed to save plan candidate: %v", err)
			}

			if err := planCandidateRepository.ReplacePlace(testContext, c.planCandidateSetId, c.planCandidateId, c.placeIdToReplace, c.placeToReplace); err != nil {
				t.Fatalf("failed to replace place: %v", err)
			}

			planCandidatePlaceEntityExist, err := generated.PlanCandidatePlaces(
				generated.PlanCandidatePlaceWhere.PlanCandidateSetID.EQ(c.planCandidateSetId),
				generated.PlanCandidatePlaceWhere.PlanCandidateID.EQ(c.planCandidateId),
				generated.PlanCandidatePlaceWhere.PlaceID.EQ(c.placeToReplace.Id),
			).Exists(testContext, testDB)
			if err != nil {
				t.Fatalf("failed to get plan candidate place: %v", err)
			}

			if !planCandidatePlaceEntityExist {
				t.Fatalf("plan candidate place should exist")
			}
		})
	}
}

func TestPlanCandidateRepository_ReplacePlace_ShouldReturnError(t *testing.T) {
	cases := []struct {
		name                  string
		planCandidateSetId    string
		planCandidateId       string
		placeIdToReplace      string
		placeToReplace        models.Place
		savedPlanCandidateSet models.PlanCandidateSet
	}{
		{
			name:               "replace with not existing place",
			planCandidateSetId: "test-plan-candidate-set",
			planCandidateId:    "test-plan-candidate",
			placeIdToReplace:   "not-existing-place",
			placeToReplace:     models.Place{Id: "place-to-replace"},
			savedPlanCandidateSet: models.PlanCandidateSet{
				Id:        "test-plan-candidate-set",
				ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
				Plans: []models.Plan{
					{
						Id: "test-plan-candidate",
						Places: []models.Place{
							{Id: "first-place"},
							{Id: "second-place"},
							{Id: "third-place"},
						},
					},
				},
			},
		},
	}

	planCandidateRepository, err := NewPlanCandidateRepository(testDB)
	if err != nil {
		t.Fatalf("failed to create plan candidate repository: %v", err)
	}

	for _, c := range cases {
		testContext := context.Background()
		t.Run(c.name, func(t *testing.T) {
			t.Cleanup(func() {
				err := cleanup(testContext, testDB)
				if err != nil {
					t.Fatalf("failed to cleanup: %v", err)
				}
			})

			// 事前に Place を作成しておく
			placesInPlanCandidates := array.Map(c.savedPlanCandidateSet.Plans, func(plan models.Plan) []models.Place { return plan.Places })
			if err := savePlaces(testContext, testDB, array.Flatten(placesInPlanCandidates)); err != nil {
				t.Fatalf("failed to save places: %v", err)
			}
			if err := savePlaces(testContext, testDB, []models.Place{c.placeToReplace}); err != nil {
				t.Fatalf("failed to save places: %v", err)
			}

			// 事前にPlanCandidateSetを作成しておく
			if err := savePlanCandidateSet(testContext, testDB, c.savedPlanCandidateSet); err != nil {
				t.Fatalf("failed to save plan candidate: %v", err)
			}

			err := planCandidateRepository.ReplacePlace(testContext, c.planCandidateSetId, c.planCandidateId, c.placeIdToReplace, c.placeToReplace)
			if err == nil {
				t.Fatalf("error should be returned")
			}
		})
	}
}

func TestPlanCandidateRepository_UpdateLikeToPlaceInPlanCandidate_Like(t *testing.T) {
	cases := []struct {
		name                                   string
		planCandidateSetId                     string
		placeId                                string
		savedPlaces                            []models.Place
		savedPlanCandidate                     models.PlanCandidateSet
		savedPlanCandidateSetLikePlaceEntities []generated.PlanCandidateSetLikePlace
	}{
		{
			name:               "like from none",
			planCandidateSetId: "test-plan-candidate-set",
			placeId:            "test-place",
			savedPlaces: []models.Place{
				{Id: "test-place"},
			},
			savedPlanCandidate: models.PlanCandidateSet{
				Id:        "test-plan-candidate-set",
				ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
				Plans: []models.Plan{
					{
						Id:     "test-plan-candidate",
						Places: []models.Place{{Id: "test-place"}},
					},
				},
			},
			savedPlanCandidateSetLikePlaceEntities: []generated.PlanCandidateSetLikePlace{},
		},
		{
			name:               "like from like",
			planCandidateSetId: "test-plan-candidate-set",
			placeId:            "test-place",
			savedPlaces: []models.Place{
				{Id: "test-place"},
			},
			savedPlanCandidate: models.PlanCandidateSet{
				Id:        "test-plan-candidate-set",
				ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
				Plans: []models.Plan{
					{
						Id:     "test-plan-candidate",
						Places: []models.Place{{Id: "test-place"}},
					},
				},
			},
			savedPlanCandidateSetLikePlaceEntities: []generated.PlanCandidateSetLikePlace{
				{
					ID:                 uuid.New().String(),
					PlanCandidateSetID: "test-plan-candidate-set",
					PlaceID:            "test-place",
				},
			},
		},
	}

	planCandidateRepository, err := NewPlanCandidateRepository(testDB)
	if err != nil {
		t.Fatalf("failed to create plan candidate repository: %v", err)
	}

	for _, c := range cases {
		testContext := context.Background()
		t.Run(c.name, func(t *testing.T) {
			t.Cleanup(func() {
				err := cleanup(testContext, testDB)
				if err != nil {
					t.Fatalf("failed to cleanup: %v", err)
				}
			})

			// 事前に Place を作成しておく
			if err := savePlaces(testContext, testDB, c.savedPlaces); err != nil {
				t.Fatalf("failed to save places: %v", err)
			}

			// 事前に PlanCandidateSet を作成しておく
			if err := savePlanCandidateSet(testContext, testDB, c.savedPlanCandidate); err != nil {
				t.Fatalf("failed to save plan candidate: %v", err)
			}

			// 事前に PlanCandidateSetLikePlace を作成しておく
			for _, planCandidateSetLikePlaceEntity := range c.savedPlanCandidateSetLikePlaceEntities {
				if err := planCandidateSetLikePlaceEntity.Insert(testContext, testDB, boil.Infer()); err != nil {
					t.Fatalf("failed to save plan candidate set like place: %v", err)
				}
			}

			err := planCandidateRepository.UpdateLikeToPlaceInPlanCandidateSet(testContext, c.planCandidateSetId, c.placeId, true)
			if err != nil {
				t.Fatalf("failed to update like to place in plan candidate: %v", err)
			}

			numPlanCandidateSetLikePlaceEntity, err := generated.PlanCandidateSetLikePlaces(
				generated.PlanCandidateSetLikePlaceWhere.PlanCandidateSetID.EQ(c.planCandidateSetId),
				generated.PlanCandidateSetLikePlaceWhere.PlaceID.EQ(c.placeId),
			).Count(testContext, testDB)
			if err != nil {
				t.Fatalf("failed to get plan candidate set like place: %v", err)
			}

			if numPlanCandidateSetLikePlaceEntity != 1 {
				t.Fatalf("wrong number of plan candidate set like place expected: %v, actual: %v", 1, numPlanCandidateSetLikePlaceEntity)
			}
		})
	}
}

func TestPlanCandidateRepository_UpdateLikeToPlaceInPlanCandidate_Unlike(t *testing.T) {
	cases := []struct {
		name                                   string
		planCandidateSetId                     string
		placeId                                string
		savedPlaces                            []models.Place
		savedPlanCandidate                     models.PlanCandidateSet
		savedPlanCandidateSetLikePlaceEntities []generated.PlanCandidateSetLikePlace
	}{
		{
			name:               "unlike from like",
			planCandidateSetId: "test-plan-candidate-set",
			placeId:            "test-place",
			savedPlaces: []models.Place{
				{Id: "test-place"},
			},
			savedPlanCandidate: models.PlanCandidateSet{
				Id:        "test-plan-candidate-set",
				ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
				Plans: []models.Plan{
					{
						Id:     "test-plan-candidate",
						Places: []models.Place{{Id: "test-place"}},
					},
				},
			},
			savedPlanCandidateSetLikePlaceEntities: []generated.PlanCandidateSetLikePlace{
				{
					ID:                 uuid.New().String(),
					PlanCandidateSetID: "test-plan-candidate-set",
					PlaceID:            "test-place",
				},
			},
		},
		{
			name:               "unlike from none",
			planCandidateSetId: "test-plan-candidate-set",
			placeId:            "test-place",
			savedPlaces: []models.Place{
				{Id: "test-place"},
			},
			savedPlanCandidate: models.PlanCandidateSet{
				Id:        "test-plan-candidate-set",
				ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
				Plans: []models.Plan{
					{
						Id:     "test-plan-candidate",
						Places: []models.Place{{Id: "test-place"}},
					},
				},
			},
			savedPlanCandidateSetLikePlaceEntities: []generated.PlanCandidateSetLikePlace{},
		},
	}

	planCandidateRepository, err := NewPlanCandidateRepository(testDB)
	if err != nil {
		t.Fatalf("failed to create plan candidate repository: %v", err)
	}

	for _, c := range cases {
		testContext := context.Background()
		t.Run(c.name, func(t *testing.T) {
			t.Cleanup(func() {
				err := cleanup(testContext, testDB)
				if err != nil {
					t.Fatalf("failed to cleanup: %v", err)
				}
			})

			// 事前に Place を作成しておく
			if err := savePlaces(testContext, testDB, c.savedPlaces); err != nil {
				t.Fatalf("failed to save places: %v", err)
			}

			// 事前に PlanCandidateSet を作成しておく
			if err := savePlanCandidateSet(testContext, testDB, c.savedPlanCandidate); err != nil {
				t.Fatalf("failed to save plan candidate: %v", err)
			}

			// 事前に PlanCandidateSetLikePlace を作成しておく
			for _, planCandidateSetLikePlaceEntity := range c.savedPlanCandidateSetLikePlaceEntities {
				if err := planCandidateSetLikePlaceEntity.Insert(testContext, testDB, boil.Infer()); err != nil {
					t.Fatalf("failed to save plan candidate set like place: %v", err)
				}
			}

			err := planCandidateRepository.UpdateLikeToPlaceInPlanCandidateSet(testContext, c.planCandidateSetId, c.placeId, false)
			if err != nil {
				t.Fatalf("failed to update like to place in plan candidate: %v", err)
			}

			isPlanCandidateSetLikePlaceEntityExist, err := generated.PlanCandidateSetLikePlaces(
				generated.PlanCandidateSetLikePlaceWhere.PlanCandidateSetID.EQ(c.planCandidateSetId),
				generated.PlanCandidateSetLikePlaceWhere.PlaceID.EQ(c.placeId),
			).Exists(testContext, testDB)
			if err != nil {
				t.Fatalf("failed to get plan candidate set like place: %v", err)
			}

			if isPlanCandidateSetLikePlaceEntityExist {
				t.Fatalf("plan candidate set like place should not exist")
			}
		})
	}
}
