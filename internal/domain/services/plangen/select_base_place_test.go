package plangen

import (
	"poroto.app/poroto/planner/internal/domain/models"
	"testing"
)

func TestIsAlreadyAdded(t *testing.T) {
	cases := []struct {
		name     string
		place    models.GooglePlace
		places   []models.GooglePlace
		expected bool
	}{
		{
			name:  "should return true when place is already added",
			place: models.GooglePlace{PlaceId: "1"},
			places: []models.GooglePlace{
				{PlaceId: "1"},
				{PlaceId: "2"},
			},
			expected: true,
		},
		{
			name:  "should return false when place is not added",
			place: models.GooglePlace{PlaceId: "3"},
			places: []models.GooglePlace{
				{PlaceId: "1"},
				{PlaceId: "2"},
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

func TestIsSameCategoryPlace(t *testing.T) {
	cases := []struct {
		name     string
		a        models.GooglePlace
		b        models.GooglePlace
		expected bool
	}{
		{
			name: "should return true when two places are same category",
			a: models.GooglePlace{
				Types: []string{models.CategoryRestaurant.SubCategories[0]},
			},
			b: models.GooglePlace{
				Types: []string{models.CategoryRestaurant.SubCategories[1]},
			},
			expected: true,
		},
		{
			name: "should return false when two places are not same category",
			a: models.GooglePlace{
				Types: []string{models.CategoryRestaurant.SubCategories[0]},
			},
			b: models.GooglePlace{
				Types: []string{models.CategoryAmusements.SubCategories[0]},
			},
			expected: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := isSameCategoryPlace(c.a, c.b)
			if actual != c.expected {
				t.Errorf("expected: %v, actual: %v", c.expected, actual)
			}
		})
	}
}
