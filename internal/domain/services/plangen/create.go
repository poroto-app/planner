package plangen

import (
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/placefilter"
)

const (
	defaultMaxPlanDuration = 180
)

type CreatePlanPlacesParams struct {
	locationStart                models.GeoLocation
	placeStart                   models.GooglePlace
	places                       []models.GooglePlace
	placesOtherPlansContain      []models.GooglePlace
	freeTime                     *int
	createBasedOnCurrentLocation bool
	shouldOpenWhileTraveling     bool
}

// createPlanPlaces プランの候補地となる場所を作成する
func (s Service) createPlanPlaces(ctx context.Context, params CreatePlanPlacesParams) ([]models.GooglePlace, error) {
	placesFiltered := params.places

	// 現在、開いている場所のみに絞る
	if params.shouldOpenWhileTraveling {
		placesFiltered = placefilter.FilterByOpeningNow(placesFiltered)
	}

	// 開始地点となる場所から1500m圏内の場所に絞る
	placesFiltered = placefilter.FilterWithinDistanceRange(
		placesFiltered,
		params.placeStart.Location,
		0,
		1500,
	)

	// 重複した場所を削除
	placesFiltered = placefilter.FilterDuplicated(placesFiltered)

	// 会社はプランに含まれないようにする
	placesFiltered = placefilter.FilterCompany(placesFiltered)

	// 場所のカテゴリによるフィルタリング
	placesFiltered = placefilter.FilterIgnoreCategory(placesFiltered)
	placesFiltered = placefilter.FilterByCategory(placesFiltered, models.GetCategoryToFilter(), true)

	// 他のプランに含まれている場所を除外する
	placesFiltered = placefilter.FilterPlaces(placesFiltered, func(place models.GooglePlace) bool {
		if params.placesOtherPlansContain == nil {
			return true
		}

		for _, placeOtherPlanContain := range params.placesOtherPlansContain {
			if place.PlaceId == placeOtherPlanContain.PlaceId {
				return false
			}
		}
		return true
	})

	// レビューの高い順でソート
	placesSorted := placesFiltered
	sort.SliceStable(placesSorted, func(i, j int) bool {
		return placesSorted[i].Rating > placesSorted[j].Rating
	})

	placesInPlan := make([]models.GooglePlace, 0)

	// 指定された場所を基準としてプランを作成するときは必ず含める
	if params.locationStart.Equal(params.placeStart.Location) {
		placesInPlan = append(placesInPlan, params.placeStart)
	}

	for _, place := range placesSorted {
		if place.PlaceId == params.placeStart.PlaceId {
			continue
		}

		// 飲食店やカフェは複数回含めない
		categoriesFood := []models.LocationCategory{
			models.CategoryRestaurant,
			models.CategoryMealTakeaway,
			models.CategoryCafe,
		}
		if isAlreadyHavePlaceCategoryOf(placesInPlan, categoriesFood) && isCategoryOf(place.Types, categoriesFood) {
			log.Printf("skip place %s because the cafe or restaurant is already in plan\n", place.Name)
			continue
		}

		// 最適経路で巡ったときの所要時間を計算
		sortedByDistance := sortPlacesByDistanceFrom(params.locationStart, append(placesInPlan, place))
		timeInPlan := planTimeFromPlaces(params.locationStart, sortedByDistance)

		// 予定の時間内に収まらない場合はスキップ
		if params.freeTime != nil && timeInPlan > uint(*params.freeTime) {
			log.Printf("skip place %s because it will be over time\n", place.Name)
			continue
		}

		// 予定の時間を指定しない場合、3時間を超える場合はスキップ
		if params.freeTime == nil && timeInPlan > defaultMaxPlanDuration {
			log.Printf("skip place %s because it will be over time\n", place.Name)
			continue
		}

		// 予定の時間内に閉まってしまう場合はスキップ
		if params.shouldOpenWhileTraveling && params.freeTime != nil && !s.isOpeningWithIn(
			ctx,
			place,
			time.Now(),
			time.Minute*time.Duration(*params.freeTime),
		) {
			log.Printf("skip place %s because it will be closed\n", place.Name)
			continue
		}

		placesInPlan = append(placesInPlan, place)
	}

	if len(placesInPlan) == 0 {
		return nil, fmt.Errorf("could not contain any places in plan")
	}

	return placesInPlan, nil
}

func isAlreadyHavePlaceCategoryOf(placesInPlan []models.GooglePlace, categories []models.LocationCategory) bool {
	var categoriesInPlan []models.LocationCategory
	for _, place := range placesInPlan {
		cs := models.GetCategoriesFromSubCategories(place.Types)
		categoriesInPlan = append(categoriesInPlan, cs...)
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
func sortPlacesByDistanceFrom(location models.GeoLocation, places []models.GooglePlace) []models.GooglePlace {
	placesSorted := make([]models.GooglePlace, len(places))
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
func planTimeFromPlaces(locationStart models.GeoLocation, places []models.GooglePlace) uint {
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
