package place

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/placefilter"
	"poroto.app/poroto/planner/internal/domain/services/placesearch"
	"time"
)

const (
	defaultMaxPlacesToSuggest = 3
)

type FetchCandidatePlacesInput struct {
	PlanCandidateSetId string
	NLimit             int
}

// FetchCandidatePlaces はプランの候補となる場所を取得する
func (s Service) FetchCandidatePlaces(
	ctx context.Context,
	input FetchCandidatePlacesInput,
) (*[]models.Place, error) {
	if input.PlanCandidateSetId == "" {
		panic("PlanCandidateSetId is empty")
	}

	if input.NLimit == 0 {
		input.NLimit = defaultMaxPlacesToSuggest
	}

	planCandidateSet, err := s.planCandidateRepository.Find(ctx, input.PlanCandidateSetId, time.Now())
	if err != nil {
		return nil, err
	}

	if planCandidateSet.MetaData.LocationStart == nil {
		s.logger.Warn(
			"plan candidate set has no start location",
			zap.String("Id", planCandidateSet.Id),
		)

		return nil, nil
	}

	// 付近の場所を検索
	placesNearby, err := s.placeSearchService.SearchNearbyPlaces(ctx, placesearch.SearchNearbyPlacesInput{
		Location:           *planCandidateSet.MetaData.LocationStart,
		PlanCandidateSetId: &planCandidateSet.Id,
	})
	if err != nil {
		return nil, fmt.Errorf("error while fetching nearby places: %v\n", err)
	}

	// 検索された場所を保存
	placesFiltered := placesNearby
	placesFiltered = placefilter.FilterDefaultIgnore(placefilter.FilterDefaultIgnoreInput{
		Places:        placesFiltered,
		StartLocation: *planCandidateSet.MetaData.LocationStart,
	})

	placesSortedByRating := models.SortPlacesByRating(placesFiltered)

	placesToSuggest := make([]models.Place, 0, len(placesSortedByRating))
	for _, place := range placesSortedByRating {
		if planCandidateSet.HasPlace(place.Id) {
			continue
		}

		placesWithPhoto := s.placeSearchService.FetchPlacesPhotosAndSave(ctx, place)
		place = placesWithPhoto[0]

		placesToSuggest = append(placesToSuggest, place)

		if len(placesToSuggest) >= input.NLimit {
			break
		}
	}

	return &placesToSuggest, nil
}
