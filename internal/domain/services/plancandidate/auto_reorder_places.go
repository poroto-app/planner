package plancandidate

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/models"
)

type AutoReorderPlacesInput struct {
	PlanCandidateId string
	PlanId          string
}

// AutoReorderPlaces はプラン候補の場所をスタート地点からの移動が最小になるように並び替える
func (s *Service) AutoReorderPlaces(ctx context.Context, input AutoReorderPlacesInput) (*models.Plan, error) {
	plan, err := s.planCandidateRepository.FindPlan(ctx, input.PlanCandidateId, input.PlanId)
	if err != nil {
		return nil, fmt.Errorf("failed to find plan: %w", err)
	}

	if plan == nil {
		return nil, fmt.Errorf("plan not found")
	}

	placesReordered := plan.PlacesReorderedToMinimizeDistance()
	plan.Places = placesReordered

	var placeIdsOrdered []string
	for _, place := range placesReordered {
		placeIdsOrdered = append(placeIdsOrdered, place.Id)
	}

	if _, err := s.planCandidateRepository.UpdatePlacesOrder(ctx, input.PlanId, input.PlanCandidateId, placeIdsOrdered); err != nil {
		return nil, fmt.Errorf("failed to update places order: %v", err)
	}

	return plan, nil
}
