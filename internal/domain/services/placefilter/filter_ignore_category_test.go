package placefilter

import (
	"poroto.app/poroto/planner/internal/domain/models"
	"reflect"
	"testing"
)

func TestFilterIgnoreCategory(t *testing.T) {
	cases := []struct {
		name           string
		placesToFilter []models.GooglePlace
		expected       []models.GooglePlace
	}{
		{
			name: "should remove ignore category",
			placesToFilter: []models.GooglePlace{
				{
					Name:  "Museo Nacional de Bellas Artes",
					Types: []string{"museum"},
				},
				{
					Name:  "ATM",
					Types: []string{"atm"},
				},
			},
			expected: []models.GooglePlace{
				{
					Name:  "Museo Nacional de Bellas Artes",
					Types: []string{"museum"},
				},
			},
		},
		{
			name: "ignore if place has at least one ignore category",
			placesToFilter: []models.GooglePlace{
				{
					Name:  "Museo Nacional de Bellas Artes",
					Types: []string{"museum", "church"},
				},
				{
					Name:  "Cafe",
					Types: []string{"cafe"},
				},
			},
			expected: []models.GooglePlace{
				{
					Name:  "Cafe",
					Types: []string{"cafe"},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := FilterIgnoreCategory(c.placesToFilter)
			if !reflect.DeepEqual(c.expected, actual) {
				t.Errorf("expected: %v\nactual: %v", c.expected, actual)
			}
		})
	}
}
