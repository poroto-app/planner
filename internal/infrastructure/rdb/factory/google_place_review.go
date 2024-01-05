package factory

import (
	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/entities"
)

func NewGooglePlaceReviewFromEntity(googlePlaceReviewEntity entities.GooglePlaceReview) models.GooglePlaceReview {
	return models.GooglePlaceReview{
		Rating:                uint(googlePlaceReviewEntity.Rating.Int),
		Text:                  googlePlaceReviewEntity.Text.Ptr(),
		Time:                  googlePlaceReviewEntity.Time.Int,
		AuthorName:            googlePlaceReviewEntity.AuthorName.String,
		AuthorUrl:             googlePlaceReviewEntity.AuthorURL.Ptr(),
		AuthorProfileImageUrl: googlePlaceReviewEntity.AuthorProfilePhotoURL.Ptr(),
		Language:              googlePlaceReviewEntity.Language.Ptr(),
	}
}

func NewGooglePlaceReviewEntityFromGooglePlaceReview(googlePlaceReview models.GooglePlaceReview) entities.GooglePlaceReview {
	return entities.GooglePlaceReview{
		ID:                    uuid.New().String(),
		Rating:                null.IntFrom(int(googlePlaceReview.Rating)),
		Text:                  null.StringFromPtr(googlePlaceReview.Text),
		Time:                  null.IntFrom(googlePlaceReview.Time),
		AuthorName:            null.StringFrom(googlePlaceReview.AuthorName),
		AuthorURL:             null.StringFromPtr(googlePlaceReview.AuthorUrl),
		AuthorProfilePhotoURL: null.StringFromPtr(googlePlaceReview.AuthorProfileImageUrl),
		Language:              null.StringFromPtr(googlePlaceReview.Language),
	}
}

func NewGooglePlaceReviewSliceFromGooglePlaceDetail(googlePlaceDetail models.GooglePlaceDetail) entities.GooglePlaceReviewSlice {
	if len(googlePlaceDetail.Reviews) == 0 {
		return nil
	}

	var googlePlaceReviewEntities entities.GooglePlaceReviewSlice
	for _, googlePlaceReview := range googlePlaceDetail.Reviews {
		gpr := NewGooglePlaceReviewEntityFromGooglePlaceReview(googlePlaceReview)
		googlePlaceReviewEntities = append(googlePlaceReviewEntities, &gpr)
	}
	return googlePlaceReviewEntities
}

func NewGooglePlaceReviewSliceFromGooglePlace(googlePlace models.GooglePlace) entities.GooglePlaceReviewSlice {
	if googlePlace.PlaceDetail == nil {
		return nil
	}
	return NewGooglePlaceReviewSliceFromGooglePlaceDetail(*googlePlace.PlaceDetail)
}
