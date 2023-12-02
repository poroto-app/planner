package placefilter

import (
	"github.com/google/go-cmp/cmp"
	"poroto.app/poroto/planner/internal/domain/models"
	"testing"
)

func TestFilterDuplicated(t *testing.T) {
	cases := []struct {
		name           string
		placesToFilter []models.Place
		expected       []models.Place
	}{
		{
			name: "no duplicated",
			placesToFilter: []models.Place{
				{Id: "1"},
				{Id: "2"},
			},
			expected: []models.Place{
				{Id: "1"},
				{Id: "2"},
			},
		},
		{
			name: "duplicated",
			placesToFilter: []models.Place{
				{Id: "1"},
				{Id: "1"},
			},
			expected: []models.Place{
				{Id: "1"},
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
