package entity

import (
	"time"

	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

type PlaceSearchResultEntity struct {
	PlanCandidateId string         `firestore:"plan_candidate_id"`
	Places          []places.Place `firestore:"places"`
	CreatedAt       time.Time      `firestore:"created_at,omitempty,serverTimestamp"`
	UpdatedAt       time.Time      `firestore:"updated_at,omitempty"`
}
