package placefilter

import (
	"reflect"
	"testing"

	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

func TestFuncFilterByCategory(t *testing.T) {
	cases := []struct {
		name                   string
		includeGivenCategories bool
		placesToFilter         []places.Place
		categories             []models.LocationCategory
		expected               []places.Place
	}{
		{
			name:                   "should filter places by category and include given categories",
			includeGivenCategories: true,
			placesToFilter: []places.Place{
				{
					PlaceID: "Place_1",
					Types:   []string{"museum"},
				},
				{
					PlaceID: "Place_2",
					Types:   []string{"atm"},
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
					PlaceID: "Place_1",
					Types:   []string{"museum"},
				},
			},
		},
		{
			name:                   "should filter places by category and exclude given categories",
			includeGivenCategories: false,
			placesToFilter: []places.Place{
				{
					PlaceID: "Place_1",
					Types: []string{
						"museum",
					},
				},
				{
					PlaceID: "Place_2",
					Types:   []string{"atm"},
				},
			},
			categories: []models.LocationCategory{
				{
					Name:          "atm",
					SubCategories: []string{"atm"},
				},
			},
			expected: []places.Place{
				{
					PlaceID: "Place_1",
					Types:   []string{"museum"},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := FilterByCategory(c.placesToFilter, c.categories, c.includeGivenCategories)
			if !reflect.DeepEqual(c.expected, actual) {
				t.Errorf("expected: %v\nactual: %v", c.expected, actual)
			}
		})
	}
}
