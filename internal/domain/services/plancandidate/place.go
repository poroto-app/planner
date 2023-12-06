package plancandidate

import (
	"context"
	"go.uber.org/zap"
	"sort"

	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/placefilter"
	googleplaces "poroto.app/poroto/planner/internal/infrastructure/api/google/places"
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

		// TODO: キャッシュする
		// TODO: 大きいサイズの写真も取得する
		thumbnail, err := s.placesApi.FetchPlacePhoto([]models.GooglePlacePhotoReference{
			{
				PhotoReference: place.Google.PhotoReferences[0],
			},
		}, googleplaces.ImageSizeSmall())
		if err != nil {
			s.logger.Warn(
				"error while fetching place photo",
				zap.String("placeId", place.Id),
				zap.String("planCandidateId", createPlanSessionId),
				zap.Error(err),
			)
			continue
		}

		place.Google.Photos = &[]models.GooglePlacePhoto{*thumbnail}

		placesToSuggest = append(placesToSuggest, place)

		if len(placesToSuggest) >= nLimit {
			break
		}
	}

	return &placesToSuggest, nil
}
