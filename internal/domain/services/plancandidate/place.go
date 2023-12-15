package plancandidate

import (
	"context"
	"sort"

	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/placefilter"
)

// FetchCandidatePlaces はプランの候補となる場所を取得する
func (s Service) FetchCandidatePlaces(
	ctx context.Context,
	createPlanSessionId string,
	nLimit int,
) (*[]models.Place, error) {
	if nLimit <= 0 {
		panic("nLimit must be greater than 0")
	}

	placesSearched, err := s.placeService.FetchSearchedPlaces(ctx, createPlanSessionId)
	if err != nil {
		return nil, err
	}

	planCandidate, err := s.planCandidateRepository.Find(ctx, createPlanSessionId)
	if err != nil {
		return nil, err
	}

	placesFiltered := placesSearched
	placesFiltered = placefilter.FilterIgnoreCategory(placesFiltered)
	placesFiltered = placefilter.FilterByCategory(placesFiltered, models.GetCategoryToFilter(), true)

	// 重複した場所を削除
	placesFiltered = placefilter.FilterDuplicated(placesFiltered)

	placesSortedByRating := placesFiltered
	sort.Slice(placesSortedByRating, func(i, j int) bool {
		return placesSortedByRating[i].Google.Rating > placesSortedByRating[j].Google.Rating
	})

	placesToSuggest := make([]models.Place, 0, len(placesSortedByRating))
	for _, place := range placesSortedByRating {
		if planCandidate.HasPlace(place.Id) {
			continue
		}

		placesWithPhoto := s.placeService.FetchPlacesPhotosAndSave(ctx, place)
		place = placesWithPhoto[0]

		placesToSuggest = append(placesToSuggest, place)

		if len(placesToSuggest) >= nLimit {
			break
		}
	}

	return &placesToSuggest, nil
}
