package place

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/placefilter"
)

const (
	defaultMaxPlacesNearPlan            = 10
	defaultRadiusToSearchPlacesNearPlan = 2000
)

type PlacesNearPlanInput struct {
	PlanID string
	Limit  int
	Radius float64
}

func (s Service) FetchPlacesNearPlan(ctx context.Context, input PlacesNearPlanInput) (*[]models.Place, error) {
	if input.Limit < 0 {
		panic("limit must be greater than 0")
	}

	if input.Radius < 0 {
		panic("radius must be greater than 0")
	}

	if input.Limit == 0 {
		input.Limit = defaultMaxPlacesNearPlan
	}

	if input.Radius == 0 {
		input.Radius = defaultRadiusToSearchPlacesNearPlan
	}

	plan, err := s.planRepository.Find(ctx, input.PlanID)
	if err != nil {
		return nil, err
	}

	if len(plan.Places) == 0 {
		return nil, nil
	}

	planLocation := plan.Places[0].Location
	places, err := s.placeRepository.FindByLocation(ctx, planLocation, input.Radius)
	if err != nil {
		return nil, fmt.Errorf("error while fetching places near plan: %v", err)
	}

	placesFiltered := placefilter.FilterDefaultIgnore(placefilter.FilterDefaultIgnoreInput{
		Places:              places,
		StartLocation:       planLocation,
		IgnoreDistanceRange: input.Radius,
	})

	// プランに含まれている場所を除外する
	placesFiltered = array.Filter(placesFiltered, func(place models.Place) bool {
		_, isPlaceInPlan := array.Find(plan.Places, func(placeInPlan models.Place) bool {
			return placeInPlan.Id == place.Id
		})
		return !isPlaceInPlan
	})

	// レビューの高い順に並び替える
	placesNearbyPlan := models.SortPlacesByRating(placesFiltered)

	if len(placesNearbyPlan) > input.Limit {
		placesNearbyPlan = placesNearbyPlan[:input.Limit]
	}

	// 写真を取得
	placesNearbyPlan = s.placeSearchService.FetchPlacesPhotosAndSave(ctx, placesNearbyPlan...)

	return &placesNearbyPlan, nil
}
