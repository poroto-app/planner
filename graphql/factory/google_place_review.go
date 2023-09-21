package factory

import (
	graphql "poroto.app/poroto/planner/graphql/model"
	"poroto.app/poroto/planner/internal/domain/models"
)

func GooglePlaceReviewFromDomainModel(review models.GooglePlaceReview) *graphql.GooglePlaceReview {
	return &graphql.GooglePlaceReview{
		Rating:           int(review.Rating),
		Text:             review.Text,
		Time:             review.Time,
		AuthorName:       review.AuthorName,
		AuthorURL:        review.AuthorUrl,
		AuthorPhotoURL:   review.AuthorProfileImageUrl,
		Language:         review.Language,
		OriginalLanguage: review.Language,
	}
}
