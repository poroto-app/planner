package plangen

import (
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/placefilter"
	api "poroto.app/poroto/planner/internal/infrastructure/api/google/places"
	"sort"
)

const (
	maxBasePlaceCount = 3
)

// selectBasePlace は，プランの起点となる場所を選択する
// baseLocation: プランを作成する基準となる場所(現在地など)
func (s Service) selectBasePlace(
	places []api.Place,
	categoryNamesPreferred *[]string,
	categoryNamesDisliked *[]string,
	shouldOpenNow bool,
) []api.Place {
	// ユーザーが拒否した場所は取り除く
	if categoryNamesDisliked != nil {
		var categoriesDisliked []models.LocationCategory
		for _, categoryName := range *categoryNamesDisliked {
			category := models.GetCategoryOfName(categoryName)
			if category != nil {
				categoriesDisliked = append(categoriesDisliked, *category)
			}
		}
		places = placefilter.FilterByCategory(places, categoriesDisliked, false)
	}

	if shouldOpenNow {
		places = placefilter.FilterByOpeningNow(places)
	}

	var categoriesPreferred []models.LocationCategory
	if categoryNamesPreferred != nil {
		for _, categoryName := range *categoryNamesPreferred {
			category := models.GetCategoryOfName(categoryName)
			if category != nil {
				categoriesPreferred = append(categoriesPreferred, *category)
			}
		}
	}

	// カテゴリごとにレビューの高い場所から選択する
	placesSelected := selectByReview(places, categoriesPreferred)
	if len(placesSelected) == maxBasePlaceCount {
		return placesSelected
	}

	// 選択された場所から遠い場所を選択する
	placesSelected = selectByDistanceFromPlaces(places, placesSelected)

	return placesSelected[:maxBasePlaceCount]
}

// selectByReview は，レビューの高い順に場所を選択する
// categoriesPreferred が指定される場合は、同じカテゴリの場所が含まれないように選択する
func selectByReview(
	places []api.Place,
	categoriesPreferred []models.LocationCategory,
) []api.Place {
	// レビューの高い順にソート
	sort.SliceStable(places, func(i, j int) bool {
		return places[i].Rating > places[j].Rating
	})

	var placesSelected []api.Place
	for _, place := range places {
		// 既に選択済みの場所は除外
		if isAlreadyAdded(place, placesSelected) {
			continue
		}

		// 既に選択された場所と異なるカテゴリの場所が選択されるようにする
		if len(categoriesPreferred) > 0 {
			isAlreadyHaveSameCategory := false
			for _, placeSelected := range placesSelected {
				if isSameCategoryPlace(place, placeSelected) {
					isAlreadyHaveSameCategory = true
					break
				}
			}
			if isAlreadyHaveSameCategory {
				continue
			}
		}

		placesSelected = append(placesSelected, place)
		if len(placesSelected) == maxBasePlaceCount {
			break
		}
	}

	return placesSelected
}

// selectByDistanceFromPlaces は，プラン間の内容が重複しないようにするため、既に選択された場所から遠い場所を選択する
func selectByDistanceFromPlaces(
	places []api.Place,
	placesSelected []api.Place,
) []api.Place {
	// 既に選択された場所から遠い順にソート
	sort.SliceStable(places, func(i, j int) bool {
		sumDistanceI := 0.0
		for _, placeSelected := range placesSelected {
			sumDistanceI += placeSelected.Location.ToGeoLocation().DistanceInMeter(places[i].Location.ToGeoLocation())
		}

		sumDistanceJ := 0.0
		for _, placeSelected := range placesSelected {
			sumDistanceJ += placeSelected.Location.ToGeoLocation().DistanceInMeter(places[j].Location.ToGeoLocation())
		}

		return sumDistanceI > sumDistanceJ
	})

	for _, place := range places {
		// 既に選択済みの場所は除外
		if isAlreadyAdded(place, placesSelected) {
			continue
		}

		placesSelected = append(placesSelected, place)
	}

	return placesSelected
}

func isAlreadyAdded(place api.Place, places []api.Place) bool {
	for _, p := range places {
		if p.PlaceID == place.PlaceID {
			return true
		}
	}
	return false
}

func isSameCategoryPlace(a, b api.Place) bool {
	categoriesOfA := categoriesOfPlace(a)
	categoriesOfB := categoriesOfPlace(b)
	for _, categoryOfA := range categoriesOfA {
		for _, categoryOfB := range categoriesOfB {
			if categoryOfA.Name == categoryOfB.Name {
				return true
			}
		}
	}
	return false
}

func categoriesOfPlace(place api.Place) []models.LocationCategory {
	var categories []models.LocationCategory
	for _, placeType := range place.Types {
		category := models.CategoryOfSubCategory(placeType)
		if category != nil {
			categories = append(categories, *category)
		}
	}
	return categories
}
