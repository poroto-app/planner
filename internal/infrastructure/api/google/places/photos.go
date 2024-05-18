package places

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"path"
	"poroto.app/poroto/planner/internal/domain/models"
)

type ImageSize struct {
	Width  uint
	Height uint
}

type PlacePhotoWithSize struct {
	photoReference models.GooglePlacePhotoReference
	image          models.Image
}

const (
	imgMaxHeight = 500
	imgMaxWidth  = 500
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

// FetchPlacePhoto は，指定された場所の画像を１件取得する
func (r PlacesApi) FetchPlacePhoto(photoReferences []models.GooglePlacePhotoReference) (*models.GooglePlacePhoto, error) {
	if len(photoReferences) == 0 {
		return nil, nil
	}

	photoReference := photoReferences[0]
	imgUrl, err := imgUrlBuilder(imgMaxWidth, imgMaxHeight, photoReference.PhotoReference, r.apiKey)
	if err != nil {
		return nil, err
	}

	publicImageUrl, err := fetchPublicImageUrl(imgUrl)
	if err != nil {
		return nil, fmt.Errorf("error while fetching public image url: %w", err)
	}

	image := models.Image{
		Width:          imgMaxWidth,
		Height:         imgMaxHeight,
		URL:            *publicImageUrl,
		IsGooglePhotos: true,
	}
	googlePhoto := photoReference.ToGooglePlacePhoto(&image, &image)
	return &googlePhoto, nil
}

// FetchPlacePhotos は，指定された場所の写真を全件取得する
// imageSizeTypes が指定されている場合は，高画質の写真を取得する
// 画像取得は呼び出し料金が高いため、複数の場所の写真を取得するときは注意
// https://developers.google.com/maps/documentation/places/web-service/usage-and-billing?hl=ja#places-photo-new
func (r PlacesApi) FetchPlacePhotos(ctx context.Context, photoReferences []models.GooglePlacePhotoReference, maxPhotoCount int) ([]models.GooglePlacePhoto, error) {
	ch := make(chan *PlacePhotoWithSize, len(photoReferences))
	for iPhoto, photoReference := range photoReferences {
		go func(ctx context.Context, photoIndex int, photoReference models.GooglePlacePhotoReference, ch chan<- *PlacePhotoWithSize) {
			// 画像取得数が上限に達した場合は、何もしない
			if photoIndex >= maxPhotoCount {
				ch <- nil
				return
			}

			// 画像サイズを指定（上限を超えている場合は、上限に合わせる）
			var imageSize ImageSize
			if photoReference.Width > imgMaxWidth || photoReference.Height > imgMaxHeight {
				imageSize = ImageSize{
					Width:  imgMaxWidth,
					Height: imgMaxHeight,
				}
			} else {
				imageSize = ImageSize{
					Width:  uint(photoReference.Width),
					Height: uint(photoReference.Height),
				}
			}

			imgUrl, err := imgUrlBuilder(imageSize.Width, imageSize.Height, photoReference.PhotoReference, r.apiKey)
			if err != nil {
				// TODO: channelにエラーを送信するようにする
				r.logger.Warn(
					"skipping photoReference because of error while building image url",
					zap.Error(err),
					zap.String("photoReference", photoReference.PhotoReference),
				)
				ch <- nil
				return
			}

			r.logger.Info(
				"Places API Fetch Place Photo",
				zap.String("photoReference", photoReference.PhotoReference),
			)
			publicImageUrl, err := fetchPublicImageUrl(imgUrl)
			if err != nil {
				// TODO: channelにエラーを送信するようにする
				r.logger.Warn(
					"skipping photoReference because of error while fetching public image url",
					zap.Error(err),
					zap.String("photoReference", photoReference.PhotoReference),
				)
				ch <- nil
				return
			}

			ch <- &PlacePhotoWithSize{
				photoReference: photoReference,
				image: models.Image{
					Width:          imageSize.Width,
					Height:         imageSize.Height,
					URL:            *publicImageUrl,
					IsGooglePhotos: true,
				},
			}
		}(ctx, iPhoto, photoReference, ch)
	}

	var placePhotos []models.GooglePlacePhoto
	for i := 0; i < len(photoReferences); i++ {
		placePhotoWithSize := <-ch
		if placePhotoWithSize == nil {
			continue
		}
		placePhotos = append(placePhotos, placePhotoWithSize.photoReference.ToGooglePlacePhoto(nil, &placePhotoWithSize.image))
	}

	// すべての写真の取得に失敗した場合は、エラーを返す
	if len(photoReferences) > 0 && len(placePhotos) == 0 {
		return nil, fmt.Errorf("could not fetch any photos")
	}

	return placePhotos, nil
}
