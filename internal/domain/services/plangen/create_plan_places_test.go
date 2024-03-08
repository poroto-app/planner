package plangen

import (
	"poroto.app/poroto/planner/internal/domain/models"
	"testing"
)

func TestSortPlacesByDistanceFrom(t *testing.T) {
	cases := []struct {
		name     string
		location models.GeoLocation
		places   []models.Place
		expected []models.Place
	}{
		{
			name: "should return Places sorted by distance from location",
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
