package placefilter

import (
	"reflect"
	"testing"

	"poroto.app/poroto/planner/internal/domain/models"
)

func TestFilterWithinDistanceRange(t *testing.T) {
	cases := []struct {
		name            string
		placesToFilter  []models.Place
		currentLocation models.GeoLocation
		startInMeter    float64
		endInMeter      float64
		expected        []models.Place
	}{
		{
			name: "should filter places by distance range",
			placesToFilter: []models.Place{
				{
					Google: models.GooglePlace{
						Name: "Tokyo Sky Tree",
						Location: models.GeoLocation{
							Latitude:  35.710063,
							Longitude: 139.8107,
						},
					},
				},
				{
					Google: models.GooglePlace{
						Name: "Tokyo Tower",
						Location: models.GeoLocation{
							Latitude:  35.658581,
							Longitude: 139.745433,
						},
					},
				},
			},
			// とうきょうスカイツリー駅
			currentLocation: models.GeoLocation{
				Latitude:  35.7104,
				Longitude: 139.8093,
			},
			startInMeter: 0,
			endInMeter:   500,
			expected: []models.Place{
				{
					Google: models.GooglePlace{
						Name: "Tokyo Sky Tree",
						Location: models.GeoLocation{
							Latitude:  35.710063,
							Longitude: 139.8107,
						},
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := FilterWithinDistanceRange(c.placesToFilter, c.currentLocation, c.startInMeter, c.endInMeter)
			if !reflect.DeepEqual(c.expected, actual) {
				t.Errorf("expected: %v\nactual: %v", c.expected, actual)
			}
		})
	}
}
