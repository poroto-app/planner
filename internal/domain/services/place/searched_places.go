package place

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/models"
)

func (s Service) SaveSearchedPlaces(ctx context.Context, planCandidateId string, googlePlaces []models.GooglePlace) ([]models.Place, error) {
	type savePlaceFromGooglePlaceResult struct {
		place *models.Place
		err   error
	}

	// models.Google を保存し，models.Place を取得する
	chPlaces := make(chan savePlaceFromGooglePlaceResult, len(googlePlaces))
	for _, googlePlace := range googlePlaces {
		go func(googlePlace models.GooglePlace) {
			place, err := s.placeRepository.SavePlacesFromGooglePlace(ctx, googlePlace)
			if err != nil {
				chPlaces <- savePlaceFromGooglePlaceResult{
					place: nil,
					err:   fmt.Errorf("error while saving place: %v\n", err),
				}
				return
			}
			chPlaces <- savePlaceFromGooglePlaceResult{
				place: place,
				err:   err,
			}
		}(googlePlace)
	}

	var places []models.Place
	for range googlePlaces {
		result := <-chPlaces
		if result.err != nil {
			return nil, result.err
		}
		places = append(places, *result.place)
	}

	// PlanCandidate と検索された場所の紐付けを行う
	var placeIds []string
	for _, place := range places {
		placeIds = append(placeIds, place.Id)
	}
	if err := s.planCandidateRepository.AddSearchedPlacesForPlanCandidate(ctx, planCandidateId, placeIds); err != nil {
		return nil, fmt.Errorf("error while adding searched places for plan candidate: %v\n", err)
	}

	return places, nil
}

func (s Service) FetchSearchedPlaces(ctx context.Context, planCandidateId string) ([]models.Place, error) {
	places, err := s.placeRepository.FindByPlanCandidateId(ctx, planCandidateId)
	if err != nil {
		return nil, fmt.Errorf("error while fetching searched places for plan candidate: %v\n", err)
	}

	return places, nil
}
