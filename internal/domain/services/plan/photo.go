package plan

import (
	"context"
	"fmt"
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
