package factory

import (
	"poroto.app/poroto/planner/internal/domain/models"
	graphql "poroto.app/poroto/planner/internal/interface/graphql/model"
)

func PlanCandidateFromDomainModel(planCandidate *models.PlanCandidate, startLocation *models.GeoLocation) *graphql.PlanCandidate {
	if planCandidate == nil {
		return nil
	}

	return &graphql.PlanCandidate{
		ID:                            planCandidate.Id,
		Plans:                         PlansFromDomainModel(&planCandidate.Plans, startLocation),
		LikedPlaceIds:                 planCandidate.LikedPlaceIds,
		CreatedBasedOnCurrentLocation: planCandidate.MetaData.CreatedBasedOnCurrentLocation,
	}
}
