package rdb

import (
	"context"
	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
	"testing"
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
