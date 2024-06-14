package factory

import (
	"fmt"
	"github.com/volatiletech/null/v8"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
	"sort"
)

func PlanCandidateEntityFromDomainModel(planCandidate models.Plan, planCandidateSetId string, sortOrder int) generated.PlanCandidate {
	return generated.PlanCandidate{
		ID:                 planCandidate.Id,
		PlanCandidateSetID: planCandidateSetId,
		Name:               planCandidate.Name,
		SortOrder:          sortOrder,
		ParentPlanID:       null.StringFromPtr(planCandidate.ParentPlanId),
	}
}

// NewPlanCandidatesFromEntities プラン候補一覧を順番を考慮して生成する
func NewPlanCandidatesFromEntities(
	planCandidateSlice generated.PlanCandidateSlice,
	planCandidatePlaces generated.PlanCandidatePlaceSlice,
	planCandidateSetId string,
	places []models.Place,
	author *models.User,
) (*[]models.Plan, error) {
	planCandidateEntities := array.MapAndFilter(planCandidateSlice, func(planCandidate *generated.PlanCandidate) (generated.PlanCandidate, bool) {
		if planCandidate == nil {
			return generated.PlanCandidate{}, false
		}

		if planCandidate.PlanCandidateSetID != planCandidateSetId {
			return generated.PlanCandidate{}, false
		}

		return *planCandidate, true
	})

	planCandidateEntitiesOrdered := planCandidateEntities
	sort.Slice(planCandidateEntitiesOrdered, func(i, j int) bool {
		return planCandidateEntitiesOrdered[i].SortOrder < planCandidateEntitiesOrdered[j].SortOrder
	})

	plans, err := array.MapWithErr(planCandidateEntitiesOrdered, func(planCandidateEntity generated.PlanCandidate) (*models.Plan, error) {
		plan, err := NewPlanCandidateFromEntity(planCandidateEntity, planCandidatePlaces, places, author)
		if err != nil {
			return nil, err
		}

		return plan, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to map plans: %w", err)
	}

	return plans, nil
}

// NewPlanCandidateFromEntity プラン候補を場所の順番を考慮して生成する
func NewPlanCandidateFromEntity(
	planCandidateEntity generated.PlanCandidate,
	planCandidatePlaces generated.PlanCandidatePlaceSlice,
	places []models.Place,
	author *models.User,
) (*models.Plan, error) {
	planCandidateEntities := array.MapAndFilter(planCandidatePlaces, func(planCandidatePlace *generated.PlanCandidatePlace) (generated.PlanCandidatePlace, bool) {
		if planCandidatePlace == nil {
			return generated.PlanCandidatePlace{}, false
		}

		if planCandidatePlace.PlanCandidateID != planCandidateEntity.ID {
			return generated.PlanCandidatePlace{}, false
		}

		return *planCandidatePlace, true
	})

	planCandidateEntitiesOrdered := planCandidateEntities
	sort.Slice(planCandidateEntitiesOrdered, func(i, j int) bool {
		return planCandidateEntitiesOrdered[i].SortOrder < planCandidateEntitiesOrdered[j].SortOrder
	})

	placesOrdered, err := array.MapWithErr(planCandidateEntitiesOrdered, func(planCandidatePlace generated.PlanCandidatePlace) (*models.Place, error) {
		place, ok := array.Find(places, func(place models.Place) bool {
			return place.Id == planCandidatePlace.PlaceID
		})
		if !ok {
			return nil, fmt.Errorf("failed to find place with id %s", planCandidatePlace.PlaceID)
		}

		return &place, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to map places: %w", err)
	}

	return &models.Plan{
		Id:           planCandidateEntity.ID,
		Name:         planCandidateEntity.Name,
		Places:       *placesOrdered,
		Author:       author,
		ParentPlanId: planCandidateEntity.ParentPlanID.Ptr(),
	}, nil
}
