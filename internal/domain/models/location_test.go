package models

import (
	"math"
	"testing"
)

func TestDistanceInMeter(t *testing.T) {
	cases := []struct {
		name      string
		locationA GeoLocation
		locationB GeoLocation
		expected  float64
	}{
		{
			name: "Distance between Tokyo Tower and Tokyo Sky Tree",
			locationA: GeoLocation{
				Latitude:  35.658581,
				Longitude: 139.745433,
			},
			locationB: GeoLocation{
				Latitude:  35.710063,
				Longitude: 139.8107,
			},
			expected: 8226,
		},
		{
			name: "Distance between Nagoya City Science Museum and Nagoya City Museum",
			locationA: GeoLocation{
				Latitude:  35.165077,
				Longitude: 136.899703,
			},
			locationB: GeoLocation{
				Latitude:  35.163926,
				Longitude: 136.901071,
			},
			expected: 178,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			resultAB := c.locationA.DistanceInMeter(c.locationB)
			if math.Abs(resultAB-c.expected) > 10 {
				t.Errorf("expected: %f\nactual: %f", c.expected, resultAB)
			}

			resultBA := c.locationB.DistanceInMeter(c.locationA)
			if math.Abs(resultBA-c.expected) > 10 {
				t.Errorf("expected: %f\nactual: %f", c.expected, resultBA)
			}
		})
	}
}
