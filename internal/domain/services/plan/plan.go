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
	"poroto.app/poroto/planner/internal/domain/services/plangenerator"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
	"poroto.app/poroto/planner/internal/infrastructure/firestore"
)

type PlanService struct {
	placesApi                   places.PlacesApi
	planRepository              repository.PlanRepository
	planCandidateRepository     repository.PlanCandidateRepository
	placeSearchResultRepository repository.PlaceSearchResultRepository
	planGeneratorService        plangenerator.Service
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

	placeSearchResultRepository, err := firestore.NewPlaceSearchResultRepository(ctx)
	if err != nil {
		return nil, err
	}

	planGeneratorService, err := plangenerator.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while initializing plan generator service: %v", err)
	}

	return &PlanService{
		placesApi:                   *placesApi,
		planRepository:              planRepository,
		planCandidateRepository:     planCandidateRepository,
		placeSearchResultRepository: placeSearchResultRepository,
		planGeneratorService:        *planGeneratorService,
	}, err
}

func (s PlanService) CreatePlanByLocation(
	ctx context.Context,
	createPlanSessionId string,
	locationStart models.GeoLocation,
	// TODO: ユーザーに却下された場所を引数にする（プランを作成時により多くの場所を取得した場合、YESと答えたカテゴリの場所からしかプランを作成できなくなるため）
	categoryNamesPreferred *[]string,
	freeTime *int,
	createBasedOnCurrentLocation bool,
) (*[]models.Plan, error) {
	// 付近の場所を検索
	var placesSearched []places.Place

	//　キャッシュがあれば利用する
	placesCached, err := s.placeSearchResultRepository.Find(ctx, createPlanSessionId)
	if err != nil {
		log.Printf("error while fetching places from cache: %v\n", err)
	} else if placesCached != nil {
		log.Printf("use cached places[%v]\n", createPlanSessionId)
		placesSearched = placesCached
	}

	if placesSearched == nil {
		placesSearched, err = s.placesApi.FindPlacesFromLocation(ctx, &places.FindPlacesFromLocationRequest{
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

		if err := s.placeSearchResultRepository.Save(ctx, createPlanSessionId, placesSearched); err != nil {
			log.Printf("error while saving places to cache: %v\n", err)
		}
		log.Printf("save places to cache[%v]\n", createPlanSessionId)
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
	placesFilter = placesFilter.FilterIgnoreCategory()
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

	// 最もおすすめ度が高い３つの場所を基準にプランを作成する
	performanceTimer := time.Now()
	chPlans := make(chan *models.Plan, len(placesRecommend))
	for _, placeRecommend := range placesRecommend {
		go func(ctx context.Context, placeRecommend places.Place, chPlan chan<- *models.Plan) {
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
				return
			}
			chPlans <- plan
		}(ctx, placeRecommend, chPlans)
	}

	plans := make([]models.Plan, 0)
	for i := 0; i < len(placesRecommend); i++ {
		plan := <-chPlans
		if plan == nil {
			continue
		}

		plans = append(plans, *plan)
	}
	log.Printf("created plans[%v]\n", time.Since(performanceTimer))

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
			Name:                  place.Name,
			GooglePlaceId:         &place.PlaceID,
			Location:              place.Location.ToGeoLocation(),
			EstimatedStayDuration: categoryMain.EstimatedStayDuration,
			Category:              categoryMain.Name,
		})
		timeInPlan += timeInPlace
		categoriesInPlan = append(categoriesInPlan, categoryMain.Name)
		previousLocation = place.Location.ToGeoLocation()
		transitions = s.planGeneratorService.AddTransition(placesInPlan, transitions, tripTime, createBasedOnCurrentLocation)
	}

	if len(placesInPlan) == 0 {
		return nil, fmt.Errorf("could not contain any places in plan")
	}

	// 場所の画像を取得
	performanceTimer := time.Now()
	placesInPlan = s.planGeneratorService.FetchPlacesPhotos(ctx, placesInPlan)
	log.Printf("fetching place photos took %v\n", time.Since(performanceTimer))

	title, err := s.planGeneratorService.GeneratePlanTitle(placesInPlan)
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

func (s PlanService) CreatePlanFromPlace(
	ctx context.Context,
	createPlanSessionId string,
	placeId string,
) (*models.Plan, error) {
	planCandidate, err := s.planCandidateRepository.Find(ctx, createPlanSessionId)
	if err != nil {
		return nil, fmt.Errorf("error while fetching plan candidate")
	}

	// TODO: ユーザーの興味等を保存しておいて、それを反映させる
	placesSearched, err := s.placeSearchResultRepository.Find(ctx, createPlanSessionId)
	if err != nil {
		return nil, err
	}

	var placeStart *places.Place
	for _, place := range placesSearched {
		if place.PlaceID == placeId {
			placeStart = &place
			break
		}
	}

	if placeStart == nil {
		return nil, fmt.Errorf("place not found")
	}

	planCreated, err := s.createPlanByLocation(
		ctx,
		placeStart.Location.ToGeoLocation(),
		*placeStart,
		placesSearched,
		// TODO: freeTimeの項目を保存し、それを反映させる
		nil,
		planCandidate.CreatedBasedOnCurrentLocation,
	)
	if err != nil {
		return nil, err
	}

	if _, err = s.planCandidateRepository.AddPlan(ctx, createPlanSessionId, planCreated); err != nil {
		return nil, err
	}

	return planCreated, nil
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
