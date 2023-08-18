package plangen

import (
	"poroto.app/poroto/planner/internal/domain/models"
	api "poroto.app/poroto/planner/internal/infrastructure/api/google/places"
	"testing"
)

func TestIsAlreadyAdded(t *testing.T) {
	cases := []struct {
		name     string
		place    api.Place
		places   []api.Place
		expected bool
	}{
		{
			name:  "should return true when place is already added",
			place: api.Place{PlaceID: "1"},
			places: []api.Place{
				{PlaceID: "1"},
				{PlaceID: "2"},
			},
			expected: true,
		},
		{
			name:  "should return false when place is not added",
			place: api.Place{PlaceID: "3"},
			places: []api.Place{
				{PlaceID: "1"},
				{PlaceID: "2"},
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
		a        api.Place
		b        api.Place
		expected bool
	}{
		{
			name: "should return true when two places are same category",
			a: api.Place{
				Types: []string{models.CategoryRestaurant.SubCategories[0]},
			},
			b: api.Place{
				Types: []string{models.CategoryRestaurant.SubCategories[1]},
			},
			expected: true,
		},
		{
			name: "should return false when two places are not same category",
			a: api.Place{
				Types: []string{models.CategoryRestaurant.SubCategories[0]},
			},
			b: api.Place{
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
