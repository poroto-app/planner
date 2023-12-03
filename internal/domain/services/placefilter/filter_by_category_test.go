package placefilter

import (
	"github.com/google/go-cmp/cmp"
	"testing"

	"poroto.app/poroto/planner/internal/domain/models"
)

func TestFuncFilterByCategory(t *testing.T) {
	cases := []struct {
		name                   string
		includeGivenCategories bool
		placesToFilter         []models.Place
		categories             []models.LocationCategory
		expected               []models.Place
	}{
		{
			name:                   "should filter places by category and include given categories",
			includeGivenCategories: true,
			placesToFilter: []models.Place{
				{
					Id:     "Place_1",
					Google: models.GooglePlace{Types: []string{"museum"}},
				},
				{
					Id:     "Place_2",
					Google: models.GooglePlace{Types: []string{"atm"}},
				},
			},
			categories: []models.LocationCategory{
				{
					Name:          "amusements",
					SubCategories: []string{"museum"},
				},
			},
			expected: []models.Place{
				{
					Id:     "Place_1",
					Google: models.GooglePlace{Types: []string{"museum"}},
				},
			},
		},
		{
			name:                   "should filter places by category and exclude given categories",
			includeGivenCategories: false,
			placesToFilter: []models.Place{
				{
					Id:     "Place_1",
					Google: models.GooglePlace{Types: []string{"museum"}},
				},
				{
					Id:     "Place_2",
					Google: models.GooglePlace{Types: []string{"atm"}},
				},
			},
			categories: []models.LocationCategory{
				{
					Name:          "atm",
					SubCategories: []string{"atm"},
				},
			},
			expected: []models.Place{
				{
					Id:     "Place_1",
					Google: models.GooglePlace{Types: []string{"museum"}},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := FilterByCategory(c.placesToFilter, c.categories, c.includeGivenCategories)
			if diff := cmp.Diff(c.expected, actual); diff != "" {
				t.Errorf("FilterByCategory() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
