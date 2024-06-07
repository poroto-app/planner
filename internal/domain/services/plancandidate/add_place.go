package plancandidate

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/models"
	"time"
)

// AddPlaceAfterPlace プランに指定された場所を追加する
// すでに指定された場所が登録されている場合は、なにもしない
func (s Service) AddPlaceAfterPlace(ctx context.Context, planCandidateSetId string, planId string, previousPlaceId string, placeId string) (*models.Plan, error) {
	planCandidateSet, err := s.planCandidateRepository.Find(ctx, planCandidateSetId, time.Now())
	if err != nil {
		return nil, fmt.Errorf("error while fetching plan candidate: %v", err)
	}

	planToUpdate := planCandidateSet.GetPlan(planId)
	if planToUpdate == nil {
		return nil, fmt.Errorf("plan not found: %v", planId)
	}

	s.logger.Debug(
		"Fetching searched places for plan candidate",
		zap.String("planCandidateSetId", planCandidateSetId),
	)

	// 追加する場所を取得
	placeToAdd, err := s.placeRepository.Find(ctx, placeId)
	if err != nil {
		return nil, fmt.Errorf("error while fetching place: %v", err)
	}

	if placeToAdd == nil {
		return nil, fmt.Errorf("place not found: %v", placeId)
	}

	// 重複して追加しないようにする
	for _, place := range planToUpdate.Places {
		if place.Id == placeToAdd.Id {
			s.logger.Debug(
				"Place is already added to plan candidate",
				zap.String("placeId", placeId),
				zap.String("planCandidateSetId", planCandidateSetId),
			)
			return planToUpdate, nil
		}
	}

	// 画像を取得
	s.logger.Info(
		"Fetching photos and reviews for places for plan candidate",
		zap.String("planCandidateSetId", planCandidateSetId),
	)
	placesWithPhoto := s.placeSearchService.FetchPlacesPhotosAndSave(ctx, *placeToAdd)
	placeToAdd = &placesWithPhoto[0]
	s.logger.Info(
		"Successfully fetched photos and reviews for places for plan candidate",
		zap.String("planCandidateSetId", planCandidateSetId),
	)

	// プランに指定された場所を追加
	s.logger.Info(
		"Adding place to plan candidate",
		zap.String("planCandidateSetId", planCandidateSetId),
	)
	if err := s.planCandidateRepository.AddPlaceToPlan(ctx, planCandidateSetId, planId, previousPlaceId, *placeToAdd); err != nil {
		return nil, fmt.Errorf("error while adding place to plan candidate: %v", err)
	}
	s.logger.Info(
		"Successfully added place to plan candidate",
		zap.String("planCandidateSetId", planCandidateSetId),
	)

	// 最新のプランの情報を取得
	s.logger.Info(
		"Fetching plan candidate",
		zap.String("planCandidateSetId", planCandidateSetId),
	)
	planCandidateSet, err = s.planCandidateRepository.Find(ctx, planCandidateSetId, time.Now())
	if err != nil {
		return nil, fmt.Errorf("error while fetching plan candidate: %v", err)
	}
	s.logger.Info(
		"Successfully fetched plan candidate",
		zap.String("planCandidateSetId", planCandidateSetId),
	)

	plan := planCandidateSet.GetPlan(planId)
	if plan == nil {
		return nil, fmt.Errorf("plan not found: %v", planId)
	}

	return plan, nil
}
