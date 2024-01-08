package place

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/factory"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

// FetchGooglePlace GooglePlace ID から場所の情報を取得する
// 過去に検索が行われている場合はキャッシュを利用する
// まだ検索が行われていない場合は、PlaceDetail APIを用いて場所の情報を取得し、保存する
func (s Service) FetchGooglePlace(ctx context.Context, googlePlaceId string) (*models.Place, error) {
	// キャッシュがある場合は取得する
	savedPlace, err := s.placeRepository.FindByGooglePlaceID(ctx, googlePlaceId)
	if err != nil {
		return nil, fmt.Errorf("could not fetch google place detail: %v", err)
	}

	if savedPlace != nil {
		return savedPlace, nil
	}

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

	// 保存する
	places, err := s.placeRepository.SavePlacesFromGooglePlaces(ctx, googlePlace)
	if err != nil {
		return nil, fmt.Errorf("could not save google place detail: %v", err)
	}
	if len(*places) == 0 {
		return nil, fmt.Errorf("could not save google place detail: %v", err)
	}

	return &(*places)[0], nil
}
