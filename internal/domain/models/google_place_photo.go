package models

type GooglePlacePhoto struct {
	PhotoReference   string
	Width            int
	Height           int
	HTMLAttributions []string
	Small            *Image
	Large            *Image
}

func (g GooglePlacePhoto) ToImage() ImageSmallLarge {
	var small, large *string
	if g.Small != nil {
		small = &g.Small.URL
	}

	if g.Large != nil {
		large = &g.Large.URL
	}
	return ImageSmallLarge{
		Small:          small,
		Large:          large,
		IsGooglePhotos: true,
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
