package factory

import (
	"github.com/google/uuid"
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
