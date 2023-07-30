package plan

import (
	"context"
	"log"

	"github.com/google/uuid"
	"poroto.app/poroto/planner/internal/domain/models"
	placesApi "poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

// FetchCandidatePlaces はプランの候補となる場所を取得する
func (s PlanService) FetchCandidatePlaces(
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

	places := make([]*models.Place, 0, len(placesSearched))
	for _, place := range placesSearched {
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
			Width:  placesApi.ImgThumbnailMaxHeight,
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
			Id:                    uuid.New().String(),
			GooglePlaceId:         &place.PlaceID,
			Name:                  place.Name,
			Location:              place.Location.ToGeoLocation(),
			EstimatedStayDuration: categoryMain.EstimatedStayDuration,
			Category:              categoryMain.Name,
			Thumbnail:             &thumbnailUrl,
		})
	}

	return places, nil
}
