package place

import (
	"context"
	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
)

// FetchGooglePlacesPhotos は，指定された場所の写真を一括で取得する
// すでに写真がある場合は，何もしない
func (s Service) FetchGooglePlacesPhotos(ctx context.Context, places []models.GooglePlace) []models.GooglePlace {
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

// FetchGooglePlacesPhotosAndSave は，指定された場所の写真を一括で取得し，保存する
// 事前に FetchPlaceDetailAndSave で models.GooglePlaceDetail を取得しておく必要がある
func (s Service) FetchGooglePlacesPhotosAndSave(ctx context.Context, planCandidateId string, places ...models.GooglePlace) []models.GooglePlace {
	// 写真が取得されていない場所のみ、画像が保存されるようにする
	var googlePlaceIdsAlreadyHasImages []string
	for _, place := range places {
		if place.Photos != nil && len(*place.Photos) > 0 {
			googlePlaceIdsAlreadyHasImages = append(googlePlaceIdsAlreadyHasImages, place.PlaceId)
		}
	}

	// 画像を取得
	places = s.FetchGooglePlacesPhotos(ctx, places)

	// 画像を保存
	for _, place := range places {
		// すでに写真が取得済みの場合は何もしない
		alreadyHasImages := array.IsContain(googlePlaceIdsAlreadyHasImages, place.PlaceId)
		if alreadyHasImages {
			continue
		}

		if place.Photos == nil || len(*place.Photos) == 0 {
			continue
		}

		if err := s.placeRepository.SaveGooglePlacePhotos(ctx, place.PlaceId, *place.Photos); err != nil {
			continue
		}
	}

	return places
}

func (s Service) FetchPlacesPhotosAndSave(ctx context.Context, planCandidateId string, places ...models.Place) []models.Place {
	googlePlaces := make([]models.GooglePlace, len(places))
	for i, place := range places {
		googlePlaces[i] = place.Google
	}

	googlePlaces = s.FetchGooglePlacesPhotosAndSave(ctx, planCandidateId, googlePlaces...)

	for i, googlePlace := range googlePlaces {
		places[i].Google = googlePlace
	}

	return places
}
