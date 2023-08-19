package models

import "time"

type PlanCandidate struct {
	Id                            string
	Plans                         []Plan
	CreatedBasedOnCurrentLocation bool
	MetaData                      PlanCandidateMetaData
	ExpiresAt                     time.Time
}

func (p PlanCandidate) HasPlace(googlePlaceId string) bool {
	for _, plan := range p.Plans {
		for _, place := range plan.Places {
			if place.GooglePlaceId == nil {
				continue
			}

			if googlePlaceId == *place.GooglePlaceId {
				return true
			}
		}
	}
	return false
}
