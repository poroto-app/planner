package placefilter

import (
	"github.com/google/go-cmp/cmp"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
	"testing"
)

func TestFilterCompany(t *testing.T) {
	cases := []struct {
		name     string
		places   []places.Place
		expected []places.Place
	}{
		{
			name: "should filter company",
			places: []places.Place{
				{
					Name: "株式会社 Example",
				},
				{
					Name: "Example(株)",
				},
				{
					Name: "Example（株）",
				},
			},
			expected: []places.Place{},
		},
		{
			name: "should not filter non-company",
			places: []places.Place{
				{
					Name: "Example Example Example",
				},
			},
			expected: []places.Place{
				{
					Name: "Example Example Example",
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := FilterCompany(c.places)
			if diff := cmp.Diff(actual, c.expected); diff != "" {
				t.Errorf("FilterCompany() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
