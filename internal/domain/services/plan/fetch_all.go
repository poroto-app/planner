package plan

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/repository"
	"poroto.app/poroto/planner/internal/domain/utils"
)

const (
	defaultPlanPageSize = 10
)

type FetchPlansInput struct {
	PageToken *string
	Limit     *int
}

func (s Service) FetchPlans(ctx context.Context, input FetchPlansInput) (*[]models.Plan, *string, error) {
	var queryCursor *repository.SortedByCreatedAtQueryCursor
	if input.PageToken != nil {
		queryCursor = utils.ToPointer(repository.SortedByCreatedAtQueryCursor(*input.PageToken))
	}

	limit := defaultPlanPageSize
	if input.Limit != nil {
		limit = *input.Limit
	}

	plans, nextQueryCursor, err := s.planRepository.SortedByCreatedAt(ctx, queryCursor, limit)
	if err != nil {
		return nil, nil, fmt.Errorf("error while fetching plans: %v", err)
	}

	var nextPageToken *string
	if nextQueryCursor != nil {
		pk := string(*nextQueryCursor)
		nextPageToken = &pk
	}

	return plans, nextPageToken, nil
}
