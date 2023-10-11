package plancandidate

import (
	"context"
	"fmt"
	"log"
	"poroto.app/poroto/planner/internal/domain/utils"
	"sort"

	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/placefilter"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

// TODO: PlanGeneratorServiceに持っていく
func (s Service) CategoriesNearLocation(
	ctx context.Context,
	location models.GeoLocation,
	createPlanSessionId string,
) ([]models.LocationCategory, error) {
	placesSearched, err := s.placeService.SearchNearbyPlaces(ctx, location)
	if err != nil {
		return nil, fmt.Errorf("error while fetching places: %v\n", err)
	}

	if err := s.placeSearchResultRepository.Save(ctx, createPlanSessionId, placesSearched); err != nil {
		return nil, fmt.Errorf("error while saving places to cache: %v\n", err)
	}

	placesFiltered := placesSearched
	placesFiltered = placefilter.FilterIgnoreCategory(placesFiltered)
	placesFiltered = placefilter.FilterByCategory(placesFiltered, models.GetCategoryToFilter(), true)
	placesFiltered = placefilter.FilterCompany(placesFiltered)

	// TODO: 現在時刻でフィルタリングするかを指定できるようにする
	placesFiltered = placefilter.FilterByOpeningNow(placesFiltered)

	// 場所をカテゴリごとにグループ化し、対応する場所の少ないカテゴリから順に写真を取得する
	placeCategoryGroups := groupPlacesByCategory(placesFiltered)
	sort.Slice(placeCategoryGroups, func(i, j int) bool {
		return len(placeCategoryGroups[i].places) < len(placeCategoryGroups[j].places)
	})

	// 検索された場所のカテゴリとその写真を取得
	categories := make([]models.LocationCategory, 0)
	placesUsedOfCategory := make([]models.GooglePlace, 0)
	for _, categoryPlaces := range placeCategoryGroups {
		category := models.GetCategoryOfName(categoryPlaces.category)
		if category == nil {
			continue
		}

		// すでに他のカテゴリで利用した場所は利用しない
		placesNotUsedInOtherCategory := placefilter.FilterPlaces(categoryPlaces.places, func(place models.GooglePlace) bool {
			return placefilter.FindById(placesUsedOfCategory, place.PlaceId) == nil
		})

		// カテゴリと関連の強い場所から順に写真を取得する
		placesSortedByCategoryIndex := placesNotUsedInOtherCategory
		sort.Slice(placesSortedByCategoryIndex, func(i, j int) bool {
			return indexOfCategory(placesSortedByCategoryIndex[i], *category) < indexOfCategory(placesSortedByCategoryIndex[j], *category)
		})

		//　カテゴリに属する場所のうち、写真が取得可能なものを取得
		for _, place := range placesSortedByCategoryIndex {
			placePhoto, err := s.placesApi.FetchPlacePhoto(place.PhotoReferences, places.ImageSizeLarge())
			if err != nil {
				log.Printf("error while fetching place photo: %v\n", err)
				continue
			}
			if placePhoto != nil {
				category.Photo = utils.StrCopyPointerValue(placePhoto)
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
	places   []models.GooglePlace
}

// groupPlacesByCategory は場所をカテゴリごとにグループ化する
// 同じ場所が複数のカテゴリに含まれることがある
func groupPlacesByCategory(placesToGroup []models.GooglePlace) []groupPlacesByCategoryResult {
	locationsGroupByCategory := make(map[string][]models.GooglePlace, 0)
	for _, location := range placesToGroup {
		for _, subCategory := range location.Types {
			category := models.CategoryOfSubCategory(subCategory)
			if category == nil {
				continue
			}

			if _, ok := locationsGroupByCategory[category.Name]; ok {
				locationsGroupByCategory[category.Name] = []models.GooglePlace{}
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

// indexOfCategory は models.GooglePlace.Types 中の`category`に対応するTypeのインデックスを返す
func indexOfCategory(place models.GooglePlace, category models.LocationCategory) int {
	for i, placeType := range place.Types {
		c := models.CategoryOfSubCategory(placeType)
		if c.Name == category.Name {
			return i
		}
	}
	return -1
}
