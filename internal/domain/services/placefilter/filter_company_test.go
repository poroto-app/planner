package placefilter

import (
	"github.com/google/go-cmp/cmp"
	"poroto.app/poroto/planner/internal/domain/models"
	"testing"
)

func TestFilterCompany(t *testing.T) {
	cases := []struct {
		name     string
		places   []models.Place
		expected []models.Place
	}{
		{
			name: "should filter company",
			places: []models.Place{
				{
					Google: models.GooglePlace{Name: "株式会社 Example"},
				},
				{
					Google: models.GooglePlace{Name: "Example(株)"},
				},
				{
					Google: models.GooglePlace{Name: "Example（株）"},
				},
			},
			expected: []models.Place{},
		},
		{
			name: "should not filter non-company",
			places: []models.Place{
				{
					Google: models.GooglePlace{Name: "Example Example Example"},
				},
			},
			expected: []models.Place{
				{
					Google: models.GooglePlace{Name: "Example Example Example"},
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
