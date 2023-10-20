package placefilter

import (
	"poroto.app/poroto/planner/internal/domain/models"
	"reflect"
	"testing"
)

func TestFilterByOpeningNow(t *testing.T) {
	cases := []struct {
		name           string
		placesToFilter []models.GooglePlace
		expected       []models.GooglePlace
	}{
		{
			name: "should filter places by opening now",
			placesToFilter: []models.GooglePlace{
				{
					Name:    "Museo Nacional de Bellas Artes",
					OpenNow: true,
				},
				{
					Name:    "Subway",
					OpenNow: false,
				},
			},
			expected: []models.GooglePlace{
				{
					Name:    "Museo Nacional de Bellas Artes",
					OpenNow: true,
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
