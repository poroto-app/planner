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

func (s Service) FetchPlacesDetail(ctx context.Context, places []models.GooglePlace) []models.GooglePlace {
	if len(places) == 0 {
		return nil
	}

	ch := make(chan *models.GooglePlace, len(places))
	for _, place := range places {
		go func(ctx context.Context, place models.GooglePlace, ch chan<- *models.GooglePlace) {
			placeDetail, err := s.FetchPlaceDetail(ctx, place)
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
