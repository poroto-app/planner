package plancandidate

import (
	"context"
	"fmt"
	"log"
	"sort"

	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/models/placefilter"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

func (s Service) CategoriesNearLocation(
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

	placesFilter := placefilter.NewPlacesFilter(placesSearched)
	placesFilter = placesFilter.FilterIgnoreCategory()
	placesFilter = placesFilter.FilterByCategory(models.GetCategoryToFilter())

	// TODO: 現在時刻でフィルタリングするかを指定できるようにする
	placesFilter = placesFilter.FilterByOpeningNow()

	// 場所をカテゴリごとにグループ化し、対応する場所の少ないカテゴリから順に写真を取得する
	placeCategoryGroups := groupPlacesByCategory(placesFilter.Places())
	sort.Slice(placeCategoryGroups, func(i, j int) bool {
		return len(placeCategoryGroups[i].places) < len(placeCategoryGroups[j].places)
	})

	// 検索された場所のカテゴリとその写真を取得
	categories := make([]models.LocationCategory, 0)
	placesUsedOfCategory := make([]places.Place, 0)
	for _, categoryPlaces := range placeCategoryGroups {
		category := models.GetCategoryOfName(categoryPlaces.category)
		if category == nil {
			continue
		}

		// すでに他のカテゴリで利用した場所は利用しない
		placesNotUsedInOtherCategory := placefilter.NewPlacesFilter(
			categoryPlaces.places,
		).FilterPlaces(func(place places.Place) bool {
			return placefilter.NewPlacesFilter(placesUsedOfCategory).FindById(place.PlaceID) == nil
		}).Places()

		// カテゴリと関連の強い場所から順に写真を取得する
		placesSortedByCategoryIndex := placesNotUsedInOtherCategory
		sort.Slice(placesSortedByCategoryIndex, func(i, j int) bool {
			return indexOfCategory(placesSortedByCategoryIndex[i], *category) < indexOfCategory(placesSortedByCategoryIndex[j], *category)
		})

		//　カテゴリに属する場所のうち、写真が取得可能なものを取得
		for _, place := range placesSortedByCategoryIndex {
			placePhoto, err := s.placesApi.FetchPlacePhoto(place, nil)
			if err != nil {
				log.Printf("error while fetching place photo: %v\n", err)
				continue
			}
			if placePhoto != nil {
				category.Photo = placePhoto.ImageUrl
				placesUsedOfCategory = append(placesUsedOfCategory, place)
				break
			}
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

// indexOfCategory は places.Place.Types 中の`category`に対応するTypeのインデックスを返す
func indexOfCategory(place places.Place, category models.LocationCategory) int {
	for i, placeType := range place.Types {
		c := models.CategoryOfSubCategory(placeType)
		if c.Name == category.Name {
			return i
		}
	}
	return -1
}
