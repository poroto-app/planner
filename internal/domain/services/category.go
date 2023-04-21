package services

import (
	"context"
	"fmt"

	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

type LocationCategory struct {
	Name          string
	SubCategories []string
	Photo         places.PlacePhoto
}

// Place.Type がどの大カテゴリに所属するか
func CategoryOfType(placeType string) models.LocationCategory {
	for _, category := range models.AllCategory {
		if array.IsContain(category.SubCategories, placeType) {
			return category
		}
	}

	return models.LocationCategory{
		Name:          "Undefined",
		SubCategories: []string{placeType},
	}
}

// 付近のPlacesTypesを呼び出して大カテゴリに集約する関数
func (s PlanService) CategoriesNearLocation(
	ctx context.Context,
	req *places.FindPlacesFromLocationRequest,
) ([]LocationCategory, error) {
	var categories []LocationCategory
	var locationCategory models.LocationCategory

	var bookedCategories = []string{}

	placesTypesSearched, err := fetchNearPlacesTypes(ctx, req, &s.placesApi)
	if err != nil {
		return nil, err
	}

	for _, placeType := range placesTypesSearched {
		locationCategory = CategoryOfType(placeType.Name)
		if array.IsContain(bookedCategories, locationCategory.Name) {
			continue
		}

		categories = append(categories, LocationCategory{
			Name:          locationCategory.Name,
			SubCategories: locationCategory.SubCategories,
			Photo:         placeType.Photo,
		})

		bookedCategories = append(bookedCategories, locationCategory.Name)
	}
	return categories, nil
}

// 付近のPlacesTypesをとってくるだけの関数（s PlanServiceはそのうちとる）
func fetchNearPlacesTypes(
	ctx context.Context,
	req *places.FindPlacesFromLocationRequest,
	placesApi *places.PlacesApi,
) ([]LocationCategory, error) {
	var nearCategories = []string{}
	var nearLocationCategories = []LocationCategory{}

	placesSearched, err := placesApi.FindPlacesFromLocation(ctx, req)
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
				nearLocationCategories = append(nearLocationCategories, LocationCategory{
					Name:          category,
					SubCategories: []string{},
					Photo: places.PlacePhoto{
						ImageUrl: "https://example.com/sample/category.jpg",
					},
				})
				continue
			}

			nearLocationCategories = append(nearLocationCategories, LocationCategory{
				Name:          category,
				SubCategories: []string{},
				Photo:         photos[0],
			})
		}
	}
	return nearLocationCategories, nil
}
