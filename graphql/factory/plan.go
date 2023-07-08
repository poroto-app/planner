package factory

import (
	"log"

	graphql "poroto.app/poroto/planner/graphql/model"
	"poroto.app/poroto/planner/internal/domain/models"
)

func PlansFromDomainModel(plans *[]models.Plan) []*graphql.Plan {
	graphqlPlans := make([]*graphql.Plan, len(*plans))

	for i, plan := range *plans {
		graphqlPlan := PlanFromDomainModel(plan)
		graphqlPlans[i] = &graphqlPlan
	}

	return graphqlPlans
}

func PlanFromDomainModel(plan models.Plan) graphql.Plan {
	places := make([]*graphql.Place, len(plan.Places))
	for i, place := range plan.Places {
		places[i] = PlaceFromDomainModel(place)
	}

	transitions := make([]*graphql.Transition, len(plan.Transitions))
	for _, place := range plan.Places {
		nextPlace, duration, err := plan.GetTransition(place.Id)
		if err != nil {
			log.Println("error while getting transition: ", err)
			continue
		}

		if nextPlace == nil {
			continue
		}

		transitions = append(transitions, &graphql.Transition{
			From:     PlaceFromDomainModel(place),
			To:       PlaceFromDomainModel(*nextPlace),
			Duration: int(duration),
		})
	}

	return graphql.Plan{
		ID:            plan.Id,
		Name:          plan.Name,
		Places:        places,
		TimeInMinutes: int(plan.TimeInMinutes),
		Transitions:   transitions,
	}
}
