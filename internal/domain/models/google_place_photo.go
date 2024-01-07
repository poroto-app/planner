package models

import "poroto.app/poroto/planner/internal/domain/utils"

type GooglePlacePhoto struct {
	PhotoReference   string
	Width            int
	Height           int
	HTMLAttributions []string
	Small            *string
	Large            *string
}

func (g GooglePlacePhoto) ToImage() ImageSmallLarge {
	return ImageSmallLarge{
		Small: utils.StrCopyPointerValue(g.Small),
		Large: utils.StrCopyPointerValue(g.Large),
	}
}

func (g GooglePlacePhoto) ToPhotoReference() GooglePlacePhotoReference {
	return GooglePlacePhotoReference{
		PhotoReference:   g.PhotoReference,
		Width:            g.Width,
		Height:           g.Height,
		HTMLAttributions: g.HTMLAttributions,
	}
}
