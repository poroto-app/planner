package plan

import (
	"context"
	"fmt"

	"poroto.app/poroto/planner/internal/domain/models"
)

// fetchPlacePhotos は，指定された場所の写真を取得する
func (s PlanService) fetchPlacePhotos(ctx context.Context, placeId string) (thumbnailUrl *string, photoUrls []string, err error) {
	placePhotos, err := s.placesApi.FetchPlacePhotos(ctx, placeId)
	if err != nil {
		return nil, nil, fmt.Errorf("error while fetching place photos: %v\n", err)
	}

	for _, photo := range placePhotos {
		photoUrls = append(photoUrls, photo.ImageUrl)

		if thumbnailUrl == nil {
			thumbnailUrl = &photo.ImageUrl
		}
	}

	return thumbnailUrl, photoUrls, nil
}

// fetchPlacesPhotos は，指定された場所の写真を一括で取得する
func (s PlanService) fetchPlacesPhotos(ctx context.Context, places []models.Place) []models.Place {
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

			thumbnailUrl, photoUrls, err := s.fetchPlacePhotos(ctx, *place.GooglePlaceId)
			if err != nil {
				ch <- place
				return
			}

			if thumbnailUrl != nil {
				place.Thumbnail = thumbnailUrl
			}

			if len(photoUrls) > 0 {
				place.Photos = photoUrls
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
