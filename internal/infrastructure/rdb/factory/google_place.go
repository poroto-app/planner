package factory

import (
	"github.com/volatiletech/null/v8"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
)

func NewGooglePlaceFromEntity(
	googlePlaceEntity generated.GooglePlace,
	googlePlaceTypeSlice generated.GooglePlaceTypeSlice,
	googlePlacePhotoReferenceSlice generated.GooglePlacePhotoReferenceSlice,
	googlePlacePhotoAttributionSlice generated.GooglePlacePhotoAttributionSlice,
	googlePlacePhotoSlice generated.GooglePlacePhotoSlice,
	googlePlaceReviewSlice generated.GooglePlaceReviewSlice,
	googlePlaceOpeningPeriodSlice generated.GooglePlaceOpeningPeriodSlice,
) (*models.GooglePlace, error) {
	googlePlaceTypes := NewGooglePlaceTypesFromEntity(googlePlaceTypeSlice)

	var googlePlacePhotoReferences []models.GooglePlacePhotoReference
	for _, googlePlacePhotoReferenceEntity := range googlePlacePhotoReferenceSlice {
		if googlePlacePhotoReferenceEntity == nil {
			continue
		}
		if googlePlacePhotoReferenceEntity.GooglePlaceID != googlePlaceEntity.GooglePlaceID {
			continue
		}
		gpr := NewGooglePlacePhotoReferenceFromEntity(*googlePlacePhotoReferenceEntity, googlePlacePhotoAttributionSlice)
		googlePlacePhotoReferences = append(googlePlacePhotoReferences, gpr)
	}

	var googlePlacePhotos *[]models.GooglePlacePhoto
	if len(googlePlacePhotoReferences) > 0 {
		gp := make([]models.GooglePlacePhoto, 0, len(googlePlacePhotoReferences))
		for _, googlePlacePhotoReferenceEntity := range googlePlacePhotoReferenceSlice {
			if googlePlacePhotoReferenceEntity == nil {
				continue
			}
			googlePlacePhoto := NewGooglePlacePhotoFromEntity(
				*googlePlacePhotoReferenceEntity,
				googlePlacePhotoSlice,
				googlePlacePhotoAttributionSlice,
			)
			if googlePlacePhoto == nil {
				continue
			}
			gp = append(gp, *googlePlacePhoto)
		}
		googlePlacePhotos = &gp
	}

	googlePlaceDetail, err := NewGooglePlaceDetailFromGooglePlaceEntity(
		googlePlaceReviewSlice,
		googlePlaceOpeningPeriodSlice,
		googlePlacePhotoReferenceSlice,
		googlePlacePhotoAttributionSlice,
		googlePlaceEntity.GooglePlaceID,
	)
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

func NewGooglePlaceEntityFromGooglePlace(googlePlace models.GooglePlace, placeId string) generated.GooglePlace {
	return generated.GooglePlace{
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
