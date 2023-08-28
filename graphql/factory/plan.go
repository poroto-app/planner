package factory

import (
	"fmt"
	"log"
	graphql "poroto.app/poroto/planner/graphql/model"
	"poroto.app/poroto/planner/internal/domain/models"
)

func PlansFromDomainModel(plans *[]models.Plan) []*graphql.Plan {
	graphqlPlans := make([]*graphql.Plan, 0)

	for _, plan := range *plans {
		graphqlPlan, err := PlanFromDomainModel(plan)
		if err != nil {
			log.Println("error while converting plan to graphql model: ", err)
			continue
		}
		graphqlPlans = append(graphqlPlans, graphqlPlan)
	}

	return graphqlPlans
}

func PlanFromDomainModel(plan models.Plan) (*graphql.Plan, error) {
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
				return nil, fmt.Errorf("could not find place %s in plan %s", *t.FromPlaceId, plan.Id)
			}
		}

		placeTo := plan.GetPlace(t.ToPlaceId)
		if placeTo == nil {
			return nil, fmt.Errorf("could not find place %s in plan %s", t.ToPlaceId, plan.Id)
		}

		transitions[i] = &graphql.Transition{
			From:     PlaceFromDomainModel(placeFrom),
			To:       PlaceFromDomainModel(placeTo),
			Duration: int(t.Duration),
		}
	}

	return &graphql.Plan{
		ID:            plan.Id,
		Name:          plan.Name,
		Places:        places,
		TimeInMinutes: int(plan.TimeInMinutes),
		Transitions:   transitions,
		AuthorID:      plan.AuthorId,
	}, nil
}
