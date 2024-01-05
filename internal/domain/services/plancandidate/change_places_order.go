package plancandidate

import (
	"context"
	"poroto.app/poroto/planner/internal/interface/graphql/model"

	"poroto.app/poroto/planner/internal/domain/models"
)

func (s Service) ChangePlacesOrderPlanCandidate(
	ctx context.Context,
	planId string,
	planCandidateId string,
	placeIdsOrdered []string,
	currentLocation *model.GeoLocation,
) (*models.Plan, error) {
	// TODO：移動時間の再計算処理を実装（latitude, longitudeがnilでなければ使う）
	if err := s.planCandidateRepository.UpdatePlacesOrder(ctx, planId, planCandidateId, placeIdsOrdered); err != nil {
		return nil, err
	}

	plan, err := s.planCandidateRepository.FindPlan(ctx, planCandidateId, planId)
	if err != nil {
		return nil, err
	}

	return plan, nil
}
