package factory

import (
	graphql "poroto.app/poroto/planner/graphql/model"
	"poroto.app/poroto/planner/internal/domain/models"
)

func PlansFromDomainModel(plans *[]models.Plan) []*graphql.Plan {
	graphqlPlans := make([]*graphql.Plan, len(*plans))

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

	for i, plan := range *plans {
		graphqlPlan := PlanFromDomainModel(plan)
		graphqlPlans[i] = &graphqlPlan
	}

	return graphqlPlans
}

func PlanFromDomainModel(plan models.Plan) graphql.Plan {
	places := make([]*graphql.Place, len(plan.Places))
	for i, place := range plan.Places {
		places[i] = &graphql.Place{
			Name:   place.Name,
			Photos: place.Photos,
			Location: &graphql.GeoLocation{
				Latitude:  place.Location.Latitude,
				Longitude: place.Location.Longitude,
			},
			EstimatedStayDuration: int(place.EstimatedStayDuration),
		}
	}
	return graphql.Plan{
		ID:            plan.Id,
		Name:          plan.Name,
		Places:        places,
		TimeInMinutes: int(plan.TimeInMinutes),
	}
}
