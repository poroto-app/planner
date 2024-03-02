package plangen

import (
	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/placefilter"
	"sort"
)

const (
	defaultMaxBasePlaceCount = 3
	defaultRadius            = 800
)

type SelectBasePlaceInput struct {
	BaseLocation           models.GeoLocation
	Places                 []models.Place
	CategoryNamesPreferred *[]string
	CategoryNamesDisliked  *[]string
	MaxBasePlaceCount      int
	Radius                 int
}

// SelectBasePlace は，プランの起点となる場所を選択する
func (s Service) SelectBasePlace(input SelectBasePlaceInput) []models.Place {
	if input.MaxBasePlaceCount == 0 {
		input.MaxBasePlaceCount = defaultMaxBasePlaceCount
	}

	if input.Radius == 0 {
		input.Radius = defaultRadius
	}

	if input.BaseLocation.IsZero() {
		panic("base location is zero value")
	}

	places := input.Places

	// レビューが低い、またはレビュー数が少ない場所を除外する
	places = placefilter.FilterByRating(places, 3.0, 10)
	s.logger.Debug("places after filtering by rating", zap.Int("places", len(places)))

	// ユーザーが拒否した場所は取り除く
	if input.CategoryNamesDisliked != nil {
		categoriesDisliked := models.GetCategoriesFromSubCategories(*input.CategoryNamesDisliked)
		places = placefilter.FilterByCategory(places, categoriesDisliked, false)
		s.logger.Debug("places after filtering by disliked categories", zap.Int("places", len(places)))
	}

	for filterDistance := 800; filterDistance < 2000; filterDistance += 200 {
		// 距離によって1件も場所が取得できない場合は、距離を広げて再度取得する
		placesFiltered := placefilter.FilterDefaultIgnore(placefilter.FilterDefaultIgnoreInput{
			Places:              places,
			StartLocation:       input.BaseLocation,
			IgnoreDistanceRange: float64(filterDistance),
		})

		if len(placesFiltered) >= 10 {
			places = placesFiltered
			break
		}
	}
	s.logger.Debug("places after filtering default ignore", zap.Int("places", len(places)))

	// カテゴリごとにレビューの高い場所から選択する
	placesSelected := selectByReview(places)
	if len(placesSelected) == input.MaxBasePlaceCount {
		return placesSelected
	}

	// 選択された場所から遠い場所を選択する
	placesSelected = selectByDistanceFromPlaces(places, placesSelected)

	if len(placesSelected) > input.MaxBasePlaceCount {
		return placesSelected[:input.MaxBasePlaceCount]
	}

	return placesSelected
}

// selectByReview は，レビューの高い順に場所を選択する
// categoriesPreferred が指定される場合は、同じカテゴリの場所が含まれないように選択する
func selectByReview(places []models.Place) []models.Place {
	// レビューの高い順にソート
	places = models.SortPlacesByRating(places)

	var placesSelected []models.Place
	for _, place := range places {
		// 既に選択済みの場所は除外
		_, isAlreadyAdded := array.Find(placesSelected, func(p models.Place) bool {
			return p.Google.PlaceId == place.Google.PlaceId
		})
		if isAlreadyAdded {
			continue
		}

		// 既に選択された場所と異なるカテゴリの場所が選択されるようにする
		isAlreadyHaveSameCategory := false
		for _, placeSelected := range placesSelected {
			if place.IsSameCategoryPlace(placeSelected) {
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
		if len(placesSelected) == defaultMaxBasePlaceCount {
			break
		}
	}

	return placesSelected
}

// selectByDistanceFromPlaces は，プラン間の内容が重複しないようにするため、既に選択された場所から遠い場所を選択する
func selectByDistanceFromPlaces(
	places []models.Place,
	placesSelected []models.Place,
) []models.Place {
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
		_, isAlreadyAdded := array.Find(placesSelected, func(p models.Place) bool {
			return p.Google.PlaceId == place.Google.PlaceId
		})
		if isAlreadyAdded {
			continue
		}

		placesSelected = append(placesSelected, place)
	}

	return placesSelected
}

// isNearFromPlaces placeBase　が placesCompare　のいずれかの場所から distance メートル以内にあるかどうかを判定する
func isNearFromPlaces(
	placeBase models.Place,
	placesCompare []models.Place,
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
