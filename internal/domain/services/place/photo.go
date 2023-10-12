package place

import (
	"context"
	"poroto.app/poroto/planner/internal/domain/array"
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

// FetchPlacesPhotosAndSave は，指定された場所の写真を一括で取得し，保存する
func (s Service) FetchPlacesPhotosAndSave(ctx context.Context, planCandidateId string, places []models.Place) []models.Place {
	// 写真が取得されていない場所のみ、画像が保存されるようにする
	var googlePlaceIdsWithPhotos []string
	for _, place := range places {
		if place.GooglePlaceId == nil {
			continue
		}

		if array.IsContain(googlePlaceIdsWithPhotos, *place.GooglePlaceId) {
			continue
		}

		googlePlaceIdsWithPhotos = append(googlePlaceIdsWithPhotos, *place.GooglePlaceId)
	}

	// 画像を取得
	places = s.FetchPlacesPhotos(ctx, places)

	// 画像を保存
	for _, place := range places {
		if place.GooglePlaceId == nil {
			continue
		}

		// すでに写真が取得済みの場合は何もしない
		if !array.IsContain(googlePlaceIdsWithPhotos, *place.GooglePlaceId) {
			continue
		}

		if err := s.placeSearchResultRepository.SaveImagesIfNotExist(ctx, planCandidateId, *place.GooglePlaceId, place.Images); err != nil {
			continue
		}
	}

	return places
}
