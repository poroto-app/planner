package models

import "testing"

func TestHasPlace(t *testing.T) {
	cases := []struct {
		name          string
		planCandidate PlanCandidate
		googlePlaceId string
		expected      bool
	}{
		{
			name: "Has place",
			planCandidate: PlanCandidate{
				Plans: []Plan{
					{
						Places: []Place{{GooglePlaceId: toStrPointer("1")}},
					},
				},
			},
			googlePlaceId: "1",
			expected:      true,
		},
		{
			name: "Does not have place",
			planCandidate: PlanCandidate{
				Plans: []Plan{
					{
						Places: []Place{{GooglePlaceId: toStrPointer("1")}},
					},
				},
			},
			googlePlaceId: "2",
			expected:      false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := c.planCandidate.HasPlace(c.googlePlaceId)
			if result != c.expected {
				t.Errorf("expected: %t\nactual: %t", c.expected, result)
			}
		})
	}
}

func toStrPointer(v string) *string {
	return &v
}
