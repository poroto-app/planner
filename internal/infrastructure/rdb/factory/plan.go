package factory

import (
	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
)

func NewPlanEntityFromDomainModel(plan models.Plan) generated.Plan {
	if plan.Id == "" {
		plan.Id = uuid.New().String()
	}
	return generated.Plan{
		ID:     plan.Id,
		UserID: null.StringFromPtr(plan.AuthorId),
		Name:   plan.Name,
	}
}

func NewPlanFromEntity(
	planEntity generated.Plan,
	planPlaceSlice generated.PlanPlaceSlice,
	places []models.Place,
) (*models.Plan, error) {
	planPlaces, err := NewPlanPlacesFromEntities(planPlaceSlice, places, planEntity.ID)
	if err != nil {
		return nil, err
	}

	return &models.Plan{
		Id:       planEntity.ID,
		Name:     planEntity.Name,
		AuthorId: planEntity.UserID.Ptr(),
		Places:   *planPlaces,
	}, nil
}
