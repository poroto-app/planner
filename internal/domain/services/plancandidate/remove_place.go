package plancandidate

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/models"
)

func (s Service) RemovePlaceFromPlan(ctx context.Context, planCandidateId string, planId string, placeId string) (*models.Plan, error) {
	err := s.planCandidateRepository.RemovePlaceFromPlan(ctx, planCandidateId, planId, placeId)
	if err != nil {
		return nil, fmt.Errorf("error while removing place from plan candidate: %v", err)
	}

	planCandidate, err := s.planCandidateRepository.Find(ctx, planCandidateId)
	if err != nil {
		return nil, fmt.Errorf("error while retrieving plan candidate: %v", err)
	}

	plan := planCandidate.GetPlan(planId)
	if err != nil {
		return nil, fmt.Errorf("plan not found in plan candidate: %v", err)
	}

	return plan, nil
}
