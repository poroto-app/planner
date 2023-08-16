package plangen

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/google/uuid"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/placefilter"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

func (s Service) CreatePlan(
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

	for _, place := range placesWithInRange {
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

		var categoryMain *models.LocationCategory
		for _, placeType := range place.Types {
			c := models.CategoryOfSubCategory(placeType)
			if c != nil {
				categoryMain = c
				break
			}
		}
		// MEMO: カテゴリが不明な場合，滞在時間が取得できない
		if categoryMain == nil {
			log.Printf("place %s has no category\n", place.Name)
			continue
		}

		tripTime := s.travelTimeBetween(
			previousLocation,
			place.Location.ToGeoLocation(),
			80.0,
		)
		timeInPlace := categoryMain.EstimatedStayDuration + tripTime
		if freeTime != nil && timeInPlan+timeInPlace > uint(*freeTime) {
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
		transitions = s.AddTransition(placesInPlan, transitions, tripTime, createBasedOnCurrentLocation)
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

// isOpeningWithIn は，指定された場所が指定された時間内に開いているかを判定する
func (s Service) isOpeningWithIn(
	ctx context.Context,
	place places.Place,
	startTime time.Time,
	duration time.Duration,
) bool {
	placeOpeningPeriods, err := s.placesApi.FetchPlaceOpeningPeriods(ctx, place)
	if err != nil {
		log.Printf("error while fetching place periods: %v\n", err)
		return false
	}
	// 時刻フィルタリング用変数
	endTime := startTime.Add(time.Minute * duration)
	today := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location())

	for _, placeOpeningPeriod := range placeOpeningPeriods {
		weekday := startTime.Weekday()
		if placeOpeningPeriod.DayOfWeek != weekday.String() {
			continue
		}

		// TODO: パース処理に関するテストを書く
		openingPeriodHour, opHourErr := strconv.Atoi(placeOpeningPeriod.OpeningTime[:2])
		openingPeriodMinute, opMinuteErr := strconv.Atoi(placeOpeningPeriod.OpeningTime[2:])
		closingPeriodHour, clHourErr := strconv.Atoi(placeOpeningPeriod.ClosingTime[:2])
		closingPeriodMinute, clMinuteErr := strconv.Atoi(placeOpeningPeriod.ClosingTime[2:])
		if opHourErr != nil || opMinuteErr != nil || clHourErr != nil || clMinuteErr != nil {
			log.Println("error while converting period [string->int]")
			continue
		}

		openingTime := today.Add(time.Hour*time.Duration(openingPeriodHour) + time.Minute*time.Duration(openingPeriodMinute))
		closingTime := today.Add(time.Hour*time.Duration(closingPeriodHour) + time.Minute*time.Duration(closingPeriodMinute))

		// 開店時刻 < 開始時刻 && 終了時刻 < 閉店時刻 の判断
		if startTime.After(openingTime) && endTime.Before(closingTime) {
			return true
		}
	}
	return false
}

func (s Service) travelTimeBetween(
	locationDeparture models.GeoLocation,
	locationDestination models.GeoLocation,
	meterPerMinutes float64,
) uint {
	var timeInMinutes uint = 0
	distance := locationDeparture.DistanceInMeter(locationDestination)
	if distance > 0.0 && meterPerMinutes > 0.0 {
		timeInMinutes = uint(distance / meterPerMinutes)
	}
	return timeInMinutes
}
