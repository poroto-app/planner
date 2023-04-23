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

type ImageSize struct {
	Width  uint
	Height uint
}

const (
	imgMaxHeight          = 1000
	imgMaxWidth           = 1000
	ImgThumbnailMaxHeight = 400
	ImgThumbnailMaxWidth  = 400
)

func imgUrlBuilder(maxWidth uint, maxHeight uint, photoReference string, apiKey string) (string, error) {
	u, err := url.Parse("https://maps.googleapis.com")
	if err != nil {
		return "", err
	}

	u.Path = path.Join(u.Path, "maps", "api", "place", "photo")
	q := u.Query()

	q.Set("maxwidth", fmt.Sprint(maxWidth))
	q.Set("maxheight", fmt.Sprint(maxHeight))
	q.Set("photo_reference", photoReference)
	q.Set("key", apiKey)

	u.RawQuery = q.Encode()
	return u.String(), nil
}

// FetchPlacePhoto は，指定された場所のサムネイル画像を１件取得する
// TODO: ImageUrlにAPIキーが含まれないように、リダイレクト先のURLを取得して返す
func (r PlacesApi) FetchPlacePhoto(place Place, imageSize *ImageSize) (*PlacePhoto, error) {
	if len(place.photoReferences) == 0 {
		return nil, nil
	}

	if imageSize == nil {
		imageSize = &ImageSize{
			Width:  imgMaxWidth,
			Height: imgMaxHeight,
		}
	}

	imgUrl, err := imgUrlBuilder(imageSize.Width, imageSize.Height, place.photoReferences[0], r.apiKey)
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
	resp, err := r.mapsClient.PlaceDetails(ctx, &maps.PlaceDetailsRequest{
		PlaceID: place.PlaceID,
		Fields: []maps.PlaceDetailsFieldMask{
			maps.PlaceDetailsFieldMaskPhotos,
		},
	})
	if err != nil {
		return nil, err
	}

	var placePhotos []PlacePhoto
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
