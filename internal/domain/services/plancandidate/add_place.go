package plancandidate

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"poroto.app/poroto/planner/internal/domain/models"
)

// AddPlace プランに指定された場所を追加する
// すでに指定された場所が登録されている場合は、なにもしない
func (s Service) AddPlace(ctx context.Context, planCandidateId string, planId string, placeId string) (*models.Plan, error) {
	planCandidate, err := s.planCandidateRepository.Find(ctx, planCandidateId)
	if err != nil {
		return nil, fmt.Errorf("error while fetching plan candidate: %v\n", err)
	}

	planToUpdate := planCandidate.GetPlan(planId)
	if planToUpdate == nil {
		return nil, fmt.Errorf("plan not found: %v\n", planId)
	}

	log.Printf("Fetching searched places for plan candidate: %v\n", planCandidateId)
	placesSearched, err := s.placeSearchResultRepository.Find(ctx, planCandidateId)
	if err != nil {
		return nil, err
	}
	log.Printf("Successfully fetched searched places for plan candidate: %v\n", planCandidateId)

	// 追加する場所を取得する
	var placeToAdd *models.Place
	for _, place := range placesSearched {
		// TODO: Google の PlaceIDでマッチさせるのではなく、plannerが作成したIDでマッチするようにする
		if place.PlaceID == placeId {
			categories := models.GetCategoriesFromSubCategories(place.Types)

			var estimatedStayDuration uint = 0
			if len(categories) > 0 {
				estimatedStayDuration = categories[0].EstimatedStayDuration
			}

			placeToAdd = &models.Place{
				Id:                    uuid.New().String(),
				GooglePlaceId:         &place.PlaceID,
				Name:                  place.Name,
				Location:              place.Location.ToGeoLocation(),
				EstimatedStayDuration: estimatedStayDuration,
				Categories:            categories,
			}
			break
		}
	}
	if placeToAdd == nil {
		return nil, fmt.Errorf("place not found: %v\n", placeId)
	}

	// 重複して追加しないようにする
	for _, place := range planToUpdate.Places {
		if place.GooglePlaceId == nil || placeToAdd.GooglePlaceId == nil {
			continue
		}

		if *place.GooglePlaceId == *placeToAdd.GooglePlaceId {
			log.Printf("Place %v is already added to plan candidate %v\n", placeId, planCandidateId)
			return planToUpdate, nil
		}
	}

	places := []models.Place{*placeToAdd}

	// 画像を取得
	log.Printf("Fetching photos and reviews for places for plan candidate: %v\n", planCandidateId)
	places = s.placeService.FetchPlacesPhotosAndSave(ctx, planCandidateId, places)
	log.Printf("Successfully fetched photos for places for plan candidate: %v\n", planCandidateId)

	// レビューを取得
	log.Printf("Fetching reviews for places for plan candidate: %v\n", planCandidateId)
	places = s.placeService.FetchPlaceReviewsAndSave(ctx, planCandidateId, places)
	log.Printf("Successfully fetched reviews for places for plan candidate: %v\n", planCandidateId)

	placeToAdd = &places[0]

	// プランに指定された場所を追加
	log.Printf("Adding place to plan candidate %v\n", planCandidateId)
	if err := s.planCandidateRepository.AddPlaceToPlan(ctx, planCandidateId, planId, *placeToAdd); err != nil {
		return nil, fmt.Errorf("error while adding place to plan candidate: %v\n", err)
	}
	log.Printf("Successfully added place to plan candidate %v\n", planCandidateId)

	// 最新のプランの情報を取得
	log.Printf("Fetching plan candidate: %v\n", planCandidateId)
	planCandidate, err = s.planCandidateRepository.Find(ctx, planCandidateId)
	if err != nil {
		return nil, fmt.Errorf("error while fetching plan candidate: %v\n", err)
	}
	log.Printf("Successfully fetched plan candidate: %v\n", planCandidateId)

	plan := planCandidate.GetPlan(planId)
	if plan == nil {
		return nil, fmt.Errorf("plan not found: %v\n", planId)
	}

	return plan, nil
}
