package plancandidate

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/models"
)

func (s Service) AddPlace(ctx context.Context, planCandidateId string, planId string, placeId string) (*models.Plan, error) {
	placesSearched, err := s.placeSearchResultRepository.Find(ctx, planCandidateId)
	if err != nil {
		return nil, err
	}

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
				Id:                    place.PlaceID,
				Name:                  place.Name,
				Location:              place.Location.ToGeoLocation(),
				EstimatedStayDuration: estimatedStayDuration,
				Categories:            categories,
			}
			break
		}
	}

	if placeToAdd == nil {
		return nil, nil
	}

	// TODO: キャッシュする
	places := []models.Place{*placeToAdd}
	places = s.planGeneratorService.FetchPlacesPhotos(ctx, places)
	places = s.planGeneratorService.FetchReviews(ctx, places)
	*placeToAdd = places[0]

	// プランに指定された場所を追加
	if err := s.planCandidateRepository.AddPlaceToPlan(ctx, planCandidateId, planId, *placeToAdd); err != nil {
		return nil, fmt.Errorf("error while adding place to plan candidate: %v\n", err)
	}

	// 最新のプランの情報を取得
	planCandidate, err := s.planCandidateRepository.Find(ctx, planCandidateId)
	if err != nil {
		return nil, fmt.Errorf("error while fetching plan candidate: %v\n", err)
	}

	plan := planCandidate.GetPlan(planId)
	if plan == nil {
		return nil, fmt.Errorf("plan not found: %v\n", planId)
	}

	return plan, nil
}
