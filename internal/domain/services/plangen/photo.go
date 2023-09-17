package plangen

import (
	"context"
	"poroto.app/poroto/planner/internal/domain/models"
	api "poroto.app/poroto/planner/internal/infrastructure/api/google/places"
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

			photos, err := s.placesApi.FetchPlacePhotos(ctx, *place.GooglePlaceId, api.ImageSizeLarge(), api.ImageSizeThumbnail())
			if err != nil {
				ch <- place
				return
			}

			for _, photo := range photos {
				if photo.ImageSize.Same(api.ImageSizeLarge()) {
					place.Photos = append(place.Photos, photo.ImageUrl)
				} else if photo.ImageSize.Same(api.ImageSizeThumbnail()) {
					place.Thumbnails = append(place.Thumbnails, photo.ImageUrl)
				}
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
