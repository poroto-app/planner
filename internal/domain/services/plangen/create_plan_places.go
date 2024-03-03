package plangen

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/placefilter"
)

const (
	defaultMaxPlanDuration = 180
	defaultMaxPlaceInPlan  = 4

	placeDistanceRangeInPlan = 500 // 徒歩5分以内
)

type CreatePlanPlacesParams struct {
	PlanCandidateId         string
	LocationStart           models.GeoLocation
	PlaceStart              models.Place
	Places                  []models.Place
	PlacesOtherPlansContain []models.Place
	CategoryNamesDisliked   *[]string
	FreeTime                *int
	MaxPlace                int
}

// createPlanPlaces プランの候補地となる場所を作成する
func (s Service) createPlanPlaces(ctx context.Context, params CreatePlanPlacesParams) ([]models.Place, error) {
	if params.PlanCandidateId == "" {
		panic("PlanCandidateId is required")
	}

	if params.MaxPlace == 0 {
		params.MaxPlace = defaultMaxPlaceInPlan
	}

	placesFiltered := params.Places
	placesFiltered = placefilter.FilterDefaultIgnore(placefilter.FilterDefaultIgnoreInput{
		Places:              placesFiltered,
		StartLocation:       params.LocationStart,
		IgnoreDistanceRange: placeDistanceRangeInPlan,
	})
	s.logger.Debug("places after filtering by distance", zap.Int("places", len(placesFiltered)))

	// ユーザーが拒否した場所は取り除く
	if params.CategoryNamesDisliked != nil {
		categoriesDisliked := models.GetCategoriesFromSubCategories(*params.CategoryNamesDisliked)
		placesFiltered = placefilter.FilterByCategory(placesFiltered, categoriesDisliked, false)
		s.logger.Debug("places after filtering by disliked categories", zap.Int("places", len(placesFiltered)))
	}

	// 他のプランに含まれている場所を除外する
	placesFiltered = placefilter.FilterPlaces(placesFiltered, func(place models.Place) bool {
		if params.PlacesOtherPlansContain == nil {
			return true
		}

		for _, placeOtherPlanContain := range params.PlacesOtherPlansContain {
			if place.Id == placeOtherPlanContain.Id {
				return false
			}
		}
		return true
	})
	s.logger.Debug("places after filtering by other plans", zap.Int("places", len(placesFiltered)))

	// レビューの高い順でソート
	placesSorted := models.SortPlacesByRating(placesFiltered)

	placesInPlan := make([]models.Place, 0)

	// 指定された場所を基準としてプランを作成するときは必ず含める
	if params.LocationStart.Equal(params.PlaceStart.Location) {
		placesInPlan = append(placesInPlan, params.PlaceStart)
	}

	for _, place := range placesSorted {
		// プランに含まれる場所の数が上限に達したら終了
		if len(placesInPlan) >= params.MaxPlace {
			s.logger.Debug(
				"skip place because the number of Places in plan is over",
				zap.String("place", place.Google.Name),
				zap.Int("MaxPlace", params.MaxPlace),
				zap.Int("placesInPlan", len(placesInPlan)),
			)
			break
		}

		if place.Id == params.PlaceStart.Id {
			continue
		}

		// 飲食店を複数回含めない
		if isAlreadyHavePlaceCategoryOf(placesInPlan, models.FoodCategories()) && isCategoryOf(place.Google.Types, models.FoodCategories()) {
			s.logger.Debug(
				"skip place because the cafe or restaurant is already in plan",
				zap.String("place", place.Google.Name),
			)
			continue
		}

		// 最適経路で巡ったときの所要時間を計算
		sortedByDistance := sortPlacesByDistanceFrom(params.LocationStart, append(placesInPlan, place))
		timeInPlan := planTimeFromPlaces(params.LocationStart, sortedByDistance)

		// 予定の時間内に収まらない場合はスキップ
		if params.FreeTime != nil && timeInPlan > uint(*params.FreeTime) {
			s.logger.Debug(
				"skip place because it will be over time",
				zap.String("place", place.Google.Name),
				zap.Uint("timeInPlan", timeInPlan),
				zap.Int("FreeTime", *params.FreeTime),
			)
			continue
		}

		// 予定の時間を指定しない場合、3時間を超える場合はスキップ
		if params.FreeTime == nil && timeInPlan > defaultMaxPlanDuration {
			s.logger.Debug(
				"skip place because it will be over time",
				zap.String("place", place.Google.Name),
				zap.Uint("timeInPlan", timeInPlan),
				zap.Int("defaultMaxPlanDuration", defaultMaxPlanDuration),
			)
			continue
		}

		placesInPlan = append(placesInPlan, place)
	}

	if len(placesInPlan) == 0 {
		return nil, fmt.Errorf("could not contain any Places in plan")
	}

	return placesInPlan, nil
}

		return false
	}

	return true
}

// sortPlacesByDistanceFrom location からplacesを巡回する最短経路をgreedy法で求める
func sortPlacesByDistanceFrom(location models.GeoLocation, places []models.Place) []models.Place {
	placesSorted := make([]models.Place, len(places))
	copy(placesSorted, places)

	prevLocation := location
	for i := 0; i < len(places); i++ {
		nearestPlaceIndex := i
		for j := i; j < len(places); j++ {
			locationCurrent := placesSorted[j].Location
			locationNearest := placesSorted[nearestPlaceIndex].Location

			distanceFromCurrent := prevLocation.DistanceInMeter(locationCurrent)
			distanceFromNearest := prevLocation.DistanceInMeter(locationNearest)
			if distanceFromCurrent < distanceFromNearest {
				nearestPlaceIndex = j
			}
		}

		placesSorted[i], placesSorted[nearestPlaceIndex] = placesSorted[nearestPlaceIndex], placesSorted[i]
		prevLocation = placesSorted[i].Location
	}
	return placesSorted
}

// planTimeFromPlaces プランの所要時間を計算する
func planTimeFromPlaces(locationStart models.GeoLocation, places []models.Place) uint {
	prevLocation := locationStart
	var planTimeInMinutes uint
	for _, place := range places {
		travelTime := prevLocation.TravelTimeTo(place.Location, 80.0)
		planTimeInMinutes += travelTime

		planTimeInMinutes += place.EstimatedStayDuration()

		prevLocation = place.Location
	}

	return planTimeInMinutes
}
