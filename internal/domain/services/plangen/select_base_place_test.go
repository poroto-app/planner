package plangen

import (
	"github.com/google/go-cmp/cmp"
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

func TestSelectByReview(t *testing.T) {
	cases := []struct {
		name                string
		places              []api.Place
		categoriesPreferred []models.LocationCategory
		expected            []api.Place
	}{
		{
			name: "should return places sorted by review",
			places: []api.Place{
				{PlaceID: "1", Rating: 4.0},
				{PlaceID: "2", Rating: 3.0},
				{PlaceID: "3", Rating: 5.0},
				{PlaceID: "4", Rating: 2.0},
			},
			categoriesPreferred: []models.LocationCategory{},
			expected: []api.Place{
				{PlaceID: "3", Rating: 5.0},
				{PlaceID: "1", Rating: 4.0},
				{PlaceID: "2", Rating: 3.0},
			},
		},
		{
			name: "should return places most higher review in preferred categories",
			places: []api.Place{
				{PlaceID: "1", Rating: 4.0, Types: []string{models.CategoryRestaurant.SubCategories[0]}},
				{PlaceID: "2", Rating: 3.0, Types: []string{models.CategoryRestaurant.SubCategories[0]}},
				{PlaceID: "3", Rating: 5.0, Types: []string{models.CategoryAmusements.SubCategories[0]}},
				{PlaceID: "4", Rating: 2.0, Types: []string{models.CategoryAmusements.SubCategories[0]}},
			},
			categoriesPreferred: []models.LocationCategory{
				models.CategoryRestaurant,
				models.CategoryAmusements,
			},
			expected: []api.Place{
				{PlaceID: "3", Rating: 5.0, Types: []string{models.CategoryAmusements.SubCategories[0]}},
				{PlaceID: "1", Rating: 4.0, Types: []string{models.CategoryRestaurant.SubCategories[0]}},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := selectByReview(c.places, c.categoriesPreferred)
			if diff := cmp.Diff(c.expected, actual); diff != "" {
				t.Errorf("(-want +got)\n%s", diff)
			}
		})
	}
}
