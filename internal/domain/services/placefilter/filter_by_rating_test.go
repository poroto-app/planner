package placefilter

import (
	"github.com/google/go-cmp/cmp"
	"poroto.app/poroto/planner/internal/domain/models"
	"testing"
)

func TestFilterByRating(t *testing.T) {
	cases := []struct {
		name                   string
		placesToFilter         []models.GooglePlace
		lowestRating           float32
		lowestUserRatingsTotal int
		expected               []models.GooglePlace
	}{
		{
			name: "should filter by rating",
			placesToFilter: []models.GooglePlace{
				{
					Rating:           2.9,
					UserRatingsTotal: 10,
				},
				{
					Rating:           3.0,
					UserRatingsTotal: 10,
				},
				{
					Rating:           3.1,
					UserRatingsTotal: 10,
				},
			},
			lowestRating:           3.0,
			lowestUserRatingsTotal: 10,
			expected: []models.GooglePlace{
				{
					Rating:           3.0,
					UserRatingsTotal: 10,
				},
				{
					Rating:           3.1,
					UserRatingsTotal: 10,
				},
			},
		},
		{
			name: "should filter by user ratings total",
			placesToFilter: []models.GooglePlace{
				{
					Rating:           3.0,
					UserRatingsTotal: 9,
				},
				{
					Rating:           3.0,
					UserRatingsTotal: 10,
				},
				{
					Rating:           3.0,
					UserRatingsTotal: 11,
				},
			},
			lowestRating:           3.0,
			lowestUserRatingsTotal: 10,
			expected: []models.GooglePlace{
				{
					Rating:           3.0,
					UserRatingsTotal: 10,
				},
				{
					Rating:           3.0,
					UserRatingsTotal: 11,
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := FilterByRating(c.placesToFilter, c.lowestRating, c.lowestUserRatingsTotal)
			if diff := cmp.Diff(c.expected, actual); diff != "" {
				t.Errorf("(-want +got):\n%s", diff)
			}
		})
	}
}
