package models

// PlanCandidateMetaData は PlanCandidate を作成するにあたって、元になった情報
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
