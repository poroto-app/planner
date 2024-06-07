package plancandidate

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/placefilter"
	"poroto.app/poroto/planner/internal/domain/services/placesearch"
	"sort"
)

// defaultMaxCategory は提示するカテゴリの種類の上限
// defaultMaxPlacesPerCategory は提示するカテゴリごとに例示する場所の上限
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

type groupPlacesByCategoryResult struct {
	category string
	places   []models.Place
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
	if err := s.CreatePlanCandidateSet(ctx, params.CreatePlanSessionId); err != nil {
		return nil, fmt.Errorf("error while creating plan candidate: %v\n", err)
	}

	// 付近の場所を検索
	placesNearby, err := s.placeSearchService.SearchNearbyPlaces(ctx, placesearch.SearchNearbyPlacesInput{
		Location:           params.Location,
		PlanCandidateSetId: &params.CreatePlanSessionId,
	})
	if err != nil {
		return nil, fmt.Errorf("error while fetching places: %v\n", err)
	}

	placesFiltered := placesNearby
	placesFiltered = placefilter.FilterDefaultIgnore(placefilter.FilterDefaultIgnoreInput{
		Places:        placesFiltered,
		StartLocation: params.Location,
	})

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

		placesInCategory := categoryPlaces.places
		if len(placesInCategory) == 0 {
			continue
		}

		// すでに他のカテゴリで利用した場所は利用しない
		placesInCategory = placefilter.FilterPlaces(categoryPlaces.places, func(place models.Place) bool {
			for _, categoryWithPlaces := range categoriesWithPlaces {
				for _, placeInOtherCategory := range categoryWithPlaces.Places {
					if placeInOtherCategory.Id == place.Id {
						return false
					}
				}
			}
			return true
		})

		// カテゴリと関連の強い場所の最初の1件の写真を取得する
		placesSortedByCategoryIndex := placesInCategory
		sort.Slice(placesSortedByCategoryIndex, func(i, j int) bool {
			return placesSortedByCategoryIndex[i].Google.IndexOfCategory(*category) < placesSortedByCategoryIndex[j].Google.IndexOfCategory(*category)
		})

		// カテゴリごとに提示される場所の数を制限する
		if len(placesSortedByCategoryIndex) > params.MaxPlacesPerCategory {
			placesSortedByCategoryIndex = placesSortedByCategoryIndex[:params.MaxPlacesPerCategory]
		}

		// カテゴリ内の場所をレビューの高い順にソート
		placesSortedByCategoryIndex = models.SortPlacesByRating(placesSortedByCategoryIndex)

		// 場所の写真を取得する
		placesWithPhotos := s.placeSearchService.FetchPlacesPhotosAndSave(ctx, placesSortedByCategoryIndex...)

		categoriesWithPlaces = append(categoriesWithPlaces, models.NewLocationCategoryWithPlaces(*category, placesWithPhotos))
	}

	return categoriesWithPlaces, nil
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
