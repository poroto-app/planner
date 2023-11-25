package factory

import (
	"fmt"
	"log"

	graphql "poroto.app/poroto/planner/graphql/model"
	"poroto.app/poroto/planner/internal/domain/models"
)

func PlansFromDomainModel(plans *[]models.Plan, startLocation *models.GeoLocation) []*graphql.Plan {
	graphqlPlans := make([]*graphql.Plan, 0)

	for _, plan := range *plans {
		graphqlPlan, err := PlanFromDomainModel(plan, startLocation)
		if err != nil {
			log.Println("error while converting plan to graphql model: ", err)
			continue
		}
		graphqlPlans = append(graphqlPlans, graphqlPlan)
	}

	return graphqlPlans
}

func PlanFromDomainModel(plan models.Plan, startLocation *models.GeoLocation) (*graphql.Plan, error) {
	places := make([]*graphql.Place, len(plan.Places))
	for i, place := range plan.Places {
		places[i] = PlaceFromDomainModel(&place)
	}

	transitions := plan.Transitions(startLocation)
	graphqlTransitionEntities := make([]*graphql.Transition, len(plan.Transitions(startLocation)))
	for i, t := range transitions {
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

		graphqlTransitionEntities[i] = &graphql.Transition{
			From:     PlaceFromDomainModel(placeFrom),
			To:       PlaceFromDomainModel(placeTo),
			Duration: int(t.Duration),
		}
	}

	return &graphql.Plan{
		ID:     plan.Id,
		Name:   plan.Name,
		Places: places,
		// MOCK: transitionを利用し時間を計算する関数を実装する
		TimeInMinutes: 0,
		Transitions:   graphqlTransitionEntities,
		AuthorID:      plan.AuthorId,
	}, nil
}
