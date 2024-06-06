package rdb

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/repository"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
	"testing"
	"time"
)

func TestPlanRepository_Save(t *testing.T) {
	cases := []struct {
		name                       string
		savedUsers                 generated.UserSlice
		savedPlaces                []models.Place
		savedPlans                 generated.PlanSlice
		savedPlanCandidateSets     generated.PlanCandidateSetSlice
		savedPlanCandidates        generated.PlanCandidateSlice
		savedPlanCandidatePlaces   generated.PlanCandidatePlaceSlice
		plan                       models.Plan
		expectedPlan               generated.Plan
		expectedPlanPlaces         generated.PlanPlaceSlice
		expectedPlanParentChildren *generated.PlanParentChild
	}{
		{
			name: "should save plan",
			savedUsers: generated.UserSlice{
				{ID: "user_id_1", FirebaseUID: "firebase_uid_1"},
			},
			savedPlaces: []models.Place{
				{Id: "place_id_1"},
				{Id: "place_id_2"},
			},
			savedPlans: generated.PlanSlice{
				{ID: "plan_parent"},
			},
			savedPlanCandidateSets: generated.PlanCandidateSetSlice{
				{ID: "plan_candidate_set_id_1", ExpiresAt: time.Now()},
			},
			savedPlanCandidates: generated.PlanCandidateSlice{
				{
					ID:                 "plan_candidate_id_1",
					PlanCandidateSetID: "plan_candidate_set_id_1",
					Name:               "plan title",
					ParentPlanID:       null.StringFrom("plan_parent"),
				},
			},
			savedPlanCandidatePlaces: generated.PlanCandidatePlaceSlice{
				{
					ID:                 uuid.New().String(),
					PlanCandidateID:    "plan_candidate_id_1",
					PlanCandidateSetID: "plan_candidate_set_id_1",
					PlaceID:            "place_id_1",
					SortOrder:          0,
				},
				{
					ID:                 uuid.New().String(),
					PlanCandidateID:    "plan_candidate_id_1",
					PlanCandidateSetID: "plan_candidate_set_id_1",
					PlaceID:            "place_id_2",
					SortOrder:          1,
				},
			},
			plan: models.Plan{
				Id:   "plan",
				Name: "plan title",
				Places: []models.Place{
					{
						Id: "place_id_1",
						Location: models.GeoLocation{
							Latitude:  35.681236,
							Longitude: 139.767125,
						},
					},
					{Id: "place_id_2"},
				},
				Author: &models.User{
					Id:          "user_id_1",
					FirebaseUID: "firebase_uid_1",
				},
				ParentPlanId: utils.ToPointer("plan_parent"),
			},
			expectedPlan: generated.Plan{
				ID:        "plan",
				Name:      "plan title",
				UserID:    null.StringFrom("user_id_1"),
				Latitude:  35.681236,
				Longitude: 139.767125,
			},
			expectedPlanPlaces: generated.PlanPlaceSlice{
				{
					PlanID:    "plan",
					PlaceID:   "place_id_1",
					SortOrder: 0,
				},
				{
					PlanID:    "plan",
					PlaceID:   "place_id_2",
					SortOrder: 1,
				},
			},
			expectedPlanParentChildren: &generated.PlanParentChild{
				ParentPlanID: "plan_parent",
				ChildPlanID:  "plan",
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

			// 事前にデータを保存
			if _, err := c.savedUsers.InsertAll(testContext, planRepository.GetDB(), boil.Infer()); err != nil {
				t.Errorf("error saving user: %v", err)
			}

			if err := savePlaces(testContext, planRepository.GetDB(), c.savedPlaces); err != nil {
				t.Errorf("error saving places: %v", err)
			}

			if _, err := c.savedPlans.InsertAll(testContext, planRepository.GetDB(), boil.Infer()); err != nil {
				t.Errorf("error saving plan: %v", err)
			}

			if _, err := c.savedPlanCandidateSets.InsertAll(testContext, planRepository.GetDB(), boil.Infer()); err != nil {
				t.Errorf("error saving plan candidate set: %v", err)
			}

			if _, err := c.savedPlanCandidates.InsertAll(testContext, planRepository.GetDB(), boil.Infer()); err != nil {
				t.Errorf("error saving plan candidate: %v", err)
			}

			if _, err := c.savedPlanCandidatePlaces.InsertAll(testContext, planRepository.GetDB(), boil.Infer()); err != nil {
				t.Errorf("error saving plan candidate place: %v", err)
			}

			if err := planRepository.Save(testContext, &c.plan); err != nil {
				t.Errorf("error saving plan: %v", err)
			}

			planEntity, err := generated.Plans(generated.PlanWhere.ID.EQ(c.plan.Id)).One(testContext, planRepository.GetDB())
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					t.Errorf("plan should be saved")
				}
				t.Errorf("error checking if plan exists: %v", err)
			}

			if diff := cmp.Diff(
				c.expectedPlan,
				*planEntity,
				cmpopts.SortSlices(func(a, b generated.PlanPlace) bool { return a.PlaceID < b.PlaceID }),
				cmpopts.IgnoreFields(generated.Plan{}, "ID", "CreatedAt", "UpdatedAt"),
			); diff != "" {
				t.Errorf("plan mismatch (-want +got):\n%s", diff)
			}

			// プランに紐づく場所を確認
			planPlaceEntities, err := generated.PlanPlaces(
				generated.PlanPlaceWhere.PlanID.EQ(c.plan.Id),
				qm.OrderBy(generated.PlanPlaceColumns.SortOrder),
			).All(testContext, planRepository.GetDB())
			if err != nil {
				t.Errorf("error fetching plan places: %v", err)
			}

			if diff := cmp.Diff(
				c.expectedPlanPlaces,
				planPlaceEntities,
				cmpopts.SortSlices(func(a, b generated.PlanPlace) bool { return a.PlaceID < b.PlaceID }),
				cmpopts.IgnoreFields(generated.PlanPlace{}, "ID", "CreatedAt", "UpdatedAt"),
			); diff != "" {
				t.Errorf("plan places mismatch (-want +got):\n%s", diff)
			}

			// プランの親子関係を確認
			planParentChildEntiy, err := generated.PlanParentChildren(
				generated.PlanParentChildWhere.ChildPlanID.EQ(c.plan.Id),
			).One(testContext, planRepository.GetDB())
			if err != nil {
				t.Errorf("error fetching plan parent children: %v", err)
			}

			if diff := cmp.Diff(
				c.expectedPlanParentChildren,
				planParentChildEntiy,
				cmpopts.IgnoreFields(generated.PlanParentChild{}, "ID", "CreatedAt", "UpdatedAt"),
			); diff != "" {
				t.Errorf("plan parent children mismatch (-want +got):\n%s", diff)
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

func TestPlanRepository_Find_WithPlaceLikeCount(t *testing.T) {
	cases := []struct {
		name                                   string
		savedPlaces                            []models.Place
		savedUsers                             generated.UserSlice
		savedPlans                             []models.Plan
		savedPlanCandidateSets                 generated.PlanCandidateSetSlice
		savedPlanCandidateSetLikePlaceEntities generated.PlanCandidateSetLikePlaceSlice
		savedUserLikePlaceEntities             generated.UserLikePlaceSlice
		planId                                 string
		expected                               models.Plan
	}{
		{
			name: "should find plan with place like count",
			savedPlaces: []models.Place{
				{Id: "test-place-1", Google: models.GooglePlace{PlaceId: "test-google-place-1"}},
				{Id: "test-place-2", Google: models.GooglePlace{PlaceId: "test-google-place-2"}},
			},
			savedUsers: generated.UserSlice{
				{ID: "test-user-1", FirebaseUID: uuid.New().String()},
				{ID: "test-user-2", FirebaseUID: uuid.New().String()},
			},
			savedPlans: []models.Plan{
				{
					Id:   "test-plan-1",
					Name: "plan title",
					Places: []models.Place{
						{Id: "test-place-1"},
						{Id: "test-place-2"},
					},
				},
			},
			savedPlanCandidateSets: generated.PlanCandidateSetSlice{
				{ID: "test-plan-candidate-set-1", ExpiresAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local)},
				{ID: "test-plan-candidate-set-2", ExpiresAt: time.Date(2020, 12, 2, 0, 0, 0, 0, time.Local)},
			},
			savedPlanCandidateSetLikePlaceEntities: generated.PlanCandidateSetLikePlaceSlice{
				{ID: uuid.New().String(), PlanCandidateSetID: "test-plan-candidate-set-1", PlaceID: "test-place-1"},
				{ID: uuid.New().String(), PlanCandidateSetID: "test-plan-candidate-set-1", PlaceID: "test-place-2"},
				{ID: uuid.New().String(), PlanCandidateSetID: "test-plan-candidate-set-2", PlaceID: "test-place-1"},
			},
			savedUserLikePlaceEntities: generated.UserLikePlaceSlice{
				{ID: uuid.New().String(), UserID: "test-user-1", PlaceID: "test-place-1"},
				{ID: uuid.New().String(), UserID: "test-user-1", PlaceID: "test-place-2"},
				{ID: uuid.New().String(), UserID: "test-user-2", PlaceID: "test-place-1"},
			},
			planId: "test-plan-1",
			expected: models.Plan{
				Id:   "test-plan-1",
				Name: "plan title",
				Places: []models.Place{
					{
						Id:        "test-place-1",
						Google:    models.GooglePlace{PlaceId: "test-google-place-1"},
						LikeCount: 4,
					},
					{
						Id:        "test-place-2",
						Google:    models.GooglePlace{PlaceId: "test-google-place-2"},
						LikeCount: 2,
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

			// 事前に User・Place・Plan・PlanCandidateSet・PlanCandidateSetLikePlace・UserLikePlace を保存
			if _, err := c.savedUsers.InsertAll(textContext, planRepository.GetDB(), boil.Infer()); err != nil {
				t.Errorf("error saving user: %v", err)
			}

			if err := savePlaces(textContext, planRepository.GetDB(), c.savedPlaces); err != nil {
				t.Errorf("error saving places: %v", err)
			}

			if err := savePlans(textContext, planRepository.GetDB(), c.savedPlans); err != nil {
				t.Errorf("error saving plan: %v", err)
			}

			if _, err := c.savedPlanCandidateSets.InsertAll(textContext, planRepository.GetDB(), boil.Infer()); err != nil {
				t.Errorf("error saving plan candidate set: %v", err)
			}

			if _, err := c.savedPlanCandidateSetLikePlaceEntities.InsertAll(textContext, planRepository.GetDB(), boil.Infer()); err != nil {
				t.Errorf("error saving plan candidate set like place: %v", err)
			}

			if _, err := c.savedUserLikePlaceEntities.InsertAll(textContext, planRepository.GetDB(), boil.Infer()); err != nil {
				t.Errorf("error saving user like place: %v", err)
			}

			plan, err := planRepository.Find(textContext, c.planId)
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

func TestPlanRepository_UpdatePlanAuthorUserByPlanCandidateSet(t *testing.T) {
	cases := []struct {
		name                   string
		userId                 string
		planCandidateSetIds    []string
		savedUsers             generated.UserSlice
		savedPlans             generated.PlanSlice
		savedPlanCandidateSets generated.PlanCandidateSetSlice
		savedPlanCandidates    generated.PlanCandidateSlice
		expectedUserPlans      generated.PlanSlice
	}{
		{
			name:   "should update plan author user",
			userId: "8fde8eff-4b18-4276-b71f-2fec30ea65c8",
			planCandidateSetIds: []string{
				"d65ecb97-99f5-474e-a349-79fa888b37f5",
			},
			savedUsers: generated.UserSlice{
				{ID: "8fde8eff-4b18-4276-b71f-2fec30ea65c8"},
			},
			savedPlans: generated.PlanSlice{
				{
					ID:     "f2c98d68-3904-455b-8832-a0f723a96735",
					UserID: null.StringFromPtr(nil),
				},
			},
			savedPlanCandidateSets: generated.PlanCandidateSetSlice{
				{ID: "d65ecb97-99f5-474e-a349-79fa888b37f5", ExpiresAt: time.Now()},
			},
			savedPlanCandidates: generated.PlanCandidateSlice{
				{
					ID:                 "f2c98d68-3904-455b-8832-a0f723a96735",
					PlanCandidateSetID: "d65ecb97-99f5-474e-a349-79fa888b37f5",
				},
			},
			expectedUserPlans: generated.PlanSlice{
				{
					ID:     "f2c98d68-3904-455b-8832-a0f723a96735",
					UserID: null.StringFrom("8fde8eff-4b18-4276-b71f-2fec30ea65c8"),
				},
			},
		},
		{
			name:   "should not update plan author user",
			userId: "8fde8eff-4b18-4276-b71f-2fec30ea65c8",
			planCandidateSetIds: []string{
				"d65ecb97-99f5-474e-a349-79fa888b37f5",
			},
			savedUsers: generated.UserSlice{
				{ID: "8fde8eff-4b18-4276-b71f-2fec30ea65c8"},
			},
			savedPlans: generated.PlanSlice{
				{
					ID:     "f2c98d68-3904-455b-8832-a0f723a96735",
					UserID: null.StringFromPtr(nil),
				},
			},
			savedPlanCandidateSets: generated.PlanCandidateSetSlice{
				{ID: "d65ecb97-99f5-474e-a349-79fa888b37f5", ExpiresAt: time.Now()},
			},
			savedPlanCandidates: generated.PlanCandidateSlice{
				{
					ID:                 "a05c61a5-2974-4a8f-9914-8639088481a8",
					PlanCandidateSetID: "d65ecb97-99f5-474e-a349-79fa888b37f5",
				},
			},
			expectedUserPlans: nil,
		},
	}

	testContext := context.Background()
	planRepository, err := NewPlanRepository(testDB)
	if err != nil {
		t.Errorf("error initializing plan repository: %v", err)
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Cleanup(func() {
				if err := cleanup(testContext, planRepository.GetDB()); err != nil {
					t.Errorf("error cleaning up: %v", err)
				}
			})

			// 事前に保存
			if _, err := c.savedUsers.InsertAll(testContext, planRepository.GetDB(), boil.Infer()); err != nil {
				t.Errorf("error saving user: %v", err)
			}

			if _, err := c.savedPlans.InsertAll(testContext, planRepository.GetDB(), boil.Infer()); err != nil {
				t.Errorf("error saving plan: %v", err)
			}

			if _, err := c.savedPlanCandidateSets.InsertAll(testContext, planRepository.GetDB(), boil.Infer()); err != nil {
				t.Errorf("error saving plan candidate set: %v", err)
			}

			if _, err := c.savedPlanCandidates.InsertAll(testContext, planRepository.GetDB(), boil.Infer()); err != nil {
				t.Errorf("error saving plan candidate: %v", err)
			}

			err := planRepository.UpdatePlanAuthorUserByPlanCandidateSet(testContext, c.userId, c.planCandidateSetIds)
			if err != nil {
				t.Errorf("error updating plan author user: %v", err)
			}

			actualUserPlans, err := generated.Plans(
				generated.PlanWhere.UserID.EQ(null.StringFrom(c.userId)),
			).All(testContext, planRepository.GetDB())
			if err != nil {
				t.Errorf("error finding user plans: %v", err)
			}

			if diff := cmp.Diff(
				c.expectedUserPlans,
				actualUserPlans,
				cmpopts.SortSlices(func(a, b *models.Plan) bool { return a.Id < b.Id }),
				cmpopts.IgnoreFields(generated.Plan{}, "CreatedAt", "UpdatedAt"),
			); diff != "" {
				t.Errorf("user plan mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
