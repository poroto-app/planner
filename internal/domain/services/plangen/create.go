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
) (*models.Plan, error) {
	// 起点となる場所との距離順でソート
	placesSortedByDistance := places
	sort.SliceStable(placesSortedByDistance, func(i, j int) bool {
		locationRecommend := placeStart.Location.ToGeoLocation()
		distanceI := locationRecommend.DistanceInMeter(placesSortedByDistance[i].Location.ToGeoLocation())
		distanceJ := locationRecommend.DistanceInMeter(placesSortedByDistance[j].Location.ToGeoLocation())
		return distanceI < distanceJ
	})

	placesWithInRange := placefilter.FilterWithinDistanceRange(
		placesSortedByDistance,
		placeStart.Location.ToGeoLocation(),
		0,
		500,
	)

	placesInPlan := make([]models.Place, 0)
	categoriesInPlan := make([]string, 0)
	transitions := make([]models.Transition, 0)
	previousLocation := locationStart
	var timeInPlan uint = 0

	for i, place := range placesWithInRange {
		var categoryMain *models.LocationCategory
		for _, placeType := range place.Types {
			c := models.CategoryOfSubCategory(placeType)
			if c != nil {
				categoryMain = c
				break
			}
		}

		// 開始位置と開始場所が同じ場合は必ず含める
		// TODO: 開始位置と開始場所が同じ場所は、場所を含めるかのチェックを飛ばすようにして、追加の処理は共通化する
		if i == 0 && locationStart.Equal(placeStart.Location.ToGeoLocation()) {
			if categoryMain == nil {
				categoryMain = &models.CategoryOther
			}
			placesInPlan = append(placesInPlan, models.Place{
				Id:                    uuid.New().String(),
				Name:                  place.Name,
				GooglePlaceId:         utils.StrPointer(place.PlaceID), // MEMO: 値コピーでないと参照が変化してしまう
				Location:              place.Location.ToGeoLocation(),
				EstimatedStayDuration: categoryMain.EstimatedStayDuration,
				Category:              categoryMain.Name,
			})

			timeInPlan += categoryMain.EstimatedStayDuration
			if categoryMain.Name != models.CategoryOther.Name {
				categoriesInPlan = append(categoriesInPlan, categoryMain.Name)
			}
			previousLocation = place.Location.ToGeoLocation()
			continue
		}

		var categoriesOfPlace []string
		for _, placeType := range place.Types {
			c := models.CategoryOfSubCategory(placeType)
			if c != nil && !array.IsContain(categoriesOfPlace, c.Name) {
				categoriesOfPlace = append(categoriesOfPlace, c.Name)
			}
		}

		// 飲食店系は複数含めない
		categoriesFood := []string{
			models.CategoryRestaurant.Name,
			models.CategoryMealTakeaway.Name,
		}
		isFoodPlace := array.HasIntersection(categoriesOfPlace, categoriesFood)
		isPlanContainsFoodPlace := array.HasIntersection(categoriesInPlan, categoriesFood)
		if isFoodPlace && isPlanContainsFoodPlace {
			log.Printf("skip place %s because plan is already has food place\n", place.Name)
			continue
		}

		// カフェを複数含めない
		isCafePlace := array.IsContain(categoriesOfPlace, models.CategoryCafe.Name)
		isPlanContainsFoodPlace = array.IsContain(categoriesInPlan, models.CategoryCafe.Name)
		if isCafePlace && isPlanContainsFoodPlace {
			log.Printf("skip place %s because plan is already has cafe place\n", place.Name)
			continue
		}

		// MEMO: カテゴリが不明な場合，滞在時間が取得できない
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
		if freeTime != nil && !s.isOpeningWithIn(
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
		})
		timeInPlan += timeInPlace
		categoriesInPlan = append(categoriesInPlan, categoryMain.Name)
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
