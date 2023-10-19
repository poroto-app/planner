package place

import (
	"context"
	"log"

	"poroto.app/poroto/planner/internal/domain/models"
)

// FetchPriceLevel は、プランに含まれるすべての場所の価格帯を一括で取得する
func (s Service) FetchPriceLevel(ctx context.Context, places []models.GooglePlace) []models.GooglePlace {
	ch := make(chan *models.GooglePlace, len(places))
	for _, place := range places {
		go func(ctx context.Context, place models.GooglePlace, ch chan<- *models.GooglePlace) {
			// すでに価格帯がある場合は何もしない
			if place.PriceLevel != nil {
				ch <- &place
				return
			}

			priceLevel, err := s.placesApi.FetchPlacePriceLevelRequest(ctx, place.PlaceId)
			if err != nil {
				log.Printf("error while fetching place price level: %v\n", err)
				ch <- nil
				return
			}
			place.PriceLevel = priceLevel
			ch <- &place
		}(ctx, place, ch)
	}

	for i := 0; i < len(places); i++ {
		placeUpdated := <-ch
		if placeUpdated == nil {
			continue
		}

		for j, place := range places {
			if placeUpdated.PlaceId == place.PlaceId {
				places[j] = *placeUpdated
				break
			}
		}
	}

	return places
}
