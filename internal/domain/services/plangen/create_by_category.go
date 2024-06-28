package plangen

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/placesearch"
)

type CreatePlanByCategoryInput struct {
	PlanCandidateSetId string
	Category           models.LocationCategoryCreatePlan
	Location           models.GeoLocation
	RadiusInKm         float64
}

func (s Service) CreatePlanByCategory(ctx context.Context, input CreatePlanByCategoryInput) (*[]models.Plan, error) {
	searchRadiusInKm := input.RadiusInKm
	if searchRadiusInKm < input.Category.SearchRadiusMinInKm {
		searchRadiusInKm = input.Category.SearchRadiusMinInKm
	} else if searchRadiusInKm < 5 {
		// 検索範囲が狭い場合は広めに検索する
		searchRadiusInKm += 5
	}

	placesOfCategory, err := s.placeRepository.FindByGooglePlaceType(
		ctx,
		input.Category.GooglePlaceTypes[0],
		input.Location,
		searchRadiusInKm*1000,
	)
	if err != nil {
		return nil, fmt.Errorf("error while fetching google Places: %v\n", err)
	}

	// 検索が大量に行われないようにするため、対象とする場所の数を20件に制限する
	if len(*placesOfCategory) > 20 {
		*placesOfCategory = (*placesOfCategory)[:20]
	}

	var createPlanParams []CreatePlanParams
	for _, placeOfCategory := range *placesOfCategory {
		if len(createPlanParams) >= 3 {
			break
		}

		placesNearby, err := s.placeSearchService.SearchNearbyPlaces(ctx, placesearch.SearchNearbyPlacesInput{
			Location: placeOfCategory.Location,
		})
		if err != nil {
			return nil, fmt.Errorf("error while fetching nearby places: %v\n", err)
		}

		planPlaces, err := s.CreatePlanPlaces(CreatePlanPlacesInput{
			PlanCandidateSetId: input.PlanCandidateSetId,
			LocationStart:      placeOfCategory.Location,
			PlaceStart:         placeOfCategory,
			Places:             placesNearby,
			PlacesOtherPlansContain: array.FlatMap(createPlanParams, func(p CreatePlanParams) []models.Place {
				return p.Places
			}),
		})
		if err != nil {
			return nil, fmt.Errorf("error while creating plan places: %v\n", err)
		}

		createPlanParams = append(createPlanParams, CreatePlanParams{
			LocationStart: placeOfCategory.Location,
			PlaceStart:    placeOfCategory,
			Places:        planPlaces,
		})
	}

	plans := s.createPlanData(ctx, input.PlanCandidateSetId, createPlanParams...)

	return &plans, nil
}
