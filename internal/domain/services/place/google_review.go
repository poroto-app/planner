package place

import (
	"context"
	"log"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/utils"
	api "poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

// FetchReviews は、プランに含まれるすべての場所のレビューを一括で取得する
// すでにレビューを取得している場合は何もしない
func (s Service) FetchReviews(ctx context.Context, places []models.Place) []models.Place {
	ch := make(chan *models.Place, len(places))
	for _, place := range places {
		go func(ctx context.Context, place models.Place, ch chan<- *models.Place) {
			if place.GooglePlaceId == nil {
				ch <- nil
				return
			}

			if place.GooglePlaceReviews != nil && len(*place.GooglePlaceReviews) > 0 {
				ch <- &place
				return
			}

			reviews, err := s.placesApi.FetchPlaceReview(ctx, api.FetchPlaceReviewRequest{
				PlaceId:  *place.GooglePlaceId,
				Language: "ja",
			})
			if err != nil {
				log.Printf("error while fetching place reviews: %v\n", err)
				ch <- nil
				return
			}

			var googlePlaceReviews []models.GooglePlaceReview
			for _, review := range *reviews {
				// SEE: https://developers.google.com/maps/documentation/places/web-service/details?hl=ja#PlaceReview
				googlePlaceReviews = append(googlePlaceReviews, models.GooglePlaceReview{
					AuthorName:            review.AuthorName,
					Rating:                uint(review.Rating),
					Time:                  review.Time,
					AuthorUrl:             utils.StrOmitEmpty(review.AuthorURL),
					Language:              utils.StrOmitEmpty(review.Language),
					OriginalLanguage:      utils.StrOmitEmpty(review.Language),
					AuthorProfileImageUrl: utils.StrOmitEmpty(review.AuthorProfilePhoto),
					Text:                  utils.StrOmitEmpty(review.Text),
				})
			}

			place.GooglePlaceReviews = &googlePlaceReviews
			ch <- &place
		}(ctx, place, ch)
	}

	for i := 0; i < len(places); i++ {
		placeUpdated := <-ch
		if placeUpdated == nil {
			continue
		}

		for j, place := range places {
			if placeUpdated.Id == place.Id {
				places[j] = *placeUpdated
				break
			}
		}
	}

	return places
}

// FetchPlaceReviewsAndSave は，指定された場所のレビューを一括で取得し、保存する
func (s Service) FetchPlaceReviewsAndSave(ctx context.Context, planCandidateId string, places []models.Place) []models.Place {
	var googlePlaceIdsWithReviews []string
	for _, place := range places {
		if place.GooglePlaceId == nil {
			continue
		}

		if array.IsContain(googlePlaceIdsWithReviews, *place.GooglePlaceId) {
			continue
		}

		googlePlaceIdsWithReviews = append(googlePlaceIdsWithReviews, *place.GooglePlaceId)
	}

	// レビューを取得
	places = s.FetchReviews(ctx, places)

	// レビューを保存
	for _, place := range places {
		if place.GooglePlaceId == nil {
			continue
		}

		// すでにレビューが取得済みの場合は何もしない
		if !array.IsContain(googlePlaceIdsWithReviews, *place.GooglePlaceId) {
			continue
		}

		if place.GooglePlaceReviews == nil || len(*place.GooglePlaceReviews) == 0 {
			continue
		}

		if err := s.placeSearchResultRepository.SaveReviewsIfNotExist(ctx, planCandidateId, *place.GooglePlaceId, *place.GooglePlaceReviews); err != nil {
			continue
		}
	}

	return places
}
