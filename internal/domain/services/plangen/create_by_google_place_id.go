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
	startPlace, err := s.placeService.FetchGooglePlace(ctx, input.GooglePlaceId)
	if err != nil {
		return nil, err
	}

	if startPlace == nil {
		return nil, fmt.Errorf("could not fetch start place: %v", input.GooglePlaceId)
	}

	s.logger.Debug(
		"successfully fetched start place by google place id",
		zap.String("PlanCandidateId", input.PlanCandidateId),
		zap.String("placeId", startPlace.Id),
		zap.String("googlePlaceId", input.GooglePlaceId),
		zap.String("name", startPlace.Name),
	)

	// 付近の場所を検索
	var places []models.Place
	places = append(places, *startPlace)

	placesSearched, err := s.placeService.FetchSearchedPlaces(ctx, input.PlanCandidateId)
	if err != nil {
		s.logger.Warn(
			"error while fetching searched Places",
			zap.String("PlanCandidateId", input.PlanCandidateId),
			zap.Error(err),
		)
	}

	if len(placesSearched) > 0 {
		// すでに検索が行われている場合はキャッシュを利用する
		s.logger.Debug(
			"Places fetched",
			zap.String("PlanCandidateId", input.PlanCandidateId),
			zap.Int("Places", len(placesSearched)),
		)
		places = placesSearched

		// 開始地点の検索結果を保存する
		if _, err := s.placeService.SaveSearchedPlaces(ctx, input.PlanCandidateId, []models.GooglePlace{startPlace.Google}); err != nil {
			return nil, fmt.Errorf("error while saving searched Places: %v\n", err)
		}
	} else {
		// 検索を行っていない場合は検索を行う
		googlePlaces, err := s.placeService.SearchNearbyPlaces(ctx, place.SearchNearbyPlacesInput{Location: startPlace.Location})
		if err != nil {
			return nil, fmt.Errorf("error while fetching google Places: %v\n", err)
		}

		// プラン候補作成において検索した場所を保存する
		placesSaved, err := s.placeService.SaveSearchedPlaces(ctx, input.PlanCandidateId, googlePlaces)
		if err != nil {
			return nil, fmt.Errorf("error while saving searched Places: %v\n", err)
		}

		places = append(places, placesSaved...)
	}

	s.logger.Debug(
		"Places searched",
		zap.String("PlanCandidateId", input.PlanCandidateId),
		zap.String("startPlace", startPlace.Name),
		zap.Int("Places", len(places)),
	)

	// プラン作成の基準となる場所を選択
	var placesRecommend []models.Place
	placesRecommend = append(placesRecommend, *startPlace)
	placesRecommend = append(placesRecommend, s.SelectBasePlace(SelectBasePlaceInput{
		BaseLocation:           startPlace.Location,
		Places:                 places,
		CategoryNamesPreferred: input.CategoryNamesPreferred,
		CategoryNamesDisliked:  input.CategoryNamesDisliked,
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
			PlanCandidateId:          input.PlanCandidateId,
			LocationStart:            startPlace.Location,
			PlaceStart:               placeRecommended,
			Places:                   places,
			PlacesOtherPlansContain:  placesAlreadyInPlan,
			FreeTime:                 input.FreeTime,
			CategoryNamesDisliked:    input.CategoryNamesDisliked,
			ShouldOpenWhileTraveling: *input.ShouldOpenNow,
		})
		if err != nil {
			s.logger.Warn(
				"error while creating plan",
				zap.String("PlanCandidateId", input.PlanCandidateId),
				zap.String("placeId", placeRecommended.Id),
				zap.Error(err),
			)
			continue
		}

		createPlanParams = append(createPlanParams, CreatePlanParams{
			locationStart: startPlace.Location,
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
		StartPlace: *startPlace,
	}, nil
}
