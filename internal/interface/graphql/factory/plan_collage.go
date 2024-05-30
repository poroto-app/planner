package factory

import (
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	graphql "poroto.app/poroto/planner/internal/interface/graphql/model"
)

func PlanCollageFromDomainModel(planCollage *models.PlanCollage) *graphql.PlanCollage {
	if planCollage == nil {
		return &graphql.PlanCollage{}
	}

	graphqlPlanCollageImages := array.Map(planCollage.Images, func(image models.PlanCollageImage) *graphql.PlanCollageImage {
		return &graphql.PlanCollageImage{
			PlaceID: image.PlaceId,
			Image:   ImageFromDomainModel(&image.Image),
		}
	})

	return &graphql.PlanCollage{
		Images: graphqlPlanCollageImages,
	}
}
