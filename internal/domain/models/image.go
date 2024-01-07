package models

import (
	"fmt"
	"poroto.app/poroto/planner/internal/domain/utils"
)

type ImageSmallLarge struct {
	Small *string
	Large *string
}

func NewImage(small, large *string) (*ImageSmallLarge, error) {
	if small == nil && large == nil {
		return nil, fmt.Errorf("small and large are both nil")
	}

	return &ImageSmallLarge{
		Small: utils.StrCopyPointerValue(small),
		Large: utils.StrCopyPointerValue(large),
	}, nil
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
