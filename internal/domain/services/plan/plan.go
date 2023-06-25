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
	location models.GeoLocation,
	preferenceCategoryNames *[]string,
	freeTime *int,
) (*[]models.Plan, error) {
	placesSearched, err := s.placesApi.FindPlacesFromLocation(ctx, &places.FindPlacesFromLocationRequest{
		Location: places.Location{
			Latitude:  location.Latitude,
			Longitude: location.Longitude,
		},
		Radius:   2000,
		Language: "ja",
	})
	if err != nil {
		return nil, fmt.Errorf("error while fetching places: %v\n", err)
	}

	var preferenceCategories []models.LocationCategory
	if preferenceCategoryNames != nil {
		for _, categoryName := range *preferenceCategoryNames {
			if category := models.GetCategoryOfName(categoryName); category != nil {
				preferenceCategories = append(preferenceCategories, *category)
			}
		}
	}

	var categoryToFiler []models.LocationCategory
	if len(*preferenceCategoryNames) > 0 {
		categoryToFiler = preferenceCategories
	} else {
		categoryToFiler = models.GetCategoryToFilter()
	}

	placesFilter := placefilter.NewPlacesFilter(placesSearched)
	placesFilter = placesFilter.FilterByCategory(categoryToFiler)

	// TODO: 現在時刻でフィルタリングするかを指定できるようにする
	placesFilter = placesFilter.FilterByOpeningNow()

	// TODO: 移動距離ではなく、移動時間でやる
	var placesRecommend []places.Place

	placesInNear := placesFilter.FilterWithinDistanceRange(location, 0, 500).Places()
	placesInMiddle := placesFilter.FilterWithinDistanceRange(location, 500, 1000).Places()
	placesInFar := placesFilter.FilterWithinDistanceRange(location, 1000, 2000).Places()
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
		// 起点となる場所との距離順でソート
		placesSortedByDistance := placesFilter.Places()
		sort.SliceStable(placesSortedByDistance, func(i, j int) bool {
			locationRecommend := placeRecommend.Location.ToGeoLocation()
			distanceI := locationRecommend.DistanceInMeter(placesSortedByDistance[i].Location.ToGeoLocation())
			distanceJ := locationRecommend.DistanceInMeter(placesSortedByDistance[j].Location.ToGeoLocation())
			return distanceI < distanceJ
		})

		//　起点となる場所から500m以内の場所を抽出
		placesWithInRange := placefilter.NewPlacesFilter(placesSortedByDistance).FilterWithinDistanceRange(
			placeRecommend.Location.ToGeoLocation(),
			0,
			500,
		).Places()

		placesInPlan := make([]models.Place, 0)
		categoriesInPlan := make([]string, 0)
		previousLocation := location
		var timeInPlan uint = 0

		for _, place := range placesWithInRange {
			// 既にプランに含まれるカテゴリの場所は無視する
			if len(place.Types) == 0 {
				continue
			}

			category := models.CategoryOfSubCategory(place.Types[0])

			// MEMO: カテゴリが不明な場合，滞在時間が取得できない
			if category == nil || array.IsContain(categoriesInPlan, category.Name) {
				continue
			}

			tripTime := s.travelTimeBetween(
				previousLocation,
				place.Location.ToGeoLocation(),
				80.0,
			)
			timeInPlace := category.EstimatedStayDuration + tripTime
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
					continue
				}
			}

			thumbnailPhoto, err := s.placesApi.FetchPlacePhoto(place, &places.ImageSize{
				Width:  places.ImgThumbnailMaxWidth,
				Height: places.ImgThumbnailMaxHeight,
			})
			if err != nil {
				log.Printf("error while fetching place thumbnail: %v\n", err)
				continue
			}
			var thumbnail *string
			if thumbnailPhoto != nil {
				thumbnail = &thumbnailPhoto.ImageUrl
			}

			placePhotos, err := s.placesApi.FetchPlacePhotos(ctx, place)
			if err != nil {
				log.Printf("error while fetching place photos: %v\n", err)
				continue
			}
			photos := make([]string, 0)
			for _, photo := range placePhotos {
				photos = append(photos, photo.ImageUrl)
			}

			placesInPlan = append(placesInPlan, models.Place{
				Name:                  place.Name,
				Photos:                photos,
				Thumbnail:             thumbnail,
				Location:              place.Location.ToGeoLocation(),
				EstimatedStayDuration: category.EstimatedStayDuration,
				Category:              category.Name,
			})
			timeInPlan += timeInPlace
			categoriesInPlan = append(categoriesInPlan, category.Name)
			previousLocation = place.Location.ToGeoLocation()
		}

		if len(placesInPlan) == 0 {
			continue
		}

		title, err := s.GeneratePlanTitle(placesInPlan)
		if err != nil {
			log.Printf("error while generating plan title: %v\n", err)
			title = &placeRecommend.Name
		}

		plans = append(plans, models.Plan{
			Id:            uuid.New().String(),
			Name:          *title,
			Places:        placesInPlan,
			TimeInMinutes: timeInPlan,
		})
	}

	return &plans, nil
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

	return false, fmt.Errorf("could not find opening period")
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
