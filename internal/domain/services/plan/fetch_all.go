package plan

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/repository"

	"poroto.app/poroto/planner/internal/domain/models"
)

func (s Service) FetchPlans(ctx context.Context, pageToken *string) (*[]models.Plan, error) {
	var queryCursor *repository.SortedByCreatedAtQueryCursor
	if pageToken != nil {
		qc := repository.SortedByCreatedAtQueryCursor(*pageToken)
		queryCursor = &qc
	}

	// TODO: PageTokenをreturnする
	plans, _, err := s.planRepository.SortedByCreatedAt(ctx, queryCursor, 10)
	if err != nil {
		return nil, fmt.Errorf("error while fetching plans: %v", err)
	}

	return plans, nil
}
