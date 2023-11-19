package places

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
)

type ImageSize struct {
	Width  uint
	Height uint
}

type ImageSizeType int

type PlacePhoto struct {
	PhotoReference string
	Small          *string
	Large          *string
}

type placePhotoWithSize struct {
	photoReference string
	imageUrl       string
	size           ImageSizeType
}

const (
	imgMaxHeightLarge = 1000
	imgMaxWidthLarge  = 1000
	imgMaxHeightSmall = 400
	imgMaxWidthSmall  = 400
)

const (
	ImageSizeTypeLarge ImageSizeType = iota
	ImageSizeTypeSmall
)

func ImageSizeLarge() ImageSize {
	return ImageSize{
		Width:  imgMaxWidthLarge,
		Height: imgMaxHeightLarge,
	}
}

func ImageSizeSmall() ImageSize {
	return ImageSize{
		Width:  imgMaxWidthSmall,
		Height: imgMaxHeightSmall,
	}
}

func (i ImageSizeType) ImageSize() ImageSize {
	switch i {
	case ImageSizeTypeLarge:
		return ImageSizeLarge()
	case ImageSizeTypeSmall:
		return ImageSizeSmall()
	default:
		panic(fmt.Sprintf("invalid image size type: %v", i))
	}
}

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
func (r PlacesApi) FetchPlacePhoto(photoReferences []string, imageSize ImageSize) (*PlacePhoto, error) {
	if len(photoReferences) == 0 {
		return nil, nil
	}

	imgUrl, err := imgUrlBuilder(imageSize.Width, imageSize.Height, photoReferences[0], r.apiKey)
	if err != nil {
		return nil, err
	}

	publicImageUrl, err := fetchPublicImageUrl(imgUrl)
	if err != nil {
		return nil, fmt.Errorf("error while fetching public image url: %w", err)
	}

	return &PlacePhoto{
		PhotoReference: photoReferences[0],
		Small:          publicImageUrl,
		Large:          publicImageUrl,
	}, nil
}

// FetchPlacePhotos は，指定された場所の写真を全件取得する
// imageSizeTypes が指定されている場合は，高画質の写真を取得する
// 画像取得は呼び出し料金が高いため、複数の場所の写真を取得するときは注意
// https://developers.google.com/maps/documentation/places/web-service/usage-and-billing?hl=ja#places-photo-new
func (r PlacesApi) FetchPlacePhotos(ctx context.Context, photoReferences []string, maxPhotoCount int, imageSizeTypes ...ImageSizeType) ([]PlacePhoto, error) {
	if len(imageSizeTypes) == 0 {
		imageSizeTypes = []ImageSizeType{ImageSizeTypeLarge}
	}

	ch := make(chan *placePhotoWithSize, len(photoReferences)*len(imageSizeTypes))
	for iPhoto, photoReference := range photoReferences {
		for _, imageSizeType := range imageSizeTypes {
			go func(ctx context.Context, photoIndex int, photoReference string, imageSizeType ImageSizeType, ch chan<- *placePhotoWithSize) {
				if photoIndex >= maxPhotoCount {
					ch <- nil
					return
				}

				imageSize := imageSizeType.ImageSize()

				imgUrl, err := imgUrlBuilder(imageSize.Width, imageSize.Height, photoReference, r.apiKey)
				if err != nil {
					log.Printf("skipping photoReference because of error while building image url: %v", err)
					ch <- nil
					return
				}

				log.Printf("Places API Fetch Place Photo: %s\n", photoReference)
				publicImageUrl, err := fetchPublicImageUrl(imgUrl)
				if err != nil {
					log.Printf("skipping photoReference because of error while fetching public image url: %v", err)
					ch <- nil
					return
				}

				ch <- &placePhotoWithSize{
					photoReference: photoReference,
					imageUrl:       *publicImageUrl,
					size:           imageSizeType,
				}
			}(ctx, iPhoto, photoReference, imageSizeType, ch)
		}
	}

	var placePhotoWithSizes []*placePhotoWithSize
	for i := 0; i < len(photoReferences)*len(imageSizeTypes); i++ {
		placePhotoWithSize := <-ch
		if placePhotoWithSize == nil {
			continue
		}
		placePhotoWithSizes = append(placePhotoWithSizes, placePhotoWithSize)
	}

	var placePhotos []PlacePhoto
	for _, photoReference := range photoReferences {
		var placePhoto PlacePhoto
		placePhoto.PhotoReference = photoReference

		for _, placePhotoWithSize := range placePhotoWithSizes {
			if placePhotoWithSize.photoReference != photoReference {
				continue
			}

			switch placePhotoWithSize.size {
			case ImageSizeTypeLarge:
				placePhoto.Large = &placePhotoWithSize.imageUrl
			case ImageSizeTypeSmall:
				placePhoto.Small = &placePhotoWithSize.imageUrl
			default:
				panic(fmt.Sprintf("invalid image size type: %v", placePhotoWithSize.size))
			}
		}

		if placePhoto.Large == nil && placePhoto.Small == nil {
			continue
		}

		placePhotos = append(placePhotos, placePhoto)
	}

	// すべての写真の取得に失敗した場合は、エラーを返す
	if len(photoReferences) > 0 && len(placePhotos) == 0 {
		return nil, fmt.Errorf("could not fetch any photos")
	}

	return placePhotos, nil
}
