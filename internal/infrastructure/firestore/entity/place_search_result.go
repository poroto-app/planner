package entity

import "poroto.app/poroto/planner/internal/infrastructure/api/google/places"

type PlaceSearchResultEntity struct {
	PlanCandidateId string `firestore:"plan_candidate_id"`
	Places          []places.Place
}
