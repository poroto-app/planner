package plangen

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"time"

	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/placefilter"
)

const (
	defaultMaxPlanDuration = 180
	defaultMaxPlace        = 4

	placeDistanceRangeInPlan = 1500
)

// TODO: 大文字スタートにする
type CreatePlanPlacesParams struct {
	planCandidateId              string
	locationStart                models.GeoLocation
	placeStart                   models.Place
	places                       []models.Place
	placesOtherPlansContain      []models.Place
	categoryNamesDisliked        *[]string
	freeTime                     *int
	createBasedOnCurrentLocation bool
	shouldOpenWhileTraveling     bool
	maxPlace                     int
}

// createPlanPlaces プランの候補地となる場所を作成する
func (s Service) createPlanPlaces(ctx context.Context, params CreatePlanPlacesParams) ([]models.Place, error) {
	if params.planCandidateId == "" {
		panic("planCandidateId is required")
	}

	if params.maxPlace == 0 {
		params.maxPlace = defaultMaxPlace
	}

	placesFiltered := params.places
	placesFiltered = placefilter.FilterDefaultIgnore(placefilter.FilterDefaultIgnoreInput{
		Places:        placesFiltered,
		StartLocation: params.locationStart,
	})

	// 現在、開いている場所のみに絞る
	if params.shouldOpenWhileTraveling {
		placesFiltered = placefilter.FilterByOpeningNow(placesFiltered)
	}

	// 開始地点となる場所から1500m圏内の場所に絞る
	placesFiltered = placefilter.FilterWithinDistanceRange(
		placesFiltered,
		params.placeStart.Location,
		0,
		placeDistanceRangeInPlan,
	)

	// ユーザーが拒否した場所は取り除く
	if params.categoryNamesDisliked != nil {
		categoriesDisliked := models.GetCategoriesFromSubCategories(*params.categoryNamesDisliked)
		placesFiltered = placefilter.FilterByCategory(placesFiltered, categoriesDisliked, false)
	}

	// 他のプランに含まれている場所を除外する
	placesFiltered = placefilter.FilterPlaces(placesFiltered, func(place models.Place) bool {
		if params.placesOtherPlansContain == nil {
			return true
		}

		for _, placeOtherPlanContain := range params.placesOtherPlansContain {
			if place.Id == placeOtherPlanContain.Id {
				return false
			}
		}
		return true
	})

	// レビューの高い順でソート
	placesSorted := models.SortPlacesByRating(placesFiltered)

	placesInPlan := make([]models.Place, 0)

	// 指定された場所を基準としてプランを作成するときは必ず含める
	if params.locationStart.Equal(params.placeStart.Location) {
		placesInPlan = append(placesInPlan, params.placeStart)
	}

	for _, place := range placesSorted {
		// プランに含まれる場所の数が上限に達したら終了
		if len(placesInPlan) >= params.maxPlace {
			s.logger.Debug(
				"skip place because the number of places in plan is over",
				zap.String("place", place.Google.Name),
				zap.Int("maxPlace", params.maxPlace),
				zap.Int("placesInPlan", len(placesInPlan)),
			)
			break
		}

		if place.Id == params.placeStart.Id {
			continue
		}

		// 飲食店やカフェは複数回含めない
		categoriesFood := []models.LocationCategory{
			models.CategoryRestaurant,
			models.CategoryMealTakeaway,
			models.CategoryCafe,
		}
		if isAlreadyHavePlaceCategoryOf(placesInPlan, categoriesFood) && isCategoryOf(place.Google.Types, categoriesFood) {
			s.logger.Debug(
				"skip place because the cafe or restaurant is already in plan",
				zap.String("place", place.Google.Name),
			)
			continue
		}

		// 最適経路で巡ったときの所要時間を計算
		sortedByDistance := sortPlacesByDistanceFrom(params.locationStart, append(placesInPlan, place))
		timeInPlan := planTimeFromPlaces(params.locationStart, sortedByDistance)

		// 予定の時間内に収まらない場合はスキップ
		if params.freeTime != nil && timeInPlan > uint(*params.freeTime) {
			s.logger.Debug(
				"skip place because it will be over time",
				zap.String("place", place.Google.Name),
				zap.Uint("timeInPlan", timeInPlan),
				zap.Int("freeTime", *params.freeTime),
			)
			continue
		}

		// 予定の時間を指定しない場合、3時間を超える場合はスキップ
		if params.freeTime == nil && timeInPlan > defaultMaxPlanDuration {
			s.logger.Debug(
				"skip place because it will be over time",
				zap.String("place", place.Google.Name),
				zap.Uint("timeInPlan", timeInPlan),
				zap.Int("defaultMaxPlanDuration", defaultMaxPlanDuration),
			)
			continue
		}

		if params.shouldOpenWhileTraveling && params.freeTime == nil {
			// 場所の詳細を取得(Place Detailリクエストが発生するため、ある程度フィルタリングしたあとに行う)
			placeDetail, err := s.placeService.FetchPlaceDetailAndSave(ctx, place.Google.PlaceId)
			if err != nil {
				s.logger.Warn(
					"error while fetching place detail",
					zap.String("place", place.Google.Name),
					zap.Error(err),
				)
				continue
			}

			place.Google.PlaceDetail = placeDetail

			// 予定の時間内に閉まってしまう場合はスキップ
			isOpeningWhilePlan, err := s.isOpeningWithIn(place, time.Now(), time.Minute*time.Duration(timeInPlan))
			if err != nil {
				s.logger.Warn(
					"error while checking opening hours",
					zap.String("place", place.Google.Name),
					zap.Error(err),
				)
				continue
			}

			if !isOpeningWhilePlan {
				s.logger.Debug(
					"skip place because it will be closed",
					zap.String("place", place.Google.Name),
				)
				continue
			}
		}

		placesInPlan = append(placesInPlan, place)
	}

	if len(placesInPlan) == 0 {
		return nil, fmt.Errorf("could not contain any places in plan")
	}

	return placesInPlan, nil
}

func isAlreadyHavePlaceCategoryOf(placesInPlan []models.Place, categories []models.LocationCategory) bool {
	var categoriesInPlan []models.LocationCategory
	for _, place := range placesInPlan {
		categoriesInPlan = append(categoriesInPlan, place.Categories()...)
	}

	for _, category := range categories {
		for _, categoryInPlan := range categoriesInPlan {
			if categoryInPlan.Name == category.Name {
				return true
			}
		}
	}
	return false
}

func isCategoryOf(placeTypes []string, categories []models.LocationCategory) bool {
	categoriesOfPlace := models.GetCategoriesFromSubCategories(placeTypes)
	for _, category := range categories {
		for _, categoryOfPlace := range categoriesOfPlace {
			if categoryOfPlace.Name == category.Name {
				return true
			}
		}
	}
	return false
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
