package factory

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
)

func NewPlanCandidateMetaDataFromDomainModel(planCandidateSetMetaData models.PlanCandidateMetaData, planCandidateSetId string) generated.PlanCandidateSetMetaDatum {
	return generated.PlanCandidateSetMetaDatum{
		ID:                           uuid.New().String(),
		PlanCandidateSetID:           planCandidateSetId,
		IsCreatedFromCurrentLocation: planCandidateSetMetaData.CreatedBasedOnCurrentLocation,
		LatitudeStart:                planCandidateSetMetaData.LocationStart.Latitude,
		LongitudeStart:               planCandidateSetMetaData.LocationStart.Longitude,
		PlanDurationMinutes:          null.IntFromPtr(planCandidateSetMetaData.FreeTime),
	}
}

func NewPlanCandidateMetaDataFromEntity(
	planCandidateSetMetaDataSlice generated.PlanCandidateSetMetaDatumSlice,
	planCandidateSetCategorySlice generated.PlanCandidateSetMetaDataCategorySlice,
	planCandidateSetMetaDataCreateByCategory generated.PlanCandidateSetMetaDataCreateByCategorySlice,
	planCandidateSetId string,
) (*models.PlanCandidateMetaData, error) {
	planCandidateSetMetaData, ok := array.Find(planCandidateSetMetaDataSlice, func(planCandidateSetMetaData *generated.PlanCandidateSetMetaDatum) bool {
		if planCandidateSetMetaData == nil {
			return false
		}
		return planCandidateSetMetaData.PlanCandidateSetID == planCandidateSetId
	})
	if !ok {
		return nil, fmt.Errorf("failed to find plan candidate set meta data")
	}

	return &models.PlanCandidateMetaData{
		CreatedBasedOnCurrentLocation: planCandidateSetMetaData.IsCreatedFromCurrentLocation,
		CategoriesPreferred:           newPlanCandidateSetMetaDataPreferredCategoriesFromEntity(planCandidateSetCategorySlice, planCandidateSetId),
		CategoriesRejected:            newPlanCandidateSetMetaDataRejectedCategoriesFromEntity(planCandidateSetCategorySlice, planCandidateSetId),
		LocationStart: &models.GeoLocation{
			Latitude:  planCandidateSetMetaData.LatitudeStart,
			Longitude: planCandidateSetMetaData.LongitudeStart,
		},
		FreeTime:                 planCandidateSetMetaData.PlanDurationMinutes.Ptr(),
		CreateByCategoryMetaData: newPlanCandidateSetMetaDataCreateByCategoryFromEntry(planCandidateSetMetaDataCreateByCategory),
	}, nil
}
