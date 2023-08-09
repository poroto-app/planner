package entity

import (
	"reflect"
	"testing"

	"poroto.app/poroto/planner/internal/domain/models"
)

func TestFromPlanEntity(t *testing.T) {
	cases := []struct {
		name            string
		planId          string
		planName        string
		places          []PlaceEntity
		timeInMinutes   int
		placeIdsOrdered []string
		transitions     *[]TransitionsEntity
		expected        models.Plan
	}{
		{
			name:     "test01",
			planId:   "I01",
			planName: "N01",
			places: []PlaceEntity{
				{
					Id: "1",
				},
				{
					Id: "2",
				},
				{
					Id: "3",
				},
				{
					Id: "4",
				},
			},
			timeInMinutes:   30,
			placeIdsOrdered: []string{"3", "2", "4", "1"},
			transitions:     &[]TransitionsEntity{},
			expected: models.Plan{
				Id:   "I01",
				Name: "N01",
				Places: []models.Place{
					{
						Id: "3",
					},
					{
						Id: "2",
					},
					{
						Id: "4",
					},
					{
						Id: "1",
					},
				},
				Transitions:   []models.Transition{},
				TimeInMinutes: 30,
			},
		},
		{
			name:     "test02",
			planId:   "I02",
			planName: "N02",
			places: []PlaceEntity{
				{
					Id: "2",
				},
				{
					Id: "4",
				},
				{
					Id: "1",
				},
				{
					Id: "3",
				},
			},
			timeInMinutes:   60,
			placeIdsOrdered: []string{"4", "1", "3", "2"},
			transitions:     &[]TransitionsEntity{},
			expected: models.Plan{
				Id:   "I02",
				Name: "N02",
				Places: []models.Place{
					{
						Id: "4",
					},
					{
						Id: "1",
					},
					{
						Id: "3",
					},
					{
						Id: "2",
					},
				},
				Transitions:   []models.Transition{},
				TimeInMinutes: 60,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := fromPlanEntity(
				c.planId,
				c.planName,
				c.places,
				c.timeInMinutes,
				c.placeIdsOrdered,
				c.transitions,
			)
			if !reflect.DeepEqual(c.expected, actual) {
				t.Errorf("expected: %v\nactual: %v", c.expected, actual)
			}
		})
	}
}
