package place

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/factory"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

// FetchGooglePlace PlaceDetail APIを用いて場所の情報を取得する
func (s Service) FetchGooglePlace(ctx context.Context, googlePlaceId string) (*models.GooglePlace, error) {
	placeDetailEntity, err := s.placesApi.FetchPlaceDetail(ctx, places.FetchPlaceDetailRequest{
		PlaceId:  googlePlaceId,
		Language: "ja",
	})
	if err != nil {
		return nil, err
	}

	if placeDetailEntity == nil {
		return nil, fmt.Errorf("could not fetch google place detail: %v", googlePlaceId)
	}

	googlePlace := factory.GooglePlaceFromPlaceEntity(*placeDetailEntity, nil)

	return &googlePlace, nil
}

// FetchPlaceDetailAndSave Place Detail　情報を取得し、保存する
func (s Service) FetchPlaceDetailAndSave(ctx context.Context, planCandidateId string, googlePlaceId string) (*models.GooglePlaceDetail, error) {
	// キャッシュがある場合は取得する
	savedPlace, err := s.placeInPlanCandidateRepository.FindByGooglePlaceId(ctx, planCandidateId, googlePlaceId)
	if err != nil {
		return nil, fmt.Errorf("could not fetch google place detail: %v", err)
	}

	if savedPlace != nil && savedPlace.Google.PlaceDetail != nil {
		return savedPlace.Google.PlaceDetail, nil
	}

	placeDetailEntity, err := s.placesApi.FetchPlaceDetail(ctx, places.FetchPlaceDetailRequest{
		PlaceId:  googlePlaceId,
		Language: "ja",
	})
	if err != nil {
		return nil, err
	}

	if placeDetailEntity.PlaceDetail == nil {
		return nil, fmt.Errorf("could not fetch google place detail: %v", googlePlaceId)
	}

	placeDetail := factory.GooglePlaceDetailFromPlaceDetailEntity(*placeDetailEntity.PlaceDetail)

	// キャッシュする
	if err := s.placeInPlanCandidateRepository.SaveGooglePlaceDetail(ctx, planCandidateId, googlePlaceId, placeDetail); err != nil {
		return nil, fmt.Errorf("could not save google place detail: %v", err)
	}

	return &placeDetail, nil
}

// FetchGooglePlacesDetailAndSave 複数の場所の Place Detail 情報を並行に取得し、保存する
func (s Service) FetchGooglePlacesDetailAndSave(ctx context.Context, planCandidateId string, places []models.GooglePlace) []models.GooglePlace {
	if len(places) == 0 {
		return nil
	}

	ch := make(chan *models.GooglePlace, len(places))
	for _, place := range places {
		go func(ctx context.Context, place models.GooglePlace, ch chan<- *models.GooglePlace) {
			if place.PlaceDetail != nil {
				ch <- &place
				return
			}

			placeDetail, err := s.FetchPlaceDetailAndSave(ctx, planCandidateId, place.PlaceId)
			if err != nil {
				ch <- nil
				return
			}

			place.PlaceDetail = placeDetail

			ch <- &place
		}(ctx, place, ch)
	}

	for i := 0; i < len(places); i++ {
		placeWithPlaceDetail := <-ch
		if placeWithPlaceDetail == nil {
			continue
		}

		for iPlace, place := range places {
			if placeWithPlaceDetail.PlaceId == place.PlaceId {
				places[iPlace] = *placeWithPlaceDetail
			}
		}
	}

	return places
}

func (s Service) FetchPlacesDetailAndSave(ctx context.Context, planCandidateId string, places []models.Place) []models.Place {
	googlePlaces := make([]models.GooglePlace, len(places))
	for i, place := range places {
		googlePlaces[i] = place.Google
	}

	googlePlaces = s.FetchGooglePlacesDetailAndSave(ctx, planCandidateId, googlePlaces)

	for i, googlePlace := range googlePlaces {
		places[i].Google = googlePlace
	}

	return places
}
