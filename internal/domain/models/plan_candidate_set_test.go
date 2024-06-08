package models

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestPlanCandidateSet_HasPlace(t *testing.T) {
	cases := []struct {
		name             string
		planCandidateSet PlanCandidateSet
		placeId          string
		expected         bool
	}{
		{
			name: "Has place of placeId",
			planCandidateSet: PlanCandidateSet{
				Plans: []Plan{
					{
						Places: []Place{{Id: "1"}},
					},
				},
			},
			placeId:  "1",
			expected: true,
		},
		{
			name: "Does not have place of placeId",
			planCandidateSet: PlanCandidateSet{
				Plans: []Plan{
					{
						Places: []Place{{Id: "1"}},
					},
				},
			},
			placeId:  "2",
			expected: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := c.planCandidateSet.HasPlace(c.placeId)
			if result != c.expected {
				t.Errorf("expected: %t\nactual: %t", c.expected, result)
			}
		})
	}
}

func TestPlanCandidateSet_GetPlan(t *testing.T) {
	cases := []struct {
		name             string
		planCandidateSet PlanCandidateSet
		planId           string
		expected         *Plan
	}{
		{
			name:             "Has plan",
			planCandidateSet: PlanCandidateSet{Plans: []Plan{{Id: "1"}}},
			planId:           "1",
			expected:         &Plan{Id: "1"},
		},
		{
			name:             "Does not have plan",
			planCandidateSet: PlanCandidateSet{Plans: []Plan{{Id: "1"}}},
			planId:           "2",
			expected:         nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := c.planCandidateSet.GetPlan(c.planId)
			if diff := cmp.Diff(result, c.expected); diff != "" {
				t.Errorf("GetPlan() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
