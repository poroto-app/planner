package entity

import "poroto.app/poroto/planner/internal/domain/models"

type PlaceInPlanCandidateEntity struct {
	Id            string `firestore:"id"`
	GooglePlaceId string `firestore:"google_place_id"`
}

func ToPlaceInPlanCandidateEntity(placeInPlanCandidate models.PlaceInPlanCandidate) PlaceInPlanCandidateEntity {
	return PlaceInPlanCandidateEntity{
		Id:            placeInPlanCandidate.Id,
		GooglePlaceId: placeInPlanCandidate.Google.PlaceId,
	}
}

func FromPlaceInPlanCandidateEntity(entity PlaceInPlanCandidateEntity, googlePlace models.GooglePlace) models.PlaceInPlanCandidate {
	return models.PlaceInPlanCandidate{
		Id:     entity.Id,
		Google: googlePlace,
	}
}
