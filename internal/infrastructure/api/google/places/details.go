package places

import (
	"context"
	"log"

	"googlemaps.github.io/maps"
)

type FetchPlaceRequest struct {
	PlaceId  string
	Language string
}

// FetchPlaceDetail は IDを指定することで、対応する場所の情報を取得する
// 取得される内容は FindPlacesFromLocation と同じ
func (r PlacesApi) FetchPlaceDetail(ctx context.Context, req FetchPlaceRequest) (*Place, error) {
	log.Println("Places API Place Details: ", req)

	resp, err := r.mapsClient.PlaceDetails(ctx, &maps.PlaceDetailsRequest{
		PlaceID:  req.PlaceId,
		Language: req.Language,
		Fields: []maps.PlaceDetailsFieldMask{
			maps.PlaceDetailsFieldMaskPlaceID,
			maps.PlaceDetailsFieldMaskName,
			maps.PlaceDetailsFieldMaskTypes,
			maps.PlaceDetailsFieldMaskGeometryLocation,
			maps.PlaceDetailsFieldMaskRatings,
			maps.PlaceDetailsFieldMaskUserRatingsTotal,
			maps.PlaceDetailsFieldMaskPriceLevel,
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
		resp.UserRatingsTotal,
		resp.PriceLevel,
	)

	placeDetail := createPlaceDetail(
		resp.Reviews,
		resp.Photos,
		resp.OpeningHours,
	)

	place.PlaceDetail = &placeDetail

	return &place, nil
}
