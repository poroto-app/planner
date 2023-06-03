package plan

import (
	"reflect"
	"testing"

	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

func TestFuncFilterByCategory(t *testing.T) {
	cases := []struct {
		name           string
		placesToFilter []places.Place
		categories     []models.LocationCategory
		expected       []places.Place
	}{
		{
			name: "should filter places by category",
			placesToFilter: []places.Place{
				{
					Name: "Museo Nacional de Bellas Artes",
					Types: []string{
						"museum",
					},
				},
				{
					Name:  "ATM",
					Types: []string{"atm"},
				},
			},
			categories: []models.LocationCategory{
				{
					Name:          "amusements",
					SubCategories: []string{"museum"},
				},
			},
			expected: []places.Place{
				{
					Name: "Museo Nacional de Bellas Artes",
					Types: []string{
						"museum",
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := PlanService{}
			actual := s.filterByCategory(c.placesToFilter, c.categories)
			if !reflect.DeepEqual(actual, c.expected) {
				t.Errorf("expected %v, got %v", c.expected, actual)
			}
		})
	}
}

func TestFilterByOpeningNow(t *testing.T) {
	cases := []struct {
		name           string
		placesToFilter []places.Place
		expected       []places.Place
	}{
		{
			name: "should filter places by opening now",
			placesToFilter: []places.Place{
				{
					Name:    "Museo Nacional de Bellas Artes",
					OpenNow: true,
				},
				{
					Name:    "Subway",
					OpenNow: false,
				},
			},
			expected: []places.Place{
				{
					Name:    "Museo Nacional de Bellas Artes",
					OpenNow: true,
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := PlanService{}
			actual := s.filterByOpeningNow(c.placesToFilter)
			if !reflect.DeepEqual(actual, c.expected) {
				t.Errorf("expected %v, got %v", c.expected, actual)
			}
		})
	}
}

func TestFilterWithinDistanceRange(t *testing.T) {
	cases := []struct {
		name            string
		placesToFilter  []places.Place
		currentLocation models.GeoLocation
		startInMeter    float64
		endInMeter      float64
		expected        []places.Place
	}{
		{
			name: "should filter places by distance range",
			placesToFilter: []places.Place{
				{
					Name: "Tokyo Sky Tree",
					Location: places.Location{
						Latitude:  35.710063,
						Longitude: 139.8107,
					},
				},
				{
					Name: "Tokyo Tower",
					Location: places.Location{
						Latitude:  35.658581,
						Longitude: 139.745433,
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
			expected: []places.Place{
				{
					Name: "Tokyo Sky Tree",
					Location: places.Location{
						Latitude:  35.710063,
						Longitude: 139.8107,
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := PlanService{}
			actual := s.filterWithinDistanceRange(c.placesToFilter, c.currentLocation, c.startInMeter, c.endInMeter)
			if !reflect.DeepEqual(actual, c.expected) {
				t.Errorf("expected %v, got %v", c.expected, actual)
			}
		})
	}
}
