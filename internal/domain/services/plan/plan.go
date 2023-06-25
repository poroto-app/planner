package plan

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/google/uuid"
	"googlemaps.github.io/maps"
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
	preferenceCategoryNames *[]string,
	freeTime *int,
) (*[]models.Plan, error) {
	//// カテゴリでフィルタ
	//var categoriesPreferred []models.LocationCategory
	//if preferenceCategoryNames != nil {
	//	for _, categoryName := range *preferenceCategoryNames {
	//		if category := models.GetCategoryOfName(categoryName); category != nil {
	//			categoriesPreferred = append(categoriesPreferred, *category)
	//		}
	//	}
	//}
	//
	//var categoriesToFilter []models.LocationCategory
	//if len(*preferenceCategoryNames) > 0 {
	//	categoriesToFilter = categoriesPreferred
	//} else {
	//	categoriesToFilter = models.GetCategoryToFilter()
	//}

	placesSearched, err := s.fetchPlacesFromLocation(ctx, locationStart, []maps.PlaceType{
		// Amusements
		maps.PlaceTypeAmusementPark,
		maps.PlaceTypeSpa,
		//	Book
		maps.PlaceTypeBookStore,
		maps.PlaceTypeLibrary,
		//	Cafe
		maps.PlaceTypeCafe,
		// Culture
		maps.PlaceTypeArtGallery,
		maps.PlaceTypeMuseum,
		// Natural
		maps.PlaceTypeAquarium,
		maps.PlaceTypeZoo,
		// Shopping
		maps.PlaceTypeStore,
	})
	if err != nil {
		return nil, fmt.Errorf("error while fetching places: %v", err)
	}

	placesFilter := placefilter.NewPlacesFilter(placesSearched)
	// 営業中の場所のみフィルタ
	// TODO: 現在時刻でフィルタリングするかを指定できるようにする
	placesFilter = placesFilter.FilterByOpeningNow()

	// カテゴリでフィルタ
	var categoriesPreferred []models.LocationCategory
	if preferenceCategoryNames != nil {
		for _, categoryName := range *preferenceCategoryNames {
			if category := models.GetCategoryOfName(categoryName); category != nil {
				categoriesPreferred = append(categoriesPreferred, *category)
			}
		}
	}

	var categoriesToFilter []models.LocationCategory
	if len(*preferenceCategoryNames) > 0 {
		categoriesToFilter = categoriesPreferred
	} else {
		categoriesToFilter = models.GetCategoryToFilter()
	}

	placesFilter = placesFilter.FilterByCategory(categoriesToFilter)

	// 起点となる場所を決める
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
		plan, err := s.createPlanFromLocation(ctx, locationStart, placeRecommend, placesFilter.Places(), freeTime)
		if err != nil {
			log.Println(err)
			continue
		}

		plans = append(plans, *plan)
	}

	return &plans, nil
}

func (s PlanService) createPlanFromLocation(
	ctx context.Context,
	locationStart models.GeoLocation,
	placeStart places.Place,
	places []places.Place,
	freeTime *int,
) (*models.Plan, error) {
	// TODO: おすすめする飲食店が決まったら、飲食店以外の場所を取得する

	placesFilter := placefilter.NewPlacesFilter(places)

	// 起点となる場所との距離順でソート
	placesSortedByDistance := placesFilter.Places()
	sort.SliceStable(placesSortedByDistance, func(i, j int) bool {
		locationRecommend := placeStart.Location.ToGeoLocation()
		distanceI := locationRecommend.DistanceInMeter(placesSortedByDistance[i].Location.ToGeoLocation())
		distanceJ := locationRecommend.DistanceInMeter(placesSortedByDistance[j].Location.ToGeoLocation())
		return distanceI < distanceJ
	})

	//　起点となる場所から1000m以内の場所を抽出
	//　MEMO: 広めの場所を取得するが、時間の上限があるため、そこまで多くならないはず
	placesWithInRange := placefilter.NewPlacesFilter(placesSortedByDistance).FilterWithinDistanceRange(
		placeStart.Location.ToGeoLocation(),
		0,
		1000,
	).Places()

	log.Printf("generate plan from %d places\n", len(placesWithInRange))

	placesInPlan := make([]models.Place, 0)
	categoriesInPlan := make([]string, 0)
	previousLocation := locationStart
	var timeInPlan uint = 0

	for _, place := range placesWithInRange {
		// 既にプランに含まれるカテゴリの場所は無視する
		if len(place.Types) == 0 {
			log.Printf("place %s has no category\n", place.Name)
			continue
		}

		// MEMO: カテゴリが不明な場合，滞在時間が取得できない
		var categoryMain *models.LocationCategory
		for _, placeType := range place.Types {
			category := models.CategoryOfSubCategory(placeType)
			if category != nil {
				categoryMain = category
				break
			}
		}
		if categoryMain == nil {
			log.Printf("place %s has no category\n", place.Name)
			continue
		}

		var categoriesOfPlace []string
		for _, placeType := range place.Types {
			category := models.CategoryOfSubCategory(placeType)
			if category != nil {
				categoriesOfPlace = append(categoriesOfPlace, category.Name)
			}
		}

		// 飲食店系は複数含めない
		categoriesFood := []string{
			models.CategoryRestaurant.Name,
			models.CategoryMealTakeaway.Name,
		}
		isFoodPlace := array.HasIntersection(categoriesOfPlace, categoriesFood)
		if isFoodPlace && array.HasIntersection(categoriesInPlan, categoriesFood) {
			log.Printf("skip place %s because plan is already has food place\n", place.Name)
			continue
		}

		// カフェを複数含めない
		isCafePlace := array.IsContain(categoriesOfPlace, models.CategoryCafe.Name)
		if isCafePlace && array.IsContain(categoriesInPlan, models.CategoryCafe.Name) {
			log.Printf("skip place %s because plan is already has cafe place\n", place.Name)
			continue
		}

		timeInPlace := categoryMain.EstimatedStayDuration

		// 開始地点から最初の場所までの移動時間はプランの時間に含めない
		// TODO: Planの変数として、移動時間を持たせる
		if !place.Location.ToGeoLocation().Equal(locationStart) {
			timeInPlace += s.travelTimeBetween(
				previousLocation,
				place.Location.ToGeoLocation(),
				80.0,
			)
		}

		// 5時間を超えたら強制的に終了
		if timeInPlan+timeInPlace > uint((time.Hour * 5).Minutes()) {
			break
		}

		if freeTime != nil && timeInPlan+timeInPlace > uint(*freeTime) {
			break
		}

		//　指定した時間内に店が閉まっている場合は無視する
		if freeTime != nil {
			isOpening, err := s.isOpeningWithIn(
				ctx,
				place,
				time.Now(),
				time.Minute*time.Duration(*freeTime),
			)

			if err != nil {
				log.Printf("error while checking opening: %v\n", err)
				continue
			}

			if !isOpening {
				log.Printf("skip place %s because it will be closed\n", place.Name)
				continue
			}
		}

		thumbnail, photos, err := s.fetchPlacePhotos(ctx, place)
		if err != nil {
			log.Printf("error while fetching place photos: %v\n", err)
			continue
		}

		placesInPlan = append(placesInPlan, models.Place{
			Name:                  place.Name,
			Photos:                photos,
			Thumbnail:             thumbnail,
			Location:              place.Location.ToGeoLocation(),
			EstimatedStayDuration: categoryMain.EstimatedStayDuration,
			Category:              categoryMain.Name,
		})

		timeInPlan += timeInPlace
		categoriesInPlan = append(categoriesInPlan, categoriesOfPlace...)
		previousLocation = place.Location.ToGeoLocation()
	}

	if len(placesInPlan) == 0 {
		return nil, fmt.Errorf("could not find places in plan")
	}

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
	}, nil
}

func (s PlanService) fetchPlacesFromLocation(
	ctx context.Context,
	locationStart models.GeoLocation,
	typesToSearch []maps.PlaceType,
) ([]places.Place, error) {
	var placesSearched []places.Place
	if placesSearchedWithAllType, err := s.placesApi.FindPlacesFromLocation(ctx, &places.FindPlacesFromLocationRequest{
		Location: places.Location{
			Latitude:  locationStart.Latitude,
			Longitude: locationStart.Longitude,
		},
		Radius:   2000,
		Language: "ja",
	}); err != nil {
		return nil, fmt.Errorf("error while fetching placesSearchedWithAllType: %v\n", err)
	} else {
		placesSearched = append(placesSearched, placesSearchedWithAllType...)
	}

	var placesTypesInSearchResult []string
	for _, place := range placesSearched {
		placesTypesInSearchResult = append(placesTypesInSearchResult, place.Types...)
	}

	// Places APIで取得すると、飲食店が多く取得されるため取得しない
	for _, locationType := range typesToSearch {
		//　検索結果に含まれるカテゴリの場所は検索しない
		if array.IsContain(placesTypesInSearchResult, string(locationType)) {
			log.Printf("skip search places with category: %s because it is already searched\n", locationType)
			continue
		}

		log.Printf("search places with category: %s\n", locationType)
		if placesSearchedWithAllType, err := s.placesApi.FindPlacesFromLocation(ctx, &places.FindPlacesFromLocationRequest{
			Location: places.Location{
				Latitude:  locationStart.Latitude,
				Longitude: locationStart.Longitude,
			},
			Radius:   2000,
			Language: "ja",
			Type:     &locationType,
		}); err != nil {
			log.Printf("error while fetching places with category %s: %v\n", locationType, err)
			continue
		} else {
			placesSearched = append(placesSearched, placesSearchedWithAllType...)
		}
	}

	return placesSearched, nil
}

// isOpeningWithIn は，指定された場所が指定された時間内に開いているかを判定する
func (s PlanService) isOpeningWithIn(ctx context.Context, place places.Place, startTime time.Time, duration time.Duration) (bool, error) {
	// 時刻フィルタリング用変数
	endTime := startTime.Add(duration)
	today := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location())

	placeOpeningPeriods, err := s.placesApi.FetchPlaceOpeningPeriods(ctx, place)
	if err != nil {
		return false, fmt.Errorf("error while fetching place periods: %v\n", err)
	}

	for _, placeOpeningPeriod := range placeOpeningPeriods {
		// startTime で指定された曜日のみを確認
		if placeOpeningPeriod.DayOfWeek != startTime.Weekday().String() {
			continue
		}

		// TODO: 変換処理の共通化 & テストを実装
		openingPeriodHour, opHourErr := strconv.Atoi(placeOpeningPeriod.OpeningTime[:2])
		openingPeriodMinute, opMinuteErr := strconv.Atoi(placeOpeningPeriod.OpeningTime[2:])
		closingPeriodHour, clHourErr := strconv.Atoi(placeOpeningPeriod.ClosingTime[:2])
		closingPeriodMinute, clMinuteErr := strconv.Atoi(placeOpeningPeriod.ClosingTime[2:])
		if opHourErr != nil || opMinuteErr != nil || clHourErr != nil || clMinuteErr != nil {
			return false, fmt.Errorf("error while converting period [string->int]")
		}

		openingTime := today.Add(time.Hour*time.Duration(openingPeriodHour) + time.Minute*time.Duration(openingPeriodMinute))
		closingTime := today.Add(time.Hour*time.Duration(closingPeriodHour) + time.Minute*time.Duration(closingPeriodMinute))

		// 開店時刻 < 開始時刻 && 終了時刻 < 閉店時刻 の判断
		if startTime.After(openingTime) && endTime.Before(closingTime) {
			return true, nil
		}
	}

	// 開店時間が指定されていない場合は空いているとして扱う
	return true, nil
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
