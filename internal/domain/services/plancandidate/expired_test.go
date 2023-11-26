package plancandidate

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/mock"
)

func TestDeleteExpiredPlanCandidates(t *testing.T) {
	cases := []struct {
		name                   string
		expiresAt              time.Time
		planCandidates         map[string]models.PlanCandidate
		expectedPlanCandidates map[string]models.PlanCandidate
	}{
		{
			name:      "expired plan candidates are deleted",
			expiresAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			planCandidates: map[string]models.PlanCandidate{
				"planCandidate1": {
					Id:        "planCandidate1",
					ExpiresAt: time.Date(2019, 12, 31, 23, 59, 59, 0, time.UTC),
				},
				"planCandidate2": {
					Id:        "planCandidate2",
					ExpiresAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				"planCandidate3": {
					Id:        "planCandidate3",
					ExpiresAt: time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC),
				},
			},
			expectedPlanCandidates: map[string]models.PlanCandidate{
				"planCandidate3": {
					Id:        "planCandidate3",
					ExpiresAt: time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC),
				},
			},
		},
	}

	for _, c := range cases {
		planCandidateRepository := mock.NewPlanCandidateRepository(c.planCandidates)

		service := Service{
			planCandidateRepository: planCandidateRepository,
		}

		err := service.DeleteExpiredPlanCandidates(context.Background(), c.expiresAt)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if diff := cmp.Diff(c.expectedPlanCandidates, c.planCandidates); diff != "" {
			t.Errorf("unexpected plan candidates (-want +got):\n%s", diff)
		}
	}
}
