package places

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/utils"
	"time"

	"googlemaps.github.io/maps"
)

type NearbySearchRequest struct {
	Location    Location
	Radius      uint
	Language    string
	Type        *maps.PlaceType
	SearchCount int
}

func (r PlacesApi) NearbySearch(ctx context.Context, req *NearbySearchRequest) ([]Place, error) {
	var placeType maps.PlaceType
	if req.Type != nil {
		placeType = *req.Type
	}

	placeSearchResults, err := r.nearBySearch(ctx, &maps.NearbySearchRequest{
		Location: &maps.LatLng{
			Lat: req.Location.Latitude,
			Lng: req.Location.Longitude,
		},
		Radius:   req.Radius,
		Language: req.Language,
		Type:     placeType,
	}, req.SearchCount)
	if err != nil {
		return nil, err
	}

	// Getting places nearby
	var places []Place
	for _, place := range placeSearchResults {
		places = append(places, createPlace(
			place.PlaceID,
			place.Name,
			place.Types,
			place.Geometry,
			place.Photos,
			place.OpeningHours != nil && place.OpeningHours.OpenNow != nil && *place.OpeningHours.OpenNow,
			place.Rating,
			place.UserRatingsTotal,
			utils.StrOmitEmpty(place.FormattedAddress),
			utils.StrOmitEmpty(place.Vicinity),
			place.PriceLevel,
		))
	}

	return places, nil
}

// nearBySearch Places API Nearby Search
// https://developers.google.com/maps/documentation/places/web-service/search-nearby
// pageCount は 1 以上の整数で、ページング処理を行う回数を指定する。
func (r PlacesApi) nearBySearch(ctx context.Context, req *maps.NearbySearchRequest, pagesCount int) ([]maps.PlacesSearchResult, error) {
	r.logger.Info(
		"Places API Nearby Search",
		zap.String("location", fmt.Sprintf("%f,%f", req.Location.Lat, req.Location.Lng)),
		zap.Uint("radius", req.Radius),
		zap.String("language", req.Language),
		zap.String("type", string(req.Type)),
		zap.Int("pagesCount", pagesCount),
	)
	if pagesCount < 1 {
		pagesCount = 1
	}

	var placeSearchResults []maps.PlacesSearchResult
	pageToken := ""
	for i := 0; i < pagesCount; i++ {
		if pageToken == "" {
			res, err := r.neaBySearchOnce(ctx, req)
			if err != nil {
				return placeSearchResults, err
			}
			placeSearchResults = append(placeSearchResults, res.Results...)
			pageToken = res.NextPageToken
			continue
		}

		// ノータイムでリクエストを送信すると、INVALID_REQUEST となってしまう。
		time.Sleep(2000 * time.Millisecond)

		res, err := r.nearBySearchWithPageToken(ctx, pageToken, req.Language)
		if err != nil {
			return placeSearchResults, err
		}
		placeSearchResults = append(placeSearchResults, res.Results...)
		pageToken = res.NextPageToken
	}

	return placeSearchResults, nil
}

func (r PlacesApi) neaBySearchOnce(ctx context.Context, req *maps.NearbySearchRequest) (*maps.PlacesSearchResponse, error) {
	res, err := r.mapsClient.NearbySearch(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error while nearby search: %v", err)
	}
	return &res, nil
}

func (r PlacesApi) nearBySearchWithPageToken(ctx context.Context, nextPageToken string, language string) (*maps.PlacesSearchResponse, error) {
	res, err := r.mapsClient.NearbySearch(ctx, &maps.NearbySearchRequest{
		PageToken: nextPageToken,
		Language:  language,
	})
	if err != nil {
		return nil, fmt.Errorf("error while nearby search with page token: %v", err)
	}
	return &res, nil
}
