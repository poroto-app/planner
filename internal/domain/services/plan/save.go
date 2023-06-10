package plan

import (
	"context"
	"fmt"

	"poroto.app/poroto/planner/internal/domain/models"
)

func (s PlanService) SavePlanFromPlanCandidate(ctx context.Context, planCandidateId string, planId string) (*models.Plan, error) {
	planCandidate, err := s.planCandidateRepository.Find(ctx, planCandidateId)
	if err != nil {
		return nil, err
	}

	var planToSave *models.Plan
	for _, plan := range planCandidate.Plans {
		if plan.Id == planId {
			planToSave = &plan
			break
		}
	}
	if planToSave == nil {
		return nil, fmt.Errorf("plan(%v) not found in plan candidate(%v)", planId, planCandidateId)
	}

	if err := s.planRepository.Save(planToSave); err != nil {
		return nil, err
	}

	return planToSave, nil
}
