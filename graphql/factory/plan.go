package factory

import (
	"poroto.app/poroto/planner/graphql/model"
	"poroto.app/poroto/planner/internal/domain/models"
)

func PlansFromDomainModel(plans *[]models.Plan) []*model.Plan {
	graphqlPlans := make([]*model.Plan, 0)

	for _, plan := range *plans {
		places := make([]*model.Place, 0)
		for _, place := range plan.Places {
			places = append(places, &model.Place{
				ID:     place.Id,
				Name:   place.Name,
				Photos: place.Photos,
				Location: &model.GeoLocation{
					Latitude:  place.Location.Latitude,
					Longitude: place.Location.Longitude,
				},
				EstimatedStayDuration: int(place.EstimatedStayDuration),
			})
		}

		graphqlPlans = append(graphqlPlans, &model.Plan{
			ID:            plan.Id,
			Name:          plan.Name,
			Places:        places,
			TimeInMinutes: int(plan.TimeInMinutes),
		})
	}

	return graphqlPlans
}
