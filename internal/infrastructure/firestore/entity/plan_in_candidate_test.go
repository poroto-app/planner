package entity

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"poroto.app/poroto/planner/internal/domain/models"
)

func TestFromPlanInCandidateEntity(t *testing.T) {
	cases := []struct {
		name     string
		entity   PlanInCandidateEntity
		places   []models.PlaceInPlanCandidate
		expected *models.Plan
	}{
		{
			name: "正常系",
			entity: PlanInCandidateEntity{
				PlaceIdsOrdered: []string{"02", "01"},
			},
			places: []models.PlaceInPlanCandidate{
				{Id: "01"},
				{Id: "02"},
			},
			expected: &models.Plan{
				Places: []models.Place{
					{Id: "02"},
					{Id: "01"},
				},
			},
		},
		{
			name: "placeIdsOrdered に重複がある場合は nil を返す",
			entity: PlanInCandidateEntity{
				PlaceIdsOrdered: []string{"01", "01"},
			},
			places: []models.PlaceInPlanCandidate{
				{Id: "01"},
				{Id: "02"},
			},
			expected: nil,
		},
		{
			name: "placeIdsOrdered と places の示す場所が一致しない場合はnilを返す",
			entity: PlanInCandidateEntity{
				PlaceIdsOrdered: []string{"10", "20"},
			},
			places: []models.PlaceInPlanCandidate{
				{Id: "01"},
				{Id: "02"},
			},
			expected: nil,
		},
		{
			name: "placeIdsOrdered と places の大きさが一致しない場合は場合はnil",
			entity: PlanInCandidateEntity{
				PlaceIdsOrdered: []string{"01", "02", "03"},
			},
			places: []models.PlaceInPlanCandidate{
				{Id: "01"},
				{Id: "02"},
			},
			expected: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result, _ := FromPlanInCandidateEntity(
				c.entity.Id,
				c.entity.Name,
				c.places,
				c.entity.PlaceIdsOrdered,
				c.entity.TimeInMinutes,
			)
			if diff := cmp.Diff(c.expected, result); diff != "" {
				t.Errorf("FromPlanInCandidateEntity() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestValidatePlanInCandidateEntity(t *testing.T) {
	cases := []struct {
		name            string
		placeIdsOrdered []string
		places          []models.Place
		valid           bool
	}{
		{
			name:            "正常系",
			placeIdsOrdered: []string{"02", "01"},
			places: []models.Place{
				{Id: "01"},
				{Id: "02"},
			},
			valid: true,
		},
		{
			name:            "placeIdsOrdered に重複がある場合は false",
			placeIdsOrdered: []string{"01", "01"},
			places: []models.Place{
				{Id: "01"},
				{Id: "02"},
			},
			valid: false,
		},
		{
			name:            "placeIdsOrdered と places の示す場所が一致しない場合は false",
			placeIdsOrdered: []string{"10", "20"},
			places: []models.Place{
				{Id: "01"},
				{Id: "02"},
			},
			valid: false,
		},
		{
			name:            "placeIdsOrdered と places の示す場所が一致しない場合は false",
			placeIdsOrdered: []string{"10", "20"},
			valid:           false,
		},
		{
			name:            "placeIdsOrdered と places の大きさが一致しない場合は false",
			placeIdsOrdered: []string{"01", "02", "03"},
			valid:           false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := validatePlaceInPlanCandidateEntity(
				c.places,
				c.placeIdsOrdered,
			)

			result := err == nil
			if diff := cmp.Diff(c.valid, result); diff != "" {
				t.Errorf("valid %v, but got %v", c.valid, result)
			}
		})
	}
}
