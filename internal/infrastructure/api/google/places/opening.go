package places

import (
	"context"

	"googlemaps.github.io/maps"
)

type PlaceOpeningPeriod struct {
	DayOfWeek   string
	OpeningTime string
	ClosingTime string
}

func (r PlacesApi) FetchPlaceOpeningPeriods(ctx context.Context, place Place) ([]PlaceOpeningPeriod, error) {
	resp, err := r.mapsClient.PlaceDetails(ctx, &maps.PlaceDetailsRequest{
		PlaceID: place.PlaceID,
		Fields: []maps.PlaceDetailsFieldMask{
			maps.PlaceDetailsFieldMaskOpeningHours,
		},
	})
	if err != nil {
		return nil, err
	}

	var placeOpeningPeriods []PlaceOpeningPeriod
	for _, period := range resp.OpeningHours.Periods {
		placeOpeningPeriods = append(placeOpeningPeriods, PlaceOpeningPeriod{
			DayOfWeek:   period.Open.Day.String(),
			OpeningTime: period.Open.Time,
			ClosingTime: period.Close.Time,
		})
	}
	return placeOpeningPeriods, nil
}
