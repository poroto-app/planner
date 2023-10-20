package plancandidate

import (
	"context"
	"fmt"
	"log"
	"poroto.app/poroto/planner/internal/domain/models"
)

func (s Service) ReplacePlace(ctx context.Context, planCandidateId string, planId string, placeIdToBeReplaced string, placeIdToReplace string) (*models.Plan, error) {
	log.Printf("start fetching plan candidate: %v\n", planCandidateId)
	planCandidate, err := s.planCandidateRepository.Find(ctx, planCandidateId)
	if err != nil {
		return nil, fmt.Errorf("error while fetching plan candidate: %v\n", err)
	}
	log.Printf("succeeded fetching plan candidate: %v\n", planCandidateId)

	planToUpdate := planCandidate.GetPlan(planId)
	if planToUpdate == nil {
		return nil, fmt.Errorf("plan not found: %v\n", planId)
	}

	log.Printf("start searching places: %v\n", planCandidateId)
	placesSearched, err := s.placeSearchResultRepository.Find(ctx, planCandidateId)
	if err != nil {
		return nil, err
	}
	log.Printf("succeeded searching places: %v\n", planCandidateId)

	log.Printf("start fetching place to be replaced: %v\n", placeIdToBeReplaced)
	placeToBeReplaced := planToUpdate.GetPlace(placeIdToBeReplaced)
	if placeToBeReplaced == nil {
		return nil, fmt.Errorf("place to be replaced not found: %v\n", placeIdToBeReplaced)
	}
	log.Printf("succeeded fetching place to be replaced: %v\n", placeIdToBeReplaced)

	// 指定された場所がすでにプランに含まれている場合は何もしない
	if planToUpdate.GetPlace(placeIdToReplace) != nil {
		return nil, fmt.Errorf("place to replace already exists: %v\n", placeIdToReplace)
	}

	var placeToReplace *models.Place
	for _, place := range placesSearched {
		// TODO: PlaceRepositoryを用いて、Planner APIが指定したPlaceIdで取得できるようにする
		if place.PlaceId == placeIdToReplace {
			p := place.ToPlace()
			placeToReplace = &p
			break
		}
	}
	if placeToReplace == nil {
		return nil, fmt.Errorf("place to replace not found: %v\n", placeIdToReplace)
	}

	log.Printf("start replacing place: %v\n", placeIdToBeReplaced)
	if err := s.planCandidateRepository.ReplacePlace(ctx, planCandidateId, planId, placeIdToBeReplaced, *placeToReplace); err != nil {
		return nil, fmt.Errorf("error while replacing place: %v\n", err)
	}
	log.Printf("succeeded replacing place: %v\n", placeIdToBeReplaced)

	planCandidateUpdated, err := s.planCandidateRepository.Find(ctx, planCandidateId)
	if err != nil {
		return nil, fmt.Errorf("error while fetching plan candidate: %v\n", err)
	}

	planUpdated := planCandidateUpdated.GetPlan(planId)
	if planUpdated == nil {
		return nil, fmt.Errorf("plan not found: %v\n", planId)
	}

	return planUpdated, nil
}
