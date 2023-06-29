package plan

import (
	"context"
	"fmt"

	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

// fetchPlacePhotos は，指定された場所の写真を取得する
func (s PlanService) fetchPlacePhotos(ctx context.Context, place places.Place) (thumbnailUrl *string, photoUrls []string, err error) {
	thumbnailPhoto, err := s.placesApi.FetchPlacePhoto(place, &places.ImageSize{
		Width:  places.ImgThumbnailMaxWidth,
		Height: places.ImgThumbnailMaxHeight,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error while fetching place thumbnail: %v\n", err)
	}

	if thumbnailPhoto != nil {
		thumbnailUrl = &thumbnailPhoto.ImageUrl
	}

	placePhotos, err := s.placesApi.FetchPlacePhotos(ctx, place)
	if err != nil {
		return nil, nil, fmt.Errorf("error while fetching place photos: %v\n", err)
	}

	for _, photo := range placePhotos {
		photoUrls = append(photoUrls, photo.ImageUrl)
	}

	return thumbnailUrl, photoUrls, nil
}
