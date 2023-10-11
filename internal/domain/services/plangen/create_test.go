package plangen

import (
	"poroto.app/poroto/planner/internal/domain/models"
	"testing"
)

func TestIsAlreadyHavePlaceCategoryOf(t *testing.T) {
	cases := []struct {
		name       string
		places     []models.GooglePlace
		categories []models.LocationCategory
		expected   bool
	}{
		{
			name: "should return true when places has a place of category",
			places: []models.GooglePlace{
				{
					Types: []string{models.CategoryAmusements.SubCategories[0]},
				},
			},
			categories: []models.LocationCategory{
				models.CategoryAmusements,
			},
			expected: true,
		},
		{
			name: "should return false when places does not have a place of category",
			places: []models.GooglePlace{
				{
					Types: []string{models.CategoryAmusements.SubCategories[0]},
				},
			},
			categories: []models.LocationCategory{
				models.CategoryRestaurant,
			},
			expected: false,
		},
		{
			name: "should return true when places has a place of category",
			places: []models.GooglePlace{
				{
					Types: []string{models.CategoryAmusements.SubCategories[0]},
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
		places   []models.GooglePlace
		expected []models.GooglePlace
	}{
		{
			name: "should return places sorted by distance from location",
			location: models.GeoLocation{
				Latitude:  0,
				Longitude: 0,
			},
			places: []models.GooglePlace{
				{
					PlaceId: "1",
					Location: models.GeoLocation{
						Latitude:  2,
						Longitude: 0,
					},
				},
				{
					PlaceId: "2",
					Location: models.GeoLocation{
						Latitude:  3,
						Longitude: 0,
					},
				},
				{
					PlaceId: "3",
					Location: models.GeoLocation{
						Latitude:  1,
						Longitude: 0,
					},
				},
			},
			expected: []models.GooglePlace{
				{
					PlaceId: "3",
					Location: models.GeoLocation{
						Latitude:  1,
						Longitude: 0,
					},
				},
				{
					PlaceId: "1",
					Location: models.GeoLocation{
						Latitude:  2,
						Longitude: 0,
					},
				},
				{
					PlaceId: "2",
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
				if result[i].PlaceId != c.expected[i].PlaceId {
					t.Errorf("expected: %v\nactual: %v", result[i].PlaceId, c.expected[i].PlaceId)
				}
			}
		})
	}
}
