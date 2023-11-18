package factory

import (
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

func GooglePlaceReviewsFromPlaceDetail(placeDetail places.PlaceDetail) []models.GooglePlaceReview {
	var reviews []models.GooglePlaceReview
	for _, review := range placeDetail.Reviews {
		reviews = append(reviews, models.GooglePlaceReview{
			Rating:                uint(review.Rating),
			Text:                  utils.StrPointer(review.Text),
			Time:                  review.Time,
			AuthorName:            review.AuthorName,
			AuthorUrl:             utils.StrPointer(review.AuthorURL),
			AuthorProfileImageUrl: utils.StrPointer(review.AuthorProfilePhoto),
			Language:              utils.StrPointer(review.Language),
			OriginalLanguage:      utils.StrPointer(review.Language),
		})
	}
	return reviews
}
