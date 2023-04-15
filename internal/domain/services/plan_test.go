package services

import (
	"reflect"
	"testing"

	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

func TestFuncFilterByCategory(t *testing.T) {
	cases := []struct {
		name           string
		placesToFilter []places.Place
		categories     []models.LocationCategory
		expected       []places.Place
	}{
		{
			name: "should filter places by category",
			placesToFilter: []places.Place{
				{
					Name: "Museo Nacional de Bellas Artes",
					Types: []string{
						"museum",
					},
				},
				{
					Name:  "ATM",
					Types: []string{"atm"},
				},
			},
			categories: []models.LocationCategory{
				{
					Name:          "amusements",
					SubCategories: []string{"museum"},
				},
			},
			expected: []places.Place{
				{
					Name: "Museo Nacional de Bellas Artes",
					Types: []string{
						"museum",
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := PlanService{}
			actual := s.filterByCategory(c.placesToFilter, c.categories)
			if !reflect.DeepEqual(actual, c.expected) {
				t.Errorf("expected %v, got %v", c.expected, actual)
			}
		})
	}
}
