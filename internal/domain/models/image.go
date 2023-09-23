package models

import (
	"fmt"
	"poroto.app/poroto/planner/internal/domain/utils"
)

type Image struct {
	Small *string
	Large *string
}

func NewImage(small, large *string) (*Image, error) {
	if small == nil && large == nil {
		return nil, fmt.Errorf("small and large are both nil")
	}

	return &Image{
		Small: utils.StrCopyPointerValue(small),
		Large: utils.StrCopyPointerValue(large),
	}, nil
}

// Default は，画像のデフォルトのURLを返す
func (i Image) Default() string {
	if i.Large != nil {
		return *i.Large
	}
	if i.Small != nil {
		return *i.Small
	}

	panic("image is empty")
}
