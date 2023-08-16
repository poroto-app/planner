package plan

import (
	"context"

	"poroto.app/poroto/planner/internal/domain/models"
)

func (s PlanService) FindPlanCandidate(ctx context.Context, planCandidateId string) (*models.PlanCandidate, error) {
	return s.planCandidateRepository.Find(ctx, planCandidateId)
}
