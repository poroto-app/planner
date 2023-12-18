package plancandidate

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/place"
	"poroto.app/poroto/planner/internal/domain/services/placefilter"
	"sort"
)

const (
	defaultMaxCategory          = 3
	defaultMaxPlacesPerCategory = 1
)

type CategoryNearLocationParams struct {
	Location             models.GeoLocation
	CreatePlanSessionId  string
	MaxCategory          int
	MaxPlacesPerCategory int
}

// TODO: PlanGeneratorServiceに持っていく
func (s Service) CategoriesNearLocation(
	ctx context.Context,
	params CategoryNearLocationParams,
) ([]models.LocationCategoryWithPlaces, error) {
	if params.MaxCategory <= 0 {
		params.MaxCategory = defaultMaxCategory
	}

	if params.MaxPlacesPerCategory <= 0 {
		params.MaxPlacesPerCategory = defaultMaxPlacesPerCategory
	}

	// プラン候補を作成
	if err := s.CreatePlanCandidate(ctx, params.CreatePlanSessionId); err != nil {
		return nil, fmt.Errorf("error while creating plan candidate: %v\n", err)
	}

	// 付近の場所を検索
	placesSearched, err := s.placeService.SearchNearbyPlaces(ctx, place.SearchNearbyPlacesInput{Location: params.Location})
	if err != nil {
		return nil, fmt.Errorf("error while fetching places: %v\n", err)
	}

	// 検索された場所を保存
	places, err := s.placeService.SaveSearchedPlaces(ctx, params.CreatePlanSessionId, placesSearched)
	if err != nil {
		return nil, fmt.Errorf("error while saving searched places: %v\n", err)
	}

	placesFiltered := places
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
	categoriesWithPlaces := make([]models.LocationCategoryWithPlaces, 0)
	for _, categoryPlaces := range placeCategoryGroups {
		// 取得したカテゴリが上限に達したら終了
		if len(categoriesWithPlaces) >= params.MaxCategory {
			break
		}

		category := models.GetCategoryOfName(categoryPlaces.category)
		if category == nil {
			continue
		}

		// すでに他のカテゴリで利用した場所は利用しない
		placesNotUsedInOtherCategory := placefilter.FilterPlaces(categoryPlaces.places, func(place models.Place) bool {
			var placesAlreadyAdded []models.Place
			for _, categoryWithPlaces := range categoriesWithPlaces {
				placesAlreadyAdded = append(placesAlreadyAdded, categoryWithPlaces.Places...)
			}

			for _, placeAlreadyAdded := range placesAlreadyAdded {
				if placeAlreadyAdded.Id == place.Id {
					return false
				}
			}

			return true
		})

		// カテゴリと関連の強い場所の最初の1件の写真を取得する
		placesSortedByCategoryIndex := placesNotUsedInOtherCategory
		sort.Slice(placesSortedByCategoryIndex, func(i, j int) bool {
			return placesSortedByCategoryIndex[i].Google.IndexOfCategory(*category) < placesSortedByCategoryIndex[j].Google.IndexOfCategory(*category)
		})
		if len(placesSortedByCategoryIndex) > params.MaxPlacesPerCategory {
			placesSortedByCategoryIndex = placesSortedByCategoryIndex[:params.MaxPlacesPerCategory]
		}

		// 場所の写真を取得する
		placesWithPhotos := s.placeService.FetchPlacesPhotosAndSave(ctx, placesSortedByCategoryIndex...)

		categoriesWithPlaces = append(categoriesWithPlaces, models.NewLocationCategoryWithPlaces(*category, placesWithPhotos))
	}

	return categoriesWithPlaces, nil
}

type groupPlacesByCategoryResult struct {
	category string
	places   []models.Place
}

// groupPlacesByCategory は場所をカテゴリごとにグループ化する
// 同じ場所が複数のカテゴリに含まれることがある
func groupPlacesByCategory(placesToGroup []models.Place) []groupPlacesByCategoryResult {
	locationsGroupByCategory := make(map[string][]models.Place, 0)
	for _, place := range placesToGroup {
		for _, subCategory := range place.Google.Types {
			category := models.CategoryOfSubCategory(subCategory)
			if category == nil {
				continue
			}

			if _, ok := locationsGroupByCategory[category.Name]; ok {
				locationsGroupByCategory[category.Name] = []models.Place{}
			}

			locationsGroupByCategory[category.Name] = append(locationsGroupByCategory[category.Name], place)
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
