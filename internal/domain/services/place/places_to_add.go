package place

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/services/placefilter"
	"time"
)

const defaultMaxPlacesToAdd = 4

type FetchPlacesToAddInput struct {
	PlanCandidateId string
	PlanId          string
	NLimit          uint
}

// FetchPlacesToAdd はプランに追加する候補となる場所一覧を取得する
// nLimit によって取得する場所の数を制限することができる
func (s Service) FetchPlacesToAdd(ctx context.Context, input FetchPlacesToAddInput) ([]models.Place, error) {
	if input.NLimit == 0 {
		input.NLimit = defaultMaxPlacesToAdd
	}

	if input.PlanCandidateId == "" {
		return nil, fmt.Errorf("plan candidate id is empty")
	}

	if input.PlanId == "" {
		return nil, fmt.Errorf("plan id is empty")
	}

	planCandidate, err := s.planCandidateRepository.Find(ctx, input.PlanCandidateId, time.Now())
	if err != nil {
		return nil, fmt.Errorf("error while fetching plan candidate: %v", err)
	}

	var plan *models.Plan
	for _, p := range planCandidate.Plans {
		if p.Id == input.PlanId {
			plan = &p
			break
		}
	}
	if plan == nil {
		return nil, fmt.Errorf("plan not found")
	}

	if len(plan.Places) == 0 {
		return nil, fmt.Errorf("plan has no places")
	}

	startPlace := plan.Places[0]

	placesSearched, err := s.placeSearchService.FetchSearchedPlaces(ctx, input.PlanCandidateId)
	if err != nil {
		return nil, fmt.Errorf("error while fetching places searched: %v", err)
	}

	placesFiltered := placesSearched
	placesFiltered = placefilter.FilterDefaultIgnore(placefilter.FilterDefaultIgnoreInput{
		Places:        placesFiltered,
		StartLocation: startPlace.Location,
	})

	// すでにプランに含まれている場所を除外する
	placesFiltered = placefilter.FilterPlaces(placesFiltered, func(place models.Place) bool {
		for _, placeInPlan := range plan.Places {
			if placeInPlan.Id == place.Id {
				return false
			}
		}
		return true
	})

	// レビューの高い順でソート
	placesFiltered = models.SortPlacesByRating(placesFiltered)

	// TODO: すべてのカテゴリの場所が表示されるようにする
	var placesToAdd []models.Place
	for _, place := range placesFiltered {
		if len(placesToAdd) >= int(input.NLimit) {
			break
		}

		placesToAdd = append(placesToAdd, place)
	}

	// 写真を取得
	placesToAdd = s.placeSearchService.FetchPlacesPhotosAndSave(ctx, placesToAdd...)

	return placesToAdd, nil
}
