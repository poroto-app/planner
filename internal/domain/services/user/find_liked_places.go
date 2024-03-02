package user

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/models"
)

type FindLikedPlacesInput struct {
	UserId            string
	FirebaseAuthToken string
}

func (s Service) FindLikePlaces(ctx context.Context, input FindLikedPlacesInput) (*[]models.Place, error) {
	checkAuthStateResult, err := s.CheckUserAuthState(ctx, CheckUserAuthStateInput{
		UserId:            input.UserId,
		FirebaseAuthToken: input.FirebaseAuthToken,
	})
	if err != nil {
		return nil, err
	}

	if !checkAuthStateResult.IsAuthenticated {
		return nil, fmt.Errorf("user is not authenticated")
	}

	likedPlaces, err := s.placeRepository.FindLikePlacesByUserId(ctx, input.UserId)
	if err != nil {
		return nil, fmt.Errorf("error while fetching liked places: %v", err)
	}

	return likedPlaces, nil
}
