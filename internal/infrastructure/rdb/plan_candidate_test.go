package rdb

import (
	"context"
	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/entities"
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

			exists, err := entities.PlanCandidateSetExists(testContext, testDB, c.planCandidateId)
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
		expected              *models.PlanCandidate
	}{
		{
			name: "plan candidate set with only id",
			now:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
			savedPlanCandidateSet: models.PlanCandidate{
				Id:        "test",
				ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
			},
			planCandidateId: "test",
			expected: &models.PlanCandidate{
				Id:        "test",
				ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
			},
		},
		{
			name: "expired plan candidate set will not be returned",
			now:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
			savedPlanCandidateSet: models.PlanCandidate{
				Id:        "test",
				ExpiresAt: time.Date(2019, 12, 1, 0, 0, 0, 0, time.Local),
			},
			planCandidateId: "test",
			expected:        nil,
		},
		{
			name: "plan candidate set with plan candidate",
			now:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
			savedPlanCandidateSet: models.PlanCandidate{
				Id:        "test",
				ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
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
			expected: &models.PlanCandidate{
				Id:        "test",
				ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
				MetaData: models.PlanCandidateMetaData{
					CreatedBasedOnCurrentLocation: true,
					LocationStart:                 &models.GeoLocation{Latitude: 139.767125, Longitude: 35.681236},
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
			name: "plan candidate set without PlanCandidateSetMetaData",
			now:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
			savedPlanCandidateSet: models.PlanCandidate{
				Id:        "test",
				ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
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
			for _, plan := range c.savedPlanCandidateSet.Plans {
				for _, place := range plan.Places {
					placeEntity := entities.Place{ID: place.Id}
					if err := placeEntity.Insert(testContext, testDB, boil.Infer()); err != nil {
						t.Fatalf("failed to insert place: %v", err)
					}

					googlePlaceEntity := entities.GooglePlace{PlaceID: place.Google.PlaceId}
					if err := placeEntity.AddGooglePlaces(testContext, testDB, true, &googlePlaceEntity); err != nil {
						t.Fatalf("failed to add google place: %v", err)
					}
				}
			}

			// 事前にPlanCandidateSetを作成しておく
			if err := planCandidateRepository.Create(testContext, c.savedPlanCandidateSet.Id, c.savedPlanCandidateSet.ExpiresAt); err != nil {
				t.Fatalf("failed to create plan candidate: %v", err)
			}

			// 事前にPlanCandidateSetMetaDataを作成しておく
			if !c.savedPlanCandidateSet.MetaData.IsZero() {
				planCandidateSetMetaDataEntity := entities.PlanCandidateSetMetaDatum{
					ID:                           uuid.New().String(),
					PlanCandidateSetID:           c.savedPlanCandidateSet.Id,
					IsCreatedFromCurrentLocation: c.savedPlanCandidateSet.MetaData.CreatedBasedOnCurrentLocation,
					LatitudeStart:                c.savedPlanCandidateSet.MetaData.LocationStart.Latitude,
					LongitudeStart:               c.savedPlanCandidateSet.MetaData.LocationStart.Longitude,
					PlanDurationMinutes:          null.IntFromPtr(c.savedPlanCandidateSet.MetaData.FreeTime),
				}
				if err := planCandidateSetMetaDataEntity.Insert(testContext, testDB, boil.Infer()); err != nil {
					t.Fatalf("failed to insert plan candidate set meta data: %v", err)
				}
			}

			// 事前にPlanCandidateを作成しておく
			if err := planCandidateRepository.AddPlan(testContext, c.savedPlanCandidateSet.Id, c.savedPlanCandidateSet.Plans...); err != nil {
				t.Fatalf("failed to add plan: %v", err)
			}

			actual, err := planCandidateRepository.Find(testContext, c.planCandidateId, c.now)
			if err != nil {
				t.Fatalf("failed to find plan candidate: %v", err)
			}

			if c.expected == nil {
				if actual != nil {
					t.Fatalf("plan candidate should not be found")
				}
				return
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
					t.Fatalf("wrong number of places expected: %v, actual: %v", len(c.expected.Plans[i].Places), len(plan.Places))
				}

				// Place の順番が一致する
				for j, place := range plan.Places {
					if place.Id != c.expected.Plans[i].Places[j].Id {
						t.Fatalf("wrong place id expected: %v, actual: %v", c.expected.Plans[i].Places[j].Id, place.Id)
					}
				}
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

			// 事前にPlanCandidateSetを作成しておく
			if err := planCandidateRepository.Create(testContext, c.planCandidateId, time.Now().Add(time.Hour)); err != nil {
				t.Fatalf("failed to create plan candidate: %v", err)
			}

			// 事前にPlaceを作成しておく
			for _, placeId := range c.placeIds {
				placeEntity := entities.Place{ID: placeId}
				if err := placeEntity.Insert(testContext, testDB, boil.Infer()); err != nil {
					t.Fatalf("failed to insert place: %v", err)
				}
			}

			if err := planCandidateRepository.AddSearchedPlacesForPlanCandidate(testContext, c.planCandidateId, c.placeIds); err != nil {
				t.Fatalf("failed to add searched places for plan candidate: %v", err)
			}

			numPlanCandidateSetSearchedPlaces, err := entities.
				PlanCandidateSetSearchedPlaces(entities.PlanCandidateSetSearchedPlaceWhere.PlanCandidateSetID.EQ(c.planCandidateId)).
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

			// 事前にPlanCandidateSetを作成しておく
			if err := planCandidateRepository.Create(testContext, c.planCandidateId, time.Now().Add(time.Hour)); err != nil {
				t.Fatalf("failed to create plan candidate: %v", err)
			}

			// 事前にPlaceを作成しておく
			for _, plan := range c.plans {
				for _, place := range plan.Places {
					placeEntity := entities.Place{ID: place.Id}
					if err := placeEntity.Insert(testContext, testDB, boil.Infer()); err != nil {
						t.Fatalf("failed to insert place: %v", err)
					}
				}
			}

			if err := planCandidateRepository.AddPlan(testContext, c.planCandidateId, c.plans...); err != nil {
				t.Fatalf("failed to add plan: %v", err)
			}

			// すべてのPlanCandidateが保存されている
			numPlanCandidates, err := entities.
				PlanCandidates(entities.PlanCandidateWhere.PlanCandidateSetID.EQ(c.planCandidateId)).
				Count(testContext, testDB)
			if err != nil {
				t.Fatalf("failed to get plan candidates: %v", err)
			}
			if int(numPlanCandidates) != len(c.plans) {
				t.Fatalf("wrong number of plan candidates expected: %v, actual: %v", len(c.plans), numPlanCandidates)
			}

			// すべてのPlanCandidateに対して、すべてのPlaceが保存されている
			for _, plan := range c.plans {
				numPlanCandidatePlaces, err := entities.
					PlanCandidatePlaces(entities.PlanCandidatePlaceWhere.PlanCandidateID.EQ(plan.Id)).
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
			var placeEntitySlice entities.PlaceSlice
			placeEntitySlice = append(placeEntitySlice, &entities.Place{ID: c.place.Id})
			for _, place := range c.savedPlanCandidatePlaces {
				placeEntitySlice = append(placeEntitySlice, &entities.Place{ID: place.Id})
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

			savedPlanCandidatePlaceSlice, err := entities.
				PlanCandidatePlaces(
					entities.PlanCandidatePlaceWhere.PlanCandidateID.EQ(c.planCandidateId),
					qm.OrderBy(entities.PlanCandidatePlaceColumns.SortOrder),
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
