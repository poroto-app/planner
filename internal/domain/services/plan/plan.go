package plan

import (
	"context"
	"fmt"
	"log"
	"sort"

	"github.com/google/uuid"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/repository"
	"poroto.app/poroto/planner/internal/domain/services/placefilter"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
	"poroto.app/poroto/planner/internal/infrastructure/firestore"
)

type PlanService struct {
	placesApi               places.PlacesApi
	planCandidateRepository repository.PlanCandidateRepository
}

func NewPlanService(ctx context.Context) (*PlanService, error) {
	placesApi, err := places.NewPlacesApi()
	if err != nil {
		return nil, fmt.Errorf("error while initizalizing places api: %v", err)
	}

	planCandidateRepository, err := firestore.NewPlanCandidateRepository(ctx)
	if err != nil {
		return nil, err
	}

	return &PlanService{
		placesApi:               *placesApi,
		planCandidateRepository: planCandidateRepository,
	}, err
}

func (s PlanService) CreatePlanByLocation(
	ctx context.Context,
	location models.GeoLocation,
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

	placesFilter := placefilter.NewPlacesFilter(placesSearched)

	placesFilter = placesFilter.FilterByCategory(models.GetCategoryToFilter())

	// TODO: 現在時刻でフィルタリングするかを指定できるようにする
	placesFilter = placesFilter.FilterByOpeningNow()

	// TODO: 移動距離ではなく、移動時間でやる
	var placesRecommend []places.Place

	placesInNear := placesFilter.FilterWithinDistanceRange(location, 0, 500).Places()
	placesInMiddle := placesFilter.FilterWithinDistanceRange(location, 500, 1000).Places()
	placesInFar := placesFilter.FilterWithinDistanceRange(location, 1000, 2000).Places()
	if len(placesInNear) > 0 {
		placesRecommend = append(placesRecommend, placesInNear[0])
	}
	if len(placesInMiddle) > 0 {
		placesRecommend = append(placesRecommend, placesInMiddle[0])
	}
	if len(placesInFar) > 0 {
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

			tripTime := s.travelTimeBetween(
				previousLocation,
				place.Location.ToGeoLocation(),
				80.0,
			)
			timeInPlace := category.EstimatedStayDuration + tripTime
			if freeTime != nil && timeInPlan+timeInPlace > uint(*freeTime) {
				break
			}
			placesInPlan = append(placesInPlan, models.Place{
				Name:                  place.Name,
				Photos:                photos,
				Thumbnail:             thumbnail,
				Location:              place.Location.ToGeoLocation(),
				EstimatedStayDuration: category.EstimatedStayDuration,
			})
			timeInPlan += timeInPlace
			categoriesInPlan = append(categoriesInPlan, place.Types[0])
			previousLocation = place.Location.ToGeoLocation()
		}

		if len(placesInPlan) == 0 {
			continue
		}
		plans = append(plans, models.Plan{
			Id:            uuid.New().String(),
			Name:          placeRecommend.Name,
			Places:        placesInPlan,
			TimeInMinutes: timeInPlan,
		})
	}

	return &plans, nil
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
