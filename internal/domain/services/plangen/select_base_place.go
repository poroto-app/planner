package plangen

import (
	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/placefilter"
)

const (
	defaultMaxBasePlaceCount = 3
	defaultRadius            = 2000
)

// SelectBasePlaceInput
// Places 選択候補となる場所
// IgnorePlaces は，選択されないようにする場所
type SelectBasePlaceInput struct {
	BaseLocation           models.GeoLocation
	Places                 []models.Place
	IgnorePlaces           []models.Place
	CategoryNamesPreferred *[]string
	CategoryNamesDisliked  *[]string
	MaxBasePlaceCount      int
	Radius                 int
}

// SelectBasePlace は，プランの起点となる場所の候補を選択する
func (s Service) SelectBasePlace(input SelectBasePlaceInput) []models.Place {
	s.logger.Debug(
		"SelectBasePlace",
		zap.Int("Places", len(input.Places)),
		zap.Int("IgnorePlaces", len(input.IgnorePlaces)),
		zap.Int("MaxBasePlaceCount", input.MaxBasePlaceCount),
		zap.Int("Radius", input.Radius),
	)
	if input.MaxBasePlaceCount == 0 {
		input.MaxBasePlaceCount = defaultMaxBasePlaceCount
	}

	if input.Radius == 0 {
		input.Radius = defaultRadius
	}

	if input.BaseLocation.IsZero() {
		panic("base location is zero value")
	}

	placesFiltered := input.Places

	// レビューが低い、またはレビュー数が少ない場所を除外する
	placesFiltered = placefilter.FilterByRating(placesFiltered, 3.0, 10)
	s.logger.Debug("Places after filtering by rating", zap.Int("Places", len(placesFiltered)))

	// ユーザーが拒否した場所は取り除く
	if input.CategoryNamesDisliked != nil {
		categoriesDisliked := models.GetCategoriesFromSubCategories(*input.CategoryNamesDisliked)
		placesFiltered = placefilter.FilterByCategory(placesFiltered, categoriesDisliked, false)
		s.logger.Debug("Places after filtering by disliked categories", zap.Int("Places", len(placesFiltered)))
	}

	// すでに選択された場所は除外
	placesFiltered = array.Filter(placesFiltered, func(place models.Place) bool {
		_, isIgnorePlace := array.Find(input.IgnorePlaces, func(p models.Place) bool {
			return p.Google.PlaceId == place.Google.PlaceId
		})
		return !isIgnorePlace
	})
	s.logger.Debug("Places after filtering ignore places", zap.Int("places", len(placesFiltered)))

	placesFiltered = placefilter.FilterDefaultIgnore(placefilter.FilterDefaultIgnoreInput{
		Places:              placesFiltered,
		StartLocation:       input.BaseLocation,
		IgnoreDistanceRange: float64(input.Radius),
	})
	s.logger.Debug("Places after filtering default ignore", zap.Int("places", len(placesFiltered)))

	placesSelected := make([]models.Place, 0, input.MaxBasePlaceCount)
	for len(placesSelected) < input.MaxBasePlaceCount || len(placesFiltered) > 0 {
		// プラン間で重複が発生しないように、すでに選択された場所から500m以内の場所は選択しない
		placesFiltered = array.Filter(placesFiltered, func(place models.Place) bool {
			_, isDistanceWithIn := array.Find(placesSelected, func(p models.Place) bool {
				return p.Location.DistanceInMeter(place.Location) < 500
			})
			return !isDistanceWithIn
		})
		s.logger.Debug("Places after filtering by distance from selected places", zap.Int("Places", len(placesFiltered)))

		// レビューの最も高い場所を選択する
		placesSortedByRating := models.SortPlacesByRating(placesFiltered)
		if len(placesSortedByRating) == 0 {
			break
		}

		placesSelected = append(placesSelected, placesSortedByRating[0])
	}

	if len(placesSelected) > input.MaxBasePlaceCount {
		return placesSelected[:input.MaxBasePlaceCount]
	}

	return placesSelected
}
