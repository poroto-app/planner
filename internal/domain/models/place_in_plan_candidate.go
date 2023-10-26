package models

import "poroto.app/poroto/planner/internal/domain/utils"

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

func (p PlaceInPlanCandidate) EstimatedStayDuration() uint {
	if len(p.Categories()) == 0 {
		return 0
	}
	return p.Categories()[0].EstimatedStayDuration
}

func (p PlaceInPlanCandidate) IsSameCategoryPlace(other PlaceInPlanCandidate) bool {
	for _, categoryOfA := range p.Categories() {
		for _, categoryOfB := range other.Categories() {
			if categoryOfA.Name == categoryOfB.Name {
				return true
			}
		}
	}
	return false
}

func (p PlaceInPlanCandidate) ToPlace() Place {
	if p.Google.Images == nil {
		p.Google.Images = new([]Image)
	}

	return Place{
		Id:                 p.Id,
		GooglePlaceId:      utils.StrPointer(p.Google.PlaceId),
		Name:               p.Google.Name,
		Location:           p.Google.Location,
		Images:             *p.Google.Images,
		Categories:         GetCategoriesFromSubCategories(p.Google.Types),
		GooglePlaceReviews: p.Google.Reviews,
		PriceLevel:         p.Google.PriceLevel,
	}
}
