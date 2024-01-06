package rdb

import (
	"context"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"poroto.app/poroto/planner/internal/domain/models"
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
				Id:       uuid.New().String(),
				Name:     "plan title",
				AuthorId: nil,
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
		savedPlaces []models.Place
		savedPlan   models.Plan
		expected    models.Plan
	}{
		{
			name: "should find plan",
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
		savedUser   models.User
		savedPlaces []models.Place
		savedPlans  []models.Plan
		authorId    string
		expected    []models.Plan
	}{
		{
			name:      "should find plans by author id",
			savedUser: models.User{Id: "c61a8b42-2c07-4957-913d-6930f0d881ec"},
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
					Id:       "f2c98d68-3904-455b-8832-a0f723a96735",
					AuthorId: toPointer("c61a8b42-2c07-4957-913d-6930f0d881ec"),
					Name:     "plan title 1",
					Places: []models.Place{
						{Id: "f2c98d68-3904-455b-8832-a0f723a96735"},
						{Id: "c61a8b42-2c07-4957-913d-6930f0d881ec"},
					},
				},
				{
					Id:       "c61a8b42-2c07-4957-913d-6930f0d881ec",
					AuthorId: toPointer("c61a8b42-2c07-4957-913d-6930f0d881ec"),
					Name:     "plan title 2",
					Places: []models.Place{
						{Id: "82cc1884-ac59-11ee-a506-0242ac120002"},
						{Id: "c72a3614-ac59-11ee-a506-0242ac120002"},
					},
				},
			},
			authorId: "c61a8b42-2c07-4957-913d-6930f0d881ec",
			expected: []models.Plan{
				{
					Id:       "f2c98d68-3904-455b-8832-a0f723a96735",
					AuthorId: toPointer("c61a8b42-2c07-4957-913d-6930f0d881ec"),
					Name:     "plan title 1",
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
				},
				{
					Id:       "c61a8b42-2c07-4957-913d-6930f0d881ec",
					AuthorId: toPointer("c61a8b42-2c07-4957-913d-6930f0d881ec"),
					Name:     "plan title 2",
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

			// 事前に Place を保存
			if err := savePlaces(textContext, planRepository.GetDB(), c.savedPlaces); err != nil {
				t.Errorf("error saving places: %v", err)
			}

			// 事前に User を保存
			if err := saveUsers(textContext, planRepository.GetDB(), []models.User{c.savedUser}); err != nil {
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
		queryCursor         *string
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
					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        "c61a8b42-2c07-4957-913d-6930f0d881ec",
					Name:      "plan title 2",
					CreatedAt: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
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
			queryCursor: toPointer(newQueryCursor(
				time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
				"c61a8b42-2c07-4957-913d-6930f0d881ec",
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

			plans, err := planRepository.SortedByCreatedAt(textContext, c.queryCursor, c.limit)
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