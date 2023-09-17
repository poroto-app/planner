package plancandidate

import (
	"context"
	"log"
	"sort"

	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/placefilter"
	placesApi "poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

const (
	maxAddablePlaces = 10
)

// FetchCandidatePlaces はプランの候補となる場所を取得する
func (s Service) FetchCandidatePlaces(
	ctx context.Context,
	createPlanSessionId string,
) ([]*models.Place, error) {
	placesSearched, err := s.placeSearchResultRepository.Find(ctx, createPlanSessionId)
	if err != nil {
		return nil, err
	}

	planCandidate, err := s.planCandidateRepository.Find(ctx, createPlanSessionId)
	if err != nil {
		return nil, err
	}

	placesFiltered := placesSearched
	placesFiltered = placefilter.FilterIgnoreCategory(placesFiltered)
	placesFiltered = placefilter.FilterByCategory(placesFiltered, models.GetCategoryToFilter(), true)

	// 重複した場所を削除
	placesFiltered = placefilter.FilterDuplicated(placesFiltered)

	placesSortedByRating := placesFiltered
	sort.Slice(placesSortedByRating, func(i, j int) bool {
		return placesSortedByRating[i].Rating > placesSortedByRating[j].Rating
	})

	places := make([]*models.Place, 0, len(placesSortedByRating))
	for _, place := range placesSortedByRating {
		if planCandidate.HasPlace(place.PlaceID) {
			continue
		}

		var categoryMain *models.LocationCategory
		for _, placeType := range place.Types {
			c := models.CategoryOfSubCategory(placeType)
			if c != nil {
				categoryMain = c
				break
			}
		}
		// MEMO: カテゴリが不明な場合，滞在時間が取得できない
		if categoryMain == nil {
			log.Printf("place %s has no category\n", place.Name)
			continue
		}

		thumbnail, err := s.placesApi.FetchPlacePhoto(place, &placesApi.ImageSize{
			Width:  placesApi.ImgThumbnailMaxWidth,
			Height: placesApi.ImgThumbnailMaxHeight,
		})
		if err != nil {
			log.Printf("error while fetching place photo: %v\n", err)
			continue
		}

		if thumbnail == nil {
			log.Printf("place %s has no thumbnail\n", place.Name)
			continue
		}

		thumbnailUrl := thumbnail.ImageUrl

		places = append(places, &models.Place{
			// TODO: Google Places APIで取得されるIDと対応関係のあるIDを別で保存する
			Id:                    place.PlaceID,
			GooglePlaceId:         &place.PlaceID,
			Name:                  place.Name,
			Location:              place.Location.ToGeoLocation(),
			EstimatedStayDuration: categoryMain.EstimatedStayDuration,
			Category:              categoryMain.Name,
			Categories:            models.GetCategoriesFromSubCategories(place.Types),
			Thumbnail:             &thumbnailUrl,
			Photos:                []string{thumbnailUrl},
		})

		if len(places) >= maxAddablePlaces {
			break
		}
	}

	return places, nil
}
