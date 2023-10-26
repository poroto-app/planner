package placefilter

import (
	"poroto.app/poroto/planner/internal/domain/models"
	"reflect"
	"testing"
)

func TestFilterByOpeningNow(t *testing.T) {
	cases := []struct {
		name           string
		placesToFilter []models.PlaceInPlanCandidate
		expected       []models.PlaceInPlanCandidate
	}{
		{
			name: "should filter places by opening now",
			placesToFilter: []models.PlaceInPlanCandidate{
				{
					Id:     "Place_1",
					Google: models.GooglePlace{OpenNow: true},
				},
				{
					Id:     "Place_2",
					Google: models.GooglePlace{OpenNow: false},
				},
			},
			expected: []models.PlaceInPlanCandidate{
				{
					Id:     "Place_1",
					Google: models.GooglePlace{OpenNow: true},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := FilterByOpeningNow(c.placesToFilter)
			if !reflect.DeepEqual(c.expected, actual) {
				t.Errorf("expected: %v\nactual: %v", c.expected, actual)
			}
		})
	}
}
