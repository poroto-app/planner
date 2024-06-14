package factory

import (
	"poroto.app/poroto/planner/internal/domain/models"
	graphql "poroto.app/poroto/planner/internal/interface/graphql/model"
)

func PlanCandidateSetFromDomainModel(planCandidateSet *models.PlanCandidateSet) *graphql.PlanCandidate {
	if planCandidateSet == nil {
		return nil
	}

	return &graphql.PlanCandidate{
		ID:                            planCandidateSet.Id,
		Plans:                         PlansFromDomainModel(&planCandidateSet.Plans, planCandidateSet.MetaData.GetLocationStart()),
		LikedPlaceIds:                 planCandidateSet.LikedPlaceIds,
		CreatedBasedOnCurrentLocation: planCandidateSet.MetaData.CreatedBasedOnCurrentLocation,
	}
}
