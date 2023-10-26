package plancandidate

import (
	"context"
	"log"
	"sort"

	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/placefilter"
	googleplaces "poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

const (
	maxAddablePlaces = 10
)

// FetchCandidatePlaces はプランの候補となる場所を取得する
func (s Service) FetchCandidatePlaces(
	ctx context.Context,
	createPlanSessionId string,
) (*[]models.Place, error) {
	placesSaved, err := s.placeInPlanCandidateRepository.FindByPlanCandidateId(ctx, createPlanSessionId)
	if err != nil {
		return nil, err
	}

	planCandidate, err := s.planCandidateRepository.Find(ctx, createPlanSessionId)
	if err != nil {
		return nil, err
	}

	placesFiltered := *placesSaved
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
		thumbnailImageUrl, err := s.placesApi.FetchPlacePhoto(place.Google.PhotoReferences, googleplaces.ImageSizeSmall())
		if err != nil {
			log.Printf("error while fetching place photo: %v\n", err)
			continue
		}

		image, err := models.NewImage(thumbnailImageUrl, nil)
		if err != nil {
			log.Printf("error while creating image: %v\n", err)
			continue
		}

		place.Google.Images = &[]models.Image{*image}

		placesToSuggest = append(placesToSuggest, place.ToPlace())

		if len(placesToSuggest) >= maxAddablePlaces {
			break
		}
	}

	return &placesToSuggest, nil
}
