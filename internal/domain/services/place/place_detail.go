package place

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/factory"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

// FetchPlaceDetailAndSave Place Detail　情報を取得し、保存する
func (s Service) FetchPlaceDetailAndSave(ctx context.Context, planCandidateId string, googlePlace models.GooglePlace) (*models.GooglePlaceDetail, error) {
	if googlePlace.PlaceDetail != nil {
		return nil, nil
	}

	// キャッシュがある場合は取得する
	savedPlace, err := s.placeInPlanCandidateRepository.FindByGooglePlaceId(ctx, planCandidateId, googlePlace.PlaceId)
	if err != nil {
		return nil, fmt.Errorf("could not fetch google place detail: %v", err)
	}

	if savedPlace != nil && savedPlace.Google.PlaceDetail != nil {
		return savedPlace.Google.PlaceDetail, nil
	}

	placeDetailEntity, err := s.placesApi.FetchPlaceDetail(ctx, places.FetchPlaceDetailRequest{
		PlaceId:  googlePlace.PlaceId,
		Language: "ja",
	})
	if err != nil {
		return nil, err
	}

	if placeDetailEntity.PlaceDetail == nil {
		return nil, fmt.Errorf("could not fetch google place detail: %v", googlePlace.PlaceId)
	}

	placeDetail := factory.GooglePlaceDetailFromPlaceDetailEntity(*placeDetailEntity.PlaceDetail)

	// キャッシュする
	if err := s.placeInPlanCandidateRepository.SaveGooglePlaceDetail(ctx, planCandidateId, googlePlace.PlaceId, placeDetail); err != nil {
		return nil, fmt.Errorf("could not save google place detail: %v", err)
	}

	return &placeDetail, nil
}

// FetchPlacesDetailAndSave 複数の場所の Place Detail 情報を並行に取得し、保存する
func (s Service) FetchPlacesDetailAndSave(ctx context.Context, planCandidateId string, places []models.GooglePlace) []models.GooglePlace {
	if len(places) == 0 {
		return nil
	}

	ch := make(chan *models.GooglePlace, len(places))
	for _, place := range places {
		go func(ctx context.Context, place models.GooglePlace, ch chan<- *models.GooglePlace) {
			placeDetail, err := s.FetchPlaceDetailAndSave(ctx, planCandidateId, place)
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

		for i, place := range places {
			if placeWithPlaceDetail.PlaceId == place.PlaceId {
				places[i] = *placeWithPlaceDetail
			}
		}
	}

	return places
}
