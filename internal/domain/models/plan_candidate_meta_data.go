package models

// PlanCandidateMetaData は PlanCandidate を作成するにあたって、元になった情報
type PlanCandidateMetaData struct {
	CreatedBasedOnCurrentLocation bool
	CategoriesPreferred           *[]LocationCategory
	CategoriesRejected            *[]LocationCategory
	LocationStart                 *GeoLocation
	FreeTime                      *int
}
