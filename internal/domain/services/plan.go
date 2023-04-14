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

	// TODO: フィルタリングするカテゴリを指定できるようにする
	placesSearched = s.filterByCategory(placesSearched)

	// TODO: 移動距離ではなく、移動時間でやる
	var placesRecommend []places.Place
	placesInNear := FilterWithinDistanceRange(location, 0, 500, placesSearched)
	placesInMiddle := FilterWithinDistanceRange(location, 500, 1000, placesSearched)
	placesInFar := FilterWithinDistanceRange(location, 1000, 2000, placesSearched)
	if len(placesInNear) > 0 {
		placesRecommend = append(placesRecommend, placesInNear[0])
	}
	if len(placesInMiddle) > 0 {
		placesRecommend = append(placesRecommend, placesInMiddle[0])
	}
	if len(placesInFar) > 0 {
		placesRecommend = append(placesRecommend, placesInFar[0])
	}

	plans := []models.Plan{} // MEMO: 空配列の時のjsonのレスポンスがnullにならないように宣言
	for _, placeSearched := range placesRecommend {
		placePhotos, err := placesApi.FetchPlacePhotos(context.Background(), placeSearched)
		if err != nil {
			continue
		}
		photos := []string{}
		for _, photo := range placePhotos {
			photos = append(photos, photo.ImageUrl)
		}

		plans = append(plans, models.Plan{
			Name: placeSearched.Name,
			Places: []models.Place{
				{
					Name:   placeSearched.Name,
					Photos: photos,
					Location: models.GeoLocation{
						Latitude:  placeSearched.Location.Latitude,
						Longitude: placeSearched.Location.Longitude,
					},
				},
			},
		})
	}

	return &plans, nil
}

func (s PlanService) filterByCategory(
	placesToFilter []places.Place,
) []places.Place {
	categories := map[string][]string{}
	categories["amusements"] = []string{"amusement_park", "aquarium", "art_gallary", "museum"}
	categories["restaurants"] = []string{"bakery", "bar", "cafe", "food", "restaurant"}

	var categoriesSlice []string
	for _, value := range categories {
		categoriesSlice = append(categoriesSlice, value...)
	}

	var placesInCategory []places.Place
	for _, place := range placesToFilter {
		if array.HasIntersection(place.Types, categoriesSlice) {
			placesInCategory = append(placesInCategory, place)
		}
	}

	return placesInCategory
}

func FilterWithinDistanceRange(
	currentLocation models.GeoLocation,
	startInMeter float64,
	endInMeter float64,
	placesToFilter []places.Place,
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
