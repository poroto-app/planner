package entity

import (
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/utils"
)

type ImageEntity struct {
	GooglePlaceId string  `firestore:"google_place_id"`
	Small         *string `firestore:"small,omitempty"`
	Large         *string `firestore:"large,omitempty"`
}

func ToImageEntity(googlePlaceId string, image models.ImageSmallLarge) ImageEntity {
	return ImageEntity{
		GooglePlaceId: googlePlaceId,
		Small:         utils.StrCopyPointerValue(image.Small),
		Large:         utils.StrCopyPointerValue(image.Large),
	}
}

func FromImageEntity(image ImageEntity) models.ImageSmallLarge {
	return models.ImageSmallLarge{
		Small: utils.StrCopyPointerValue(image.Small),
		Large: utils.StrCopyPointerValue(image.Large),
	}
}
