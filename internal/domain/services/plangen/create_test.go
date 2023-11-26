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
					Google: models.GooglePlace{
						Types: []string{models.CategoryAmusements.SubCategories[0]},
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
					Google: models.GooglePlace{
						Types: []string{models.CategoryAmusements.SubCategories[0]},
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

					Google: models.GooglePlace{
						Types: []string{models.CategoryAmusements.SubCategories[0]},
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

func TestSortPlacesByDistanceFrom(t *testing.T) {
	cases := []struct {
		name     string
		location models.GeoLocation
		places   []models.Place
		expected []models.Place
	}{
		{
			name: "should return places sorted by distance from location",
			location: models.GeoLocation{
				Latitude:  0,
				Longitude: 0,
			},
			places: []models.Place{
				{
					Id: "1",
					Location: models.GeoLocation{
						Latitude:  2,
						Longitude: 0,
					},
				},
				{
					Id: "2",
					Location: models.GeoLocation{
						Latitude:  3,
						Longitude: 0,
					},
				},
				{
					Id: "3",
					Location: models.GeoLocation{
						Latitude:  1,
						Longitude: 0,
					},
				},
			},
			expected: []models.Place{
				{
					Id: "3",
					Location: models.GeoLocation{
						Latitude:  1,
						Longitude: 0,
					},
				},
				{
					Id: "1",
					Location: models.GeoLocation{
						Latitude:  2,
						Longitude: 0,
					},
				},
				{
					Id: "2",
					Location: models.GeoLocation{
						Latitude:  3,
						Longitude: 0,
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := sortPlacesByDistanceFrom(c.location, c.places)
			for i := 0; i < len(result); i++ {
				if result[i].Id != c.expected[i].Id {
					t.Errorf("expected: %v\nactual: %v", result[i].Id, c.expected[i].Id)
				}
			}
		})
	}
}
