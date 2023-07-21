package placefilter

import (
	"reflect"
	"testing"

	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

func TestFilterIgnoreCategory(t *testing.T) {
	cases := []struct {
		name           string
		placesToFilter []places.Place
		expected       []places.Place
	}{
		{
			name: "should remove ignore category",
			placesToFilter: []places.Place{
				{
					Name:  "Museo Nacional de Bellas Artes",
					Types: []string{"museum"},
				},
				{
					Name:  "ATM",
					Types: []string{"atm"},
				},
			},
			expected: []places.Place{
				{
					Name:  "Museo Nacional de Bellas Artes",
					Types: []string{"museum"},
				},
			},
		},
		{
			name: "ignore if place has at least one ignore category",
			placesToFilter: []places.Place{
				{
					Name:  "Museo Nacional de Bellas Artes",
					Types: []string{"museum", "church"},
				},
				{
					Name:  "Cafe",
					Types: []string{"cafe"},
				},
			},
			expected: []places.Place{
				{
					Name:  "Cafe",
					Types: []string{"cafe"},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			placesFilter := NewPlacesFilter(c.placesToFilter)
			actual := placesFilter.FilterIgnoreCategory()
			if !reflect.DeepEqual(c.expected, actual.Places()) {
				t.Errorf("expected: %v\nactual: %v", c.expected, actual.Places())
			}
		})
	}
}
