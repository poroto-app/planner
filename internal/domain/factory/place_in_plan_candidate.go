package factory

import "poroto.app/poroto/planner/internal/domain/models"

func PlaceInPlanCandidateFromGooglePlace(id string, googlePlace models.GooglePlace) models.PlaceInPlanCandidate {
	return models.PlaceInPlanCandidate{
		Id:     id,
		Google: googlePlace,
	}
}
