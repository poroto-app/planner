package place

import (
	"context"
	"fmt"
)

type UploadPlacePhotoInPlanInput struct {
	UserId   string
	PlaceId  string
	PhotoUrl string
	Width    int
	Height   int
}

func (s Service) UploadPlacePhotoInPlan(
	ctx context.Context,
	input UploadPlacePhotoInPlanInput,
) error {
	err := s.placeRepository.SavePlacePhotos(ctx, input.UserId, input.PlaceId, input.PhotoUrl, input.Width, input.Height)
	if err != nil {
		return fmt.Errorf("error while saving place photos: %v", err)
	}
	return nil
}
