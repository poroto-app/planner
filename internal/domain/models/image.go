package models

type Image struct {
	Width  uint
	Height uint
	URL    string
}

type ImageSmallLarge struct {
	Small *string
	Large *string
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
