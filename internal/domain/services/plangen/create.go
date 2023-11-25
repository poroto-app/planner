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
	defaultMaxPlace        = 4
)

type CreatePlanPlacesParams struct {
	planCandidateId              string
	locationStart                models.GeoLocation
	placeStart                   models.PlaceInPlanCandidate
	places                       []models.PlaceInPlanCandidate
	placesOtherPlansContain      []models.PlaceInPlanCandidate
	freeTime                     *int
	createBasedOnCurrentLocation bool
	shouldOpenWhileTraveling     bool
	maxPlace                     int
}

// createPlanPlaces プランの候補地となる場所を作成する
func (s Service) createPlanPlaces(ctx context.Context, params CreatePlanPlacesParams) ([]models.PlaceInPlanCandidate, error) {
	if params.planCandidateId == "" {
		panic("planCandidateId is required")
	}

	if params.maxPlace == 0 {
		params.maxPlace = defaultMaxPlace
	}

	placesFiltered := params.places

	// 現在、開いている場所のみに絞る
	if params.shouldOpenWhileTraveling {
		placesFiltered = placefilter.FilterByOpeningNow(placesFiltered)
	}

	// 開始地点となる場所から1500m圏内の場所に絞る
	placesFiltered = placefilter.FilterWithinDistanceRange(
		placesFiltered,
		params.placeStart.Location(),
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
	placesFiltered = placefilter.FilterPlaces(placesFiltered, func(place models.PlaceInPlanCandidate) bool {
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
	placesSorted := placesFiltered
	sort.SliceStable(placesSorted, func(i, j int) bool {
		return placesSorted[i].Google.Rating > placesSorted[j].Google.Rating
	})

	placesInPlan := make([]models.PlaceInPlanCandidate, 0)

	// 指定された場所を基準としてプランを作成するときは必ず含める
	if params.locationStart.Equal(params.placeStart.Location()) {
		placesInPlan = append(placesInPlan, params.placeStart)
	}

	for _, place := range placesSorted {
		// プランに含まれる場所の数が上限に達したら終了
		if len(placesInPlan) >= params.maxPlace {
			log.Printf("skip place %s because the number of places in plan is over\n", place.Google.Name)
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
			log.Printf("skip place %s because the cafe or restaurant is already in plan\n", place.Google.Name)
			continue
		}

		// 最適経路で巡ったときの所要時間を計算
		sortedByDistance := sortPlacesByDistanceFrom(params.locationStart, append(placesInPlan, place))
		timeInPlan := planTimeFromPlaces(params.locationStart, sortedByDistance)

		// 予定の時間内に収まらない場合はスキップ
		if params.freeTime != nil && timeInPlan > uint(*params.freeTime) {
			log.Printf("skip place %s because it will be over time\n", place.Google.Name)
			continue
		}

		// 予定の時間を指定しない場合、3時間を超える場合はスキップ
		if params.freeTime == nil && timeInPlan > defaultMaxPlanDuration {
			log.Printf("skip place %s because it will be over time\n", place.Google.Name)
			continue
		}

		if params.shouldOpenWhileTraveling && params.freeTime == nil {
			// 場所の詳細を取得(Place Detailリクエストが発生するため、ある程度フィルタリングしたあとに行う)
			placeDetail, err := s.placeService.FetchPlaceDetailAndSave(ctx, params.planCandidateId, place.Google.PlaceId)
			if err != nil {
				log.Printf("error while fetching place detail: %v\n", err)
				continue
			}

			place.Google.PlaceDetail = placeDetail

			// 予定の時間内に閉まってしまう場合はスキップ
			isOpeningWhilePlan, err := s.isOpeningWithIn(place, time.Now(), time.Minute*time.Duration(timeInPlan))
			if err != nil {
				log.Printf("error while checking opening hours: %v\n", err)
				continue
			}

			if !isOpeningWhilePlan {
				log.Printf("skip place %s because it will be closed\n", place.Google.Name)
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

func isAlreadyHavePlaceCategoryOf(placesInPlan []models.PlaceInPlanCandidate, categories []models.LocationCategory) bool {
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
func sortPlacesByDistanceFrom(location models.GeoLocation, places []models.PlaceInPlanCandidate) []models.PlaceInPlanCandidate {
	placesSorted := make([]models.PlaceInPlanCandidate, len(places))
	copy(placesSorted, places)

	prevLocation := location
	for i := 0; i < len(places); i++ {
		nearestPlaceIndex := i
		for j := i; j < len(places); j++ {
			locationCurrent := placesSorted[j].Location()
			locationNearest := placesSorted[nearestPlaceIndex].Location()

			distanceFromCurrent := prevLocation.DistanceInMeter(locationCurrent)
			distanceFromNearest := prevLocation.DistanceInMeter(locationNearest)
			if distanceFromCurrent < distanceFromNearest {
				nearestPlaceIndex = j
			}
		}

		placesSorted[i], placesSorted[nearestPlaceIndex] = placesSorted[nearestPlaceIndex], placesSorted[i]
		prevLocation = placesSorted[i].Location()
	}
	return placesSorted
}

// planTimeFromPlaces プランの所要時間を計算する
func planTimeFromPlaces(locationStart models.GeoLocation, places []models.PlaceInPlanCandidate) uint {
	prevLocation := locationStart
	var planTimeInMinutes uint
	for _, place := range places {
		travelTime := prevLocation.TravelTimeTo(place.Location(), 80.0)
		planTimeInMinutes += travelTime

		planTimeInMinutes += place.EstimatedStayDuration()

		prevLocation = place.Location()
	}

	return planTimeInMinutes
}
