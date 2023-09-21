package entity

import "poroto.app/poroto/planner/internal/domain/models"

type GooglePlaceReviewEntity struct {
	Rating           int     `firestore:"rating"`
	Text             *string `firestore:"text,omitempty"`
	Time             int     `firestore:"time"`
	AuthorName       string  `firestore:"author_name"`
	AuthorUrl        *string `firestore:"author_url,omitempty"`
	AuthorProfileUrl *string `firestore:"author_profile_url,omitempty"`
	Language         *string `firestore:"language,omitempty"`
	OriginalLanguage *string `firestore:"original_language,omitempty"`
}

func ToGooglePlaceReviewEntity(review models.GooglePlaceReview) GooglePlaceReviewEntity {
	return GooglePlaceReviewEntity{
		Rating:           int(review.Rating),
		Text:             review.Text,
		Time:             review.Time,
		AuthorName:       review.AuthorName,
		AuthorUrl:        review.AuthorUrl,
		AuthorProfileUrl: review.AuthorProfileImageUrl,
		Language:         review.Language,
		OriginalLanguage: review.OriginalLanguage,
	}
}

func FromGooglePlaceReviewEntity(entity GooglePlaceReviewEntity) models.GooglePlaceReview {
	return models.GooglePlaceReview{
		Rating:                uint(entity.Rating),
		Text:                  entity.Text,
		Time:                  entity.Time,
		AuthorName:            entity.AuthorName,
		AuthorUrl:             entity.AuthorUrl,
		AuthorProfileImageUrl: entity.AuthorProfileUrl,
		Language:              entity.Language,
		OriginalLanguage:      entity.OriginalLanguage,
	}
}
