package place

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/user"
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
		checkAuthStateResult, err := s.userService.CheckUserAuthState(ctx, user.CheckUserAuthStateInput{
			UserId:            input.UserId,
			FirebaseAuthToken: input.FirebaseAuthToken,
		})
		if err != nil {
			s.logger.Error("error while checking user auth state", zap.Error(err))
			continue
		}

		if !checkAuthStateResult.IsAuthenticated {
			s.logger.Error("user is not authenticated", zap.String("userId", input.UserId), zap.String("firebaseAuthToken", input.FirebaseAuthToken))
			continue
		}

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
