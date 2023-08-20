package plan

import (
	"context"
	"fmt"

	"poroto.app/poroto/planner/internal/domain/models"
)

func (s Service) FetchPlans(ctx context.Context, nextPageToken *string) (*[]models.Plan, error) {
	plans, err := s.planRepository.SortedByCreatedAt(ctx, nextPageToken, 10)
	if err != nil {
		return nil, fmt.Errorf("error while fetching plans: %v", err)
	}
	return plans, nil
}
