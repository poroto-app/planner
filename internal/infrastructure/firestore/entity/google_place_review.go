package entity

import "poroto.app/poroto/planner/internal/domain/models"

type GooglePlaceReviewEntity struct {
	GooglePlaceId    string  `firestore:"google_place_id"`
	Rating           int     `firestore:"rating"`
	Text             *string `firestore:"text,omitempty"`
	Time             int     `firestore:"time"`
	AuthorName       string  `firestore:"author_name"`
	AuthorUrl        *string `firestore:"author_url,omitempty"`
	AuthorProfileUrl *string `firestore:"author_profile_url,omitempty"`
	Language         *string `firestore:"language,omitempty"`
	OriginalLanguage *string `firestore:"original_language,omitempty"`
}

func GooglePlaceReviewEntityFromGooglePlaceReview(review models.GooglePlaceReview, googlePlaceId string) GooglePlaceReviewEntity {
	return GooglePlaceReviewEntity{
		GooglePlaceId:    googlePlaceId,
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

func (g GooglePlaceReviewEntity) ToGooglePlaceReview() models.GooglePlaceReview {
	return models.GooglePlaceReview{
		Rating:                uint(g.Rating),
		Text:                  g.Text,
		Time:                  g.Time,
		AuthorName:            g.AuthorName,
		AuthorUrl:             g.AuthorUrl,
		AuthorProfileImageUrl: g.AuthorProfileUrl,
		Language:              g.Language,
		OriginalLanguage:      g.OriginalLanguage,
	}
}
