package factory

import (
	"poroto.app/poroto/planner/internal/domain/models"
	graphql "poroto.app/poroto/planner/internal/infrastructure/graphql/model"
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
