package plancandidate

import (
	"context"
	"poroto.app/poroto/planner/internal/domain/models"
)

func (s Service) ChangePlacesOrderPlanCandidateSet(ctx context.Context, planId string, planCandidateSetId string, placeIdsOrdered []string) (*models.Plan, error) {
	// TODO：移動時間の再計算処理を実装（latitude, longitudeがnilでなければ使う）
	if err := s.planCandidateRepository.UpdatePlacesOrder(ctx, planId, planCandidateSetId, placeIdsOrdered); err != nil {
		return nil, err
	}

	plan, err := s.planCandidateRepository.FindPlan(ctx, planCandidateSetId, planId)
	if err != nil {
		return nil, err
	}

	return plan, nil
}
