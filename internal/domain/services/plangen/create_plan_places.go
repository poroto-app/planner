package plangen

import (
	"fmt"
	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/placefilter"
)

const (
	defaultMaxPlanDuration = 180
	defaultMaxPlaceInPlan  = 4

	placeDistanceRangeInPlan = 500 // 徒歩5分以内
)

type CreatePlanPlacesInput struct {
	PlanCandidateId         string
	LocationStart           models.GeoLocation
	PlaceStart              models.Place
	Places                  []models.Place
	PlacesOtherPlansContain []models.Place
	CategoryNamesDisliked   *[]string
	FreeTime                *int
	MaxPlace                int
}

// CreatePlanPlaces プランの候補地となる場所を作成する
func (s Service) CreatePlanPlaces(input CreatePlanPlacesInput) ([]models.Place, error) {
	if input.PlanCandidateId == "" {
		panic("PlanCandidateId is required")
	}

	if input.MaxPlace == 0 {
		input.MaxPlace = defaultMaxPlaceInPlan
	}

	/**
	* プラン作成の方針
	* 1. スタート地点から近い場所の中で、レビューの高い場所を選択
	* 2. その場所から近い場所の中で、レビューの高い場所を選択
	* 3. 1, 2を繰り返し、プランに含まれる場所がMaxPlaceに達するまで続ける
	 */
	placesInPlan := make([]models.Place, 0)
	placesInPlan = append(placesInPlan, input.PlaceStart)
	for len(placesInPlan) < input.MaxPlace {
		prevPlace := placesInPlan[len(placesInPlan)-1]
		nextPlace := s.getNextPlaceForPlan(prevPlace, placesInPlan, input, placeDistanceRangeInPlan)
		if nextPlace == nil {
			break
		}
		placesInPlan = append(placesInPlan, *nextPlace)
	}

	if len(placesInPlan) == 0 {
		return nil, fmt.Errorf("could not contain any Places in plan")
	}

	return placesInPlan, nil
}

func (s Service) getNextPlaceForPlan(prevPlace models.Place, placesInPlan []models.Place, input CreatePlanPlacesInput, placeDistanceRangeInPlan float64) *models.Place {
	// 最後に追加した場所から近い場所を選択
	placesFiltered := placefilter.FilterDefaultIgnore(placefilter.FilterDefaultIgnoreInput{
		Places:              input.Places,
		StartLocation:       prevPlace.Location,
		IgnoreDistanceRange: placeDistanceRangeInPlan,
	})
	s.logger.Debug("Places after filtering by default ignore", zap.Int("Places", len(placesFiltered)))

	// ユーザーが拒否した場所は取り除く
	if input.CategoryNamesDisliked != nil {
		categoriesDisliked := models.GetCategoriesFromSubCategories(*input.CategoryNamesDisliked)
		placesFiltered = placefilter.FilterByCategory(placesFiltered, categoriesDisliked, false)
		s.logger.Debug("Places after filtering by disliked categories", zap.Int("Places", len(placesFiltered)))
	}

	// 他のプランに含まれている場所を除外する
	placesFiltered = placefilter.FilterPlaces(placesFiltered, func(place models.Place) bool {
		if input.PlacesOtherPlansContain == nil {
			return true
		}

		for _, placeOtherPlanContain := range input.PlacesOtherPlansContain {
			if place.Id == placeOtherPlanContain.Id {
				return false
			}
		}
		return true
	})
	s.logger.Debug("places after filtering by other plans", zap.Int("places", len(placesFiltered)))

	if len(placesFiltered) == 0 {
		return nil
	}

	// レビューの高い場所からプランに含められる場所を選択
	for _, place := range models.SortPlacesByRating(placesFiltered) {
		if s.checkForIncludeForPlan(place, placesInPlan, input) {
			return &place
		}
	}

	return nil
}

func (s Service) checkForIncludeForPlan(
	place models.Place,
	placesInPlan []models.Place,
	input CreatePlanPlacesInput,
) bool {
	// すでにプランに含まれている場所はスキップ
	if _, isAlreadyInPlan := array.Find(placesInPlan, func(p models.Place) bool {
		return p.Id == place.Id
	}); isAlreadyInPlan {
		return false
	}

	// メインカテゴリが飲食店の場所が、一定数以上含まれないようにする
	for _, condition := range []struct {
		category            models.LocationCategory
		numPlacesCanContain int
	}{
		{models.CategoryRestaurant, 1},
		{models.CategoryCafe, 2},
		{models.CategoryBakery, 2},
	} {
		if place.MainCategory() == nil || !place.MainCategory().IsCategoryOf(condition.category) {
			continue
		}

		placesInPlanWithCategory := array.Filter(placesInPlan, func(placeInPlan models.Place) bool {
			if placeInPlan.MainCategory() == nil {
				return false
			}
			return placeInPlan.MainCategory().IsCategoryOf(condition.category)
		})
		if len(placesInPlanWithCategory) >= condition.numPlacesCanContain {
			s.logger.Debug(
				fmt.Sprintf("skip place because the %d %s places are already in plan", condition.numPlacesCanContain, condition.category.Name),
				zap.String("place", place.Google.Name),
			)
			return false
		}
	}

	// 最適経路で巡ったときの所要時間が予定の時間を超える場合はスキップ
	sortedByDistance := sortPlacesByDistanceFrom(input.LocationStart, append(placesInPlan, place))
	timeInPlan := planTimeFromPlaces(input.LocationStart, sortedByDistance)
	if input.FreeTime != nil && timeInPlan > uint(*input.FreeTime) {
		s.logger.Debug(
			"skip place because it will be over time",
			zap.String("place", place.Google.Name),
			zap.Uint("timeInPlan", timeInPlan),
			zap.Int("FreeTime", *input.FreeTime),
		)

		return false
	}

	// 予定の時間を指定しない場合、3時間を超える場合はスキップ
	if input.FreeTime == nil && timeInPlan > defaultMaxPlanDuration {
		s.logger.Debug(
			"skip place because it will be over time",
			zap.String("place", place.Google.Name),
			zap.Uint("timeInPlan", timeInPlan),
			zap.Int("defaultMaxPlanDuration", defaultMaxPlanDuration),
		)

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
