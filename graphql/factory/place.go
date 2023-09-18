package factory

import (
	graphql "poroto.app/poroto/planner/graphql/model"
	"poroto.app/poroto/planner/internal/domain/models"
)

func PlaceFromDomainModel(place *models.Place) *graphql.Place {
	if place == nil {
		return nil
	}

	var photos []*graphql.Image
	for _, photo := range place.Photos {
		photos = append(photos, &graphql.Image{
			Default: photo,
		})
	}

	return &graphql.Place{
		ID:     place.Id,
		Name:   place.Name,
		Photos: photos,
		Location: &graphql.GeoLocation{
			Latitude:  place.Location.Latitude,
			Longitude: place.Location.Longitude,
		},
		EstimatedStayDuration: int(place.EstimatedStayDuration),
	}
}
