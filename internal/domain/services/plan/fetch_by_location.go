package plan

import (
	"context"

	"poroto.app/poroto/planner/internal/domain/models"
)

func (s PlanService) FetchPlansByLocation(
	ctx context.Context,
	location models.GeoLocation,
	pageToken *string,
) (plans *[]models.Plan, nextPageToken *string, err error) {
	plans, nextPageToken, err = s.planRepository.SortedByLocation(ctx, location, pageToken, 10)
	if err != nil {
		return nil, nil, err
	}

	return plans, nextPageToken, nil
}
