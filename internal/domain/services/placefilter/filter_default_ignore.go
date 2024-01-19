package placefilter

import "poroto.app/poroto/planner/internal/domain/models"

const (
	ignorePlaceDistanceRange = 1500
)

type FilterDefaultIgnoreInput struct {
	Places              []models.Place
	StartLocation       models.GeoLocation
	IgnoreDistanceRange float64
}

// FilterDefaultIgnore はプラン作成時に共通して無視する場所をフィルタリングする
func FilterDefaultIgnore(input FilterDefaultIgnoreInput) []models.Place {
	if input.StartLocation.IsZero() {
		panic("StartLocation is empty")
	}

	if input.IgnoreDistanceRange == 0 {
		input.IgnoreDistanceRange = ignorePlaceDistanceRange
	}

	placesFiltered := input.Places

	//　重複した場所を削除
	placesFiltered = FilterDuplicated(placesFiltered)

	// 特定のカテゴリは無視する
	placesFiltered = FilterIgnoreCategory(placesFiltered)
	placesFiltered = FilterByCategory(placesFiltered, models.GetCategoryToFilter(), true)

	// 会社は無視する
	placesFiltered = FilterCompany(placesFiltered)

	// 1.5km圏外の場所は無視する
	placesFiltered = FilterWithinDistanceRange(placesFiltered, input.StartLocation, 0, input.IgnoreDistanceRange)

	// 画像がない場所は無視する
	placesFiltered = FilterByHasPhoto(placesFiltered)

	return placesFiltered

}
