package factory

import (
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
	"testing"
)

func TestNewPlanPlaceSliceFromDomainModel(t *testing.T) {
	cases := []struct {
		name       string
		planPlaces []models.Place
		planId     string
		expected   []generated.PlanPlace
	}{
		{
			name: "should return a valid slice",
			planPlaces: []models.Place{
				{
					Id: "ec7c607d-454a-4644-929a-c3b1e078842d",
				},
				{
					Id: "339809cf-d515-4a64-bbcd-c6a899051273",
				},
			},
			expected: []generated.PlanPlace{
				{
					PlaceID:   "ec7c607d-454a-4644-929a-c3b1e078842d",
					SortOrder: 0,
				},
				{
					PlaceID:   "339809cf-d515-4a64-bbcd-c6a899051273",
					SortOrder: 1,
				},
			},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			actual := NewPlanPlaceSliceFromDomainMode(c.planPlaces, c.planId)
			if len(actual) != len(c.expected) {
				t.Errorf("wrong plan place slice length, want: %d, got: %d", len(c.expected), len(actual))
			}

			for i, expected := range c.expected {
				if actual[i].PlaceID != expected.PlaceID {
					t.Errorf("wrong place id, want: %s, got: %s", expected.PlaceID, actual[i].PlaceID)
				}
				if actual[i].SortOrder != expected.SortOrder {
					t.Errorf("wrong sort order, want: %d, got: %d", expected.SortOrder, actual[i].SortOrder)
				}
			}
		})
	}
}
