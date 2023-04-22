package places

import (
	"context"
	"fmt"
	"net/url"
	"path"

	"googlemaps.github.io/maps"
)

type PlacePhoto struct {
	ImageUrl string
}

const (
	imgMaxHeight          = 1000
	imgMaxWidth           = 1000
	imgThumbnailMaxHeight = 400
	imgThumbnailMaxWidth  = 400
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

// FetchPlaceThumbnail は，指定された場所のサムネイル画像を１件取得する
// TODO: ImageUrlにAPIキーが含まれないように、リダイレクト先のURLを取得して返す
func (r PlacesApi) FetchPlaceThumbnail(place Place) (*PlacePhoto, error) {
	if len(place.photoReferences) == 0 {
		return nil, nil
	}

	imgUrl, err := imgUrlBuilder(imgThumbnailMaxWidth, imgThumbnailMaxHeight, place.photoReferences[0], r.apiKey)
	if err != nil {
		return nil, err
	}

	return &PlacePhoto{
		ImageUrl: imgUrl,
	}, nil
}

// FetchPlacePhotos は，指定された場所の写真を全件取得する
// TODO: ImageUrlにAPIキーが含まれないように、リダイレクト先のURLを取得して返す
func (r PlacesApi) FetchPlacePhotos(ctx context.Context, place Place) ([]PlacePhoto, error) {
	var placePhotos []PlacePhoto

	resp, err := r.mapsClient.PlaceDetails(ctx, &maps.PlaceDetailsRequest{
		PlaceID: place.PlaceID,
		Fields: []maps.PlaceDetailsFieldMask{
			maps.PlaceDetailsFieldMaskPhotos,
		},
	})
	if err != nil {
		return nil, err
	}

	for _, photo := range resp.Photos {
		imgUrl, err := imgUrlBuilder(imgMaxWidth, imgMaxHeight, photo.PhotoReference, r.apiKey)
		if err != nil {
			return nil, err
		}

		placePhotos = append(placePhotos, PlacePhoto{
			ImageUrl: imgUrl,
		})
	}

	return placePhotos, nil
}
