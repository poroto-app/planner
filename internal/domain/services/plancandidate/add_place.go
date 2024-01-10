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
func (s Service) AddPlaceAfterPlace(ctx context.Context, planCandidateId string, planId string, previousPlaceId string, placeId string) (*models.Plan, error) {
	planCandidate, err := s.planCandidateRepository.Find(ctx, planCandidateId, time.Now())
	if err != nil {
		return nil, fmt.Errorf("error while fetching plan candidate: %v", err)
	}

	planToUpdate := planCandidate.GetPlan(planId)
	if planToUpdate == nil {
		return nil, fmt.Errorf("plan not found: %v", planId)
	}

	s.logger.Debug(
		"Fetching searched places for plan candidate",
		zap.String("planCandidateId", planCandidateId),
	)
	places, err := s.placeSearchService.FetchSearchedPlaces(ctx, planCandidateId)
	if err != nil {
		return nil, err
	}
	s.logger.Debug(
		"Successfully fetched searched places for plan candidate",
		zap.String("planCandidateId", planCandidateId),
	)

	// 追加する場所を検索された場所一覧から取得する
	var placeToAdd *models.Place
	for _, place := range places {
		if place.Id == placeId {
			placeToAdd = &place
			break
		}
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
				zap.String("planCandidateId", planCandidateId),
			)
			return planToUpdate, nil
		}
	}

	// 画像を取得
	s.logger.Info(
		"Fetching photos and reviews for places for plan candidate",
		zap.String("planCandidateId", planCandidateId),
	)
	placesWithPhoto := s.placeSearchService.FetchPlacesPhotosAndSave(ctx, *placeToAdd)
	placeToAdd = &placesWithPhoto[0]
	s.logger.Info(
		"Successfully fetched photos and reviews for places for plan candidate",
		zap.String("planCandidateId", planCandidateId),
	)

	// プランに指定された場所を追加
	s.logger.Info(
		"Adding place to plan candidate",
		zap.String("planCandidateId", planCandidateId),
	)
	if err := s.planCandidateRepository.AddPlaceToPlan(ctx, planCandidateId, planId, previousPlaceId, *placeToAdd); err != nil {
		return nil, fmt.Errorf("error while adding place to plan candidate: %v", err)
	}
	s.logger.Info(
		"Successfully added place to plan candidate",
		zap.String("planCandidateId", planCandidateId),
	)

	// 最新のプランの情報を取得
	s.logger.Info(
		"Fetching plan candidate",
		zap.String("planCandidateId", planCandidateId),
	)
	planCandidate, err = s.planCandidateRepository.Find(ctx, planCandidateId, time.Now())
	if err != nil {
		return nil, fmt.Errorf("error while fetching plan candidate: %v", err)
	}
	s.logger.Info(
		"Successfully fetched plan candidate",
		zap.String("planCandidateId", planCandidateId),
	)

	plan := planCandidate.GetPlan(planId)
	if plan == nil {
		return nil, fmt.Errorf("plan not found: %v", planId)
	}

	return plan, nil
}
