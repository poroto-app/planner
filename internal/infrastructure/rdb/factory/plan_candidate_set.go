package factory

import (
	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/array"
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
	planCandidateSetLikePlaceSlice generated.PlanCandidateSetLikePlaceSlice,
	places []models.Place,
	author *models.User,
) (*models.PlanCandidateSet, error) {
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
		author,
	)
	if err != nil {
		return nil, err
	}

	likedPlaceIds := array.MapAndFilter(planCandidateSetLikePlaceSlice, func(planCandidateSetLikePlace *generated.PlanCandidateSetLikePlace) (string, bool) {
		if planCandidateSetLikePlace == nil {
			return "", false
		}
		if planCandidateSetLikePlace.PlanCandidateSetID != planCandidateSetEntity.ID {
			return "", false
		}
		return planCandidateSetLikePlace.PlaceID, true
	})

	return &models.PlanCandidateSet{
		Id:              planCandidateSetEntity.ID,
		Plans:           *plans,
		MetaData:        *planCandidateSetMetaData,
		IsPlaceSearched: planCandidateSetEntity.IsPlaceSearched,
		ExpiresAt:       planCandidateSetEntity.ExpiresAt,
		LikedPlaceIds:   likedPlaceIds,
	}, nil
}
