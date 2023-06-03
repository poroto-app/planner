package plan

import (
	"context"
	"fmt"
	"log"

	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

func (s PlanService) CategoriesNearLocation(
	ctx context.Context,
	location models.GeoLocation,
) ([]models.LocationCategory, error) {
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

	placesSearched = s.filterByCategory(placesSearched, models.GetCategoryToFilter())

	// TODO: 現在時刻でフィルタリングするかを指定できるようにする
	placesSearched = s.filterByOpeningNow(placesSearched)

	// 検索された場所のカテゴリとその写真を取得
	categories := make([]models.LocationCategory, 0)
	for _, categoryPlaces := range groupPlacesByCategory(placesSearched) {
		category := models.GetCategoryOfName(categoryPlaces.category)
		if category == nil {
			continue
		}

		var placePhoto *places.PlacePhoto
		for _, place := range categoryPlaces.places {
			placePhoto, err = s.placesApi.FetchPlacePhoto(place, nil)
			if err != nil {
				log.Printf("error while fetching place photo: %v\n", err)
				continue
			}
			if placePhoto != nil {
				break
			}
		}

		if placePhoto != nil {
			category.Photo = placePhoto.ImageUrl
		}
		categories = append(categories, *category)
	}

	return categories, nil
}

type groupPlacesByCategoryResult struct {
	category string
	places   []places.Place
}

// groupPlacesByCategory は場所をカテゴリごとにグループ化する
// 同じ場所が複数のカテゴリに含まれることがある
func groupPlacesByCategory(placesToGroup []places.Place) []groupPlacesByCategoryResult {
	locationsGroupByCategory := make(map[string][]places.Place, 0)
	for _, location := range placesToGroup {
		for _, subCategory := range location.Types {
			category := models.CategoryOfSubCategory(subCategory)
			if category == nil {
				continue
			}

			if _, ok := locationsGroupByCategory[category.Name]; ok {
				locationsGroupByCategory[category.Name] = []places.Place{}
			}

			locationsGroupByCategory[category.Name] = append(locationsGroupByCategory[category.Name], location)
		}
	}

	var result []groupPlacesByCategoryResult
	for category, placesOfCategory := range locationsGroupByCategory {
		result = append(result, groupPlacesByCategoryResult{
			category: category,
			places:   placesOfCategory,
		})
	}

	return result
}
