package user

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/models"
)

type BindPlanCandidateSetToUserInput struct {
	PlanCandidateSetIds []string
	UserId              string
	FirebaseAuthToken   string
}

// BindPlanCandidateSetToUser 未ログイン時に作成されたプランや、いいねした場所の情報をユーザーと紐づける
func (s Service) BindPlanCandidateSetToUser(ctx context.Context, input BindPlanCandidateSetToUserInput) (*models.User, error) {
	checkAuthStateResult, err := s.CheckUserAuthState(ctx, CheckUserAuthStateInput{
		UserId:            input.UserId,
		FirebaseAuthToken: input.FirebaseAuthToken,
	})
	if err != nil {
		return nil, fmt.Errorf("error while checking user auth state: %v", err)
	}

	if !checkAuthStateResult.IsAuthenticated {
		return nil, fmt.Errorf("user is not authenticated")
	}

	if err := s.placeRepository.UpdateLikeByPlanCandidateSetToUser(ctx, input.UserId, input.PlanCandidateSetIds); err != nil {
		return nil, fmt.Errorf("error while updating like to place in plan: %v", err)
	}

	if err := s.planRepository.UpdatePlanAuthorUserByPlanCandidateSet(ctx, input.UserId, input.PlanCandidateSetIds); err != nil {
		return nil, fmt.Errorf("error while updating author of plan candidate to user: %v", err)
	}

	return &checkAuthStateResult.User, nil
}
