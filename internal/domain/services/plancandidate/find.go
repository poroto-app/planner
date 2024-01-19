package plancandidate

import (
	"context"
	"time"

	"poroto.app/poroto/planner/internal/domain/models"
)

func (s Service) FindPlanCandidate(ctx context.Context, planCandidateId string) (*models.PlanCandidate, error) {
	return s.planCandidateRepository.Find(ctx, planCandidateId, time.Now())
}
