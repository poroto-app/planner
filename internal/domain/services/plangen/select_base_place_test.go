package plangen

import (
	"poroto.app/poroto/planner/internal/domain/models"
	"testing"
)

func TestIsAlreadyAdded(t *testing.T) {
	cases := []struct {
		name     string
		place    models.PlaceInPlanCandidate
		places   []models.PlaceInPlanCandidate
		expected bool
	}{
		{
			name:  "should return true when place is already added",
			place: models.PlaceInPlanCandidate{Id: "1"},
			places: []models.PlaceInPlanCandidate{
				{Id: "1"},
				{Id: "2"},
			},
			expected: true,
		},
		{
			name:  "should return false when place is not added",
			place: models.PlaceInPlanCandidate{Id: "3"},
			places: []models.PlaceInPlanCandidate{
				{Id: "1"},
				{Id: "2"},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := isAlreadyAdded(c.place, c.places)
			if actual != c.expected {
				t.Errorf("expected: %v, actual: %v", c.expected, actual)
			}
		})
	}
}
