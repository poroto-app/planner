package plangen

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"googlemaps.github.io/maps"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/placesearch"
)

// CreatePlanByLocationInput
// GooglePlaceId が指定された場合は、その場所を起点としてプランを作成する
type CreatePlanByLocationInput struct {
	PlanCandidateId              string
	LocationStart                models.GeoLocation
	GooglePlaceId                *string
	CategoryNamesPreferred       *[]string
	CategoryNamesDisliked        *[]string
	FreeTime                     *int
	CreateBasedOnCurrentLocation bool
	ShouldOpenWhileTraveling     bool
}

func (s Service) CreatePlanByLocation(ctx context.Context, input CreatePlanByLocationInput) (*[]models.Plan, error) {
	// 付近の場所を検索
	var places []models.Place

	// すでに検索を行っている場合はその結果を取得
	placesSearched, err := s.placeSearchService.FetchSearchedPlaces(ctx, input.PlanCandidateId)
	if err != nil {
		s.logger.Warn(
			"error while fetching searched Places",
			zap.String("PlanCandidateId", input.PlanCandidateId),
			zap.Error(err),
		)
	} else if placesSearched != nil {
		s.logger.Debug(
			"Places fetched",
			zap.String("PlanCandidateId", input.PlanCandidateId),
			zap.Int("Places", len(placesSearched)),
		)
		places = placesSearched
	}

	// 検索を行っていない場合は検索を行う
	if places == nil {
		googlePlaces, err := s.placeSearchService.SearchNearbyPlaces(ctx, placesearch.SearchNearbyPlacesInput{Location: input.LocationStart})
		if err != nil {
			return nil, fmt.Errorf("error while fetching google Places: %v\n", err)
		}

		placesSaved, err := s.placeSearchService.SaveSearchedPlaces(ctx, input.PlanCandidateId, googlePlaces)
		if err != nil {
			return nil, fmt.Errorf("error while saving searched Places: %v\n", err)
		}

		places = placesSaved
	}

	s.logger.Debug(
		"Places searched",
		zap.String("PlanCandidateId", input.PlanCandidateId),
		zap.Float64("lat", input.LocationStart.Latitude),
		zap.Float64("lng", input.LocationStart.Longitude),
		zap.Int("Places", len(places)),
	)

	// プラン作成の基準となる場所を選択
	var placesRecommend []models.Place

	// 指定された場所の情報を取得する
	if input.GooglePlaceId != nil {
		// TODO: 他のplacesRecommendが指定された場所と近くならないようにする
		place, found, err := s.findOrFetchPlaceById(ctx, input.PlanCandidateId, places, *input.GooglePlaceId)
		if err != nil {
			s.logger.Warn(
				"error while fetching place",
				zap.String("place", *input.GooglePlaceId),
				zap.Error(err),
			)
		}

		// 開始地点となる場所が建物であれば、そこを基準としたプランを作成する
		if place != nil && array.IsContain(place.Google.Types, string(maps.AutocompletePlaceTypeEstablishment)) {
			placesRecommend = append(placesRecommend, *place)
			if !found {
				places = append(places, *place)
			}
		}
	}

	placesRecommend = append(placesRecommend, s.SelectBasePlace(SelectBasePlaceInput{
		BaseLocation:           input.LocationStart,
		Places:                 places,
		CategoryNamesPreferred: input.CategoryNamesPreferred,
		CategoryNamesDisliked:  input.CategoryNamesDisliked,
		MaxBasePlaceCount:      10, // 選択した場所からプランが作成できないこともあるため、多めに取得する
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
		createPlanParam := s.createPlan(ctx, input, places, placeRecommend, createPlanParams, input.ShouldOpenWhileTraveling)
		if createPlanParam != nil {
			createPlanParams = append(createPlanParams, *createPlanParam)
		}

		if len(createPlanParams) >= 3 {
			break
		}
	}

	plans := s.createPlanData(ctx, input.PlanCandidateId, createPlanParams...)

	// 場所を指定してプランを作成した場合、その場所を起点としたプランを最初に表示する
	if input.GooglePlaceId != nil {
		for i, plan := range plans {
			if len(plan.Places) == 0 {
				continue
			}

			firstPlace := plan.Places[0]
			if firstPlace.Google.PlaceId == *input.GooglePlaceId {
				plans[0], plans[i] = plans[i], plans[0]
				break
			}
		}
	}

	return &plans, nil
}

// findOrFetchPlaceById は、googlePlaceId に対応する場所を
// placesSearched から探し、なければAPIを使って取得する
func (s Service) findOrFetchPlaceById(
	ctx context.Context,
	planCandidateId string,
	placesSearched []models.Place,
	googlePlaceId string,
) (*models.Place, bool, error) {
	for _, placeSearched := range placesSearched {
		if placeSearched.Google.PlaceId == googlePlaceId {
			// すでに取得されている場合はそれを返す
			return &placeSearched, true, nil
		}
	}

	place, err := s.placeSearchService.FetchGooglePlace(ctx, googlePlaceId)
	if err != nil {
		return nil, false, fmt.Errorf("error while fetching place: %v", err)
	}

	if place == nil {
		return nil, false, nil
	}

	// キャッシュする
	if _, err := s.placeSearchService.SaveSearchedPlaces(ctx, planCandidateId, []models.GooglePlace{place.Google}); err != nil {
		return nil, false, fmt.Errorf("error while saving searched Places: %v", err)
	}

	return place, false, nil
}

func (s Service) createPlan(ctx context.Context, input CreatePlanByLocationInput, places []models.Place, placeRecommend models.Place, createdPlanParams []CreatePlanParams, shouldOpenWhileTraveling bool) *CreatePlanParams {
	var placesInPlan []models.Place
	for _, createPlanParam := range createdPlanParams {
		placesInPlan = append(placesInPlan, createPlanParam.places...)
	}

	planPlaces, err := s.createPlanPlaces(
		ctx,
		CreatePlanPlacesParams{
			PlanCandidateId:          input.PlanCandidateId,
			LocationStart:            input.LocationStart,
			PlaceStart:               placeRecommend,
			Places:                   places,
			PlacesOtherPlansContain:  placesInPlan,
			FreeTime:                 input.FreeTime,
			CategoryNamesDisliked:    input.CategoryNamesDisliked,
			ShouldOpenWhileTraveling: shouldOpenWhileTraveling,
		},
	)
	if err != nil {
		s.logger.Warn(
			"error while creating plan",
			zap.String("place", placeRecommend.Google.Name),
			zap.Error(err),
		)
		return nil
	}

	return &CreatePlanParams{
		locationStart: input.LocationStart,
		placeStart:    placeRecommend,
		places:        planPlaces,
	}
}
