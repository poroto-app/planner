package models

type PlaceInPlanCandidate struct {
	Id     string
	Google GooglePlace
}

func (p PlaceInPlanCandidate) Location() GeoLocation {
	return p.Google.Location
}

func (p PlaceInPlanCandidate) Categories() []LocationCategory {
	return GetCategoriesFromSubCategories(p.Google.Types)
}
