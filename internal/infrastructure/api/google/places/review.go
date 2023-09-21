package places

import (
	"context"
	"googlemaps.github.io/maps"
)

type FetchPlaceReviewRequest struct {
	PlaceId  string
	Language string
}

func (r PlacesApi) FetchPlaceReview(ctx context.Context, req FetchPlaceReviewRequest) (*[]maps.PlaceReview, error) {
	resp, err := r.mapsClient.PlaceDetails(ctx, &maps.PlaceDetailsRequest{
		PlaceID:  req.PlaceId,
		Language: req.Language,
		Fields: []maps.PlaceDetailsFieldMask{
			maps.PlaceDetailsFieldMaskReviews,
		},
	})
	if err != nil {
		return nil, err
	}

	return &resp.Reviews, nil
}
