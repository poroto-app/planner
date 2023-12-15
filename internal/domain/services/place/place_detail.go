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
