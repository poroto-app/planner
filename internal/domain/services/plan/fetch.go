package plan

import (
	"context"

	"poroto.app/poroto/planner/internal/domain/models"
)

func (s Service) FetchPlan(ctx context.Context, planId string) (*models.Plan, error) {
	plan, err := s.planRepository.Find(ctx, planId)
	if err != nil {
		return nil, err
	}

	return plan, nil
}
