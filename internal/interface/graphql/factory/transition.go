package factory

import (
	"fmt"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	graphql "poroto.app/poroto/planner/internal/interface/graphql/model"
)

func TransitionFromDomainModel(transition models.Transition, places []models.Place) (*graphql.Transition, error) {
	var placeFrom *models.Place
	if transition.FromPlaceId != nil {
		p, found := array.Find(places, func(p models.Place) bool {
			return p.Id == *transition.FromPlaceId
		})
		if !found {
			return nil, fmt.Errorf("could not find place %s", *transition.FromPlaceId)
		}
		placeFrom = &p
	}

	placeTo, found := array.Find(places, func(p models.Place) bool {
		return p.Id == transition.ToPlaceId
	})
	if !found {
		return nil, fmt.Errorf("could not find place %s", transition.ToPlaceId)
	}

	return &graphql.Transition{
		From:     PlaceFromDomainModel(placeFrom),
		To:       PlaceFromDomainModel(&placeTo),
		Duration: int(transition.Duration),
	}, nil
}
