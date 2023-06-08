package models

import "time"

type PlanCandidate struct {
	Id                            string
	Plans                         []Plan
	CreatedBasedOnCurrentLocation bool
	ExpiresAt                     time.Time
}
