package place

import (
	"context"
	"fmt"
	"go.uber.org/zap"
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
func (s Service) FetchPlaceDetailAndSave(ctx context.Context, googlePlaceId string) (*models.GooglePlaceDetail, error) {
	// キャッシュがある場合は取得する
	savedPlace, err := s.placeRepository.FindByGooglePlaceID(ctx, googlePlaceId)
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
	if err := s.placeRepository.SaveGooglePlaceDetail(ctx, googlePlaceId, placeDetail); err != nil {
		return nil, fmt.Errorf("could not save google place detail: %v", err)
	}

	return &placeDetail, nil
}

// FetchPlacesDetailAndSave 複数の場所の Place Detail 情報を並行に取得し、保存する
func (s Service) FetchPlacesDetailAndSave(ctx context.Context, places []models.Place) []models.Place {
	if len(places) == 0 {
		return nil
	}

	googlePlaces := make([]models.GooglePlace, len(places))
	for i, place := range places {
		googlePlaces[i] = place.Google
	}

	// Place Detailを並行に取得する
	ch := make(chan *models.GooglePlace, len(googlePlaces))
	for _, googlePlace := range googlePlaces {
		go func(ctx context.Context, googlePlace models.GooglePlace, ch chan<- *models.GooglePlace) {
			if googlePlace.PlaceDetail != nil {
				s.logger.Info(
					"skip fetching place detail because place detail already exist",
					zap.String("placeId", googlePlace.PlaceId),
				)
				ch <- &googlePlace
				return
			}

			// 取得と保存を行う
			placeDetail, err := s.FetchPlaceDetailAndSave(ctx, googlePlace.PlaceId)
			if err != nil {
				ch <- nil
				return
			}

			googlePlace.PlaceDetail = placeDetail

			ch <- &googlePlace
		}(ctx, googlePlace, ch)
	}

	for i := 0; i < len(googlePlaces); i++ {
		placeWithPlaceDetail := <-ch
		if placeWithPlaceDetail == nil {
			continue
		}

		for iPlace, googlePlace := range googlePlaces {
			if placeWithPlaceDetail.PlaceId == googlePlace.PlaceId {
				googlePlaces[iPlace] = *placeWithPlaceDetail
			}
		}
	}

	// Place DetailとGoogle Placeを紐付ける
	for i, googlePlace := range googlePlaces {
		places[i].Google = googlePlace
	}

	return places
}
