package plancandidate

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/services/user"

	"poroto.app/poroto/planner/internal/domain/models"
)

type LikeToPlaceInPlanCandidateInput struct {
	PlanCandidateId   string
	PlaceId           string
	Like              bool
	UserId            *string
	FirebaseAuthToken *string
}

func (s Service) LikeToPlaceInPlanCandidate(
	ctx context.Context,
	input LikeToPlaceInPlanCandidateInput,
) (*models.PlanCandidateSet, error) {
	if input.UserId != nil && input.FirebaseAuthToken == nil {
		return nil, fmt.Errorf("firebase auth token is required")
	}

	if input.UserId != nil {
		// ログインの場合はユーザーとしてLikeを更新する
		checkAuthStateResult, err := s.userService.CheckUserAuthState(ctx, user.CheckUserAuthStateInput{
			UserId:            *input.UserId,
			FirebaseAuthToken: *input.FirebaseAuthToken,
		})
		if err != nil {
			return nil, fmt.Errorf("error while checking user auth state: %v", err)
		}

		if !checkAuthStateResult.IsAuthenticated {
			return nil, fmt.Errorf("user is not authenticated")
		}

		err = s.placeRepository.UpdateLikeByUserId(ctx, *input.UserId, input.PlaceId, input.Like)
		if err != nil {
			return nil, fmt.Errorf("error while updating like to place in plan candidate: %v", err)
		}
	} else {
		// ゲストの場合はプラン候補作成セッションとしてLikeを更新する
		err := s.planCandidateRepository.UpdateLikeToPlaceInPlanCandidate(ctx, input.PlanCandidateId, input.PlaceId, input.Like)
		if err != nil {
			return nil, fmt.Errorf("error while updating like to place in plan candidate: %v", err)
		}
	}

	planCandidate, err := s.FindPlanCandidate(ctx, FindPlanCandidateInput{
		PlanCandidateId:   input.PlanCandidateId,
		UserId:            input.UserId,
		FirebaseAuthToken: input.FirebaseAuthToken,
	})
	if err != nil {
		return nil, fmt.Errorf("error while fetching plan candidate after updating: %v", err)
	}

	return planCandidate, nil
}
