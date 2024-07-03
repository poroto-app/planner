package factory

import (
	"github.com/google/uuid"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
)

func NewPlanCandidateMetaDataCreateByCategoryFromDomainModel(planCandidateSetId string, createByPlanMetaData models.CreateByCategoryMetaData) generated.PlanCandidateSetMetaDataCreateByCategory {
	return generated.PlanCandidateSetMetaDataCreateByCategory{
		ID:                 uuid.New().String(),
		PlanCandidateSetID: planCandidateSetId,
		CategoryID:         createByPlanMetaData.Category.Id,
		Latitude:           createByPlanMetaData.Location.Latitude,
		Longitude:          createByPlanMetaData.Location.Longitude,
		RangeInMeters:      int(createByPlanMetaData.RadiusInKm * 1000),
	}
}

func newPlanCandidateSetMetaDataCreateByCategoryFromEntry(planCandidateSetMetaDataCreateByCategorySlice generated.PlanCandidateSetMetaDataCreateByCategorySlice) *models.CreateByCategoryMetaData {
	if len(planCandidateSetMetaDataCreateByCategorySlice) == 0 {
		return nil
	}

	category, found := array.Find(
		array.FlatMap(
			models.GetAllLocationCategorySetCreatePlan(),
			func(categorySetCreatePlan models.LocationCategorySetCreatePlan) []models.LocationCategoryCreatePlan {
				return categorySetCreatePlan.Categories
			},
		),
		func(category models.LocationCategoryCreatePlan) bool {
			return category.Id == planCandidateSetMetaDataCreateByCategorySlice[0].CategoryID
		},
	)
	if !found {
		return nil
	}

	return &models.CreateByCategoryMetaData{
		Category: category,
		Location: models.GeoLocation{
			Latitude:  planCandidateSetMetaDataCreateByCategorySlice[0].Latitude,
			Longitude: planCandidateSetMetaDataCreateByCategorySlice[0].Longitude,
		},
		RadiusInKm: float64(planCandidateSetMetaDataCreateByCategorySlice[0].RangeInMeters) / 1000,
	}
}
