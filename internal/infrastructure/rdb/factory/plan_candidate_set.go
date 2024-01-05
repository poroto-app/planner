package factory

import (
	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
)

func NewPlanCandidateSetFromEntity(
	planCandidateSetEntity generated.PlanCandidateSet,
	planCandidateSlice generated.PlanCandidateSlice,
	planCandidateSetMetaDataSlice generated.PlanCandidateSetMetaDatumSlice,
	planCandidateSetCategorySlice generated.PlanCandidateSetMetaDataCategorySlice,
	planCandidatePlaces generated.PlanCandidatePlaceSlice,
	places []models.Place,
) (*models.PlanCandidate, error) {
	planCandidateSetMetaData, err := NewPlanCandidateMetaDataFromEntity(planCandidateSetMetaDataSlice, planCandidateSetCategorySlice, planCandidateSetEntity.ID)
	if err != nil {
		// PlanCandidateSetMetaDataがない場合はエラーにしない
		logger, err := utils.NewLogger(utils.LoggerOption{
			Tag: "NewPlanCandidateSetFromEntity",
		})
		if err != nil {
			return nil, err
		}
		logger.Debug("skip to create PlanCandidateSetMetaData", zap.Error(err))
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
