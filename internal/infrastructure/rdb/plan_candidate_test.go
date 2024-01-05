package rdb

import (
	"context"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
	"testing"
	"time"
)

func TestPlanCandidateRepository_Create(t *testing.T) {
	cases := []struct {
		name            string
		planCandidateId string
		expiresAt       time.Time
	}{
		{
			name:            "success",
			planCandidateId: uuid.New().String(),
			expiresAt:       time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
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

			if err := planCandidateRepository.Create(testContext, c.planCandidateId, c.expiresAt); err != nil {
				t.Fatalf("failed to create plan candidate: %v", err)
			}

			exists, err := generated.PlanCandidateSetExists(testContext, testDB, c.planCandidateId)
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
		savedPlanCandidateSet models.PlanCandidate
		planCandidateId       string
		expected              models.PlanCandidate
	}{
		{
			name: "plan candidate set with only id",
			now:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
			savedPlanCandidateSet: models.PlanCandidate{
				Id:        "test",
				ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
			},
			planCandidateId: "test",
			expected: models.PlanCandidate{
				Id:        "test",
				ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
			},
		},
		{
			name: "plan candidate set with plan candidate",
			now:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
			savedPlanCandidateSet: models.PlanCandidate{
				Id:        "test",
				ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
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
					},
				},
			},
			planCandidateId: "test",
			expected: models.PlanCandidate{
				Id:        "test",
				ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
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
							{Id: "test-place"},
						},
					},
				},
			},
		},
		{
			name:            "plan candidate set without PlanCandidateSetMetaData",
			now:             time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
			planCandidateId: "test",
			savedPlanCandidateSet: models.PlanCandidate{
				Id:        "test",
				ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
				Plans: []models.Plan{
					{
						Id: "test-plan",
						Places: []models.Place{
							{Id: "test-place", Google: models.GooglePlace{PlaceId: "test-google-place"}},
						},
					},
				},
			},
			expected: models.PlanCandidate{
				Id:        "test",
				ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
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
			placesInPlanCandidates := array.Flatten(array.Map(c.savedPlanCandidateSet.Plans, func(plan models.Plan) []models.Place { return plan.Places }))
			if err := savePlaces(testContext, testDB, placesInPlanCandidates); err != nil {
				t.Fatalf("failed to save places: %v", err)
			}

			// 事前にPlanCandidateSetを作成しておく
			if err := savePlanCandidate(testContext, testDB, c.savedPlanCandidateSet); err != nil {
				t.Fatalf("failed to save plan candidate: %v", err)
			}

			actual, err := planCandidateRepository.Find(testContext, c.planCandidateId, c.now)
			if err != nil {
				t.Fatalf("failed to find plan candidate: %v", err)
			}

			if actual == nil {
				t.Fatalf("plan candidate should be found")
			}

			// Id の値が一致する
			if actual.Id != c.expected.Id {
				t.Fatalf("wrong plan candidate id expected: %v, actual: %v", c.expected.Id, actual.Id)
			}

			// ExpiresAt の値が一致する
			if !actual.ExpiresAt.Equal(c.expected.ExpiresAt) {
				t.Fatalf("wrong plan candidate expires at expected: %v, actual: %v", c.expected.ExpiresAt, actual.ExpiresAt)
			}

			// Plans の数が一致する
			if len(actual.Plans) != len(c.expected.Plans) {
				t.Fatalf("wrong number of plans expected: %v, actual: %v", len(c.expected.Plans), len(actual.Plans))
			}

			// Plan の順番が一致する
			for i, plan := range actual.Plans {
				if plan.Id != c.expected.Plans[i].Id {
					t.Fatalf("wrong plan id expected: %v, actual: %v", c.expected.Plans[i].Id, plan.Id)
				}

				// Place の数が一致する
				if len(plan.Places) != len(c.expected.Plans[i].Places) {
					t.Fatalf("wrong number of placesInPlanCandidates expected: %v, actual: %v", len(c.expected.Plans[i].Places), len(plan.Places))
				}

				// Place の順番が一致する
				for j, place := range plan.Places {
					if place.Id != c.expected.Plans[i].Places[j].Id {
						t.Fatalf("wrong place id expected: %v, actual: %v", c.expected.Plans[i].Places[j].Id, place.Id)
					}
				}
			}

			// MetaData の値が一致する
			if actual.MetaData.CreatedBasedOnCurrentLocation != c.expected.MetaData.CreatedBasedOnCurrentLocation {
				t.Fatalf("wrong created based on current location expected: %v, actual: %v", c.expected.MetaData.CreatedBasedOnCurrentLocation, actual.MetaData.CreatedBasedOnCurrentLocation)
			}

			if valueOrZero(actual.MetaData.FreeTime) != valueOrZero(c.expected.MetaData.FreeTime) {
				t.Fatalf("wrong free time expected: %v, actual: %v", valueOrZero(c.expected.MetaData.FreeTime), valueOrZero(actual.MetaData.FreeTime))
			}

			if !valueOrZero(actual.MetaData.LocationStart).Equal(valueOrZero(c.expected.MetaData.LocationStart)) {
				t.Fatalf("wrong location start expected: %v, actual: %v", valueOrZero(c.expected.MetaData.LocationStart), valueOrZero(actual.MetaData.LocationStart))
			}

			if len(valueOrZero(actual.MetaData.CategoriesPreferred)) != len(valueOrZero(c.expected.MetaData.CategoriesPreferred)) {
				t.Fatalf("wrong number of categories preferred expected: %v, actual: %v", len(valueOrZero(c.expected.MetaData.CategoriesPreferred)), len(valueOrZero(actual.MetaData.CategoriesPreferred)))
			}

			if len(valueOrZero(actual.MetaData.CategoriesRejected)) != len(valueOrZero(c.expected.MetaData.CategoriesRejected)) {
				t.Fatalf("wrong number of categories rejected expected: %v, actual: %v", len(valueOrZero(c.expected.MetaData.CategoriesRejected)), len(valueOrZero(actual.MetaData.CategoriesRejected)))
			}
		})
	}
}

func TestPlanCandidateRepository_Find_ShouldReturnNil(t *testing.T) {
	cases := []struct {
		name                  string
		now                   time.Time
		savedPlanCandidateSet models.PlanCandidate
		planCandidateId       string
	}{
		{
			name: "expired plan candidate set will not be returned",
			now:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
			savedPlanCandidateSet: models.PlanCandidate{
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
			if err := savePlanCandidate(testContext, testDB, c.savedPlanCandidateSet); err != nil {
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
		savedPlanCandidateSets                 []models.PlanCandidate
		savedPlanCandidateSetLikePlaceEntities []generated.PlanCandidateSetLikePlace
		planCandidateId                        string
		expected                               models.PlanCandidate
	}{
		{
			name: "plan candidate set with place like count",
			now:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
			savedPlaces: []models.Place{
				{Id: "test-place-1", Google: models.GooglePlace{PlaceId: "test-google-place-1"}},
				{Id: "test-place-2", Google: models.GooglePlace{PlaceId: "test-google-place-2"}},
			},
			savedPlanCandidateSets: []models.PlanCandidate{
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
				{ID: uuid.New().String(), PlanCandidateSetID: "plan-candidate-set-1", PlaceID: "test-place-2"},
				{ID: uuid.New().String(), PlanCandidateSetID: "plan-candidate-set-2", PlaceID: "test-place-1"},
			},
			planCandidateId: "plan-candidate-set-1",
			expected: models.PlanCandidate{
				Id:        "plan-candidate-set-1",
				ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
				Plans: []models.Plan{
					{
						Id: "plan-candidate-1",
						Places: []models.Place{
							{Id: "test-place-1", Google: models.GooglePlace{PlaceId: "test-google-place-1"}, LikeCount: 2},
							{Id: "test-place-2", Google: models.GooglePlace{PlaceId: "test-google-place-2"}, LikeCount: 1},
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
		textContext := context.Background()
		t.Run(c.name, func(t *testing.T) {
			t.Cleanup(func() {
				err := cleanup(textContext, testDB)
				if err != nil {
					t.Fatalf("failed to cleanup: %v", err)
				}
			})

			// 事前にPlaceを作成しておく
			if err := savePlaces(textContext, testDB, c.savedPlaces); err != nil {
				t.Fatalf("failed to save places: %v", err)
			}

			// 事前にPlanCandidateSetを作成しておく
			for _, planCandidateSet := range c.savedPlanCandidateSets {
				if err := savePlanCandidate(textContext, testDB, planCandidateSet); err != nil {
					t.Fatalf("failed to save plan candidate: %v", err)
				}
			}

			// 事前にPlanCandidateSetLikePlaceを作成しておく
			for _, planCandidateSetLikePlaceEntity := range c.savedPlanCandidateSetLikePlaceEntities {
				if err := planCandidateSetLikePlaceEntity.Insert(textContext, testDB, boil.Infer()); err != nil {
					t.Fatalf("failed to save plan candidate set like place: %v", err)
				}
			}

			actual, err := planCandidateRepository.Find(textContext, c.planCandidateId, c.now)
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
		savedPlanCandidateSet models.PlanCandidate
		expected              models.Plan
	}{
		{
			name:               "success",
			planCandidateSetId: "test-plan-candidate-set",
			planCandidateId:    "test-plan-candidate",
			savedPlanCandidateSet: models.PlanCandidate{
				Id:        "test-plan-candidate-set",
				ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
				Plans: []models.Plan{
					{
						Id: "test-plan-candidate",
						Places: []models.Place{
							{Id: "test-place", Google: models.GooglePlace{PlaceId: "test-google-place"}},
						},
					},
				},
			},
			expected: models.Plan{
				Id: "test-plan-candidate",
				Places: []models.Place{
					{Id: "test-place", Google: models.GooglePlace{PlaceId: "test-google-place"}},
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
			placesInPlans := array.Map(c.savedPlanCandidateSet.Plans, func(plan models.Plan) []models.Place { return plan.Places })
			if err := savePlaces(testContext, testDB, array.Flatten(placesInPlans)); err != nil {
				t.Fatalf("failed to save places: %v", err)
			}

			// 事前にPlanCandidateSetを作成しておく
			if err := savePlanCandidate(testContext, testDB, c.savedPlanCandidateSet); err != nil {
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

			// Place の数が一致する
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

func TestPlanCandidateRepository_FindExpiredBefore(t *testing.T) {
	cases := []struct {
		name                   string
		expiresAt              time.Time
		savedPlanCandidateSets []models.PlanCandidate
		expected               []string
	}{
		{
			name:      "success",
			expiresAt: time.Date(2020, 12, 1, 12, 0, 0, 0, time.Local),
			savedPlanCandidateSets: []models.PlanCandidate{
				{
					Id:        "test-plan-candidate-set-1",
					ExpiresAt: time.Date(2020, 12, 1, 12, 0, 0, 0, time.Local),
				},
				{
					Id:        "test-plan-candidate-set-2",
					ExpiresAt: time.Date(2020, 12, 1, 11, 59, 59, 0, time.Local),
				},
			},
			expected: []string{"test-plan-candidate-set-1"},
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
			placesInPlans := array.Flatten(array.Flatten(array.Map(c.savedPlanCandidateSets, func(planCandidate models.PlanCandidate) [][]models.Place {
				return array.Map(planCandidate.Plans, func(plan models.Plan) []models.Place { return plan.Places })
			})))
			if err := savePlaces(testContext, testDB, placesInPlans); err != nil {
				t.Fatalf("failed to save places: %v", err)
			}

			// 事前にPlanCandidateSetを作成しておく
			for _, planCandidateSet := range c.savedPlanCandidateSets {
				if err := savePlanCandidate(testContext, testDB, planCandidateSet); err != nil {
					t.Fatalf("failed to save plan candidate: %v", err)
				}
			}

			actual, err := planCandidateRepository.FindExpiredBefore(testContext, c.expiresAt)
			if err != nil {
				t.Fatalf("failed to find expired plan candidates: %v", err)
			}

			if actual == nil {
				t.Fatalf("expired plan candidates should be found")
			}

			if len(*actual) != len(c.expected) {
				t.Fatalf("wrong number of expired plan candidates expected: %v, actual: %v", len(c.expected), len(*actual))
			}
		})

	}
}

func TestPlanCandidateRepository_AddSearchedPlacesForPlanCandidate(t *testing.T) {
	cases := []struct {
		name            string
		planCandidateId string
		placeIds        []string
	}{
		{
			name:            "success",
			planCandidateId: uuid.New().String(),
			placeIds:        []string{uuid.New().String(), uuid.New().String()},
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

			// 事前にPlanCandidateSetを作成しておく]
			if err := savePlanCandidate(testContext, testDB, models.PlanCandidate{Id: c.planCandidateId, ExpiresAt: time.Now().Add(time.Hour)}); err != nil {
				t.Fatalf("failed to create plan candidate: %v", err)
			}

			// 事前にPlaceを作成しておく
			for _, placeId := range c.placeIds {
				placeEntity := generated.Place{ID: placeId}
				if err := placeEntity.Insert(testContext, testDB, boil.Infer()); err != nil {
					t.Fatalf("failed to insert place: %v", err)
				}
			}

			if err := planCandidateRepository.AddSearchedPlacesForPlanCandidate(testContext, c.planCandidateId, c.placeIds); err != nil {
				t.Fatalf("failed to add searched places for plan candidate: %v", err)
			}

			numPlanCandidateSetSearchedPlaces, err := generated.
				PlanCandidateSetSearchedPlaces(generated.PlanCandidateSetSearchedPlaceWhere.PlanCandidateSetID.EQ(c.planCandidateId)).
				Count(testContext, testDB)
			if err != nil {
				t.Fatalf("failed to get plan candidate places: %v", err)
			}

			if int(numPlanCandidateSetSearchedPlaces) != len(c.placeIds) {
				t.Fatalf("number of plan candidate places is not expected: %v", numPlanCandidateSetSearchedPlaces)
			}
		})
	}
}

func TestPlanCandidateRepository_AddPlan(t *testing.T) {
	cases := []struct {
		name            string
		planCandidateId string
		plans           []models.Plan
	}{
		{
			name:            "success",
			planCandidateId: uuid.New().String(),
			plans: []models.Plan{
				{
					Id: uuid.New().String(),
					Places: []models.Place{
						{Id: "tokyo-station"},
						{Id: "shinagawa-station"},
					},
				},
				{
					Id: uuid.New().String(),
					Places: []models.Place{
						{Id: "yokohama-station"},
						{Id: "shin-yokohama-station"},
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

			// 事前にPlaceを作成しておく
			placesInPlans := array.Map(c.plans, func(plan models.Plan) []models.Place { return plan.Places })
			if err := savePlaces(testContext, testDB, array.Flatten(placesInPlans)); err != nil {
				t.Fatalf("failed to save places: %v", err)
			}

			// 事前にPlanCandidateSetを作成しておく
			if err := savePlanCandidate(testContext, testDB, models.PlanCandidate{Id: c.planCandidateId, ExpiresAt: time.Now().Add(time.Hour)}); err != nil {
				t.Fatalf("failed to create plan candidate: %v", err)
			}

			if err := planCandidateRepository.AddPlan(testContext, c.planCandidateId, c.plans...); err != nil {
				t.Fatalf("failed to add plan: %v", err)
			}

			// すべてのPlanCandidateが保存されている
			numPlanCandidates, err := generated.
				PlanCandidates(generated.PlanCandidateWhere.PlanCandidateSetID.EQ(c.planCandidateId)).
				Count(testContext, testDB)
			if err != nil {
				t.Fatalf("failed to get plan candidates: %v", err)
			}
			if int(numPlanCandidates) != len(c.plans) {
				t.Fatalf("wrong number of plan candidates expected: %v, actual: %v", len(c.plans), numPlanCandidates)
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
		savedPlanCandidateSet models.PlanCandidate
	}{
		{
			name:               "success",
			planCandidateSetId: "test-plan-candidate-set",
			planCandidateId:    "test-plan-candidate",
			placeIdToDelete:    "second-place",
			savedPlanCandidateSet: models.PlanCandidate{
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
			savedPlanCandidateSet: models.PlanCandidate{
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
			if err := savePlanCandidate(testContext, testDB, c.savedPlanCandidateSet); err != nil {
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
		savedPlanCandidateSet models.PlanCandidate
	}{
		{
			name:               "success",
			planCandidateSetId: "test-plan-candidate-set",
			planCandidateId:    "test-plan-candidate",
			placeIdsOrdered:    []string{"third-place", "first-place", "second-place"},
			savedPlanCandidateSet: models.PlanCandidate{
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
			if err := savePlanCandidate(testContext, testDB, c.savedPlanCandidateSet); err != nil {
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
		savedPlanCandidateSet models.PlanCandidate
	}{
		{
			name:               "reorder with not existing place",
			planCandidateSetId: "test-plan-candidate-set",
			planCandidateId:    "test-plan-candidate",
			placeIdsOrdered:    []string{"third-place", "first-place", "not-existing-place"},
			savedPlanCandidateSet: models.PlanCandidate{
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
			if err := savePlanCandidate(testContext, testDB, c.savedPlanCandidateSet); err != nil {
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
		name                  string
		planCandidateSetId    string
		savedPlanCandidateSet models.PlanCandidate
		metaData              models.PlanCandidateMetaData
	}{
		{
			name:               "save plan candidate meta data",
			planCandidateSetId: "test-plan-candidate-set",
			savedPlanCandidateSet: models.PlanCandidate{
				Id:        "test-plan-candidate-set",
				ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
				MetaData:  models.PlanCandidateMetaData{},
			},
			metaData: models.PlanCandidateMetaData{
				CreatedBasedOnCurrentLocation: true,
				CategoriesPreferred:           &[]models.LocationCategory{models.CategoryRestaurant},
				CategoriesRejected:            &[]models.LocationCategory{models.CategoryBookStore},
				LocationStart:                 &models.GeoLocation{Latitude: 35.681236, Longitude: 139.767125},
				FreeTime:                      toPointer(60),
			},
		},
		{
			name:               "update plan candidate meta data",
			planCandidateSetId: "test-plan-candidate-set",
			savedPlanCandidateSet: models.PlanCandidate{
				Id:        "test-plan-candidate-set",
				ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
				MetaData: models.PlanCandidateMetaData{
					CreatedBasedOnCurrentLocation: false,
					CategoriesPreferred:           &[]models.LocationCategory{models.CategoryRestaurant},
					CategoriesRejected:            &[]models.LocationCategory{models.CategoryBookStore},
					LocationStart:                 &models.GeoLocation{Latitude: 35.681236, Longitude: 139.767125},
					FreeTime:                      toPointer(60),
				},
			},
			metaData: models.PlanCandidateMetaData{
				CreatedBasedOnCurrentLocation: true,
				CategoriesPreferred:           &[]models.LocationCategory{models.CategoryRestaurant, models.CategoryBookStore},
				CategoriesRejected:            &[]models.LocationCategory{models.CategoryShopping, models.CategoryAmusements},
				LocationStart:                 &models.GeoLocation{Latitude: 36.681236, Longitude: 140.767125},
				FreeTime:                      toPointer(120),
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

			// 事前にPlanCandidateSetを作成しておく
			if err := savePlanCandidate(testContext, testDB, c.savedPlanCandidateSet); err != nil {
				t.Fatalf("failed to save plan candidate: %v", err)
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

			if planCandidateSetMetaDataEntity.IsCreatedFromCurrentLocation != c.metaData.CreatedBasedOnCurrentLocation {
				t.Fatalf("wrong is created from current location expected: %v, actual: %v", c.metaData.CreatedBasedOnCurrentLocation, planCandidateSetMetaDataEntity.IsCreatedFromCurrentLocation)
			}

			if planCandidateSetMetaDataEntity.LatitudeStart != c.metaData.LocationStart.Latitude {
				t.Fatalf("wrong latitude start expected: %v, actual: %v", c.metaData.LocationStart.Latitude, planCandidateSetMetaDataEntity.LatitudeStart)
			}

			if planCandidateSetMetaDataEntity.LongitudeStart != c.metaData.LocationStart.Longitude {
				t.Fatalf("wrong longitude start expected: %v, actual: %v", c.metaData.LocationStart.Longitude, planCandidateSetMetaDataEntity.LongitudeStart)
			}

			if planCandidateSetMetaDataEntity.PlanDurationMinutes.Int != *c.metaData.FreeTime {
				t.Fatalf("wrong plan duration minutes expected: %v, actual: %v", c.metaData.FreeTime, planCandidateSetMetaDataEntity.PlanDurationMinutes.Int)
			}

			// CategoriesPreferred が一致する
			numCategoriesPreferred, err := generated.
				PlanCandidateSetMetaDataCategories(
					generated.PlanCandidateSetMetaDataCategoryWhere.PlanCandidateSetID.EQ(c.planCandidateSetId),
					generated.PlanCandidateSetMetaDataCategoryWhere.IsSelected.EQ(true),
				).Count(testContext, testDB)
			if err != nil {
				t.Fatalf("failed to get plan candidate set meta data categories: %v", err)
			}
			if int(numCategoriesPreferred) != len(*c.metaData.CategoriesPreferred) {
				t.Fatalf("wrong number of plan candidate set meta data categories expected: %v, actual: %v", len(*c.metaData.CategoriesPreferred), numCategoriesPreferred)
			}

			// CategoriesRejected が一致する
			numCategoriesRejected, err := generated.
				PlanCandidateSetMetaDataCategories(
					generated.PlanCandidateSetMetaDataCategoryWhere.PlanCandidateSetID.EQ(c.planCandidateSetId),
					generated.PlanCandidateSetMetaDataCategoryWhere.IsSelected.EQ(false),
				).Count(testContext, testDB)
			if err != nil {
				t.Fatalf("failed to get plan candidate set meta data categories: %v", err)
			}
			if int(numCategoriesRejected) != len(*c.metaData.CategoriesRejected) {
				t.Fatalf("wrong number of plan candidate set meta data categories expected: %v, actual: %v", len(*c.metaData.CategoriesRejected), numCategoriesRejected)
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
		savedPlanCandidateSet models.PlanCandidate
	}{
		{
			name:               "success",
			planCandidateSetId: "test-plan-candidate-set",
			planCandidateId:    "test-plan-candidate",
			placeIdToReplace:   "second-place",
			placeToReplace:     models.Place{Id: "replaced-place"},
			savedPlanCandidateSet: models.PlanCandidate{
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
			if err := savePlanCandidate(testContext, testDB, c.savedPlanCandidateSet); err != nil {
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
		savedPlanCandidateSet models.PlanCandidate
	}{
		{
			name:               "replace with not existing place",
			planCandidateSetId: "test-plan-candidate-set",
			planCandidateId:    "test-plan-candidate",
			placeIdToReplace:   "not-existing-place",
			placeToReplace:     models.Place{Id: "place-to-replace"},
			savedPlanCandidateSet: models.PlanCandidate{
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
			if err := savePlanCandidate(testContext, testDB, c.savedPlanCandidateSet); err != nil {
				t.Fatalf("failed to save plan candidate: %v", err)
			}

			err := planCandidateRepository.ReplacePlace(testContext, c.planCandidateSetId, c.planCandidateId, c.placeIdToReplace, c.placeToReplace)
			if err == nil {
				t.Fatalf("error should be returned")
			}
		})
	}
}

func TestPlanCandidateRepository_DeleteAll(t *testing.T) {
	cases := []struct {
		name                        string
		savedPlanCandidateSets      []models.PlanCandidate
		planCandidateIdsToDelete    []string
		planCandidateIdsNotToDelete []string
	}{
		{
			name: "success",
			savedPlanCandidateSets: []models.PlanCandidate{
				{
					Id:        "test-plan-candidate-set-1",
					ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
					MetaData: models.PlanCandidateMetaData{
						CreatedBasedOnCurrentLocation: true,
						CategoriesPreferred:           &[]models.LocationCategory{models.CategoryRestaurant},
						CategoriesRejected:            &[]models.LocationCategory{models.CategoryBookStore},
						LocationStart:                 &models.GeoLocation{Latitude: 35.681236, Longitude: 139.767125},
						FreeTime:                      toPointer(60),
					},
					Plans: []models.Plan{
						{
							Id: "test-plan-candidate-1",
							Places: []models.Place{
								{Id: "first-place"},
								{Id: "second-place"},
								{Id: "third-place"},
							},
						},
					},
				},
				{
					Id:        "test-plan-candidate-set-2",
					ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
					MetaData: models.PlanCandidateMetaData{
						CreatedBasedOnCurrentLocation: true,
						CategoriesPreferred:           &[]models.LocationCategory{models.CategoryRestaurant},
						CategoriesRejected:            &[]models.LocationCategory{models.CategoryBookStore},
						LocationStart:                 &models.GeoLocation{Latitude: 35.681236, Longitude: 139.767125},
						FreeTime:                      toPointer(60),
					},
					Plans: []models.Plan{
						{
							Id: "test-plan-candidate-2",
							Places: []models.Place{
								{Id: "first-place"},
							},
						},
					},
				},
			},
			planCandidateIdsToDelete:    []string{"test-plan-candidate-set-1"},
			planCandidateIdsNotToDelete: []string{"test-plan-candidate-set-2"},
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
			var placesInPlanCandidates []models.Place
			for _, planCandidateSet := range c.savedPlanCandidateSets {
				for _, plan := range planCandidateSet.Plans {
					placesInPlanCandidates = append(placesInPlanCandidates, plan.Places...)
				}
			}
			if err := savePlaces(testContext, testDB, placesInPlanCandidates); err != nil {
				t.Fatalf("failed to save places: %v", err)
			}

			// 事前にPlanCandidateSetを作成しておく
			for _, planCandidateSet := range c.savedPlanCandidateSets {
				if err := savePlanCandidate(testContext, testDB, planCandidateSet); err != nil {
					t.Fatalf("failed to save plan candidate: %v", err)
				}
			}

			if err := planCandidateRepository.DeleteAll(testContext, c.planCandidateIdsToDelete); err != nil {
				t.Fatalf("failed to delete plan candidate: %v", err)
			}

			// 削除されていないことを確認
			for _, planCandidateId := range c.planCandidateIdsNotToDelete {
				planCandidateSetEntityExist, err := generated.PlanCandidateSets(
					generated.PlanCandidateSetWhere.ID.EQ(planCandidateId),
				).Exists(testContext, testDB)
				if err != nil {
					t.Fatalf("failed to get plan candidate set: %v", err)
				}

				if !planCandidateSetEntityExist {
					t.Fatalf("plan candidate set should exist")
				}
			}

			// 削除されていることを確認
			for _, planCandidateId := range c.planCandidateIdsToDelete {
				// PlanCandidateSet が削除されていることを確認
				planCandidateSetEntityExist, err := generated.PlanCandidateSets(
					generated.PlanCandidateSetWhere.ID.EQ(planCandidateId),
				).Exists(testContext, testDB)
				if err != nil {
					t.Fatalf("failed to get plan candidate set: %v", err)
				}
				if planCandidateSetEntityExist {
					t.Fatalf("plan candidate set should not exist")
				}

				// PlanCandidateSetMetaData が削除されていることを確認
				planCandidateSetMetaDataEntityExist, err := generated.PlanCandidateSetMetaData(
					generated.PlanCandidateSetMetaDatumWhere.PlanCandidateSetID.EQ(planCandidateId),
				).Exists(testContext, testDB)
				if err != nil {
					t.Fatalf("failed to get plan candidate set meta data: %v", err)
				}
				if planCandidateSetMetaDataEntityExist {
					t.Fatalf("plan candidate set meta data should not exist")
				}

				// PlanCandidateSetMetaDataCategory が削除されていることを確認
				planCandidateSetMetaDataCategoryEntityExist, err := generated.PlanCandidateSetMetaDataCategories(
					generated.PlanCandidateSetMetaDataCategoryWhere.PlanCandidateSetID.EQ(planCandidateId),
				).Exists(testContext, testDB)
				if err != nil {
					t.Fatalf("failed to get plan candidate set meta data category: %v", err)
				}
				if planCandidateSetMetaDataCategoryEntityExist {
					t.Fatalf("plan candidate set meta data category should not exist")
				}

				// PlanCandidate が削除されていることを確認
				planCandidateEntityExist, err := generated.PlanCandidates(
					generated.PlanCandidateWhere.ID.EQ(planCandidateId),
				).Exists(testContext, testDB)
				if err != nil {
					t.Fatalf("failed to get plan candidate: %v", err)
				}
				if planCandidateEntityExist {
					t.Fatalf("plan candidate should not exist")
				}

				// PlanCandidatePlace が削除されていることを確認
				planCandidatePlaceEntityExist, err := generated.PlanCandidatePlaces(
					generated.PlanCandidatePlaceWhere.PlanCandidateSetID.EQ(planCandidateId),
				).Exists(testContext, testDB)
				if err != nil {
					t.Fatalf("failed to get plan candidate place: %v", err)
				}
				if planCandidatePlaceEntityExist {
					t.Fatalf("plan candidate place should not exist")
				}

				// PlanCandidateSetSearchedPlace が削除されていることを確認
				planCandidateSetSearchedPlaceEntityExist, err := generated.PlanCandidateSetSearchedPlaces(
					generated.PlanCandidateSetSearchedPlaceWhere.PlanCandidateSetID.EQ(planCandidateId),
				).Exists(testContext, testDB)
				if err != nil {
					t.Fatalf("failed to get plan candidate set searched place: %v", err)
				}
				if planCandidateSetSearchedPlaceEntityExist {
					t.Fatalf("plan candidate set searched place should not exist")
				}
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
		savedPlanCandidate                     models.PlanCandidate
		savedPlanCandidateSetLikePlaceEntities []generated.PlanCandidateSetLikePlace
	}{
		{
			name:               "like from none",
			planCandidateSetId: "test-plan-candidate-set",
			placeId:            "test-place",
			savedPlaces: []models.Place{
				{Id: "test-place"},
			},
			savedPlanCandidate: models.PlanCandidate{
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
			savedPlanCandidate: models.PlanCandidate{
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
			if err := savePlanCandidate(testContext, testDB, c.savedPlanCandidate); err != nil {
				t.Fatalf("failed to save plan candidate: %v", err)
			}

			// 事前に PlanCandidateSetLikePlace を作成しておく
			for _, planCandidateSetLikePlaceEntity := range c.savedPlanCandidateSetLikePlaceEntities {
				if err := planCandidateSetLikePlaceEntity.Insert(testContext, testDB, boil.Infer()); err != nil {
					t.Fatalf("failed to save plan candidate set like place: %v", err)
				}
			}

			err := planCandidateRepository.UpdateLikeToPlaceInPlanCandidate(testContext, c.planCandidateSetId, c.placeId, true)
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
		savedPlanCandidate                     models.PlanCandidate
		savedPlanCandidateSetLikePlaceEntities []generated.PlanCandidateSetLikePlace
	}{
		{
			name:               "unlike from like",
			planCandidateSetId: "test-plan-candidate-set",
			placeId:            "test-place",
			savedPlaces: []models.Place{
				{Id: "test-place"},
			},
			savedPlanCandidate: models.PlanCandidate{
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
			savedPlanCandidate: models.PlanCandidate{
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
			if err := savePlanCandidate(testContext, testDB, c.savedPlanCandidate); err != nil {
				t.Fatalf("failed to save plan candidate: %v", err)
			}

			// 事前に PlanCandidateSetLikePlace を作成しておく
			for _, planCandidateSetLikePlaceEntity := range c.savedPlanCandidateSetLikePlaceEntities {
				if err := planCandidateSetLikePlaceEntity.Insert(testContext, testDB, boil.Infer()); err != nil {
					t.Fatalf("failed to save plan candidate set like place: %v", err)
				}
			}

			err := planCandidateRepository.UpdateLikeToPlaceInPlanCandidate(testContext, c.planCandidateSetId, c.placeId, false)
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

func toPointer[T any](value T) *T {
	return &value
}

func valueOrZero[T any](value *T) T {
	var zero T
	if value == nil {
		return zero
	}
	return *value
}
