package places

import (
	"context"
	"fmt"
	"googlemaps.github.io/maps"
)

type FetchPlaceRequest struct {
	PlaceId  string
	Language string
}

func (r PlacesApi) placeDetails(ctx context.Context, req *maps.PlaceDetailsRequest) (*maps.PlaceDetailsResult, error) {
	response, err := r.mapsClient.PlaceDetails(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error while requesting place details: %w", err)
	}

	return &response, nil
}

// FetchPlace は IDを指定することで、対応する場所の情報を取得する
// 取得される内容は FindPlacesFromLocation と同じ
func (r PlacesApi) FetchPlace(ctx context.Context, req FetchPlaceRequest) (*Place, error) {
	resp, err := r.mapsClient.PlaceDetails(ctx, &maps.PlaceDetailsRequest{
		PlaceID:  req.PlaceId,
		Language: req.Language,
		Fields: []maps.PlaceDetailsFieldMask{
			maps.PlaceDetailsFieldMaskName,
			maps.PlaceDetailsFieldMaskTypes,
			maps.PlaceDetailsFieldMaskGeometryLocation,
			maps.PlaceDetailsFieldMaskOpeningHours,
			maps.PlaceDetailsFieldMaskPhotos,
			maps.PlaceDetailsFieldMaskRatings,
		},
	})
	if err != nil {
		return nil, err
	}

	var photoReferences []string
	if resp.Photos != nil {
		for _, photo := range resp.Photos {
			photoReferences = append(photoReferences, photo.PhotoReference)
		}
	}
	place := createPlace(
		resp.PlaceID,
		resp.Name,
		resp.Types,
		resp.Geometry,
		photoReferences,
		resp.OpeningHours != nil && resp.OpeningHours.OpenNow != nil && *resp.OpeningHours.OpenNow,
		resp.Rating,
	)

	return &place, nil
}
