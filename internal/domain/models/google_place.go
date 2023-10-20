package models

import (
	"poroto.app/poroto/planner/internal/domain/utils"
)

type GooglePlace struct {
	PlaceId          string
	Name             string
	Types            []string
	Location         GeoLocation
	PhotoReferences  []string
	OpenNow          bool
	Rating           float32
	UserRatingsTotal int
	Images           *[]Image
	Reviews          *[]GooglePlaceReview
}

func (g GooglePlace) ToPlace() Place {
	if g.Images == nil {
		g.Images = new([]Image)
	}

	// TODO: planner api が生成したIDと対応させる
	return Place{
		Id:                 g.PlaceId,
		GooglePlaceId:      utils.StrPointer(g.PlaceId),
		Name:               g.Name,
		Location:           g.Location,
		Images:             *g.Images,
		Categories:         GetCategoriesFromSubCategories(g.Types),
		GooglePlaceReviews: g.Reviews,
	}
}

func (g GooglePlace) EstimatedStayDuration() uint {
	categories := GetCategoriesFromSubCategories(g.Types)
	if len(categories) == 0 {
		return 0
	}
	return categories[0].EstimatedStayDuration
}
