package services

import (
	"context"
	"fmt"

	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

type PlanService struct {
	placesApi places.PlacesApi
}

func NewPlanService() (*PlanService, error) {
	placesApi, err := places.NewPlacesApi()
	if err != nil {
		return nil, fmt.Errorf("error while initizalizing places api: %v", err)
	}
	return &PlanService{
		placesApi: *placesApi,
	}, err
}

func (s PlanService) CreatePlanByLocation(
	ctx context.Context,
	location models.GeoLocation,
) (*[]models.Plan, error) {
	placesSearched, err := s.placesApi.FindPlacesFromLocation(ctx, &places.FindPlacesFromLocationRequest{
		Location: places.Location{
			Latitude:  location.Latitude,
			Longitude: location.Longitude,
		},
		Radius: 2000,
	})
	if err != nil {
		return nil, fmt.Errorf("error while fetching places: %v\n", err)
	}

	placesSearched = s.filterByCategory(placesSearched, []models.LocationCategory{
		models.CategoryAmusements,
		models.CategoryBook,
		models.CategoryCamp,
		models.CategoryCafe,
		models.CategoryCulture,
		models.CategoryNatural,
		models.CategoryPark,
		models.CategoryRestaurant,
		models.CategoryShopping,
	})

	// TODO: 現在時刻でフィルタリングするかを指定できるようにする
	placesSearched = s.filterByOpeningNow(placesSearched)

	// TODO: 移動距離ではなく、移動時間でやる
	var placesRecommend []places.Place
	placesInNear := s.filterWithinDistanceRange(placesSearched, location, 0, 500)
	placesInMiddle := s.filterWithinDistanceRange(placesSearched, location, 500, 1000)
	placesInFar := s.filterWithinDistanceRange(placesSearched, location, 1000, 2000)
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
	for _, place := range placesRecommend {
		placePhotos, err := s.placesApi.FetchPlacePhotos(context.Background(), place)
		if err != nil {
			continue
		}
		photos := make([]string, 0)
		for _, photo := range placePhotos {
			photos = append(photos, photo.ImageUrl)
		}

		plans = append(plans, models.Plan{
			Name: place.Name,
			Places: []models.Place{
				{
					Name:   place.Name,
					Photos: photos,
					Location: models.GeoLocation{
						Latitude:  place.Location.Latitude,
						Longitude: place.Location.Longitude,
					},
				},
			},
			TimeInMinutes: s.travelTimeFromCurrent(
				location,
				models.GeoLocation{
					Latitude:  place.Location.Latitude,
					Longitude: place.Location.Longitude,
				},
				80.0,
			),
		})
	}

	return &plans, nil
}

func (s PlanService) CategoriesNearLocation(
	ctx context.Context,
	location models.GeoLocation,
) ([]models.LocationCategory, error) {
	placesSearched, err := s.placesApi.FindPlacesFromLocation(ctx, &places.FindPlacesFromLocationRequest{
		Location: places.Location{
			Latitude:  location.Latitude,
			Longitude: location.Longitude,
		},
		Radius: 2000,
	})
	if err != nil {
		return nil, fmt.Errorf("error while fetching places: %v\n", err)
	}

	placesSearched = s.filterByCategory(placesSearched, []models.LocationCategory{
		models.CategoryAmusements,
		models.CategoryBook,
		models.CategoryCamp,
		models.CategoryCafe,
		models.CategoryCulture,
		models.CategoryNatural,
		models.CategoryPark,
		models.CategoryRestaurant,
		models.CategoryShopping,
	})

	// TODO: 現在時刻でフィルタリングするかを指定できるようにする
	placesSearched = s.filterByOpeningNow(placesSearched)

	// 検索された場所のカテゴリとその写真を取得
	categoryPhotos := make(map[string]string)
	for _, place := range placesSearched {
		photos, err := s.placesApi.FetchPlacePhotos(ctx, place)
		if err != nil {
			continue
		}

		// 対応するLocationCategoryを取得（重複処理および写真保存のためmapを採用）
		for _, subCategory := range place.Types {
			category := models.CategoryOfSubCategory(subCategory)
			if category == nil {
				continue
			}

			if _, ok := categoryPhotos[category.Name]; ok {
				continue
			}

			// 場所の写真を取得（取得できなかった場合はデフォルトの画像を利用）
			categoryPhotos[category.Name] = category.Photo
			if len(photos) > 0 {
				categoryPhotos[category.Name] = photos[0].ImageUrl
			}
		}
	}

	categories := make([]models.LocationCategory, 0)
	for key, value := range categoryPhotos {
		categories = append(categories, models.GetCategoryOfName(key))
		categories[len(categories)-1].Photo = value
	}

	return categories, nil
}

func (s PlanService) filterByCategory(
	placesToFilter []places.Place,
	categories []models.LocationCategory,
) []places.Place {
	var categoriesSlice []string
	for _, category := range categories {
		categoriesSlice = append(categoriesSlice, category.SubCategories...)
	}

	var placesInCategory []places.Place
	for _, place := range placesToFilter {
		if array.HasIntersection(place.Types, categoriesSlice) {
			placesInCategory = append(placesInCategory, place)
		}
	}

	return placesInCategory
}

func (s PlanService) filterByOpeningNow(
	placesToFilter []places.Place,
) []places.Place {
	var placesOpeningNow []places.Place
	for _, place := range placesToFilter {
		if place.OpenNow {
			placesOpeningNow = append(placesOpeningNow, place)
		}
	}
	return placesOpeningNow
}

func (s PlanService) filterWithinDistanceRange(
	placesToFilter []places.Place,
	currentLocation models.GeoLocation,
	startInMeter float64,
	endInMeter float64,
) []places.Place {
	var placesWithInDistance []places.Place
	for _, place := range placesToFilter {
		distance := currentLocation.DistanceInMeter(models.GeoLocation{
			Latitude:  place.Location.Latitude,
			Longitude: place.Location.Longitude,
		})
		if startInMeter <= distance && distance < endInMeter {
			placesWithInDistance = append(placesWithInDistance, place)
		}
	}
	return placesWithInDistance
}

func (s PlanService) travelTimeFromCurrent(
	currentLocation models.GeoLocation,
	targetLocation models.GeoLocation,
	meterPerMinutes float64,
) float64 {
	timeInMinutes := 0.0
	distance := currentLocation.DistanceInMeter(targetLocation)
	if distance > 0.0 && meterPerMinutes > 0.0 {
		timeInMinutes = distance / meterPerMinutes
	}
	return timeInMinutes
}
