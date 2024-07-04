package place

import (
	"context"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"time"
)

const (
	desfaultFetchDestinationPlacesForPlanCandidateLimit = 10
)

type FetchDestinationPlacesForPlanCandidateInput struct {
	PlanCandidateSetId string
	Limit              int
}

type FetchDestinationPlacesForPlanCandidateOutput struct {
	PlacesForPlanCandidates []PlacesForPlanCandidate
}

type PlacesForPlanCandidate struct {
	PlanCandidateId string
	Places          []models.Place
}

// FetchDestinationPlacesForPlanCandidate カテゴリからプランを作成したときに、その条件をもとに他の目的地を提示する
func (s Service) FetchDestinationPlacesForPlanCandidate(ctx context.Context, input FetchDestinationPlacesForPlanCandidateInput) (*FetchDestinationPlacesForPlanCandidateOutput, error) {
	if input.Limit == 0 {
		input.Limit = desfaultFetchDestinationPlacesForPlanCandidateLimit
	}

	planCandidate, err := s.planCandidateRepository.Find(ctx, input.PlanCandidateSetId, time.Now())
	if err != nil {
		return nil, err
	}

	// カテゴリからプランを作成したときのみ取得できるようにする
	if planCandidate.MetaData.CreateByCategoryMetaData == nil {
		return &FetchDestinationPlacesForPlanCandidateOutput{
			PlacesForPlanCandidates: array.Map(planCandidate.Plans, func(plan models.Plan) PlacesForPlanCandidate {
				return PlacesForPlanCandidate{
					PlanCandidateId: plan.Id,
					Places:          []models.Place{},
				}
			}),
		}, nil
	}

	googleCategoryTypes := planCandidate.MetaData.CreateByCategoryMetaData.Category.GooglePlaceTypes
	if len(googleCategoryTypes) == 0 {
		return &FetchDestinationPlacesForPlanCandidateOutput{
			PlacesForPlanCandidates: array.Map(planCandidate.Plans, func(plan models.Plan) PlacesForPlanCandidate {
				return PlacesForPlanCandidate{
					PlanCandidateId: plan.Id,
					Places:          []models.Place{},
				}
			}),
		}, nil
	}

	placesInPlans := array.FlatMap(planCandidate.Plans, func(plan models.Plan) []models.Place {
		return plan.Places
	})

	placesOfCategory, err := s.placeRepository.FindByGooglePlaceType(
		ctx,
		googleCategoryTypes[0],
		planCandidate.MetaData.CreateByCategoryMetaData.Location,
		planCandidate.MetaData.CreateByCategoryMetaData.RadiusInKm*1000,
	)
	if err != nil {
		return nil, err
	}

	// すでにプランに含まれている場所を提示しない
	placesDestination := array.Filter(*placesOfCategory, func(place models.Place) bool {
		_, found := array.Find(placesInPlans, func(p models.Place) bool {
			return p.Id == place.Id
		})
		return !found
	})

	if len(placesDestination) > input.Limit {
		placesDestination = placesDestination[:input.Limit]
	}

	return &FetchDestinationPlacesForPlanCandidateOutput{
		PlacesForPlanCandidates: array.Map(planCandidate.Plans, func(plan models.Plan) PlacesForPlanCandidate {
			return PlacesForPlanCandidate{
				PlanCandidateId: plan.Id,
				Places:          placesDestination,
			}
		}),
	}, nil
}
