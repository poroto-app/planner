package plangen

import (
	"poroto.app/poroto/planner/internal/domain/models"
	"testing"
)

func TestIsAlreadyHavePlaceCategoryOf(t *testing.T) {
	cases := []struct {
		name       string
		places     []models.Place
		categories []models.LocationCategory
		expected   bool
	}{
		{
			name: "should return true when places has a place of category",
			places: []models.Place{
				{
					Categories: []models.LocationCategory{
						models.CategoryAmusements,
					},
				},
			},
			categories: []models.LocationCategory{
				models.CategoryAmusements,
			},
			expected: true,
		},
		{
			name: "should return false when places does not have a place of category",
			places: []models.Place{
				{
					Categories: []models.LocationCategory{
						models.CategoryAmusements,
					},
				},
			},
			categories: []models.LocationCategory{
				models.CategoryRestaurant,
			},
			expected: false,
		},
		{
			name: "should return true when places has a place of category",
			places: []models.Place{
				{
					Categories: []models.LocationCategory{
						models.CategoryAmusements,
					},
				},
			},
			categories: []models.LocationCategory{
				models.CategoryAmusements,
				models.CategoryRestaurant,
			},
			expected: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := isAlreadyHavePlaceCategoryOf(c.places, c.categories)
			if result != c.expected {
				t.Errorf("expected: %v\nactual: %v", result, c.expected)
			}
		})
	}
}
