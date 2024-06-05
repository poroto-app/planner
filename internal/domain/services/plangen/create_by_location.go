package plangen

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"googlemaps.github.io/maps"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/placesearch"
	"sort"
)

const (
	defaultMaxDistanceFromStart = 1500
)

// CreatePlanByLocationInput
// GooglePlaceId が指定された場合は、その場所を起点としてプランを作成する
// MaxDistanceFromStart は、プランの起点となる場所を選択するときの LocationStart からの最大距離
type CreatePlanByLocationInput struct {
	PlanCandidateId              string
	LocationStart                models.GeoLocation
	GooglePlaceId                *string
	CategoryNamesPreferred       *[]string
	CategoryNamesDisliked        *[]string
	FreeTime                     *int
	CreateBasedOnCurrentLocation bool
	ShouldOpenWhileTraveling     bool
	MaxDistanceFromStart         int
}

// CreatePlanByLocation は指定した位置から近い場所を起点として複数のプランを作成する
func (s Service) CreatePlanByLocation(ctx context.Context, input CreatePlanByLocationInput) (*[]models.Plan, error) {
	if input.MaxDistanceFromStart == 0 {
		input.MaxDistanceFromStart = defaultMaxDistanceFromStart
	}

	// 付近の場所を検索
	placesNearby, err := s.placeSearchService.SearchNearbyPlaces(ctx, placesearch.SearchNearbyPlacesInput{
		Location:           input.LocationStart,
		PlanCandidateSetId: &input.PlanCandidateId,
	})
	if err != nil {
		return nil, fmt.Errorf("error while fetching google Places: %v\n", err)
	}

	s.logger.Debug(
		"Places searched",
		zap.String("PlanCandidateId", input.PlanCandidateId),
		zap.Float64("lat", input.LocationStart.Latitude),
		zap.Float64("lng", input.LocationStart.Longitude),
		zap.Int("placesCount", len(placesNearby)),
	)

	var createPlanParams []CreatePlanParams

	// 開始地点となる場所が建物であれば、そこを基準としたプランを作成する
	if input.GooglePlaceId != nil {
		place, _, err := s.findOrFetchPlaceById(ctx, placesNearby, *input.GooglePlaceId)
		if err != nil {
			s.logger.Warn(
				"error while fetching place",
				zap.String("place", *input.GooglePlaceId),
				zap.Error(err),
			)
		}

		if place != nil && array.IsContain(place.Google.Types, string(maps.AutocompletePlaceTypeEstablishment)) {
			createPlanParam := s.CreatePlan(input, placesNearby, *place, createPlanParams)
			if createPlanParam != nil {
				createPlanParams = append(createPlanParams, *createPlanParam)
			}
		}
	}

	for filterDistance := 500; filterDistance <= input.MaxDistanceFromStart; filterDistance += 400 {
		if len(createPlanParams) >= 3 {
			break
		}

		placesAlreadyAdded := array.FlatMap(createPlanParams, func(param CreatePlanParams) []models.Place {
			return param.Places
		})

		placesForPlanStart := s.SelectBasePlace(SelectBasePlaceInput{
			BaseLocation:      input.LocationStart,
			Places:            placesNearby,
			IgnorePlaces:      placesAlreadyAdded,
			MaxBasePlaceCount: 10,
			Radius:            filterDistance,
		})

		var createPlanParamsInRange []CreatePlanParams
		for _, placeForPlanStart := range placesForPlanStart {
			createPlanParam := s.CreatePlan(input, placesNearby, placeForPlanStart, createPlanParams)
			if createPlanParam != nil {
				createPlanParamsInRange = append(createPlanParamsInRange, *createPlanParam)
			}
		}

		if len(createPlanParamsInRange) == 0 {
			s.logger.Debug(
				"no plan created",
				zap.Int("filterDistance", filterDistance),
			)
			continue
		}

		// もっとも場所の数が多いプランを追加する
		sort.SliceStable(createPlanParamsInRange, func(i, j int) bool {
			return len(createPlanParamsInRange[i].Places) > len(createPlanParamsInRange[j].Places)
		})
		createPlanParams = append(createPlanParams, createPlanParamsInRange[0])
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
func (s Service) findOrFetchPlaceById(ctx context.Context, placesSearched []models.Place, googlePlaceId string) (*models.Place, bool, error) {
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

	return place, false, nil
}

func (s Service) CreatePlan(input CreatePlanByLocationInput, places []models.Place, placeRecommend models.Place, createdPlanParams []CreatePlanParams) *CreatePlanParams {
	var placesInPlan []models.Place
	for _, createPlanParam := range createdPlanParams {
		placesInPlan = append(placesInPlan, createPlanParam.Places...)
	}

	planPlaces, err := s.CreatePlanPlaces(CreatePlanPlacesInput{
		PlanCandidateId:         input.PlanCandidateId,
		LocationStart:           placeRecommend.Location,
		PlaceStart:              placeRecommend,
		Places:                  places,
		PlacesOtherPlansContain: placesInPlan,
		FreeTime:                input.FreeTime,
		CategoryNamesDisliked:   input.CategoryNamesDisliked,
	})
	if err != nil {
		s.logger.Warn(
			"error while creating plan",
			zap.String("place", placeRecommend.Google.Name),
			zap.Error(err),
		)
		return nil
	}

	return &CreatePlanParams{
		LocationStart: input.LocationStart,
		PlaceStart:    placeRecommend,
		Places:        planPlaces,
	}
}
