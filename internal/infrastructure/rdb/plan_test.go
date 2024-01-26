package rdb

import (
	"context"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/repository"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
	"testing"
	"time"
)

func TestPlanRepository_Save(t *testing.T) {
	cases := []struct {
		name        string
		savedPlaces []models.Place
		plan        models.Plan
	}{
		{
			name: "should save plan",
			savedPlaces: []models.Place{
				{
					Id: "f2c98d68-3904-455b-8832-a0f723a96735",
				},
				{
					Id: "c61a8b42-2c07-4957-913d-6930f0d881ec",
				},
			},
			plan: models.Plan{
				Id:   uuid.New().String(),
				Name: "plan title",
				Places: []models.Place{
					{
						Id: "f2c98d68-3904-455b-8832-a0f723a96735",
					},
					{
						Id: "c61a8b42-2c07-4957-913d-6930f0d881ec",
					},
				},
			},
		},
	}

	planRepository, err := NewPlanRepository(testDB)
	if err != nil {
		t.Errorf("error initializing plan repository: %v", err)
	}

	for _, c := range cases {
		c := c
		testContext := context.Background()
		t.Run(c.name, func(t *testing.T) {
			t.Cleanup(func() {
				if err := cleanup(testContext, planRepository.GetDB()); err != nil {
					t.Errorf("error cleaning up: %v", err)
				}
			})

			// 事前に Place を保存
			if err := savePlaces(testContext, planRepository.GetDB(), c.savedPlaces); err != nil {
				t.Errorf("error saving places: %v", err)
			}

			if err := planRepository.Save(testContext, &c.plan); err != nil {
				t.Errorf("error saving plan: %v", err)
			}

			isPlanExists, err := generated.Plans(generated.PlanWhere.ID.EQ(c.plan.Id)).Exists(testContext, planRepository.GetDB())
			if err != nil {
				t.Errorf("error checking if plan exists: %v", err)
			}
			if !isPlanExists {
				t.Errorf("plan should be saved")
			}

			planPlaceEntities, err := generated.PlanPlaces(
				generated.PlanPlaceWhere.PlanID.EQ(c.plan.Id),
				qm.OrderBy(generated.PlanPlaceColumns.SortOrder),
			).All(testContext, planRepository.GetDB())
			if err != nil {
				t.Errorf("error fetching plan places: %v", err)
			}

			if len(planPlaceEntities) != len(c.plan.Places) {
				t.Errorf("wrong plan place slice length, want: %d, got: %d", len(c.plan.Places), len(planPlaceEntities))
			}

			for i, expected := range c.plan.Places {
				if planPlaceEntities[i].PlaceID != expected.Id {
					t.Errorf("wrong place id, want: %s, got: %s", expected.Id, planPlaceEntities[i].PlaceID)
				}
				if planPlaceEntities[i].SortOrder != i {
					t.Errorf("wrong sort order, want: %d, got: %d", i, planPlaceEntities[i].SortOrder)
				}
			}
		})
	}
}

func TestPlanRepository_Find(t *testing.T) {
	cases := []struct {
		name        string
		savedUsers  generated.UserSlice
		savedPlaces []models.Place
		savedPlan   models.Plan
		expected    models.Plan
	}{
		{
			name: "should find plan",
			savedUsers: generated.UserSlice{
				{
					ID:          "8fde8eff-4b18-4276-b71f-2fec30ea65c8",
					FirebaseUID: "firebase_uid_1",
				},
			},
			savedPlaces: []models.Place{
				{
					Id: "f2c98d68-3904-455b-8832-a0f723a96735",
					Google: models.GooglePlace{
						PlaceId: "ChIJN1t_tDeuEmsRUsoyG83frY4",
					},
				},
				{
					Id: "c61a8b42-2c07-4957-913d-6930f0d881ec",
					Google: models.GooglePlace{
						PlaceId: "CwVXAAAAQwXg3w8QKxQZ6Q0X3Z4",
					},
				},
			},
			savedPlan: models.Plan{
				Id:   "f2c98d68-3904-455b-8832-a0f723a96735",
				Name: "plan title",
				Places: []models.Place{
					{Id: "f2c98d68-3904-455b-8832-a0f723a96735"},
					{Id: "c61a8b42-2c07-4957-913d-6930f0d881ec"},
				},
				Author: &models.User{
					Id:          "8fde8eff-4b18-4276-b71f-2fec30ea65c8",
					FirebaseUID: "firebase_uid_1",
				},
			},
			expected: models.Plan{
				Id:   "f2c98d68-3904-455b-8832-a0f723a96735",
				Name: "plan title",
				Places: []models.Place{
					{
						Id: "f2c98d68-3904-455b-8832-a0f723a96735",
						Google: models.GooglePlace{
							PlaceId: "ChIJN1t_tDeuEmsRUsoyG83frY4",
						},
					},
					{
						Id: "c61a8b42-2c07-4957-913d-6930f0d881ec",
						Google: models.GooglePlace{
							PlaceId: "CwVXAAAAQwXg3w8QKxQZ6Q0X3Z4",
						},
					},
				},
				Author: &models.User{
					Id:          "8fde8eff-4b18-4276-b71f-2fec30ea65c8",
					FirebaseUID: "firebase_uid_1",
				},
			},
		},
		{
			name: "plan with empty places should be found",
			savedPlan: models.Plan{
				Id:     "f2c98d68-3904-455b-8832-a0f723a96735",
				Name:   "plan title",
				Places: []models.Place{},
			},
			expected: models.Plan{
				Id:     "f2c98d68-3904-455b-8832-a0f723a96735",
				Name:   "plan title",
				Places: []models.Place{},
			},
		},
	}

	planRepository, err := NewPlanRepository(testDB)
	if err != nil {
		t.Errorf("error initializing plan repository: %v", err)
	}

	for _, c := range cases {
		c := c
		textContext := context.Background()
		t.Run(c.name, func(t *testing.T) {
			t.Cleanup(func() {
				if err := cleanup(textContext, planRepository.GetDB()); err != nil {
					t.Errorf("error cleaning up: %v", err)
				}
			})

			// 事前に User を保存
			if _, err := c.savedUsers.InsertAll(textContext, planRepository.GetDB(), boil.Infer()); err != nil {
				t.Errorf("error saving user: %v", err)
			}

			// 事前に Place を保存
			if err := savePlaces(textContext, planRepository.GetDB(), c.savedPlaces); err != nil {
				t.Errorf("error saving places: %v", err)
			}

			// 事前に Plan を保存
			if err := savePlans(textContext, planRepository.GetDB(), []models.Plan{c.savedPlan}); err != nil {
				t.Errorf("error saving plan: %v", err)
			}

			plan, err := planRepository.Find(textContext, c.savedPlan.Id)
			if err != nil {
				t.Errorf("error finding plan: %v", err)
			}

			if diff := cmp.Diff(c.expected, *plan); diff != "" {
				t.Errorf("plan mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPlanRepository_FindByAuthorId(t *testing.T) {
	cases := []struct {
		name        string
		savedUsers  generated.UserSlice
		savedPlaces []models.Place
		savedPlans  []models.Plan
		authorId    string
		expected    []models.Plan
	}{
		{
			name: "should find plans by author id",
			savedUsers: generated.UserSlice{
				{
					ID:          "8fde8eff-4b18-4276-b71f-2fec30ea65c8",
					FirebaseUID: "firebase_uid_1",
				},
			},
			savedPlaces: []models.Place{
				{
					Id:     "f2c98d68-3904-455b-8832-a0f723a96735",
					Google: models.GooglePlace{PlaceId: "ChIJN1t_tDeuEmsRUsoyG83frY4"},
				},
				{
					Id:     "c61a8b42-2c07-4957-913d-6930f0d881ec",
					Google: models.GooglePlace{PlaceId: "CwVXAAAAQwXg3w8QKxQZ6Q0X3Z4"},
				},
				{
					Id:     "82cc1884-ac59-11ee-a506-0242ac120002",
					Google: models.GooglePlace{PlaceId: "CAIR4P4SN0tdMjtLOLpNBa0HM2X4"},
				},
				{
					Id:     "c72a3614-ac59-11ee-a506-0242ac120002",
					Google: models.GooglePlace{PlaceId: "ChQFzO1BvHaLJXCyWVZ8Uie3smn2wQ"},
				},
			},
			savedPlans: []models.Plan{
				{
					Id:   "f2c98d68-3904-455b-8832-a0f723a96735",
					Name: "plan title 1",
					Places: []models.Place{
						{Id: "f2c98d68-3904-455b-8832-a0f723a96735"},
						{Id: "c61a8b42-2c07-4957-913d-6930f0d881ec"},
					},
					Author: &models.User{
						Id:          "8fde8eff-4b18-4276-b71f-2fec30ea65c8",
						FirebaseUID: "firebase_uid_1",
					},
				},
				{
					Id:   "c61a8b42-2c07-4957-913d-6930f0d881ec",
					Name: "plan title 2",
					Places: []models.Place{
						{Id: "82cc1884-ac59-11ee-a506-0242ac120002"},
						{Id: "c72a3614-ac59-11ee-a506-0242ac120002"},
					},
					Author: &models.User{
						Id:          "8fde8eff-4b18-4276-b71f-2fec30ea65c8",
						FirebaseUID: "firebase_uid_1",
					},
				},
			},
			authorId: "8fde8eff-4b18-4276-b71f-2fec30ea65c8",
			expected: []models.Plan{
				{
					Id:   "f2c98d68-3904-455b-8832-a0f723a96735",
					Name: "plan title 1",
					Places: []models.Place{
						{
							Id:     "f2c98d68-3904-455b-8832-a0f723a96735",
							Google: models.GooglePlace{PlaceId: "ChIJN1t_tDeuEmsRUsoyG83frY4"},
						},
						{
							Id:     "c61a8b42-2c07-4957-913d-6930f0d881ec",
							Google: models.GooglePlace{PlaceId: "CwVXAAAAQwXg3w8QKxQZ6Q0X3Z4"},
						},
					},
					Author: &models.User{
						Id:          "8fde8eff-4b18-4276-b71f-2fec30ea65c8",
						FirebaseUID: "firebase_uid_1",
					},
				},
				{
					Id:   "c61a8b42-2c07-4957-913d-6930f0d881ec",
					Name: "plan title 2",
					Places: []models.Place{
						{
							Id:     "82cc1884-ac59-11ee-a506-0242ac120002",
							Google: models.GooglePlace{PlaceId: "CAIR4P4SN0tdMjtLOLpNBa0HM2X4"},
						},
						{
							Id:     "c72a3614-ac59-11ee-a506-0242ac120002",
							Google: models.GooglePlace{PlaceId: "ChQFzO1BvHaLJXCyWVZ8Uie3smn2wQ"},
						},
					},
					Author: &models.User{
						Id:          "8fde8eff-4b18-4276-b71f-2fec30ea65c8",
						FirebaseUID: "firebase_uid_1",
					},
				},
			},
		},
		{
			name: "should not find plans by author id",
			savedUsers: generated.UserSlice{
				{
					ID:          "8fde8eff-4b18-4276-b71f-2fec30ea65c8",
					FirebaseUID: "firebase_uid_1",
				},
			},
			savedPlaces: []models.Place{
				{
					Id:     "f2c98d68-3904-455b-8832-a0f723a96735",
					Google: models.GooglePlace{PlaceId: "ChIJN1t_tDeuEmsRUsoyG83frY4"},
				},
			},
			savedPlans: []models.Plan{
				{
					Id:   "f2c98d68-3904-455b-8832-a0f723a96735",
					Name: "plan title 1",
					Places: []models.Place{
						{Id: "f2c98d68-3904-455b-8832-a0f723a96735"},
					},
				},
			},
			authorId: "28a52fdd-c252-4e32-a918-fcab5ed88ad8",
			expected: []models.Plan{},
		},
	}

	planRepository, err := NewPlanRepository(testDB)
	if err != nil {
		t.Errorf("error initializing plan repository: %v", err)
	}

	for _, c := range cases {
		c := c
		textContext := context.Background()
		t.Run(c.name, func(t *testing.T) {
			t.Cleanup(func() {
				if err := cleanup(textContext, planRepository.GetDB()); err != nil {
					t.Errorf("error cleaning up: %v", err)
				}
			})

			// 事前に Place を保存
			if err := savePlaces(textContext, planRepository.GetDB(), c.savedPlaces); err != nil {
				t.Errorf("error saving places: %v", err)
			}

			// 事前に User を保存
			if _, err := c.savedUsers.InsertAll(textContext, planRepository.GetDB(), boil.Infer()); err != nil {
				t.Errorf("error saving user: %v", err)
			}

			// 事前に Plan を保存
			if err := savePlans(textContext, planRepository.GetDB(), c.savedPlans); err != nil {
				t.Errorf("error saving plan: %v", err)
			}

			plans, err := planRepository.FindByAuthorId(textContext, c.authorId)
			if err != nil {
				t.Errorf("error finding plans: %v", err)
			}

			if diff := cmp.Diff(
				c.expected,
				*plans,
				cmpopts.SortSlices(func(a, b models.Plan) bool { return a.Id < b.Id }),
			); diff != "" {
				t.Errorf("plan mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPlanRepository_SortedByCreatedAt(t *testing.T) {
	cases := []struct {
		name                string
		savedPlanSlice      generated.PlanSlice
		savedPlanPlaceSlice generated.PlanPlaceSlice
		savedPlaces         []models.Place
		savedUsers          generated.UserSlice
		queryCursor         *repository.SortedByCreatedAtQueryCursor
		limit               int
		expected            []models.Plan
	}{
		{
			name: "should find plans sorted by created_at",
			savedPlaces: []models.Place{
				{
					Id:     "f2c98d68-3904-455b-8832-a0f723a96735",
					Google: models.GooglePlace{PlaceId: "ChIJN1t_tDeuEmsRUsoyG83frY4"},
				},
				{
					Id:     "c61a8b42-2c07-4957-913d-6930f0d881ec",
					Google: models.GooglePlace{PlaceId: "CwVXAAAAQwXg3w8QKxQZ6Q0X3Z4"},
				},
			},
			savedPlanSlice: []*generated.Plan{
				{
					ID:        "f2c98d68-3904-455b-8832-a0f723a96735",
					Name:      "plan title 1",
					UserID:    null.StringFrom("8fde8eff-4b18-4276-b71f-2fec30ea65c8"),
					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        "c61a8b42-2c07-4957-913d-6930f0d881ec",
					Name:      "plan title 2",
					UserID:    null.StringFrom("28a52fdd-c252-4e32-a918-fcab5ed88ad8"),
					CreatedAt: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			savedUsers: generated.UserSlice{
				{
					ID:          "8fde8eff-4b18-4276-b71f-2fec30ea65c8",
					FirebaseUID: "firebase_uid_1",
				},
				{
					ID:          "28a52fdd-c252-4e32-a918-fcab5ed88ad8",
					FirebaseUID: "firebase_uid_2",
				},
			},
			savedPlanPlaceSlice: generated.PlanPlaceSlice{
				{
					ID:        uuid.New().String(),
					PlanID:    "f2c98d68-3904-455b-8832-a0f723a96735",
					PlaceID:   "f2c98d68-3904-455b-8832-a0f723a96735",
					SortOrder: 0,
				},
				{
					ID:        uuid.New().String(),
					PlanID:    "c61a8b42-2c07-4957-913d-6930f0d881ec",
					PlaceID:   "c61a8b42-2c07-4957-913d-6930f0d881ec",
					SortOrder: 0,
				},
			},
			queryCursor: nil,
			limit:       10,
			expected: []models.Plan{
				{
					Id:   "c61a8b42-2c07-4957-913d-6930f0d881ec",
					Name: "plan title 2",
					Places: []models.Place{
						{
							Id:     "c61a8b42-2c07-4957-913d-6930f0d881ec",
							Google: models.GooglePlace{PlaceId: "CwVXAAAAQwXg3w8QKxQZ6Q0X3Z4"},
						},
					},
					Author: &models.User{
						Id:          "28a52fdd-c252-4e32-a918-fcab5ed88ad8",
						FirebaseUID: "firebase_uid_2",
					},
				},
				{
					Id:   "f2c98d68-3904-455b-8832-a0f723a96735",
					Name: "plan title 1",
					Places: []models.Place{
						{
							Id:     "f2c98d68-3904-455b-8832-a0f723a96735",
							Google: models.GooglePlace{PlaceId: "ChIJN1t_tDeuEmsRUsoyG83frY4"},
						},
					},
					Author: &models.User{
						Id:          "8fde8eff-4b18-4276-b71f-2fec30ea65c8",
						FirebaseUID: "firebase_uid_1",
					},
				},
			},
		},
		{
			name: "should find plans sorted by created_at with query cursor",
			savedPlanSlice: []*generated.Plan{
				{
					ID:        "f2c98d68-3904-455b-8832-a0f723a96735",
					Name:      "plan title 1",
					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        "c61a8b42-2c07-4957-913d-6930f0d881ec",
					Name:      "plan title 2",
					CreatedAt: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			queryCursor: toPointer(newSortByCreatedAtQueryCursor(
				time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
			)),
			limit: 10,
			expected: []models.Plan{
				{
					Id:     "f2c98d68-3904-455b-8832-a0f723a96735",
					Name:   "plan title 1",
					Places: []models.Place{},
				},
			},
		},
	}

	planRepository, err := NewPlanRepository(testDB)
	if err != nil {
		t.Errorf("error initializing plan repository: %v", err)
	}

	for _, c := range cases {
		c := c
		textContext := context.Background()
		t.Run(c.name, func(t *testing.T) {
			t.Cleanup(func() {
				if err := cleanup(textContext, planRepository.GetDB()); err != nil {
					t.Errorf("error cleaning up: %v", err)
				}
			})

			// 事前に User を保存
			if _, err := c.savedUsers.InsertAll(textContext, planRepository.GetDB(), boil.Infer()); err != nil {
				t.Errorf("error saving user: %v", err)
			}

			// 事前に Place を保存
			if err := savePlaces(textContext, planRepository.GetDB(), c.savedPlaces); err != nil {
				t.Errorf("error saving places: %v", err)
			}

			// 事前に Plan を保存
			if _, err := c.savedPlanSlice.InsertAll(textContext, planRepository.GetDB(), boil.Infer()); err != nil {
				t.Errorf("error saving plan: %v", err)
			}

			// 事前に PlanPlace を保存
			if _, err := c.savedPlanPlaceSlice.InsertAll(textContext, planRepository.GetDB(), boil.Infer()); err != nil {
				t.Errorf("error saving plan place: %v", err)
			}

			plans, _, err := planRepository.SortedByCreatedAt(textContext, c.queryCursor, c.limit)
			if err != nil {
				t.Errorf("error finding plans: %v", err)
			}

			if diff := cmp.Diff(
				c.expected,
				*plans,
				cmpopts.SortSlices(func(a, b models.Plan) bool { return a.Id < b.Id }),
			); diff != "" {
				t.Errorf("plan mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPlanRepository_SortedByLocation(t *testing.T) {
	cases := []struct {
		name        string
		savedUsers  generated.UserSlice
		savedPlaces []models.Place
		savedPlans  []models.Plan
		location    models.GeoLocation
		limit       int
		expected    []models.Plan
	}{
		{
			name: "should find plans sorted by location",
			savedUsers: generated.UserSlice{
				{
					ID:          "8fde8eff-4b18-4276-b71f-2fec30ea65c8",
					FirebaseUID: "firebase_uid_1",
				},
				{
					ID:          "28a52fdd-c252-4e32-a918-fcab5ed88ad8",
					FirebaseUID: "firebase_uid_2",
				},
			},
			savedPlaces: []models.Place{
				{
					Id:   "f2c98d68-3904-455b-8832-a0f723a96735",
					Name: "高島屋新宿店",
					Google: models.GooglePlace{
						PlaceId:  "ChIJN1t_tDeuEmsRUsoyG83frY4",
						Location: models.GeoLocation{Latitude: 35.687684359569, Longitude: 139.70220602474},
					},
				},
				{
					Id:   "c61a8b42-2c07-4957-913d-6930f0d881ec",
					Name: "札幌市時計台",
					Google: models.GooglePlace{
						PlaceId:  "CwVXAAAAQwXg3w8QKxQZ6Q0X3Z4",
						Location: models.GeoLocation{Latitude: 43.062558697622, Longitude: 141.35355044447},
					},
				},
			},
			savedPlans: []models.Plan{
				{
					Id:   "9c93c944-ac8e-11ee-a506-0242ac120002",
					Name: "新宿",
					Places: []models.Place{
						{
							Id:       "f2c98d68-3904-455b-8832-a0f723a96735",
							Name:     "高島屋新宿店",
							Location: models.GeoLocation{Latitude: 35.687684359569, Longitude: 139.70220602474},
						},
					},
					Author: &models.User{
						Id:          "8fde8eff-4b18-4276-b71f-2fec30ea65c8",
						FirebaseUID: "firebase_uid_1",
					},
				},
				{
					Id:   "9c93c944-ac8e-11ee-a506-0242ac120003",
					Name: "札幌",
					Places: []models.Place{
						{
							Id:       "c61a8b42-2c07-4957-913d-6930f0d881ec",
							Name:     "札幌市時計台",
							Location: models.GeoLocation{Latitude: 43.062558697622, Longitude: 141.35355044447},
						},
					},
					Author: &models.User{
						Id:          "28a52fdd-c252-4e32-a918-fcab5ed88ad8",
						FirebaseUID: "firebase_uid_2",
					},
				},
			},
			location: models.GeoLocation{Latitude: 35.6905, Longitude: 139.6995},
			limit:    10,
			expected: []models.Plan{
				{
					Id:   "9c93c944-ac8e-11ee-a506-0242ac120002",
					Name: "新宿",
					Places: []models.Place{
						{
							Id:       "f2c98d68-3904-455b-8832-a0f723a96735",
							Name:     "高島屋新宿店",
							Location: models.GeoLocation{Latitude: 35.687684359569, Longitude: 139.70220602474},
							Google: models.GooglePlace{
								PlaceId:  "ChIJN1t_tDeuEmsRUsoyG83frY4",
								Location: models.GeoLocation{Latitude: 35.687684359569, Longitude: 139.70220602474},
							},
						},
					},
					Author: &models.User{
						Id:          "8fde8eff-4b18-4276-b71f-2fec30ea65c8",
						FirebaseUID: "firebase_uid_1",
					},
				},
			},
		},
	}

	planRepository, err := NewPlanRepository(testDB)
	if err != nil {
		t.Errorf("error initializing plan repository: %v", err)
	}

	for _, c := range cases {
		textContext := context.Background()
		t.Run(c.name, func(t *testing.T) {
			t.Cleanup(func() {
				if err := cleanup(textContext, planRepository.GetDB()); err != nil {
					t.Errorf("error cleaning up: %v", err)
				}
			})

			// 事前に User を保存
			if _, err := c.savedUsers.InsertAll(textContext, planRepository.GetDB(), boil.Infer()); err != nil {
				t.Errorf("error saving user: %v", err)
			}

			// 事前に Place を保存
			if err := savePlaces(textContext, planRepository.GetDB(), c.savedPlaces); err != nil {
				t.Errorf("error saving places: %v", err)
			}

			// 事前に Plan を保存
			if err := savePlans(textContext, planRepository.GetDB(), c.savedPlans); err != nil {
				t.Errorf("error saving plan: %v", err)
			}

			plans, _, err := planRepository.SortedByLocation(textContext, c.location, nil, c.limit)
			if err != nil {
				t.Errorf("error finding plans: %v", err)
			}

			if diff := cmp.Diff(
				c.expected,
				*plans,
				cmpopts.SortSlices(func(a, b models.Plan) bool { return a.Id < b.Id }),
			); diff != "" {
				t.Errorf("plan mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
