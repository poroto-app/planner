package place

import (
	"context"
	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
)

// FetchPlacesPhotosAndSave は，指定された場所の写真を一括で取得し，保存する
func (s Service) FetchPlacesPhotosAndSave(ctx context.Context, places ...models.Place) []models.Place {
	var googlePlaces []models.GooglePlace
	for _, place := range places {
		googlePlaces = append(googlePlaces, place.Google)
	}

	// 写真が取得されていない場所のみ、画像が保存されるようにする
	var googlePlaceIdsAlreadyHasImages []string
	for _, googlePlace := range googlePlaces {
		if googlePlace.Photos != nil && len(*googlePlace.Photos) > 0 {
			googlePlaceIdsAlreadyHasImages = append(googlePlaceIdsAlreadyHasImages, googlePlace.PlaceId)
		}
	}

	// 画像を取得
	googlePlaces = s.fetchGooglePlacesPhotos(ctx, googlePlaces)

	// 画像を保存
	for _, googlePlace := range googlePlaces {
		// すでに写真が取得済みの場合は何もしない
		alreadyHasImages := array.IsContain(googlePlaceIdsAlreadyHasImages, googlePlace.PlaceId)
		if alreadyHasImages {
			continue
		}

		if googlePlace.Photos == nil || len(*googlePlace.Photos) == 0 {
			continue
		}

		// 新しく画像を取得した場合は，保存する
		if err := s.placeRepository.SaveGooglePlacePhotos(ctx, googlePlace.PlaceId, *googlePlace.Photos); err != nil {
			continue
		}
	}

	for i, googlePlace := range googlePlaces {
		places[i].Google = googlePlace
	}

	return places
}

// fetchGooglePlacesPhotos は，指定された場所の写真を一括で取得する
// すでに写真がある場合は，何もしない
func (s Service) fetchGooglePlacesPhotos(ctx context.Context, places []models.GooglePlace) []models.GooglePlace {
	if len(places) == 0 {
		return places
	}

	ch := make(chan models.GooglePlace, len(places))
	for _, place := range places {
		go func(ctx context.Context, place models.GooglePlace, ch chan<- models.GooglePlace) {
			// すでに写真がある場合は，何もしない
			if place.Photos != nil && len(*place.Photos) > 0 {
				s.logger.Info(
					"skip fetching place photos because photos already exist",
					zap.String("placeId", place.PlaceId),
				)
				ch <- place
				return
			}

			if place.PlaceDetail == nil || len(place.PlaceDetail.PhotoReferences) == 0 {
				s.logger.Info(
					"skip fetching place photos because photo references not found",
					zap.String("placeId", place.PlaceId),
				)
				ch <- place
				return
			}

			photos, err := s.placesApi.FetchPlacePhotos(ctx, place.PlaceDetail.PhotoReferences, 1)
			if err != nil {
				// TODO: channelを用いてエラーハンドリングする
				s.logger.Warn(
					"error while fetching place photos",
					zap.String("placeId", place.PlaceId),
					zap.Error(err),
				)
				ch <- place
				return
			}

			place.Photos = &photos
			ch <- place
		}(ctx, place, ch)
	}

	for i := 0; i < len(places); i++ {
		placeUpdated := <-ch

		for i, place := range places {
			if place.PlaceId != placeUpdated.PlaceId {
				continue
			}
			places[i] = placeUpdated
			break
		}
	}

	return places
}
