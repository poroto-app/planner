package models

import "time"

type PlanCandidate struct {
	Id            string
	Plans         []Plan
	MetaData      PlanCandidateMetaData
	ExpiresAt     time.Time
	LikedPlaceIds []string
}

func (p PlanCandidate) HasPlace(placeId string) bool {
	for _, plan := range p.Plans {
		for _, place := range plan.Places {
			if placeId == place.Id {
				return true
			}
		}
	}
	return false
}

func (p PlanCandidate) GetPlan(planId string) *Plan {
	for _, plan := range p.Plans {
		if plan.Id == planId {
			return &plan
		}
	}
	return nil
}
