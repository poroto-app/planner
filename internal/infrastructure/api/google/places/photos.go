package places

import (
	"context"
	"fmt"
)

type PlacePhoto struct {
	ImageUrl string
}

func (r PlacesApi) FetchPlacePhotos(ctx context.Context, place Place) ([]PlacePhoto, error) {
	const maxWidth int = 400
	const maxHeight int = 400
	const placePhotoApi string = "https://maps.googleapis.com/maps/api/place/photo?maxwidth=%d&maxheight=%d&photo_reference=%s&key=%s"
	var placePhotos []PlacePhoto

	for _, photoReference := range place.photoReferences {
		placePhotos = append(placePhotos, PlacePhoto{
			ImageUrl: fmt.Sprintf(placePhotoApi, maxWidth, maxHeight, photoReference, r.apiKey),
		})
	}

	return placePhotos, nil
}
