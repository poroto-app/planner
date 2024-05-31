package plancandidate

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"poroto.app/poroto/planner/internal/domain/models"
)

type CreatePlanCandidateSetFromSavedPlanInput struct {
	PlanId            string
	UserId            *string
	FirebaseAuthToken *string
}

type CreatePlanCandidateSetFromSavedPlanOutput struct {
	PlanCandidateSet models.PlanCandidate
}

func (s Service) CreatePlanCandidateSetFromSavedPlan(ctx context.Context, input CreatePlanCandidateSetFromSavedPlanInput) (*CreatePlanCandidateSetFromSavedPlanOutput, error) {
	plan, err := s.planRepository.Find(ctx, input.PlanId)
	if err != nil {
		return nil, fmt.Errorf("error while fetching plan: %v", err)
	}

	if plan == nil {
		return nil, fmt.Errorf("plan not found")
	}

	// 保存されるときにもとのプランと別のIDになるようにする
	copiedPlan := plan
	copiedPlan.Id = uuid.New().String()

	newPlanCandidateSetId := uuid.New().String()

	if err := s.CreatePlanCandidate(ctx, newPlanCandidateSetId); err != nil {
		return nil, fmt.Errorf("error while creating plan candidate: %v", err)
	}

	if err := s.planCandidateRepository.AddPlan(ctx, newPlanCandidateSetId, *copiedPlan); err != nil {
		return nil, fmt.Errorf("error while adding plan to plan candidate: %v", err)
	}

	planCandidateSet, err := s.FindPlanCandidate(ctx, FindPlanCandidateInput{
		PlanCandidateId:   newPlanCandidateSetId,
		UserId:            input.UserId,
		FirebaseAuthToken: input.FirebaseAuthToken,
	})

	if err != nil {
		return nil, fmt.Errorf("error while fetching plan candidate: %v", err)
	}

	return &CreatePlanCandidateSetFromSavedPlanOutput{
		PlanCandidateSet: *planCandidateSet,
	}, nil
}
