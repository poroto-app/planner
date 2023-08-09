package places

import (
	"context"
	"fmt"
	"log"
	"net/http"
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

// fetchPublicImageUrl は、Place Photos API によって提供される公開可能なURLを取得する
// imgUrlBuilder が生成するURLは、APIキーを含むため、この関数によってリダイレクト先のURLを取得する必要がある
func fetchPublicImageUrl(photoUrl string) (*string, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("GET", photoUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("error while creating request: %w", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while requesting: %w", err)
	}

	publicImageUrl := res.Header.Get("Location")
	return &publicImageUrl, nil
}

// FetchPlacePhoto は，指定された場所のサムネイル画像を１件取得する
// imageSize が nilの場合は、最大1000x1000の画像を取得する
func (r PlacesApi) FetchPlacePhoto(place Place, imageSize *ImageSize) (*PlacePhoto, error) {
	if len(place.PhotoReferences) == 0 {
		return nil, nil
	}

	if imageSize == nil {
		imageSize = &ImageSize{
			Width:  imgMaxWidth,
			Height: imgMaxHeight,
		}
	}

	imgUrl, err := imgUrlBuilder(imageSize.Width, imageSize.Height, place.PhotoReferences[0], r.apiKey)
	if err != nil {
		return nil, err
	}

	publicImageUrl, err := fetchPublicImageUrl(imgUrl)
	if err != nil {
		return nil, fmt.Errorf("error while fetching public image url: %w", err)
	}

	return &PlacePhoto{
		ImageUrl: *publicImageUrl,
	}, nil
}

// FetchPlacePhotos は，指定された場所の写真を全件取得する
func (r PlacesApi) FetchPlacePhotos(ctx context.Context, placeId string) ([]PlacePhoto, error) {
	resp, err := r.mapsClient.PlaceDetails(ctx, &maps.PlaceDetailsRequest{
		PlaceID: placeId,
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
			log.Printf("skipping photo because of error while building image url: %v", err)
			continue
		}

		publicImageUrl, err := fetchPublicImageUrl(imgUrl)
		if err != nil {
			log.Printf("skipping photo because of error while fetching public image url: %v", err)
			continue
		}

		placePhotos = append(placePhotos, PlacePhoto{
			ImageUrl: *publicImageUrl,
		})
	}

	// すべての写真の取得に失敗した場合は、エラーを返す
	if len(resp.Photos) > 0 && len(placePhotos) == 0 {
		return nil, fmt.Errorf("could not fetch any photos")
	}

	return placePhotos, nil
}
