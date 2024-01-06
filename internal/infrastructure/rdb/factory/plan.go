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
