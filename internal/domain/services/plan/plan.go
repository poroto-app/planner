package plan

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
	"poroto.app/poroto/planner/internal/domain/repository"
	"poroto.app/poroto/planner/internal/domain/services/placefilter"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
	"poroto.app/poroto/planner/internal/infrastructure/api/openai"
	"poroto.app/poroto/planner/internal/infrastructure/firestore"
)

type PlanService struct {
	placesApi                  places.PlacesApi
	planRepository             repository.PlanRepository
	planCandidateRepository    repository.PlanCandidateRepository
	openaiChatCompletionClient openai.ChatCompletionClient
}

func NewPlanService(ctx context.Context) (*PlanService, error) {
	placesApi, err := places.NewPlacesApi()
	if err != nil {
		return nil, fmt.Errorf("error while initizalizing places api: %v", err)
	}

	planRepository, err := firestore.NewPlanRepository(ctx)
	if err != nil {
		return nil, err
	}

	planCandidateRepository, err := firestore.NewPlanCandidateRepository(ctx)
	if err != nil {
		return nil, err
	}

	openaiChatCompletionClient, err := openai.NewChatCompletionClient()
	if err != nil {
		return nil, fmt.Errorf("error while initializing openai chat completion client: %v", err)
	}

	return &PlanService{
		placesApi:                  *placesApi,
		planRepository:             planRepository,
		planCandidateRepository:    planCandidateRepository,
		openaiChatCompletionClient: *openaiChatCompletionClient,
	}, err
}

func (s PlanService) CreatePlanByLocation(
	ctx context.Context,
	locationStart models.GeoLocation,
	// TODO: ユーザーに却下された場所を引数にする（プランを作成時により多くの場所を取得した場合、YESと答えたカテゴリの場所からしかプランを作成できなくなるため）
	categoryNamesPreferred *[]string,
	freeTime *int,
	createBasedOnCurrentLocation bool,
) (*[]models.Plan, error) {
	placesSearched, err := s.placesApi.FindPlacesFromLocation(ctx, &places.FindPlacesFromLocationRequest{
		Location: places.Location{
			Latitude:  locationStart.Latitude,
			Longitude: locationStart.Longitude,
		},
		Radius:   2000,
		Language: "ja",
	})
	if err != nil {
		return nil, fmt.Errorf("error while fetching places: %v\n", err)
	}

	var categoriesPreferred []models.LocationCategory
	if categoryNamesPreferred != nil {
		for _, categoryName := range *categoryNamesPreferred {
			if category := models.GetCategoryOfName(categoryName); category != nil {
				categoriesPreferred = append(categoriesPreferred, *category)
			}
		}
	}

	var categoriesToFiler []models.LocationCategory
	if len(*categoryNamesPreferred) > 0 {
		categoriesToFiler = categoriesPreferred
	} else {
		categoriesToFiler = models.GetCategoryToFilter()
	}

	placesFilter := placefilter.NewPlacesFilter(placesSearched)
	placesFilter = placesFilter.FilterByCategory(categoriesToFiler)

	// TODO: 現在時刻でフィルタリングするかを指定できるようにする
	placesFilter = placesFilter.FilterByOpeningNow()

	// TODO: 移動距離ではなく、移動時間でやる
	var placesRecommend []places.Place

	placesInNear := placesFilter.FilterWithinDistanceRange(locationStart, 0, 500).Places()
	placesInMiddle := placesFilter.FilterWithinDistanceRange(locationStart, 500, 1000).Places()
	placesInFar := placesFilter.FilterWithinDistanceRange(locationStart, 1000, 2000).Places()
	if len(placesInNear) > 0 {
		// TODO: 0 ~ 500mで最もレビューの高い場所を選ぶ
		placesRecommend = append(placesRecommend, placesInNear[0])
	}
	if len(placesInMiddle) > 0 {
		// TODO: 500 ~ 1000mで最もレビューの高い場所を選ぶ
		placesRecommend = append(placesRecommend, placesInMiddle[0])
	}
	if len(placesInFar) > 0 {
		// TODO: 1000 ~ 2000mで最もレビューの高い場所を選ぶ
		placesRecommend = append(placesRecommend, placesInFar[0])
	}

	plans := make([]models.Plan, 0) // MEMO: 空配列の時のjsonのレスポンスがnullにならないように宣言

	for _, placeRecommend := range placesRecommend {
		plan, err := s.createPlanByLocation(
			ctx,
			locationStart,
			placeRecommend,
			placesFilter.Places(),
			freeTime,
			createBasedOnCurrentLocation,
		)
		if err != nil {
			log.Println(err)
			continue
		}
		plans = append(plans, *plan)
	}

	return &plans, nil
}

func (s PlanService) createPlanByLocation(
	ctx context.Context,
	locationStart models.GeoLocation,
	placeStart places.Place,
	places []places.Place,
	freeTime *int,
	createBasedOnCurrentLocation bool,
) (*models.Plan, error) {
	placesFilter := placefilter.NewPlacesFilter(places)

	// 起点となる場所との距離順でソート
	placesSortedByDistance := placesFilter.Places()
	sort.SliceStable(placesSortedByDistance, func(i, j int) bool {
		locationRecommend := placeStart.Location.ToGeoLocation()
		distanceI := locationRecommend.DistanceInMeter(placesSortedByDistance[i].Location.ToGeoLocation())
		distanceJ := locationRecommend.DistanceInMeter(placesSortedByDistance[j].Location.ToGeoLocation())
		return distanceI < distanceJ
	})

	placesWithInRange := placefilter.NewPlacesFilter(placesSortedByDistance).FilterWithinDistanceRange(
		placeStart.Location.ToGeoLocation(),
		0,
		500,
	).Places()

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
			GooglePlaceId:         &place.PlaceID,
			Location:              place.Location.ToGeoLocation(),
			EstimatedStayDuration: categoryMain.EstimatedStayDuration,
			Category:              categoryMain.Name,
		})
		timeInPlan += timeInPlace
		categoriesInPlan = append(categoriesInPlan, categoryMain.Name)
		previousLocation = place.Location.ToGeoLocation()
		transitions = s.addTransition(placesInPlan, transitions, tripTime, createBasedOnCurrentLocation)
	}

	if len(placesInPlan) == 0 {
		return nil, fmt.Errorf("could not contain any places in plan")
	}

	// 場所の画像を取得
	performanceTimer := time.Now()
	for i, place := range placesInPlan {
		if place.GooglePlaceId == nil {
			continue
		}

		thumbnail, photos, err := s.fetchPlacePhotos(ctx, *place.GooglePlaceId)
		if err != nil {
			log.Printf("error while fetching place photos: %v\n", err)
			continue
		}

		if thumbnail != nil {
			placesInPlan[i].Thumbnail = thumbnail
		}

		if photos != nil {
			placesInPlan[i].Photos = photos
		}
	}
	log.Printf("fetching place photos took %v\n", time.Now().Sub(performanceTimer))

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
func (s PlanService) isOpeningWithIn(
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

func (s PlanService) travelTimeBetween(
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
