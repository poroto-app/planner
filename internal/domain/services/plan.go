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
	placesApi, err := places.NewPlacesApi()
	if err != nil {
		return nil, err
	}

	placesSearched, err := placesApi.FindPlacesFromLocation(ctx, &places.FindPlacesFromLocationRequest{
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
		{
			Name: "amusements",
			SubCategories: []string{
				"amusement_park", "aquarium", "art_gallery", "museum",
			},
		},
		{
			Name: "restaurants",
			SubCategories: []string{
				"bakery", "bar", "cafe", "food", "restaurant",
			},
		},
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
		placePhotos, err := placesApi.FetchPlacePhotos(context.Background(), place)
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
