package factory

import (
	"github.com/google/uuid"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/entities"
)

func NewPlanCandidateSetMetaDataCategorySliceFromDomainModel(
	locationCategoriesPreferred *[]models.LocationCategory,
	locationCategoriesRejected *[]models.LocationCategory,
	planCandidateSetId string,
) entities.PlanCandidateSetMetaDataCategorySlice {
	var planCandidateSetMetaDataCategorySlice entities.PlanCandidateSetMetaDataCategorySlice
	for _, locationCategory := range *locationCategoriesPreferred {
		planCandidateSetMetaDataCategoryEntity := newPlanCandidateSetMetaDataCategoryEntityFromDomainModel(locationCategory, true, planCandidateSetId)
		planCandidateSetMetaDataCategorySlice = append(planCandidateSetMetaDataCategorySlice, &planCandidateSetMetaDataCategoryEntity)
	}
	for _, locationCategory := range *locationCategoriesRejected {
		planCandidateSetMetaDataCategoryEntity := newPlanCandidateSetMetaDataCategoryEntityFromDomainModel(locationCategory, false, planCandidateSetId)
		planCandidateSetMetaDataCategorySlice = append(planCandidateSetMetaDataCategorySlice, &planCandidateSetMetaDataCategoryEntity)
	}
	return planCandidateSetMetaDataCategorySlice
}

func newPlanCandidateSetMetaDataCategoryEntityFromDomainModel(locationCategory models.LocationCategory, isSelected bool, planCandidateSetId string) entities.PlanCandidateSetMetaDataCategory {
	return entities.PlanCandidateSetMetaDataCategory{
		ID:                 uuid.New().String(),
		PlanCandidateSetID: planCandidateSetId,
		Category:           locationCategory.Name,
		IsSelected:         isSelected,
	}
}

func newPlanCandidateSetMetaDataPreferredCategoriesFromEntity(planCandidateSetMetaDataCategorySlice entities.PlanCandidateSetMetaDataCategorySlice, planCandidateSetId string) *[]models.LocationCategory {
	planCandidateSetMetaDataPreferredCategorySlice := array.Filter(planCandidateSetMetaDataCategorySlice, func(planCandidateSetCategory *entities.PlanCandidateSetMetaDataCategory) bool {
		if planCandidateSetCategory == nil {
			return false
		}
		if !planCandidateSetCategory.IsSelected {
			return false
		}
		return planCandidateSetCategory.PlanCandidateSetID == planCandidateSetId
	})

	var categoriesPreferred []string
	for _, planCandidateSetCategory := range planCandidateSetMetaDataPreferredCategorySlice {
		categoriesPreferred = append(categoriesPreferred, planCandidateSetCategory.Category)
	}

	if len(categoriesPreferred) == 0 {
		return nil
	}

	locationCategoriesPreferred := models.GetCategoriesFromSubCategories(categoriesPreferred)

	return &locationCategoriesPreferred
}

func newPlanCandidateSetMetaDataRejectedCategoriesFromEntity(planCandidateSetMetaDataCategorySlice entities.PlanCandidateSetMetaDataCategorySlice, planCandidateSetId string) *[]models.LocationCategory {
	planCandidateSetMetaDataRejectedCategorySlice := array.Filter(planCandidateSetMetaDataCategorySlice, func(planCandidateSetCategory *entities.PlanCandidateSetMetaDataCategory) bool {
		if planCandidateSetCategory == nil {
			return false
		}
		if planCandidateSetCategory.IsSelected {
			return false
		}
		return planCandidateSetCategory.PlanCandidateSetID == planCandidateSetId
	})

	var categoriesRequired []string
	for _, planCandidateSetCategory := range planCandidateSetMetaDataRejectedCategorySlice {
		categoriesRequired = append(categoriesRequired, planCandidateSetCategory.Category)
	}

	if len(categoriesRequired) == 0 {
		return nil
	}

	locationCategoriesRejected := models.GetCategoriesFromSubCategories(categoriesRequired)

	return &locationCategoriesRejected
}
