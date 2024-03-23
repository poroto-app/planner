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

	var userId *string
	if plan.Author != nil {
		userId = &plan.Author.Id
	}

	var startLocation models.GeoLocation
	if len(plan.Places) > 0 {
		startLocation = plan.Places[0].Location
	}

	return generated.Plan{
		ID:        plan.Id,
		UserID:    null.StringFromPtr(userId),
		Name:      plan.Name,
		Latitude:  startLocation.Latitude,
		Longitude: startLocation.Longitude,
	}
}

func NewPlanFromEntity(
	planEntity generated.Plan,
	planPlaceSlice generated.PlanPlaceSlice,
	places []models.Place,
	author *models.User,
) (*models.Plan, error) {
	planPlaces, err := NewPlanPlacesFromEntities(planPlaceSlice, places, planEntity.ID)
	if err != nil {
		return nil, err
	}

	return &models.Plan{
		Id:     planEntity.ID,
		Name:   planEntity.Name,
		Places: *planPlaces,
		Author: author,
	}, nil
}
