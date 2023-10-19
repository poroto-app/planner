package places

import (
	"context"

	"googlemaps.github.io/maps"
)

func (r PlacesApi) FetchPlacePriceLevelRequest(ctx context.Context, googlePlaceId string) (*int, error) {
	resp, err := r.mapsClient.PlaceDetails(ctx, &maps.PlaceDetailsRequest{
		PlaceID: googlePlaceId,
		Fields: []maps.PlaceDetailsFieldMask{
			maps.PlaceDetailsFieldMaskPriceLevel,
		},
	})
	if err != nil {
		return nil, err
	}

	return &resp.PriceLevel, nil
}
