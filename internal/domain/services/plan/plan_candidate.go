package plan

import (
	"context"

	"poroto.app/poroto/planner/graphql/model"
	"poroto.app/poroto/planner/internal/domain/models"
)

func (s PlanService) ChangePlacesOrderPlanCandidate(
	ctx context.Context,
	planId string,
	planCandidateId string,
	placeIdsOrdered []string,
	currentLocation *model.GeoLocation,
) (*models.Plan, error) {
	// MOCK：移動時間の再計算処理を実装（latitude, longitudeがnilでなければ使う）
	return s.planCandidateRepository.UpdatePlacesOrder(ctx, planId, planCandidateId, placeIdsOrdered)
}
