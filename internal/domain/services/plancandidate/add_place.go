package plancandidate

import (
	"context"
	"fmt"
	"log"

	"poroto.app/poroto/planner/internal/domain/models"
)

// AddPlaceAfterPlace プランに指定された場所を追加する
// すでに指定された場所が登録されている場合は、なにもしない
func (s Service) AddPlaceAfterPlace(ctx context.Context, planCandidateId string, planId string, previousPlaceId string, placeId string) (*models.Plan, error) {
	planCandidate, err := s.planCandidateRepository.Find(ctx, planCandidateId)
	if err != nil {
		return nil, fmt.Errorf("error while fetching plan candidate: %v", err)
	}

	planToUpdate := planCandidate.GetPlan(planId)
	if planToUpdate == nil {
		return nil, fmt.Errorf("plan not found: %v", planId)
	}

	log.Printf("Fetching searched places for plan candidate: %v\n", planCandidateId)
	places, err := s.placeInPlanCandidateRepository.FindByPlanCandidateId(ctx, planCandidateId)
	if err != nil {
		return nil, err
	}
	log.Printf("Successfully fetched searched places for plan candidate: %v\n", planCandidateId)

	// 追加する場所を取得する
	var placeToAdd *models.PlaceInPlanCandidate
	for _, place := range *places {
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
			log.Printf("Place %v is already added to plan candidate %v", placeId, planCandidateId)
			return planToUpdate, nil
		}
	}

	googlePlaces := []models.GooglePlace{placeToAdd.Google}

	// 画像を取得
	log.Printf("Fetching photos and reviews for places for plan candidate: %v\n", planCandidateId)
	googlePlaces = s.placeService.FetchPlacesPhotosAndSave(ctx, planCandidateId, googlePlaces...)
	log.Printf("Successfully fetched photos for places for plan candidate: %v\n", planCandidateId)

	// レビューを取得
	log.Printf("Fetching reviews for places for plan candidate: %v\n", planCandidateId)
	googlePlaces = s.placeService.FetchPlaceReviewsAndSave(ctx, planCandidateId, googlePlaces...)
	log.Printf("Successfully fetched reviews for places for plan candidate: %v\n", planCandidateId)

	placeToAdd.Google = googlePlaces[0]

	// プランに指定された場所を追加
	log.Printf("Adding place to plan candidate %v\n", planCandidateId)
	if err := s.planCandidateRepository.AddPlaceToPlan(ctx, planCandidateId, planId, previousPlaceId, placeToAdd.ToPlace()); err != nil {
		return nil, fmt.Errorf("error while adding place to plan candidate: %v", err)
	}
	log.Printf("Successfully added place to plan candidate %v\n", planCandidateId)

	// 最新のプランの情報を取得
	log.Printf("Fetching plan candidate: %v\n", planCandidateId)
	planCandidate, err = s.planCandidateRepository.Find(ctx, planCandidateId)
	if err != nil {
		return nil, fmt.Errorf("error while fetching plan candidate: %v", err)
	}
	log.Printf("Successfully fetched plan candidate: %v\n", planCandidateId)

	plan := planCandidate.GetPlan(planId)
	if plan == nil {
		return nil, fmt.Errorf("plan not found: %v", planId)
	}

	return plan, nil
}
