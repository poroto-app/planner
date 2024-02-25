package plan

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/user"
)

type LikeToPlaceInput struct {
	PlanId            string
	PlaceId           string
	Like              bool
	UserId            string
	FirebaseAuthToken string
}

// LikeToPlace はプランにいいねをする
func (s Service) LikeToPlace(
	ctx context.Context,
	input LikeToPlaceInput,
) (*models.Plan, error) {
	checkAuthStateResult, err := s.userService.CheckUserAuthState(ctx, user.CheckUserAuthStateInput{
		UserId:            input.UserId,
		FirebaseAuthToken: input.FirebaseAuthToken,
	})
	if err != nil {
		return nil, fmt.Errorf("error while checking user auth state: %v", err)
	}

	if !checkAuthStateResult.IsAuthenticated {
		return nil, fmt.Errorf("user is not authenticated")
	}

	err = s.placeRepository.UpdateLikeByUserId(ctx, input.UserId, input.PlaceId, input.Like)
	if err != nil {
		return nil, fmt.Errorf("error while updating like to place in plan: %v", err)
	}

	// TODO: ユーザーがいいねしたプランを取得できるようにする
	plan, err := s.FetchPlan(ctx, input.PlanId)
	if err != nil {
		return nil, fmt.Errorf("error while fetching plan after updating: %v", err)
	}

	return plan, nil
}
