package factory

import (
	"github.com/google/uuid"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
)

func NewPlacePhotoSliceFromDomainModel(placePhoto []models.PlacePhoto) generated.PlacePhotoSlice {
	var placePhotoSlice generated.PlacePhotoSlice
	for _, photo := range placePhoto {
		placePhotoSlice = append(placePhotoSlice, &generated.PlacePhoto{
			ID:       uuid.New().String(),
			PlaceID:  photo.PlaceId,
			UserID:   photo.UserId,
			PhotoURL: photo.PhotoUrl,
			Width:    photo.Width,
			Height:   photo.Height,
		})
	}
	return placePhotoSlice
}

func NewPlacePhotosFromEntities(placeId string, placePhotoSlice generated.PlacePhotoSlice) []models.PlacePhoto {
	return array.MapAndFilter(placePhotoSlice, func(placePhoto *generated.PlacePhoto) (models.PlacePhoto, bool) {
		return models.PlacePhoto{
			PlaceId:  placePhoto.PlaceID,
			UserId:   placePhoto.UserID,
			PhotoUrl: placePhoto.PhotoURL,
			Width:    placePhoto.Width,
			Height:   placePhoto.Height,
		}, placePhoto.PlaceID == placeId
	})
}
