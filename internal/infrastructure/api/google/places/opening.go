package places

import (
	"context"
	"log"

	"googlemaps.github.io/maps"
)

type PlaceOpeningPeriod struct {
	DayOfWeek   string
	OpeningTime string
	ClosingTime string
}

func (r PlacesApi) FetchPlaceOpeningPeriods(ctx context.Context, googlePlaceId string) ([]PlaceOpeningPeriod, error) {
	log.Printf("Places API Fetch Place Opening Periods: %+v\n", googlePlaceId)

	resp, err := r.mapsClient.PlaceDetails(ctx, &maps.PlaceDetailsRequest{
		PlaceID: googlePlaceId,
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
