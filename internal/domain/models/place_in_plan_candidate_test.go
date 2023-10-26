package models

import "testing"

func TestPlaceInPlanCandidate_IsSameCategoryPlace(t *testing.T) {
	cases := []struct {
		name     string
		a        PlaceInPlanCandidate
		b        PlaceInPlanCandidate
		expected bool
	}{
		{
			name: "should return true when two places are same category",
			a: PlaceInPlanCandidate{
				Google: GooglePlace{
					Types: []string{CategoryRestaurant.SubCategories[0]},
				},
			},
			b: PlaceInPlanCandidate{
				Google: GooglePlace{
					Types: []string{CategoryRestaurant.SubCategories[1]},
				},
			},
			expected: true,
		},
		{
			name: "should return false when two places are not same category",
			a: PlaceInPlanCandidate{
				Google: GooglePlace{
					Types: []string{CategoryRestaurant.SubCategories[0]},
				},
			},
			b: PlaceInPlanCandidate{
				Google: GooglePlace{
					Types: []string{CategoryAmusements.SubCategories[0]},
				},
			},
			expected: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.a.IsSameCategoryPlace(c.b)
			if actual != c.expected {
				t.Errorf("expected: %v, actual: %v", c.expected, actual)
			}
		})
	}
}
