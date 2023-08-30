package models

import (
	"testing"
)

func GetPlaceTest(t *testing.T) {
	cases := []struct {
		name     string
		plan     Plan
		placeId  string
		expected *Place
	}{
		{
			name: "Valid place ID",
			plan: Plan{
				Places: []Place{
					{
						Id:   "1",
						Name: "place1",
					},
				},
			},
			placeId: "1",
			expected: &Place{
				Id:   "1",
				Name: "place1",
			},
		},
		{
			name: "Invalid place ID",
			plan: Plan{
				Places: []Place{
					{
						Id:   "1",
						Name: "place1",
					},
				},
			},
			placeId:  "2",
			expected: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.plan.GetPlace(c.placeId)
			if actual != c.expected {
				t.Errorf("expected: %v, actual: %v", c.expected, actual)
			}
		})
	}
}
