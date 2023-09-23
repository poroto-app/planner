package entity

import (
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/utils"
)

type ImageEntity struct {
	Small *string `firestore:"small,omitempty"`
	Large *string `firestore:"large,omitempty"`
}

func ToImageEntity(image models.Image) ImageEntity {
	return ImageEntity{
		Small: utils.StrCopyPointerValue(image.Small),
		Large: utils.StrCopyPointerValue(image.Large),
	}
}

func FromImageEntity(image ImageEntity) models.Image {
	return models.Image{
		Small: utils.StrCopyPointerValue(image.Small),
		Large: utils.StrCopyPointerValue(image.Large),
	}
}
