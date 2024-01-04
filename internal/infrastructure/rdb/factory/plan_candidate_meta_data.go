package factory

import (
	"fmt"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/entities"
)

func NewPlanCandidateMetaDataFromEntity(
	planCandidateSetMetaDataSlice entities.PlanCandidateSetMetaDatumSlice,
	planCandidateSetCategorySlice entities.PlanCandidateSetCategorySlice,
	planCandidateSetId string,
) (*models.PlanCandidateMetaData, error) {
	planCandidateSetMetaData, ok := array.Find(planCandidateSetMetaDataSlice, func(planCandidateSetMetaData *entities.PlanCandidateSetMetaDatum) bool {
		if planCandidateSetMetaData == nil {
			return false
		}
		return planCandidateSetMetaData.PlanCandidateSetID == planCandidateSetId
	})
	if !ok {
		return nil, fmt.Errorf("failed to find plan candidate set meta data")
	}

	planCandidateSetCategorySlice = array.Filter(planCandidateSetCategorySlice, func(planCandidateSetCategory *entities.PlanCandidateSetCategory) bool {
		if planCandidateSetCategory == nil {
			return false
		}
		return planCandidateSetCategory.PlanCandidateSetID == planCandidateSetId
	})

	var categoriesPreferred, categoriesRequired []string
	for _, planCandidateSetCategory := range planCandidateSetCategorySlice {
		if planCandidateSetCategory.IsSelected {
			categoriesPreferred = append(categoriesPreferred, planCandidateSetCategory.Category)
		} else {
			categoriesRequired = append(categoriesRequired, planCandidateSetCategory.Category)
		}
	}

	locationCategoriesPreferred := models.GetCategoriesFromSubCategories(categoriesPreferred)
	locationCategoriesRejected := models.GetCategoriesFromSubCategories(categoriesRequired)

	return &models.PlanCandidateMetaData{
		CreatedBasedOnCurrentLocation: planCandidateSetMetaData.IsCreatedFromCurrentLocation,
		CategoriesPreferred:           &locationCategoriesPreferred,
		CategoriesRejected:            &locationCategoriesRejected,
		LocationStart: &models.GeoLocation{
			Latitude:  planCandidateSetMetaData.LatitudeStart,
			Longitude: planCandidateSetMetaData.LongitudeStart,
		},
		FreeTime: planCandidateSetMetaData.PlanDurationMinutes.Ptr(),
	}, nil
}
