package models

type PlanCollage struct {
	Images []PlanCollageImage
}

type PlanCollageImage struct {
	PlaceId  string
	ImageUrl string
}
