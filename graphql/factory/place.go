package factory

import (
	graphql "poroto.app/poroto/planner/graphql/model"
	"poroto.app/poroto/planner/internal/domain/models"
)

func PlaceFromDomainModel(place *models.Place) *graphql.Place {
	if place == nil {
		return nil
	}

	return &graphql.Place{
		ID:            place.Id,
		GooglePlaceID: place.GooglePlaceId,
		Name:          place.Name,
		Photos:        place.Photos,
		Thumbnails:    place.Thumbnails,
		Location: &graphql.GeoLocation{
			Latitude:  place.Location.Latitude,
			Longitude: place.Location.Longitude,
		},
		EstimatedStayDuration: int(place.EstimatedStayDuration),
	}
}
