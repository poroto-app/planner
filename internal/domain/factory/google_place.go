package factory

import (
	"poroto.app/poroto/planner/internal/domain/models"
	googleplaces "poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

func GooglePlaceFromPlaceEntity(place googleplaces.Place, photos *[]models.GooglePlacePhoto) models.GooglePlace {
	var placeDetail *models.GooglePlaceDetail
	if place.PlaceDetail != nil {
		d := GooglePlaceDetailFromPlaceDetailEntity(*place.PlaceDetail)
		placeDetail = &d
	}

	var photoReferences []models.GooglePlacePhotoReference
	if photos != nil {
		photoReferences = make([]models.GooglePlacePhotoReference, len(*photos))
		for i, photo := range *photos {
			photoReferences[i] = photo.ToPhotoReference()
		}
	} else if place.PhotoReferences != nil && len(place.PhotoReferences) > 0 {
		// Nearby Search で取得した場合は PhotoReference がある
		photoReferences = make([]models.GooglePlacePhotoReference, len(place.PhotoReferences))
		for i, photo := range place.PhotoReferences {
			photoReferences[i] = GooglePlacePhotoReferenceFromPhoto(photo)
		}
	}

	return models.GooglePlace{
		PlaceId: place.PlaceID,
		Name:    place.Name,
		Types:   place.Types,
		Location: models.GeoLocation{
			Latitude:  place.Location.Latitude,
			Longitude: place.Location.Longitude,
		},
		PhotoReferences:  photoReferences,
		Rating:           place.Rating,
		UserRatingsTotal: place.UserRatingsTotal,
		PriceLevel:       place.PriceLevel,
		FormattedAddress: place.FormattedAddress,
		Vicinity:         place.Vicinity,
		Photos:           photos,
		PlaceDetail:      placeDetail,
	}
}
