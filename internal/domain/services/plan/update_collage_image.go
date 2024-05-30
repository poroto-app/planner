package plan

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/user"
)

type UpdatePlanCollageImageInput struct {
	PlanId            string
	PlaceId           string
	ImageUrl          string
	UserId            string
	FirebaseAuthToken string
}

type UpdatePlanCollageImageOutput struct {
	Plan models.Plan
}

// TODO: 画像URLの代わりに画像IDを指定させる
func (s Service) UpdatePlanCollageImage(ctx context.Context, input UpdatePlanCollageImageInput) (*UpdatePlanCollageImageOutput, error) {
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

	plan, err := s.FetchPlan(ctx, input.PlanId)
	if err != nil {
		return nil, fmt.Errorf("error while fetching plan after updating: %v", err)
	}

	// プランの作者のみがプランの画像を更新できる
	if plan.Author != nil || plan.Author.Id != input.UserId {
		return nil, fmt.Errorf("user is not author of the plan")
	}

	if err = s.planRepository.UpdateCollageImage(ctx, input.PlanId, input.PlaceId, input.ImageUrl); err != nil {
		return nil, fmt.Errorf("error while updating like to place in plan: %v", err)
	}

	planCollage, err := s.planRepository.FindCollage(ctx, input.PlanId)
	if err != nil {
		return nil, fmt.Errorf("error while fetching plan collage after updating: %v", err)
	}

	plan.Collage = planCollage

	return &UpdatePlanCollageImageOutput{
		Plan: *plan,
	}, nil
}
