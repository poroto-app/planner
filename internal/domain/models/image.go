package models

// Image
// IsGooglePhotos は，Google Photos から取得した画像かどうかを示す
type Image struct {
	Width          uint
	Height         uint
	URL            string
	IsGooglePhotos bool
}

// ImageSmallLarge
// IsGooglePhotos は，Google Photos から取得した画像かどうかを示す
type ImageSmallLarge struct {
	Small          *string
	Large          *string
	IsGooglePhotos bool
}

// Default は，画像のデフォルトのURLを返す
func (i ImageSmallLarge) Default() string {
	if i.Large != nil {
		return *i.Large
	}
	if i.Small != nil {
		return *i.Small
	}

	panic("image is empty")
}
