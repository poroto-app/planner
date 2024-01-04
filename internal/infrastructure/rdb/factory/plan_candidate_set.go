package factory

import (
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/entities"
)

func NewPlanCandidateSetFromEntity(
	planCandidateSetEntity entities.PlanCandidateSet,
	planCandidateSlice entities.PlanCandidateSlice,
	planCandidateSetMetaDataSlice entities.PlanCandidateSetMetaDatumSlice,
	planCandidateSetCategorySlice entities.PlanCandidateSetCategorySlice,
	planCandidatePlaces entities.PlanCandidatePlaceSlice,
	places []models.Place,
) (*models.PlanCandidate, error) {
	planCandidateSetMetaData, err := NewPlanCandidateMetaDataFromEntity(planCandidateSetMetaDataSlice, planCandidateSetCategorySlice, planCandidateSetEntity.ID)
	if err != nil {
		// PlanCandidateSetMetaDataがない場合はエラーにしない
		planCandidateSetMetaData = &models.PlanCandidateMetaData{}
	}

	plans, err := NewPlanCandidatesFromEntities(
		planCandidateSlice,
		planCandidatePlaces,
		planCandidateSetEntity.ID,
		places,
	)
	if err != nil {
		return nil, err
	}

	return &models.PlanCandidate{
		Id:            planCandidateSetEntity.ID,
		Plans:         *plans,
		MetaData:      *planCandidateSetMetaData,
		ExpiresAt:     planCandidateSetEntity.ExpiresAt,
		LikedPlaceIds: nil, //TODO: implement me!
	}, nil
}
