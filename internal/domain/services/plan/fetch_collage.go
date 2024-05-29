package plan

import (
	"context"
	"poroto.app/poroto/planner/internal/domain/models"
)

type FetchPlanCollageInput struct {
	PlanId string
}

func (s Service) FetchPlanCollage(ctx context.Context, input FetchPlanCollageInput) (*models.PlanCollage, error) {
	return s.planRepository.FindCollage(ctx, input.PlanId)
}
