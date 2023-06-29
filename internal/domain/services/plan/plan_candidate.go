package plan

import (
	"context"
	"time"

	"poroto.app/poroto/planner/graphql/model"
	"poroto.app/poroto/planner/internal/domain/models"
)

func (s PlanService) CachePlanCandidate(ctx context.Context, session string, plans []models.Plan, createdBasedOnCurrentLocation bool) error {
	return s.planCandidateRepository.Save(ctx, &models.PlanCandidate{
		Id:                            session,
		Plans:                         plans,
		CreatedBasedOnCurrentLocation: createdBasedOnCurrentLocation,
		ExpiresAt:                     time.Now().Add(7 * 24 * time.Hour),
	})
}

func (s PlanService) FindPlanCandidate(ctx context.Context, planCandidateId string) (*models.PlanCandidate, error) {
	return s.planCandidateRepository.Find(ctx, planCandidateId)
}

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
