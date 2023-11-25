package models

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestPlace_EstimatedStayDuration(t *testing.T) {
	cases := []struct {
		name     string
		place    Place
		expected uint
	}{
		{
			name: "place has no category",
			place: Place{
				Google: GooglePlace{
					Types: []string{},
				},
			},
			expected: 0,
		},
		{
			name: "place has one category",
			place: Place{
				Google: GooglePlace{
					Types: []string{CategoryAmusements.SubCategories[0]},
				},
			},
			expected: CategoryAmusements.EstimatedStayDuration,
		},
		{
			name: "place has two categories and return the estimated stay duration of the first one",
			place: Place{
				Google: GooglePlace{
					Types: []string{
						CategoryBookStore.SubCategories[0],
						CategoryAmusements.SubCategories[1],
					},
				},
			},
			expected: CategoryBookStore.EstimatedStayDuration,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.place.EstimatedStayDuration()
			if diff := cmp.Diff(c.expected, actual); diff != "" {
				t.Errorf("EstimatedStayDuration() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
