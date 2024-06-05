package factory

import (
	"poroto.app/poroto/planner/internal/domain/models"
	graphql "poroto.app/poroto/planner/internal/interface/graphql/model"
)

func PlanCandidateFromDomainModel(planCandidate *models.PlanCandidateSet) *graphql.PlanCandidate {
	if planCandidate == nil {
		return nil
	}

	return &graphql.PlanCandidate{
		ID:                            planCandidate.Id,
		Plans:                         PlansFromDomainModel(&planCandidate.Plans, planCandidate.MetaData.LocationStart),
		LikedPlaceIds:                 planCandidate.LikedPlaceIds,
		CreatedBasedOnCurrentLocation: planCandidate.MetaData.CreatedBasedOnCurrentLocation,
	}
}
