package models

type GooglePlacePhotoReference struct {
	PhotoReference   string
	Width            int
	Height           int
	HTMLAttributions []string
}

func (g GooglePlacePhotoReference) ToGooglePlacePhoto(small, large *Image) GooglePlacePhoto {
	return GooglePlacePhoto{
		PhotoReference:   g.PhotoReference,
		Width:            g.Width,
		Height:           g.Height,
		HTMLAttributions: g.HTMLAttributions,
		Small:            small,
		Large:            large,
	}
}
