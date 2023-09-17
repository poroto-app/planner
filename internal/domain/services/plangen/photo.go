package plangen

import (
	"context"
	"poroto.app/poroto/planner/internal/domain/models"
)

// FetchPlacesPhotos は，指定された場所の写真を一括で取得すR
func (s Service) FetchPlacesPhotos(ctx context.Context, places []models.Place) []models.Place {
	if len(places) == 0 {
		return places
	}

	ch := make(chan models.Place, len(places))
	for _, place := range places {
		go func(ctx context.Context, place models.Place, ch chan<- models.Place) {
			if place.GooglePlaceId == nil {
				ch <- place
				return
			}

			placeThumbnails, placePhotos, err := s.placesApi.FetchPlacePhotos(ctx, *place.GooglePlaceId)
			if err != nil {
				ch <- place
				return
			}

			for _, photo := range placePhotos {
				place.Photos = append(place.Photos, photo.ImageUrl)
			}

			for _, photo := range placeThumbnails {
				place.Thumbnails = append(place.Thumbnails, photo.ImageUrl)
			}

			ch <- place
		}(ctx, place, ch)
	}

	for i := 0; i < len(places); i++ {
		placeUpdated := <-ch

		for i, place := range places {
			if place.Id != placeUpdated.Id {
				continue
			}

			places[i] = placeUpdated
			break
		}
	}

	return places
}
