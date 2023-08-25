package plangen

import (
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/google/uuid"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/placefilter"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

const (
	defaultMaxPlanDuration = 180
)

func (s Service) createPlan(
	ctx context.Context,
	locationStart models.GeoLocation,
	placeStart places.Place,
	places []places.Place,
	freeTime *int,
	createBasedOnCurrentLocation bool,
	shouldOpenWhileTraveling bool,
) (*models.Plan, error) {
	placesFiltered := places

	// 現在、開いている場所のみに絞る
	if shouldOpenWhileTraveling {
		placesFiltered = placefilter.FilterByOpeningNow(placesFiltered)
	}

	// 開始地点となる場所から500m圏内の場所に絞る
	placesFiltered = placefilter.FilterWithinDistanceRange(
		placesFiltered,
		placeStart.Location.ToGeoLocation(),
		0,
		500,
	)

	// 会社はプランに含まれないようにする
	placesFiltered = placefilter.FilterCompany(placesFiltered)

	// 場所のカテゴリによるフィルタリング
	placesFiltered = placefilter.FilterIgnoreCategory(placesFiltered)
	placesFiltered = placefilter.FilterByCategory(placesFiltered, models.GetCategoryToFilter(), true)

	// 起点となる場所との距離順でソート
	placesSortedByDistance := placesFiltered
	sort.SliceStable(placesSortedByDistance, func(i, j int) bool {
		locationRecommend := placeStart.Location.ToGeoLocation()
		distanceI := locationRecommend.DistanceInMeter(placesSortedByDistance[i].Location.ToGeoLocation())
		distanceJ := locationRecommend.DistanceInMeter(placesSortedByDistance[j].Location.ToGeoLocation())
		return distanceI < distanceJ
	})

	placesInPlan := make([]models.Place, 0)
	transitions := make([]models.Transition, 0)
	previousLocation := locationStart
	var timeInPlan uint = 0

	// 指定された場所を基準としてプランを作成するときは必ず含める
	if locationStart.Equal(placeStart.Location.ToGeoLocation()) {
		categoryMain := categoryMainOfPlace(placeStart)
		if categoryMain == nil {
			categoryMain = &models.CategoryOther
		}

		placesInPlan = append(placesInPlan, models.Place{
			Id:                    uuid.New().String(),
			Name:                  placeStart.Name,
			GooglePlaceId:         utils.StrPointer(placeStart.PlaceID), // MEMO: 値コピーでないと参照が変化してしまう
			Location:              placeStart.Location.ToGeoLocation(),
			EstimatedStayDuration: categoryMain.EstimatedStayDuration,
			Category:              categoryMain.Name,
		})

		timeInPlan += categoryMain.EstimatedStayDuration
		previousLocation = placeStart.Location.ToGeoLocation()
	}

	for _, place := range placesSortedByDistance {
		var categoriesOfPlace []string
		for _, placeType := range place.Types {
			c := models.CategoryOfSubCategory(placeType)
			if c != nil && !array.IsContain(categoriesOfPlace, c.Name) {
				categoriesOfPlace = append(categoriesOfPlace, c.Name)
			}
		}

		// 飲食店やカフェは複数回含めない
		if isAlreadyHavePlaceCategoryOf(placesInPlan, []models.LocationCategory{
			models.CategoryRestaurant,
			models.CategoryMealTakeaway,
			models.CategoryCafe,
		}) {
			log.Printf("skip place %s because the cafe or restaurant is already in plan\n", place.Name)
			continue
		}

		// MEMO: カテゴリが不明な場合，滞在時間が取得できない
		categoryMain := categoryMainOfPlace(place)
		if categoryMain == nil {
			log.Printf("place %s has no category\n", place.Name)
			continue
		}

		// 予定の時間内に収まらない場合はスキップ
		travelTime := previousLocation.TravelTimeTo(place.Location.ToGeoLocation(), 80.0)
		timeInPlace := categoryMain.EstimatedStayDuration + travelTime
		if freeTime != nil && timeInPlan+timeInPlace > uint(*freeTime) {
			break
		}

		// 予定の時間を指定しない場合、3時間を超えたら終了
		if freeTime == nil && timeInPlan+timeInPlace > defaultMaxPlanDuration {
			break
		}

		// 予定の時間内に閉まってしまう場合はスキップ
		if shouldOpenWhileTraveling && freeTime != nil && !s.isOpeningWithIn(
			ctx,
			place,
			time.Now(),
			time.Minute*time.Duration(*freeTime),
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
			Category:              categoryMain.Name,
			Categories:            models.GetCategoriesFromSubCategories(place.Types),
		})
		timeInPlan += timeInPlace
		previousLocation = place.Location.ToGeoLocation()
		transitions = s.AddTransition(placesInPlan, transitions, travelTime, createBasedOnCurrentLocation)
	}

	if len(placesInPlan) == 0 {
		return nil, fmt.Errorf("could not contain any places in plan")
	}

	// 場所の画像を取得
	performanceTimer := time.Now()
	placesInPlan = s.FetchPlacesPhotos(ctx, placesInPlan)
	log.Printf("fetching place photos took %v\n", time.Since(performanceTimer))

	title, err := s.GeneratePlanTitle(placesInPlan)
	if err != nil {
		log.Printf("error while generating plan title: %v\n", err)
		title = &placeStart.Name
	}

	return &models.Plan{
		Id:            uuid.New().String(),
		Name:          *title,
		Places:        placesInPlan,
		TimeInMinutes: timeInPlan,
		Transitions:   transitions,
	}, nil
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

func categoryMainOfPlace(place places.Place) *models.LocationCategory {
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
