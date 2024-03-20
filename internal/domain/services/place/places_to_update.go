package place

import (
	"context"
	"fmt"

	"poroto.app/poroto/planner/internal/domain/models"
)

type UploadPlacePhotoInPlanInput struct {
	PlaceId           string
	UserId            string
	PhotoUrl          string
	Width             int
	Height            int
	FirebaseAuthToken string
}

func (s Service) UploadPlacePhotoInPlan(
	ctx context.Context,
	inputs []UploadPlacePhotoInPlanInput,
) error {
	var placePhotos []models.PlacePhoto
	for _, input := range inputs {
		placePhotos = append(placePhotos, models.PlacePhoto{
			PlaceId:  input.PlaceId,
			UserId:   input.UserId,
			PhotoUrl: input.PhotoUrl,
			Width:    input.Width,
			Height:   input.Height,
		})
	}
	err := s.placeRepository.SavePlacePhotos(ctx, placePhotos)
	if err != nil {
		return fmt.Errorf("error while saving place photos: %v", err)
	}
	return nil
}
