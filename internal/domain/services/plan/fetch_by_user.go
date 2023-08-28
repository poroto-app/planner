package plan

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/models"
)

func (s Service) PlansByUser(ctx context.Context, userId string) (*[]models.Plan, error) {
	plans, err := s.planRepository.FindByAuthorId(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("error while finding plans by user: %v", err)
	}

	return plans, nil
}
