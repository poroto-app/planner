package places

import (
	"context"
	"fmt"
	"time"

	"googlemaps.github.io/maps"
)

func (r PlacesApi) nearBySearch(ctx context.Context, req *maps.NearbySearchRequest) ([]maps.PlacesSearchResult, error) {
	var placeSearchResults []maps.PlacesSearchResult

	pageToken := ""
	for i := 0; i < 3; i++ {
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

		res, err := r.nearBySearchWithPageToken(ctx, pageToken)
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

func (r PlacesApi) nearBySearchWithPageToken(ctx context.Context, nextPageToken string) (*maps.PlacesSearchResponse, error) {
	res, err := r.mapsClient.NearbySearch(ctx, &maps.NearbySearchRequest{
		PageToken: nextPageToken,
	})
	if err != nil {
		return nil, fmt.Errorf("error while nearby search with page token: %v", err)
	}
	return &res, nil
}
