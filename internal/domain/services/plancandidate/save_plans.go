package plancandidate

import (
	"context"
	"fmt"
	"poroto.app/poroto/planner/internal/domain/models"
)

type SavePlansInput struct {
	PlanCandidateSetId           string
	Plans                        []models.Plan
	LocationStart                *models.GeoLocation
	CategoryNamesPreferred       *[]string
	CategoryNamesRejected        *[]string
	FreeTime                     *int
	CreateBasedOnCurrentLocation bool
}

// SavePlans 作成されたプランとプラン候補のメタデータを保存する
// はじめてプランを作成したときに呼び出す
func (s Service) SavePlans(ctx context.Context, input SavePlansInput) (err error) {
	// プランをプラン候補に追加して保存する
	if err := s.planCandidateRepository.AddPlan(ctx, input.PlanCandidateSetId, input.Plans...); err != nil {
		return fmt.Errorf("error while adding plan to plan candidate: %v\n", err)
	}

	// プラン作成時の情報を保存
	var categoriesPreferred, categoriesDisliked *[]models.LocationCategory
	if input.CategoryNamesPreferred != nil {
		var categories []models.LocationCategory
		for _, categoryName := range *input.CategoryNamesPreferred {
			category := models.GetCategoryOfName(categoryName)
			if category != nil {
				categories = append(categories, *category)
			}
		}
		categoriesPreferred = &categories
	}

	if input.CategoryNamesRejected != nil {
		var categories []models.LocationCategory
		for _, categoryName := range *input.CategoryNamesRejected {
			category := models.GetCategoryOfName(categoryName)
			if category != nil {
				categories = append(categories, *category)
			}
		}
		categoriesDisliked = &categories
	}

	if err := s.planCandidateRepository.UpdatePlanCandidateMetaData(ctx, input.PlanCandidateSetId, models.PlanCandidateMetaData{
		LocationStart:                 input.LocationStart,
		CategoriesPreferred:           categoriesPreferred,
		CategoriesRejected:            categoriesDisliked,
		FreeTime:                      input.FreeTime,
		CreatedBasedOnCurrentLocation: input.CreateBasedOnCurrentLocation,
	}); err != nil {
		return fmt.Errorf("error while updating plan candidate metadata: %v\n", err)
	}

	return nil
}
