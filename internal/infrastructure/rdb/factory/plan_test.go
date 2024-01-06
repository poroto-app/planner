package factory

import (
	"github.com/google/go-cmp/cmp"
	"github.com/volatiletech/null/v8"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
	"testing"
)

func TestNewPlanEntityFromDomainModel(t *testing.T) {
	tests := []struct {
		name     string
		plan     models.Plan
		expected generated.Plan
	}{
		{
			name: "should return a valid entity",
			plan: models.Plan{
				Id:       "ec7c607d-454a-4644-929a-c3b1e078842d",
				AuthorId: utils.ToPointer("339809cf-d515-4a64-bbcd-c6a899051273"),
				Name:     "plan title",
			},
			expected: generated.Plan{
				ID:     "ec7c607d-454a-4644-929a-c3b1e078842d",
				UserID: null.StringFrom("339809cf-d515-4a64-bbcd-c6a899051273"),
				Name:   "plan title",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := NewPlanEntityFromDomainModel(tt.plan)
			if diff := cmp.Diff(tt.expected, actual); diff != "" {
				t.Errorf("wrong plan entity (-want +got):\n%s", diff)
			}
		})
	}
}

func TestNewPlanEntityFromDomainModel_EmptyID(t *testing.T) {
	cases := []struct {
		name string
		plan models.Plan
	}{
		{
			name: "should generate valid id if id is empty",
			plan: models.Plan{
				Id:       "",
				AuthorId: utils.ToPointer("339809cf-d515-4a64-bbcd-c6a899051273"),
				Name:     "plan title",
			},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			actual := NewPlanEntityFromDomainModel(c.plan)
			if actual.ID == "" {
				t.Errorf("expected: %v, actual: %v", c.plan.Id, actual.ID)
			}
		})
	}
}
