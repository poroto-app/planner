package place

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/user"
)

type UploadPlacePhotoInPlanInput struct {
	PlaceId  string
	PhotoUrl string
	Width    int
	Height   int
}

func (s Service) UploadPlacePhotoInPlan(
	ctx context.Context,
	userId string,
	firebaseAuthToken string,
	inputs []UploadPlacePhotoInPlanInput,
) error {
	var placePhotos []models.PlacePhoto
	for _, input := range inputs {
		checkAuthStateResult, err := s.userService.CheckUserAuthState(ctx, user.CheckUserAuthStateInput{
			UserId:            userId,
			FirebaseAuthToken: firebaseAuthToken,
		})
		if err != nil {
			s.logger.Error("error while checking user auth state", zap.Error(err))
			continue
		}

		if !checkAuthStateResult.IsAuthenticated {
			s.logger.Error("user is not authenticated", zap.String("userId", userId), zap.String("firebaseAuthToken", firebaseAuthToken))
			continue
		}

		placePhotos = append(placePhotos, models.PlacePhoto{
			PlaceId:  input.PlaceId,
			UserId:   userId,
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
