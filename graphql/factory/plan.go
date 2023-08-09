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
		places[i] = PlaceFromDomainModel(&place)
	}

	transitions := make([]*graphql.Transition, len(plan.Transitions))
	for i, t := range plan.Transitions {
		var placeFrom *models.Place
		if t.FromPlaceId != nil {
			placeFrom = plan.GetPlace(*t.FromPlaceId)
			if placeFrom == nil {
				log.Printf("could not find place %s in plan", *t.FromPlaceId)
				continue
			}
		}

		placeTo := plan.GetPlace(t.ToPlaceId)
		if placeTo == nil {
			log.Printf("could not find place %s in plan", t.ToPlaceId)
			continue
		}

		transitions[i] = &graphql.Transition{
			From:     PlaceFromDomainModel(placeFrom),
			To:       PlaceFromDomainModel(placeTo),
			Duration: int(t.Duration),
		}
	}

	return graphql.Plan{
		ID:            plan.Id,
		Name:          plan.Name,
		Places:        places,
		TimeInMinutes: int(plan.TimeInMinutes),
		Transitions:   transitions,
	}
}
