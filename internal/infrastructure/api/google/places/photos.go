package places

import "context"

type PlacePhoto struct {
	ImageUrl string
}

func (r PlacesApi) FetchPlacePhotos(ctx context.Context, place Place) ([]PlacePhoto, error) {
	// TODO: implement me
	return []PlacePhoto{}, nil
}
