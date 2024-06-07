package models

import "time"

type PlanCandidateSet struct {
	Id              string
	Plans           []Plan
	MetaData        PlanCandidateMetaData
	IsPlaceSearched bool
	ExpiresAt       time.Time
	LikedPlaceIds   []string
}

func (p PlanCandidateSet) HasPlace(placeId string) bool {
	for _, plan := range p.Plans {
		for _, place := range plan.Places {
			if placeId == place.Id {
				return true
			}
		}
	}
	return false
}

func (p PlanCandidateSet) GetPlan(planId string) *Plan {
	for _, plan := range p.Plans {
		if plan.Id == planId {
			return &plan
		}
	}
	return nil
}
