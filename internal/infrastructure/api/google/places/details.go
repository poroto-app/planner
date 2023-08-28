package places

import (
	"context"
	"googlemaps.github.io/maps"
)

type FetchPlaceRequest struct {
	PlaceId  string
	Language string
}

// FetchPlace は IDを指定することで、対応する場所の情報を取得する
// 取得される内容は FindPlacesFromLocation と同じ
func (r PlacesApi) FetchPlace(ctx context.Context, req FetchPlaceRequest) (*Place, error) {
	resp, err := r.mapsClient.PlaceDetails(ctx, &maps.PlaceDetailsRequest{
		PlaceID:  req.PlaceId,
		Language: req.Language,
		Fields: []maps.PlaceDetailsFieldMask{
			maps.PlaceDetailsFieldMaskPlaceID,
			maps.PlaceDetailsFieldMaskName,
			maps.PlaceDetailsFieldMaskTypes,
			maps.PlaceDetailsFieldMaskGeometryLocation,
			maps.PlaceDetailsFieldMaskOpeningHours,
			maps.PlaceDetailsFieldMaskPhotos,
			maps.PlaceDetailsFieldMaskRatings,
			maps.PlaceDetailsFieldMaskUserRatingsTotal,
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
		resp.UserRatingsTotal,
	)

	return &place, nil
}
