package placefilter

import (
	"reflect"
	"testing"

	"poroto.app/poroto/planner/internal/domain/models"
)

func TestFuncFilterByCategory(t *testing.T) {
	cases := []struct {
		name                   string
		includeGivenCategories bool
		placesToFilter         []models.GooglePlace
		categories             []models.LocationCategory
		expected               []models.GooglePlace
	}{
		{
			name:                   "should filter places by category and include given categories",
			includeGivenCategories: true,
			placesToFilter: []models.GooglePlace{
				{
					PlaceId: "Place_1",
					Types:   []string{"museum"},
				},
				{
					PlaceId: "Place_2",
					Types:   []string{"atm"},
				},
			},
			categories: []models.LocationCategory{
				{
					Name:          "amusements",
					SubCategories: []string{"museum"},
				},
			},
			expected: []models.GooglePlace{
				{
					PlaceId: "Place_1",
					Types:   []string{"museum"},
				},
			},
		},
		{
			name:                   "should filter places by category and exclude given categories",
			includeGivenCategories: false,
			placesToFilter: []models.GooglePlace{
				{
					PlaceId: "Place_1",
					Types: []string{
						"museum",
					},
				},
				{
					PlaceId: "Place_2",
					Types:   []string{"atm"},
				},
			},
			categories: []models.LocationCategory{
				{
					Name:          "atm",
					SubCategories: []string{"atm"},
				},
			},
			expected: []models.GooglePlace{
				{
					PlaceId: "Place_1",
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
