package factory

import (
	"fmt"
	"github.com/volatiletech/null/v8"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/entities"
)

func NewGooglePlaceFromEntity(googlePlaceEntity entities.GooglePlace) (*models.GooglePlace, error) {
	if googlePlaceEntity.R == nil {
		return nil, fmt.Errorf("googlePlaceEntity.R is nil")
	}

	googlePlaceTypes := NewGooglePlaceTypesFromEntity(googlePlaceEntity.R.GetGooglePlaceTypes())

	var googlePlacePhotoReferences []models.GooglePlacePhotoReference
	for _, googlePlacePhotoReferenceEntity := range googlePlaceEntity.R.GetGooglePlacePhotoReferences() {
		if googlePlacePhotoReferenceEntity == nil {
			continue
		}
		gpr := NewGooglePlacePhotoReferenceFromEntity(*googlePlacePhotoReferenceEntity, googlePlaceEntity.R.GetGooglePlacePhotoAttributions())
		googlePlacePhotoReferences = append(googlePlacePhotoReferences, gpr)
	}

	var googlePlacePhotos *[]models.GooglePlacePhoto
	if len(googlePlacePhotoReferences) > 0 {
		gp := make([]models.GooglePlacePhoto, 0, len(googlePlacePhotoReferences))
		for _, googlePlacePhotoReferenceEntity := range googlePlaceEntity.R.GetGooglePlacePhotoReferences() {
			if googlePlacePhotoReferenceEntity == nil {
				continue
			}
			googlePlacePhoto := NewGooglePlacePhotoFromEntity(
				*googlePlacePhotoReferenceEntity,
				googlePlaceEntity.R.GooglePlacePhotos,
				googlePlaceEntity.R.GooglePlacePhotoAttributions,
			)
			if googlePlacePhoto == nil {
				continue
			}
			gp = append(gp, *googlePlacePhoto)
		}
		googlePlacePhotos = &gp
	}

	var googlePlaceReviews []models.GooglePlaceReview
	for _, googlePlaceReviewEntity := range googlePlaceEntity.R.GetGooglePlaceReviews() {
		if googlePlaceReviewEntity == nil {
			continue
		}
		googlePlaceReviews = append(googlePlaceReviews, NewGooglePlaceReviewFromEntity(*googlePlaceReviewEntity))
	}

	googlePlaceDetail, err := NewGooglePlaceDetailFromGooglePlaceEntity(googlePlaceEntity)
	if err != nil {
		return nil, err
	}

	geolocation := models.GeoLocation{
		Latitude:  googlePlaceEntity.Latitude,
		Longitude: googlePlaceEntity.Longitude,
	}

	return &models.GooglePlace{
		PlaceId:          googlePlaceEntity.GooglePlaceID,
		Name:             googlePlaceEntity.Name,
		Types:            googlePlaceTypes,
		Location:         geolocation,
		PhotoReferences:  googlePlacePhotoReferences,
		OpenNow:          false, // TODO: DELETE ME
		PriceLevel:       googlePlaceEntity.PriceLevel.Int,
		Rating:           googlePlaceEntity.Rating.Float32,
		UserRatingsTotal: googlePlaceEntity.UserRatingsTotal.Int,
		Vicinity:         googlePlaceEntity.Vicinity.Ptr(),
		FormattedAddress: googlePlaceEntity.FormattedAddress.Ptr(),
		Photos:           googlePlacePhotos,
		PlaceDetail:      googlePlaceDetail,
	}, nil
}

func NewGooglePlaceEntityFromGooglePlace(googlePlace models.GooglePlace, placeId string) entities.GooglePlace {
	return entities.GooglePlace{
		GooglePlaceID:    googlePlace.PlaceId,
		PlaceID:          placeId,
		Name:             googlePlace.Name,
		Latitude:         googlePlace.Location.Latitude,
		Longitude:        googlePlace.Location.Longitude,
		PriceLevel:       null.IntFrom(googlePlace.PriceLevel),
		Rating:           null.Float32From(googlePlace.Rating),
		UserRatingsTotal: null.IntFrom(googlePlace.UserRatingsTotal),
		Vicinity:         null.StringFromPtr(googlePlace.Vicinity),
		FormattedAddress: null.StringFromPtr(googlePlace.FormattedAddress),
	}
}
