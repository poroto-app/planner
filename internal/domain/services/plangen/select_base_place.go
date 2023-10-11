package plangen

import (
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/placefilter"
	"sort"
)

const (
	maxBasePlaceCount = 3
)

// selectBasePlace は，プランの起点となる場所を選択する
func (s Service) selectBasePlace(
	places []models.GooglePlace,
	categoryNamesPreferred *[]string,
	categoryNamesDisliked *[]string,
	shouldOpenNow bool,
) []models.GooglePlace {
	// ユーザーが拒否した場所は取り除く
	if categoryNamesDisliked != nil {
		categoriesDisliked := models.GetCategoriesFromSubCategories(*categoryNamesDisliked)
		places = placefilter.FilterByCategory(places, categoriesDisliked, false)
	}

	if shouldOpenNow {
		places = placefilter.FilterByOpeningNow(places)
	}

	// カテゴリごとにレビューの高い場所から選択する
	placesSelected := selectByReview(places)
	if len(placesSelected) == maxBasePlaceCount {
		return placesSelected
	}

	// 選択された場所から遠い場所を選択する
	placesSelected = selectByDistanceFromPlaces(places, placesSelected)

	if len(placesSelected) > maxBasePlaceCount {
		return placesSelected[:maxBasePlaceCount]
	}

	return placesSelected
}

// selectByReview は，レビューの高い順に場所を選択する
// categoriesPreferred が指定される場合は、同じカテゴリの場所が含まれないように選択する
func selectByReview(places []models.GooglePlace) []models.GooglePlace {
	// レビューの高い順にソート
	sort.SliceStable(places, func(i, j int) bool {
		return places[i].Rating > places[j].Rating
	})

	var placesSelected []models.GooglePlace
	for _, place := range places {
		// 既に選択済みの場所は除外
		if isAlreadyAdded(place, placesSelected) {
			continue
		}

		// 既に選択された場所と異なるカテゴリの場所が選択されるようにする
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

		// 既に選択された場所から500m以内の場所は選択しない(プランの内容が重複する可能性が高いため)
		if isNearFromPlaces(place, placesSelected, 500) {
			continue
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
	places []models.GooglePlace,
	placesSelected []models.GooglePlace,
) []models.GooglePlace {
	// 既に選択された場所から遠い順にソート
	sort.SliceStable(places, func(i, j int) bool {
		sumDistanceI := 0.0
		for _, placeSelected := range placesSelected {
			sumDistanceI += placeSelected.Location.DistanceInMeter(places[i].Location)
		}

		sumDistanceJ := 0.0
		for _, placeSelected := range placesSelected {
			sumDistanceJ += placeSelected.Location.DistanceInMeter(places[j].Location)
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

func isAlreadyAdded(place models.GooglePlace, places []models.GooglePlace) bool {
	for _, p := range places {
		if p.PlaceId == place.PlaceId {
			return true
		}
	}
	return false
}

func isSameCategoryPlace(a, b models.GooglePlace) bool {
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

func categoriesOfPlace(place models.GooglePlace) []models.LocationCategory {
	var categories []models.LocationCategory
	for _, placeType := range place.Types {
		category := models.CategoryOfSubCategory(placeType)
		if category != nil {
			categories = append(categories, *category)
		}
	}
	return categories
}

// isNearFromPlaces placeBase　が placesCompare　のいずれかの場所から distance メートル以内にあるかどうかを判定する
func isNearFromPlaces(
	placeBase models.GooglePlace,
	placesCompare []models.GooglePlace,
	distance int,
) bool {
	for _, placeCompare := range placesCompare {
		locationOfPlaceBase := placeBase.Location
		locationOfPlaceCompare := placeCompare.Location
		distanceFromSelectedPlace := locationOfPlaceCompare.DistanceInMeter(locationOfPlaceBase)
		if int(distanceFromSelectedPlace) < distance {
			return true
		}
	}
	return false
}
