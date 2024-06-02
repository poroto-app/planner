package plan

import (
	"context"

	"poroto.app/poroto/planner/internal/domain/models"
)

const (
	defaultLimit = 10
)

func (s Service) FetchPlansByLocation(
	ctx context.Context,
	location models.GeoLocation,
	limit *int,
	pageToken *string,
) (plans *[]models.Plan, nextPageToken *string, err error) {
	if limit == nil {
		value := defaultLimit
		limit = &value
	}

	plans, nextPageToken, err = s.planRepository.FindByLocation(ctx, location, pageToken, *limit)
	if err != nil {
		return nil, nil, err
	}

	return plans, nextPageToken, nil
}
