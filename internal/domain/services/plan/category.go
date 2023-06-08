package plan

import (
	"context"
	"fmt"

	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

func (s PlanService) CategoriesNearLocation(
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

	placesSearched = s.filterByCategory(placesSearched, models.GetCategoryToFilter())

	// TODO: 現在時刻でフィルタリングするかを指定できるようにする
	placesSearched = s.filterByOpeningNow(placesSearched)

	// 検索された場所のカテゴリとその写真を取得
	categoryPhotos := make(map[string]string)
	for _, place := range placesSearched {
		// 対応するLocationCategoryを取得（重複処理および写真保存のためmapを採用）
		for _, subCategory := range place.Types {
			category := models.CategoryOfSubCategory(subCategory)
			if category == nil {
				continue
			}

			if _, ok := categoryPhotos[category.Name]; ok {
				continue
			}

			photo, err := s.placesApi.FetchPlacePhoto(place, nil)
			if err != nil {
				continue
			}

			// 場所の写真を取得（取得できなかった場合はデフォルトの画像を利用）
			categoryPhotos[category.Name] = category.Photo
			if photo != nil {
				categoryPhotos[category.Name] = photo.ImageUrl
			}
		}
	}

	categories := make([]models.LocationCategory, 0)
	for categoryName, categoryPhoto := range categoryPhotos {
		category := models.GetCategoryOfName(categoryName)
		if category == nil {
			continue
		}

		category.Photo = categoryPhoto
		categories = append(categories, *category)
	}

	return categories, nil
}
