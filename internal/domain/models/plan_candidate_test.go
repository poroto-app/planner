package models

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

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

func TestPlanCandidate_GetPlan(t *testing.T) {
	cases := []struct {
		name          string
		planCandidate PlanCandidate
		planId        string
		expected      *Plan
	}{
		{
			name:          "Has plan",
			planCandidate: PlanCandidate{Plans: []Plan{{Id: "1"}}},
			planId:        "1",
			expected:      &Plan{Id: "1"},
		},
		{
			name:          "Does not have plan",
			planCandidate: PlanCandidate{Plans: []Plan{{Id: "1"}}},
			planId:        "2",
			expected:      nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := c.planCandidate.GetPlan(c.planId)
			if diff := cmp.Diff(result, c.expected); diff != "" {
				t.Errorf("GetPlan() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func toStrPointer(v string) *string {
	return &v
}
