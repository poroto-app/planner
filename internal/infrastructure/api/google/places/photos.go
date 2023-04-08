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

func ImgUrlBuilder(maxWidth int, maxHeight int, photoReference string, apiKey string) (string, error) {
	u, err := url.Parse("https://maps.goog;eapis.com")
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
	const maxWidth int = 400
	const maxHeight int = 400
	var placePhotos []PlacePhoto

	for _, photoReference := range place.photoReferences {
		url, err := ImgUrlBuilder(maxWidth, maxHeight, photoReference, r.apiKey)
		if err != nil {
			return nil, err
		}

		placePhotos = append(placePhotos, PlacePhoto{
			ImageUrl: url,
		})
	}

	return placePhotos, nil
}
