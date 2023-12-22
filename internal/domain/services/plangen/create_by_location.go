package plangen

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/place"
)

// CreatePlanByLocation 指定された位置情報を基準とするプランを作成する
func (s Service) CreatePlanByLocation(
	ctx context.Context,
	createPlanSessionId string,
	baseLocation models.GeoLocation,
	categoryNamesPreferred *[]string,
	categoryNamesDisliked *[]string,
	freeTime *int,
	createBasedOnCurrentLocation bool,
) (*[]models.Plan, error) {
	// 付近の場所を検索
	var places []models.Place

	// すでに検索を行っている場合はその結果を取得
	placesSearched, err := s.placeService.FetchSearchedPlaces(ctx, createPlanSessionId)
	if err != nil {
		s.logger.Warn(
			"error while fetching searched Places",
			zap.String("PlanCandidateId", createPlanSessionId),
			zap.Error(err),
		)
	} else if placesSearched != nil {
		s.logger.Debug(
			"Places fetched",
			zap.String("PlanCandidateId", createPlanSessionId),
			zap.Int("Places", len(placesSearched)),
		)
		places = placesSearched
	}

	// 検索を行っていない場合は検索を行う
	if places == nil {
		googlePlaces, err := s.placeService.SearchNearbyPlaces(ctx, place.SearchNearbyPlacesInput{Location: baseLocation})
		if err != nil {
			return nil, fmt.Errorf("error while fetching google Places: %v\n", err)
		}

		placesSaved, err := s.placeService.SaveSearchedPlaces(ctx, createPlanSessionId, googlePlaces)
		if err != nil {
			return nil, fmt.Errorf("error while saving searched Places: %v\n", err)
		}

		places = placesSaved
	}

	s.logger.Debug(
		"Places searched",
		zap.String("PlanCandidateId", createPlanSessionId),
		zap.Float64("lat", baseLocation.Latitude),
		zap.Float64("lng", baseLocation.Longitude),
		zap.Int("Places", len(places)),
	)

	// プラン作成の基準となる場所を選択
	var placesRecommend []models.Place

	placesRecommend = append(placesRecommend, s.SelectBasePlace(SelectBasePlaceInput{
		BaseLocation:           baseLocation,
		Places:                 places,
		CategoryNamesPreferred: categoryNamesPreferred,
		CategoryNamesDisliked:  categoryNamesDisliked,
		ShouldOpenNow:          false,
	})...)
	for _, place := range placesRecommend {
		s.logger.Debug(
			"place recommended",
			zap.String("place", place.Google.Name),
		)
	}

	// 最もおすすめ度が高い３つの場所を基準にプランを作成する
	var createPlanParams []CreatePlanParams
	for _, placeRecommend := range placesRecommend {
		var placesInPlan []models.Place
		for _, createPlanParam := range createPlanParams {
			placesInPlan = append(placesInPlan, createPlanParam.places...)
		}

		planPlaces, err := s.createPlanPlaces(
			ctx,
			CreatePlanPlacesParams{
				PlanCandidateId:              createPlanSessionId,
				LocationStart:                baseLocation,
				PlaceStart:                   placeRecommend,
				Places:                       places,
				PlacesOtherPlansContain:      placesInPlan,
				FreeTime:                     freeTime,
				CategoryNamesDisliked:        categoryNamesDisliked,
				CreateBasedOnCurrentLocation: createBasedOnCurrentLocation,
				ShouldOpenWhileTraveling:     createBasedOnCurrentLocation, // 現在地からプランを作成した場合は、今から出発した場合に閉まってしまうお店は含めない
			},
		)
		if err != nil {
			s.logger.Warn(
				"error while creating plan",
				zap.String("place", placeRecommend.Google.Name),
				zap.Error(err),
			)
			continue
		}

		createPlanParams = append(createPlanParams, CreatePlanParams{
			locationStart: baseLocation,
			placeStart:    placeRecommend,
			places:        planPlaces,
		})
	}

	plans := s.createPlanData(ctx, createPlanSessionId, createPlanParams...)

	return &plans, nil
}
