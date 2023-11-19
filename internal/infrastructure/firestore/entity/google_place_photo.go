package entity

import (
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/utils"
)

// GooglePlacePhotoEntity 場所の写真
// models.GooglePlacePhoto だけでなく、 models.GooglePlacePhotoReference としても扱う（画像を取得しているかどうかの差のみのため）
type GooglePlacePhotoEntity struct {
	GooglePlaceId    string   `firestore:"google_place_id"`
	PhotoReference   string   `firestore:"photo_reference"`
	Width            int      `firestore:"width"`
	Height           int      `firestore:"height"`
	HTMLAttributions []string `firestore:"html_attributions"`
	Small            *string  `firestore:"small,omitempty"`
	Large            *string  `firestore:"large,omitempty"`
}

func GooglePlacePhotoEntityFromGooglePlacePhoto(googlePlacePhoto models.GooglePlacePhoto, googlePlaceId string) GooglePlacePhotoEntity {
	return GooglePlacePhotoEntity{
		GooglePlaceId:    googlePlaceId,
		PhotoReference:   googlePlacePhoto.PhotoReference,
		Width:            googlePlacePhoto.Width,
		Height:           googlePlacePhoto.Height,
		HTMLAttributions: googlePlacePhoto.HTMLAttributions,
		Small:            utils.StrCopyPointerValue(googlePlacePhoto.Small),
		Large:            utils.StrCopyPointerValue(googlePlacePhoto.Large),
	}
}

func (g GooglePlacePhotoEntity) ToGooglePlacePhoto() *models.GooglePlacePhoto {
	// SmallもLargeもnilの場合は、画像を取得していないため nil にする
	if g.Small == nil && g.Large == nil {
		return nil
	}

	return &models.GooglePlacePhoto{
		PhotoReference:   g.PhotoReference,
		Width:            g.Width,
		Height:           g.Height,
		HTMLAttributions: g.HTMLAttributions,
		Small:            utils.StrCopyPointerValue(g.Small),
		Large:            utils.StrCopyPointerValue(g.Large),
	}
}

func (g GooglePlacePhotoEntity) ToGooglePlacePhotoReference() models.GooglePlacePhotoReference {
	return models.GooglePlacePhotoReference{
		PhotoReference:   g.PhotoReference,
		Width:            g.Width,
		Height:           g.Height,
		HTMLAttributions: g.HTMLAttributions,
	}
}
