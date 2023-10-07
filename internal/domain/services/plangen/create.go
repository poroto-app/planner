package plangen

import (
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/google/uuid"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/placefilter"
	"poroto.app/poroto/planner/internal/domain/utils"
	api "poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

const (
	defaultMaxPlanDuration = 180
)

type CreatePlanPlacesParams struct {
	locationStart                models.GeoLocation
	placeStart                   api.Place
	places                       []api.Place
	placesOtherPlansContain      []models.Place
	freeTime                     *int
	createBasedOnCurrentLocation bool
	shouldOpenWhileTraveling     bool
}

// createPlanPlaces プランの候補地となる場所を作成する
func (s Service) createPlanPlaces(ctx context.Context, params CreatePlanPlacesParams) ([]models.Place, error) {
	placesFiltered := params.places

	// 現在、開いている場所のみに絞る
	if params.shouldOpenWhileTraveling {
		placesFiltered = placefilter.FilterByOpeningNow(placesFiltered)
	}

	// 開始地点となる場所から1500m圏内の場所に絞る
	placesFiltered = placefilter.FilterWithinDistanceRange(
		placesFiltered,
		params.placeStart.Location.ToGeoLocation(),
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
	placesFiltered = placefilter.FilterPlaces(placesFiltered, func(place api.Place) bool {
		if params.placesOtherPlansContain == nil {
			return true
		}

		for _, placeOtherPlanContain := range params.placesOtherPlansContain {
			if placeOtherPlanContain.GooglePlaceId != nil && place.PlaceID == *placeOtherPlanContain.GooglePlaceId {
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

	placesInPlan := make([]models.Place, 0)

	// 指定された場所を基準としてプランを作成するときは必ず含める
	if params.locationStart.Equal(params.placeStart.Location.ToGeoLocation()) {
		categoryMain := categoryMainOfPlace(params.placeStart)
		if categoryMain == nil {
			categoryMain = &models.CategoryOther
		}

		placesInPlan = append(placesInPlan, models.Place{
			Id:                    uuid.New().String(),
			Name:                  params.placeStart.Name,
			GooglePlaceId:         utils.StrPointer(params.placeStart.PlaceID), // MEMO: 値コピーでないと参照が変化してしまう
			Location:              params.placeStart.Location.ToGeoLocation(),
			EstimatedStayDuration: categoryMain.EstimatedStayDuration,
		})
	}

	for _, place := range placesSorted {
		if place.PlaceID == params.placeStart.PlaceID {
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

		// MEMO: カテゴリが不明な場合，滞在時間が取得できない
		categoryMain := categoryMainOfPlace(place)
		if categoryMain == nil {
			log.Printf("place %s has no category\n", place.Name)
			continue
		}

		// 最適経路で巡ったときの所要時間を計算
		sortedByDistance := sortPlacesByDistanceFrom(params.locationStart, append(placesInPlan, models.Place{
			Location:              place.Location.ToGeoLocation(),
			EstimatedStayDuration: categoryMain.EstimatedStayDuration,
			Categories:            models.GetCategoriesFromSubCategories(place.Types),
		}))
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

		placesInPlan = append(placesInPlan, models.Place{
			Id:                    uuid.New().String(),
			Name:                  place.Name,
			GooglePlaceId:         utils.StrPointer(place.PlaceID), // MEMO: 値コピーでないと参照が変化してしまう
			Location:              place.Location.ToGeoLocation(),
			EstimatedStayDuration: categoryMain.EstimatedStayDuration,
			Categories:            models.GetCategoriesFromSubCategories(place.Types),
		})
	}

	if len(placesInPlan) == 0 {
		return nil, fmt.Errorf("could not contain any places in plan")
	}

	return placesInPlan, nil
}

type CreatePlanParams struct {
	locationStart models.GeoLocation
	placeStart    api.Place
	places        []models.Place
}

func (s Service) createPlans(ctx context.Context, params ...CreatePlanParams) []models.Plan {
	ch := make(chan *models.Plan, len(params))
	for _, param := range params {
		go func(ctx context.Context, param CreatePlanParams, ch chan<- *models.Plan) {
			places := param.places

			chPlaceWithPhotos := make(chan []models.Place, 1)
			go func(ctx context.Context, places []models.Place, chPlaceWithPhotos chan<- []models.Place) {
				// 場所の画像を取得
				performanceTimer := time.Now()
				places = s.FetchPlacesPhotos(ctx, places)
				log.Printf("fetching place photos took %v\n", time.Since(performanceTimer))
				chPlaceWithPhotos <- places
			}(ctx, places, chPlaceWithPhotos)

			// プランのタイトルを生成
			chPlanTitle := make(chan string, 1)
			go func(ctx context.Context, places []models.Place, chPlanTitle chan<- string) {
				performanceTimer := time.Now()
				title, err := s.GeneratePlanTitle(param.places)
				if err != nil {
					log.Printf("error while generating plan title: %v\n", err)
					title = &param.placeStart.Name
				}
				log.Printf("generating plan title took %v\n", time.Since(performanceTimer))
				chPlanTitle <- *title
			}(ctx, places, chPlanTitle)

			// 場所のレビューを取得
			chPlansWithReviews := make(chan []models.Place, 1)
			go func(ctx context.Context, places []models.Place, chPlansWithReviews chan<- []models.Place) {
				performanceTimer := time.Now()
				places = s.FetchReviews(ctx, places)
				log.Printf("fetching place reviews took %v\n", time.Since(performanceTimer))
				chPlansWithReviews <- places
			}(ctx, places, chPlansWithReviews)

			// タイトル生成には2秒以上かかる場合があるため、タイムアウト処理を行う
			var title string
			chTitleTimeOut := time.NewTimer(2 * time.Second)
			select {
			case title = <-chPlanTitle:
				chTitleTimeOut.Stop()
			case <-chTitleTimeOut.C:
				log.Printf("timeout while generating plan title\n")
				title = param.placeStart.Name
			}

			placesWithPhotos := <-chPlaceWithPhotos
			placesWithReviews := <-chPlansWithReviews
			for i := 0; i < len(places); i++ {
				places[i].Images = placesWithPhotos[i].Images
				places[i].GooglePlaceReviews = placesWithReviews[i].GooglePlaceReviews
			}

			places = sortPlacesByDistanceFrom(param.locationStart, places)
			timeInPlan := planTimeFromPlaces(param.locationStart, places)

			ch <- &models.Plan{
				Id:            uuid.New().String(),
				Name:          title,
				Places:        places,
				TimeInMinutes: timeInPlan,
			}
		}(ctx, param, ch)
	}

	plans := make([]models.Plan, 0)
	for i := 0; i < len(params); i++ {
		plan := <-ch
		if plan == nil {
			continue
		}
		plans = append(plans, *plan)
	}

	return plans
}

func isAlreadyHavePlaceCategoryOf(placesInPlan []models.Place, categories []models.LocationCategory) bool {
	var categoriesInPlan []models.LocationCategory
	for _, place := range placesInPlan {
		categoriesInPlan = append(categoriesInPlan, place.Categories...)
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

func categoryMainOfPlace(place api.Place) *models.LocationCategory {
	var categoryMain *models.LocationCategory
	for _, placeType := range place.Types {
		c := models.CategoryOfSubCategory(placeType)
		if c != nil {
			categoryMain = c
			break
		}
	}
	return categoryMain
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

		if len(place.Categories) > 0 {
			categoryMain := place.Categories[0]
			planTimeInMinutes += categoryMain.EstimatedStayDuration
		}

		prevLocation = place.Location
	}

	return planTimeInMinutes
}
