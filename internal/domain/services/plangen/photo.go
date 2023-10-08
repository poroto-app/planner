package plangen

import (
	"context"
	"poroto.app/poroto/planner/internal/domain/models"
	api "poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

// FetchPlacesPhotos は，指定された場所の写真を一括で取得する
// すでに写真がある場合は，何もしない
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

			// すでに写真がある場合は，何もしない
			if place.Images != nil && len(place.Images) > 0 {
				ch <- place
				return
			}

			photos, err := s.placesApi.FetchPlacePhotos(
				ctx,
				*place.GooglePlaceId,
				api.ImageSizeTypeSmall,
				api.ImageSizeTypeLarge,
			)
			if err != nil {
				ch <- place
				return
			}

			images := make([]models.Image, 0, len(photos))
			for _, photo := range photos {
				image, err := models.NewImage(photo.Small, photo.Large)
				if err != nil {
					continue
				}

				images = append(images, *image)
			}

			place.Images = images

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
