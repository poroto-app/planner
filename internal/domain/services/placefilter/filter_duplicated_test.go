package placefilter

import (
	"github.com/google/go-cmp/cmp"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
	"testing"
)

func TestFilterDuplicated(t *testing.T) {
	cases := []struct {
		name           string
		placesToFilter []places.Place
		expected       []places.Place
	}{
		{
			name: "no duplicated",
			placesToFilter: []places.Place{
				{PlaceID: "1"},
				{PlaceID: "2"},
			},
			expected: []places.Place{
				{PlaceID: "1"},
				{PlaceID: "2"},
			},
		},
		{
			name: "duplicated",
			placesToFilter: []places.Place{
				{PlaceID: "1"},
				{PlaceID: "1"},
			},
			expected: []places.Place{
				{PlaceID: "1"},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := FilterDuplicated(c.placesToFilter)
			if diff := cmp.Diff(c.expected, result); diff != "" {
				t.Errorf("FilterDuplicated() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
