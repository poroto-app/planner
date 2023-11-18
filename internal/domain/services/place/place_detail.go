package place

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/factory"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

func (s Service) FetchPlaceDetail(ctx context.Context, place models.GooglePlace) (*models.GooglePlaceDetail, error) {
	if place.PlaceDetail != nil {
		return nil, nil
	}

	// TODO: キャッシュが有る場合は取得する

	placeDetailEntity, err := s.placesApi.FetchPlaceDetail(ctx, places.FetchPlaceDetailRequest{
		PlaceId:  place.PlaceId,
		Language: "ja",
	})
	if err != nil {
		return nil, err
	}

	if placeDetailEntity.PlaceDetail == nil {
		return nil, fmt.Errorf("could not fetch place detail: %v", place.PlaceId)
	}

	placeDetail := factory.GooglePlaceDetailFromPlaceDetailEntity(*placeDetailEntity.PlaceDetail)

	// TODO: キャッシュする

	return &placeDetail, nil
}
