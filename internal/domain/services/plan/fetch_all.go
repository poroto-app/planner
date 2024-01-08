package plan

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/repository"

	"poroto.app/poroto/planner/internal/domain/models"
)

func (s Service) FetchPlans(ctx context.Context, pageToken *string) (*[]models.Plan, *string, error) {
	var queryCursor *repository.SortedByCreatedAtQueryCursor
	if pageToken != nil {
		qc := repository.SortedByCreatedAtQueryCursor(*pageToken)
		queryCursor = &qc
	}

	plans, nextQueryCursor, err := s.planRepository.SortedByCreatedAt(ctx, queryCursor, 10)
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
