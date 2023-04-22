package places

import (
	"context"
	"fmt"
	"net/url"
	"path"
)

type PlacePhoto struct {
	ImageUrl string
}

const (
	imgMaxHeight = 1000
	imgMaxWidth  = 1000
)

func imgUrlBuilder(maxWidth int, maxHeight int, photoReference string, apiKey string) (string, error) {
	u, err := url.Parse("https://maps.googleapis.com")
	if err != nil {
		return "", err
	}

	u.Path = path.Join(u.Path, "maps", "api", "place", "photo")
	q := u.Query()

	q.Set("maxwidth", fmt.Sprint(maxWidth))
	q.Set("maxHeight", fmt.Sprint(maxHeight))
	q.Set("photo_reference", photoReference)
	q.Set("key", apiKey)

	u.RawQuery = q.Encode()
	return u.String(), nil
}

func (r PlacesApi) FetchPlacePhotos(ctx context.Context, place Place) ([]PlacePhoto, error) {
	var placePhotos []PlacePhoto

	for _, photoReference := range place.photoReferences {
		url, err := imgUrlBuilder(imgMaxWidth, imgMaxHeight, photoReference, r.apiKey)
		if err != nil {
			return nil, err
		}

		placePhotos = append(placePhotos, PlacePhoto{
			ImageUrl: url,
		})
	}

	return placePhotos, nil
}
