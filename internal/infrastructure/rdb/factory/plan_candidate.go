package factory

import (
	"fmt"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/entities"
	"sort"
)

func PlanCandidateEntityFromDomainModel(planCandidate models.Plan, planCandidateSetId string, sortOrder int) entities.PlanCandidate {
	return entities.PlanCandidate{
		ID:                 planCandidate.Id,
		PlanCandidateSetID: planCandidateSetId,
		Name:               planCandidate.Name,
		SortOrder:          sortOrder,
	}
}

// NewPlanCandidatesFromEntities プラン候補一覧を順番を考慮して生成する
func NewPlanCandidatesFromEntities(
	planCandidateSlice entities.PlanCandidateSlice,
	planCandidatePlaces entities.PlanCandidatePlaceSlice,
	planCandidateSetId string,
	places []models.Place,
) (*[]models.Plan, error) {
	planCandidateEntities := array.MapAndFilter(planCandidateSlice, func(planCandidate *entities.PlanCandidate) (entities.PlanCandidate, bool) {
		if planCandidate == nil {
			return entities.PlanCandidate{}, false
		}

		if planCandidate.PlanCandidateSetID != planCandidateSetId {
			return entities.PlanCandidate{}, false
		}

		return *planCandidate, true
	})

	planCandidateEntitiesOrdered := planCandidateEntities
	sort.Slice(planCandidateEntitiesOrdered, func(i, j int) bool {
		return planCandidateEntitiesOrdered[i].SortOrder < planCandidateEntitiesOrdered[j].SortOrder
	})

	var plans []models.Plan
	for _, planCandidateEntity := range planCandidateEntitiesOrdered {
		plan, err := NewPlanCandidateFromEntity(planCandidateEntity, planCandidatePlaces, places)
		if err != nil {
			return nil, err
		}
		plans = append(plans, *plan)
	}

	return &plans, nil
}

// NewPlanCandidateFromEntity プラン候補を場所の順番を考慮して生成する
func NewPlanCandidateFromEntity(
	planCandidateEntity entities.PlanCandidate,
	planCandidatePlaces entities.PlanCandidatePlaceSlice,
	places []models.Place,
) (*models.Plan, error) {
	planCandidateEntities := array.MapAndFilter(planCandidatePlaces, func(planCandidatePlace *entities.PlanCandidatePlace) (entities.PlanCandidatePlace, bool) {
		if planCandidatePlace == nil {
			return entities.PlanCandidatePlace{}, false
		}

		if planCandidatePlace.PlanCandidateID != planCandidateEntity.ID {
			return entities.PlanCandidatePlace{}, false
		}

		return *planCandidatePlace, true
	})

	planCandidateEntitiesOrdered := planCandidateEntities
	sort.Slice(planCandidateEntitiesOrdered, func(i, j int) bool {
		return planCandidateEntitiesOrdered[i].SortOrder < planCandidateEntitiesOrdered[j].SortOrder
	})

	placesOrdered, err := array.MapWithErr(planCandidateEntitiesOrdered, func(planCandidatePlace entities.PlanCandidatePlace) (*models.Place, error) {
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
		Id:       planCandidateEntity.ID,
		Name:     planCandidateEntity.Name,
		Places:   *placesOrdered,
		AuthorId: nil, // TODO: implement me
	}, nil
}
