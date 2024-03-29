package places

import (
	"context"
	"go.uber.org/zap"
	"googlemaps.github.io/maps"
	"poroto.app/poroto/planner/internal/domain/utils"
)

type FetchPlaceDetailRequest struct {
	PlaceId  string
	Language string
}

// FetchPlaceDetail は IDを指定することで、対応する場所の情報を取得する
// 取得される内容は NearbySearch と同じ
func (r PlacesApi) FetchPlaceDetail(ctx context.Context, req FetchPlaceDetailRequest) (*Place, error) {
	r.logger.Info(
		"Places API Place Details",
		zap.String("placeId", req.PlaceId),
		zap.String("language", req.Language),
	)

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
			maps.PlaceDetailsFieldMaskReviews,
			maps.PlaceDetailsFieldMaskPhotos,
			maps.PlaceDetailsFieldMaskOpeningHours,
		},
	})
	if err != nil {
		return nil, err
	}

	place := createPlace(
		resp.PlaceID,
		resp.Name,
		resp.Types,
		resp.Geometry,
		resp.Photos,
		resp.OpeningHours != nil && resp.OpeningHours.OpenNow != nil && *resp.OpeningHours.OpenNow,
		resp.Rating,
		resp.UserRatingsTotal,
		utils.StrOmitEmpty(resp.FormattedAddress),
		utils.StrOmitEmpty(resp.Vicinity),
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
