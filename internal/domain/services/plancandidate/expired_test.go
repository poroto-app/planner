package plancandidate

import (
	"context"
	"go.uber.org/zap"
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
		planCandidates         map[string]models.PlanCandidateSet
		expectedPlanCandidates map[string]models.PlanCandidateSet
	}{
		{
			name:      "expired plan candidates are deleted",
			expiresAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			planCandidates: map[string]models.PlanCandidateSet{
				"planCandidateSet1": {
					Id:        "planCandidateSet1",
					ExpiresAt: time.Date(2019, 12, 31, 23, 59, 59, 0, time.UTC),
				},
				"planCandidateSet2": {
					Id:        "planCandidateSet2",
					ExpiresAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				"planCandidateSet3": {
					Id:        "planCandidateSet3",
					ExpiresAt: time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC),
				},
			},
			expectedPlanCandidates: map[string]models.PlanCandidateSet{
				"planCandidateSet3": {
					Id:        "planCandidateSet3",
					ExpiresAt: time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC),
				},
			},
		},
	}

	for _, c := range cases {
		planCandidateRepository := mock.NewPlanCandidateRepository(c.planCandidates)

		logger, _ := zap.NewDevelopment()
		service := Service{
			planCandidateRepository: planCandidateRepository,
			logger:                  logger,
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
