package placefilter

import (
	"reflect"
	"testing"

	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

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
			placesFilter := NewPlacesFilter(c.placesToFilter)
			actual := placesFilter.FilterByOpeningNow()
			if !reflect.DeepEqual(c.expected, actual.Places()) {
				t.Errorf("expected: %v\nactual: %v", c.expected, actual.Places())
			}
		})
	}
}
