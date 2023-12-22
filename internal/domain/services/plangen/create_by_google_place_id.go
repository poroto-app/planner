package plangen

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/place"
)

const (
	defaultShouldOpenNow = false
)

type CreatePlanByGooglePlaceIdInput struct {
	PlanCandidateId        string
	GooglePlaceId          string
	CategoryNamesPreferred *[]string
	CategoryNamesDisliked  *[]string
	FreeTime               *int
	ShouldOpenNow          *bool
}

type CreatePlanByGooglePlaceIdOutput struct {
	Plans      []models.Plan
	StartPlace models.Place
}

func (s Service) CreatePlanByGooglePlaceId(ctx context.Context, input CreatePlanByGooglePlaceIdInput) (*CreatePlanByGooglePlaceIdOutput, error) {
	if input.ShouldOpenNow == nil {
		v := defaultShouldOpenNow
		input.ShouldOpenNow = &v
	}

	// 開始地点となる場所を検索
	startGooglePlace, err := s.placeService.FetchGooglePlace(ctx, input.GooglePlaceId)
	if err != nil {
		return nil, err
	}

	if startGooglePlace == nil {
		return nil, fmt.Errorf("could not fetch google place: %v", input.GooglePlaceId)
	}

	// キャッシュする
	placesSaved, err := s.placeService.SaveSearchedPlaces(ctx, input.PlanCandidateId, []models.GooglePlace{*startGooglePlace})
	if err != nil {
		return nil, fmt.Errorf("error while saving searched places: %v", err)
	}
	if len(placesSaved) == 0 {
		return nil, fmt.Errorf("could not save searched places")
	}
	startPlace := placesSaved[0]

	s.logger.Debug(
		"successfully fetched start place by google place id",
		zap.String("planCandidateId", input.PlanCandidateId),
		zap.String("placeId", startPlace.Id),
		zap.String("googlePlaceId", input.GooglePlaceId),
		zap.String("name", startPlace.Name),
	)

	// 付近の場所を検索
	var places []models.Place
	places = append(places, startPlace)

	placesSearched, err := s.placeService.FetchSearchedPlaces(ctx, input.PlanCandidateId)
	if err != nil {
		s.logger.Warn(
			"error while fetching searched places",
			zap.String("planCandidateId", input.PlanCandidateId),
			zap.Error(err),
		)
	}

	if len(placesSearched) > 1 {
		// すでに検索が行われている場合はキャッシュを利用する（開始地点は除く）
		s.logger.Debug(
			"places fetched",
			zap.String("planCandidateId", input.PlanCandidateId),
			zap.Int("places", len(placesSearched)),
		)
		places = placesSearched
	} else {
		// 検索を行っていない場合は検索を行う
		googlePlaces, err := s.placeService.SearchNearbyPlaces(ctx, place.SearchNearbyPlacesInput{Location: startGooglePlace.Location})
		if err != nil {
			return nil, fmt.Errorf("error while fetching google places: %v\n", err)
		}

		placesSaved, err := s.placeService.SaveSearchedPlaces(ctx, input.PlanCandidateId, googlePlaces)
		if err != nil {
			return nil, fmt.Errorf("error while saving searched places: %v\n", err)
		}

		places = append(places, placesSaved...)
	}

	s.logger.Debug(
		"places searched",
		zap.String("planCandidateId", input.PlanCandidateId),
		zap.String("startPlace", startGooglePlace.Name),
		zap.Int("places", len(places)),
	)

	// プラン作成の基準となる場所を選択
	var placesRecommend []models.Place
	placesRecommend = append(placesRecommend, startPlace)
	placesRecommend = append(placesRecommend, s.SelectBasePlace(SelectBasePlaceInput{
		BaseLocation:           startPlace.Location,
		Places:                 places,
		CategoryNamesPreferred: input.CategoryNamesPreferred,
		CategoryNamesDisliked:  input.CategoryNamesDisliked,
		ShouldOpenNow:          *input.ShouldOpenNow,
		MaxBasePlaceCount:      defaultMaxBasePlaceCount - 1,
	})...)
	for _, place := range placesRecommend {
		s.logger.Debug(
			"place recommended",
			zap.String("placeId", place.Id),
			zap.String("name", place.Name),
		)
	}

	// プランを作成
	var createPlanParams []CreatePlanParams
	for _, placeRecommended := range placesRecommend {
		var placesAlreadyInPlan []models.Place
		for _, createPlanParam := range createPlanParams {
			placesAlreadyInPlan = append(placesAlreadyInPlan, createPlanParam.places...)
		}

		// フィルタ処理は select base place などの中で行う
		placesInPlan, err := s.createPlanPlaces(ctx, CreatePlanPlacesParams{
			planCandidateId:              input.PlanCandidateId,
			locationStart:                startGooglePlace.Location,
			placeStart:                   placeRecommended,
			places:                       places,
			placesOtherPlansContain:      placesAlreadyInPlan,
			freeTime:                     input.FreeTime,
			createBasedOnCurrentLocation: false,
			categoryNamesDisliked:        input.CategoryNamesDisliked,
			shouldOpenWhileTraveling:     *input.ShouldOpenNow,
		})
		if err != nil {
			s.logger.Warn(
				"error while creating plan",
				zap.String("planCandidateId", input.PlanCandidateId),
				zap.String("placeId", placeRecommended.Id),
				zap.Error(err),
			)
			continue
		}

		createPlanParams = append(createPlanParams, CreatePlanParams{
			locationStart: startGooglePlace.Location,
			placeStart:    placeRecommended,
			places:        placesInPlan,
		})
	}

	plans := s.createPlanData(ctx, input.PlanCandidateId, createPlanParams...)

	// 指定された場所を起点としたプランを最初に表示する
	for i, plan := range plans {
		if len(plan.Places) == 0 {
			continue
		}

		if plan.Places[0].Google.PlaceId == input.GooglePlaceId {
			plans[0], plans[i] = plans[i], plans[0]
		}
	}

	return &CreatePlanByGooglePlaceIdOutput{
		Plans:      plans,
		StartPlace: startPlace,
	}, nil
}
