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

func (g GooglePlacePhoto) ToImage() Image {
	return Image{
		Small: utils.StrCopyPointerValue(g.Small),
		Large: utils.StrCopyPointerValue(g.Large),
	}
}
