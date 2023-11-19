package place

import (
	"context"
	"log"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	api "poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

// FetchPlacesPhotos は，指定された場所の写真を一括で取得する
// すでに写真がある場合は，何もしない
func (s Service) FetchPlacesPhotos(ctx context.Context, places []models.GooglePlace) []models.GooglePlace {
	if len(places) == 0 {
		return places
	}

	ch := make(chan models.GooglePlace, len(places))
	for _, place := range places {
		go func(ctx context.Context, place models.GooglePlace, ch chan<- models.GooglePlace) {
			// すでに写真がある場合は，何もしない
			if place.Images != nil && len(*place.Images) > 0 {
				log.Printf("skip fetching place photos because images already exist: %v\n", place.PlaceId)
				ch <- place
				return
			}

			photoReferences := make([]string, len(place.PlaceDetail.PhotoReferences))
			for i, photoReference := range place.PlaceDetail.PhotoReferences {
				photoReferences[i] = photoReference.PhotoReference
			}

			photos, err := s.placesApi.FetchPlacePhotos(
				ctx,
				photoReferences,
				1,
				api.ImageSizeTypeSmall,
				api.ImageSizeTypeLarge,
			)
			if err != nil {
				log.Printf("error while fetching place photos: %v\n", err)
				ch <- place
				return
			}

			var images []models.Image
			for _, photo := range photos {
				image, err := models.NewImage(photo.Small, photo.Large)
				if err != nil {
					log.Printf("error while creating image: %v\n", err)
					continue
				}

				images = append(images, *image)
			}

			place.Images = &images
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

// FetchPlacesPhotosAndSave は，指定された場所の写真を一括で取得し，保存する
// 事前に FetchPlaceDetail で models.GooglePlaceDetail を取得しておく必要がある
func (s Service) FetchPlacesPhotosAndSave(ctx context.Context, planCandidateId string, places ...models.GooglePlace) []models.GooglePlace {
	// 写真が取得されていない場所のみ、画像が保存されるようにする
	var googlePlaceIdsAlreadyHasImages []string
	for _, place := range places {
		if place.Images != nil && len(*place.Images) > 0 {
			googlePlaceIdsAlreadyHasImages = append(googlePlaceIdsAlreadyHasImages, place.PlaceId)
		}
	}

	// 画像を取得
	places = s.FetchPlacesPhotos(ctx, places)

	// 画像を保存
	for _, place := range places {
		// すでに写真が取得済みの場合は何もしない
		alreadyHasImages := array.IsContain(googlePlaceIdsAlreadyHasImages, place.PlaceId)
		if alreadyHasImages {
			continue
		}

		if place.Images == nil || len(*place.Images) == 0 {
			continue
		}

		if err := s.placeInPlanCandidateRepository.SaveGoogleImages(ctx, planCandidateId, place.PlaceId, *place.Images); err != nil {
			continue
		}
	}

	return places
}

func (s Service) FetchPlacesInPlanCandidatePhotosAndSave(ctx context.Context, planCandidateId string, places ...models.PlaceInPlanCandidate) []models.PlaceInPlanCandidate {
	googlePlaces := make([]models.GooglePlace, len(places))
	for i, place := range places {
		googlePlaces[i] = place.Google
	}

	googlePlaces = s.FetchPlacesPhotosAndSave(ctx, planCandidateId, googlePlaces...)

	for i, googlePlace := range googlePlaces {
		places[i].Google = googlePlace
	}

	return places
}
