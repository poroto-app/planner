package models

// PlanCandidateMetaData は PlanCandidateSet を作成するにあたって、元になった情報
type PlanCandidateMetaData struct {
	CreatedBasedOnCurrentLocation bool
	CategoriesPreferred           *[]LocationCategory
	CategoriesRejected            *[]LocationCategory
	LocationStart                 *GeoLocation
	FreeTime                      *int
}

func (p PlanCandidateMetaData) IsZero() bool {
	return p.CategoriesPreferred == nil &&
		p.CategoriesRejected == nil &&
		p.LocationStart == nil &&
		p.FreeTime == nil
}

func (p PlanCandidateMetaData) GetLocationStart() *GeoLocation {
	if p.CreatedBasedOnCurrentLocation {
		return p.LocationStart
	}
	return nil
}
