package models

import (
	"math"
	"testing"
)

func TestGeoLocation_IsZero(t *testing.T) {
	cases := []struct {
		name     string
		location GeoLocation
		expected bool
	}{
		{
			name:     "Zero location",
			location: GeoLocation{},
			expected: true,
		},
		{
			name: "Non-zero location",
			location: GeoLocation{
				Latitude:  35.658581,
				Longitude: 139.745433,
			},
			expected: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := c.location.IsZero()
			if result != c.expected {
				t.Errorf("expected: %v\nactual: %v", c.expected, result)
			}
		})
	}
}

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

func TestGeoLocation_CalculateMBR(t *testing.T) {
	cases := []struct {
		name     string
		location GeoLocation
		distance float64
	}{
		{
			name: "Calculate MBR for Tokyo Tower",
			location: GeoLocation{
				Latitude:  35.658581,
				Longitude: 139.745433,
			},
			distance: 1000,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			c := c

			minLocation, maxLocation := c.location.CalculateMBR(c.distance)
			distanceToMinLocation := c.location.DistanceInMeter(minLocation)

			// 対角線距離を計算する
			diagonalDistance := c.distance * math.Pow(2, 0.5)
			allowableError := 100.0

			if distanceToMinLocation < diagonalDistance-allowableError || distanceToMinLocation > diagonalDistance+allowableError {
				t.Errorf("expected: %f\nactual: %f", c.distance, distanceToMinLocation)
			}

			distanceToMaxLocation := c.location.DistanceInMeter(maxLocation)
			if distanceToMaxLocation < diagonalDistance-allowableError || distanceToMaxLocation > diagonalDistance+allowableError {
				t.Errorf("expected: %f\nactual: %f", c.distance, distanceToMaxLocation)
			}
		})
	}
}

func TestEqual(t *testing.T) {
	cases := []struct {
		name      string
		locationA GeoLocation
		locationB GeoLocation
		expected  bool
	}{
		{
			name: "Same location",
			locationA: GeoLocation{
				Latitude:  35.658581,
				Longitude: 139.745433,
			},
			locationB: GeoLocation{
				Latitude:  35.658581,
				Longitude: 139.745433,
			},
			expected: true,
		},
		{
			name: "Different location",
			locationA: GeoLocation{
				Latitude:  35.658581,
				Longitude: 139.745433,
			},
			locationB: GeoLocation{
				Latitude:  35.710063,
				Longitude: 139.8107,
			},
			expected: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			resultAB := c.locationA.Equal(c.locationB)
			if resultAB != c.expected {
				t.Errorf("expected: %v\nactual: %v", c.expected, resultAB)
			}

			resultBA := c.locationB.Equal(c.locationA)
			if resultBA != c.expected {
				t.Errorf("expected: %v\nactual: %v", c.expected, resultBA)
			}
		})
	}
}
