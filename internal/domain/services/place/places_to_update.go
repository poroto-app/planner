package place

import (
	"context"
	"fmt"

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
	checkAuthStateResult, err := s.userService.CheckUserAuthState(ctx, user.CheckUserAuthStateInput{
		UserId:            userId,
		FirebaseAuthToken: firebaseAuthToken,
	})
	if err != nil {
		return fmt.Errorf("error while checking user auth state: %v", err)
	}
	if !checkAuthStateResult.IsAuthenticated {
		return fmt.Errorf("user is not authenticated: user id: %s", userId)
	}
	for _, input := range inputs {
		placePhotos = append(placePhotos, models.PlacePhoto{
			PlaceId:  input.PlaceId,
			UserId:   userId,
			PhotoUrl: input.PhotoUrl,
			Width:    input.Width,
			Height:   input.Height,
		})
	}

	if err = s.placeRepository.SavePlacePhotos(ctx, placePhotos); err != nil {
		return fmt.Errorf("error while saving place photos: %v", err)
	}
	return nil
}
