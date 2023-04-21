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

type SubCategories struct {
	name   string
	imgUrl string
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

// 付近のPlacesTypesを呼び出して大カテゴリに集約する関数
func (s PlanService) CategoriesNearLocation(
	ctx context.Context,
	location models.GeoLocation,
) ([]models.LocationCategory, error) {
	var categories []models.LocationCategory
	var locationCategory *models.LocationCategory

	var bookedCategories = []string{}

	subCategoriesSearched, err := fetchNearSubCategories(ctx, location, &s.placesApi)
	if err != nil {
		return nil, err
	}

	for _, placeType := range subCategoriesSearched {
		locationCategory = models.CategoryOfSubCategory(placeType.name)
		if locationCategory == nil {
			continue
		}
		if array.IsContain(bookedCategories, locationCategory.Name) {
			continue
		}

		categories = append(categories, models.LocationCategory{
			Name:          locationCategory.Name,
			SubCategories: locationCategory.SubCategories,
			Photo:         locationCategory.Photo,
		})

		bookedCategories = append(bookedCategories, locationCategory.Name)
	}
	return categories, nil
}

// 付近のSubCategoriesをとってくるだけの関数
func fetchNearSubCategories(
	ctx context.Context,
	location models.GeoLocation,
	placesApi *places.PlacesApi,
) ([]SubCategories, error) {
	var nearCategories = []string{}
	var nearLocationCategories = []SubCategories{}

	placesSearched, err := placesApi.FindPlacesFromLocation(ctx, &places.FindPlacesFromLocationRequest{
		Location: places.Location{
			Latitude:  location.Latitude,
			Longitude: location.Longitude,
		},
		Radius: 2000,
	})
	if err != nil {
		return nearLocationCategories, fmt.Errorf("error while fetching places: %v\n", err)
	}
	for _, place := range placesSearched {
		for _, category := range place.Types {
			if array.IsContain(nearCategories, category) {
				continue
			}

			photos, err := placesApi.FetchPlacePhotos(ctx, place)
			if err != nil {
				continue
			}
			nearCategories = append(nearCategories, category)

			if len(photos) == 0 {
				nearLocationCategories = append(nearLocationCategories, SubCategories{
					name: category,
					// TODO: implement me!
					imgUrl: "https://placehold.jp/0a0a0a/ffffff/300x500.png?text=SubCategory",
				})
				continue
			}

			nearLocationCategories = append(nearLocationCategories, SubCategories{
				name:   category,
				imgUrl: photos[0].ImageUrl,
			})
		}
	}
	return nearLocationCategories, nil
}
