package plan

import (
	"context"
	"poroto.app/poroto/planner/internal/domain/utils"

	"poroto.app/poroto/planner/internal/domain/models"
)

const (
	defaultLimit = 10

	// 半径2km圏内のプランを検索する
	defaultDistanceToSearchPlan = 2 * 1000
)

type FetchPlansByLocationInput struct {
	Location    models.GeoLocation
	Limit       *int
	SearchRange *int
}

func (s Service) FetchPlansByLocation(ctx context.Context, input FetchPlansByLocationInput) (plans *[]models.Plan, nextPageToken *string, err error) {
	if input.Limit == nil {
		input.Limit = utils.ToPointer(defaultLimit)
	}

	if input.SearchRange == nil {
		input.SearchRange = utils.ToPointer(defaultDistanceToSearchPlan)
	}

	plans, nextPageToken, err = s.planRepository.FindByLocation(ctx, input.Location, *input.Limit, *input.SearchRange)
	if err != nil {
		return nil, nil, err
	}

	return plans, nextPageToken, nil
}
