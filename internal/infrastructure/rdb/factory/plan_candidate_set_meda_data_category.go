package factory

import (
	"github.com/google/uuid"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
)

func NewPlanCandidateSetMetaDataCategorySliceFromDomainModel(
	locationCategoriesPreferred *[]models.LocationCategory,
	locationCategoriesRejected *[]models.LocationCategory,
	planCandidateSetId string,
) generated.PlanCandidateSetMetaDataCategorySlice {
	var planCandidateSetMetaDataCategorySlice generated.PlanCandidateSetMetaDataCategorySlice

	if locationCategoriesPreferred != nil {
		for _, locationCategory := range *locationCategoriesPreferred {
			planCandidateSetMetaDataCategoryEntity := newPlanCandidateSetMetaDataCategoryEntityFromDomainModel(locationCategory, true, planCandidateSetId)
			planCandidateSetMetaDataCategorySlice = append(planCandidateSetMetaDataCategorySlice, &planCandidateSetMetaDataCategoryEntity)
		}
	}

	if locationCategoriesRejected != nil {
		for _, locationCategory := range *locationCategoriesRejected {
			planCandidateSetMetaDataCategoryEntity := newPlanCandidateSetMetaDataCategoryEntityFromDomainModel(locationCategory, false, planCandidateSetId)
			planCandidateSetMetaDataCategorySlice = append(planCandidateSetMetaDataCategorySlice, &planCandidateSetMetaDataCategoryEntity)
		}
	}

	return planCandidateSetMetaDataCategorySlice
}

func newPlanCandidateSetMetaDataCategoryEntityFromDomainModel(locationCategory models.LocationCategory, isSelected bool, planCandidateSetId string) generated.PlanCandidateSetMetaDataCategory {
	return generated.PlanCandidateSetMetaDataCategory{
		ID:                 uuid.New().String(),
		PlanCandidateSetID: planCandidateSetId,
		Category:           locationCategory.Name,
		IsSelected:         isSelected,
	}
}

func newPlanCandidateSetMetaDataPreferredCategoriesFromEntity(planCandidateSetMetaDataCategorySlice generated.PlanCandidateSetMetaDataCategorySlice, planCandidateSetId string) *[]models.LocationCategory {
	planCandidateSetMetaDataPreferredCategorySlice := array.Filter(planCandidateSetMetaDataCategorySlice, func(planCandidateSetCategory *generated.PlanCandidateSetMetaDataCategory) bool {
		if planCandidateSetCategory == nil {
			return false
		}
		if !planCandidateSetCategory.IsSelected {
			return false
		}
		return planCandidateSetCategory.PlanCandidateSetID == planCandidateSetId
	})

	categoriesPreferred := array.MapAndFilter(planCandidateSetMetaDataPreferredCategorySlice, func(planCandidateSetCategory *generated.PlanCandidateSetMetaDataCategory) (string, bool) {
		if planCandidateSetCategory == nil {
			return "", false
		}
		return planCandidateSetCategory.Category, true
	})

	if len(categoriesPreferred) == 0 {
		return nil
	}

	locationCategoriesPreferred := models.GetCategoriesFromSubCategories(categoriesPreferred)

	return &locationCategoriesPreferred
}

func newPlanCandidateSetMetaDataRejectedCategoriesFromEntity(planCandidateSetMetaDataCategorySlice generated.PlanCandidateSetMetaDataCategorySlice, planCandidateSetId string) *[]models.LocationCategory {
	planCandidateSetMetaDataRejectedCategorySlice := array.Filter(planCandidateSetMetaDataCategorySlice, func(planCandidateSetCategory *generated.PlanCandidateSetMetaDataCategory) bool {
		if planCandidateSetCategory == nil {
			return false
		}
		if planCandidateSetCategory.IsSelected {
			return false
		}
		return planCandidateSetCategory.PlanCandidateSetID == planCandidateSetId
	})

	categoriesRequired := array.MapAndFilter(planCandidateSetMetaDataRejectedCategorySlice, func(planCandidateSetCategory *generated.PlanCandidateSetMetaDataCategory) (string, bool) {
		if planCandidateSetCategory == nil {
			return "", false
		}
		return planCandidateSetCategory.Category, true
	})

	if len(categoriesRequired) == 0 {
		return nil
	}

	locationCategoriesRejected := models.GetCategoriesFromSubCategories(categoriesRequired)

	return &locationCategoriesRejected
}
