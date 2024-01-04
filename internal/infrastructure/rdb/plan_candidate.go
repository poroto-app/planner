package rdb

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/entities"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/factory"
	"time"
)

type PlanCandidateRepository struct {
	db     *sql.DB
	logger zap.Logger
}

func NewPlanCandidateRepository(db *sql.DB) (*PlanCandidateRepository, error) {
	logger, err := utils.NewLogger(utils.LoggerOption{
		Tag: "PlanCandidateRepository",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	return &PlanCandidateRepository{db: db, logger: *logger}, nil
}

func (p PlanCandidateRepository) GetDB() *sql.DB {
	return p.db
}

func (p PlanCandidateRepository) Create(cxt context.Context, planCandidateId string, expiresAt time.Time) error {
	if err := runTransaction(cxt, p, func(ctx context.Context, tx *sql.Tx) error {
		planCandidateEntity := entities.PlanCandidateSet{ID: planCandidateId, ExpiresAt: expiresAt}
		if err := planCandidateEntity.Insert(ctx, tx, boil.Infer()); err != nil {
			return fmt.Errorf("failed to insert plan candidate: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to run transaction: %w", err)
	}
	return nil
}

func (p PlanCandidateRepository) Find(ctx context.Context, planCandidateId string) (*models.PlanCandidate, error) {
	//TODO implement me
	panic("implement me")
}

func (p PlanCandidateRepository) FindPlan(ctx context.Context, planCandidateId string, planId string) (*models.Plan, error) {
	//TODO implement me
	panic("implement me")
}

func (p PlanCandidateRepository) FindExpiredBefore(ctx context.Context, expiresAt time.Time) (*[]string, error) {
	//TODO implement me
	panic("implement me")
}

func (p PlanCandidateRepository) AddSearchedPlacesForPlanCandidate(ctx context.Context, planCandidateId string, placeIds []string) error {
	if err := runTransaction(ctx, p, func(ctx context.Context, tx *sql.Tx) error {
		// TODO: BatchInsertする
		for _, placeId := range placeIds {
			planCandidatePlace := entities.PlanCandidateSetSearchedPlace{ID: uuid.New().String(), PlanCandidateSetID: planCandidateId, PlaceID: placeId}
			if err := planCandidatePlace.Insert(ctx, tx, boil.Infer()); err != nil {
				return fmt.Errorf("failed to insert plan candidate place: %w", err)
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to run transaction: %w", err)
	}
	return nil
}

func (p PlanCandidateRepository) AddPlan(ctx context.Context, planCandidateId string, plans ...models.Plan) error {
	if err := runTransaction(ctx, p, func(ctx context.Context, tx *sql.Tx) error {
		// TODO: BatchInsertする
		for iPlan, plan := range plans {
			planCandidateEntity := factory.PlanCandidateEntityFromDomainModel(plan, planCandidateId, iPlan)
			if err := planCandidateEntity.Insert(ctx, tx, boil.Infer()); err != nil {
				return fmt.Errorf("failed to insert plan candidate: %w", err)
			}

			planCandidatePlaceSlice := factory.NewPlanCandidatePlaceSliceFromDomainModel(plan.Places, planCandidateId, plan.Id)
			for _, planCandidatePlace := range planCandidatePlaceSlice {
				if err := planCandidatePlace.Insert(ctx, tx, boil.Infer()); err != nil {
					return fmt.Errorf("failed to insert plan candidate place: %w", err)
				}
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to run transaction: %w", err)
	}

	return nil
}

func (p PlanCandidateRepository) AddPlaceToPlan(ctx context.Context, planCandidateId string, planId string, previousPlaceId string, place models.Place) error {
	if err := runTransaction(ctx, p, func(ctx context.Context, tx *sql.Tx) error {
		planCandidatePlaceSlice, err := entities.
			PlanCandidatePlaces(entities.PlanCandidatePlaceWhere.PlanCandidateSetID.EQ(planCandidateId)).
			All(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to get plan candidate places: %w", err)
		}

		newOrder := 0
		for _, planCandidatePlace := range planCandidatePlaceSlice {
			if planCandidatePlace.PlaceID == previousPlaceId {
				// 挿入する場所の順序を決定
				newOrder = planCandidatePlace.SortOrder + 1
			} else if planCandidatePlace.SortOrder >= newOrder {
				// 後続の場所の順序を更新
				planCandidatePlace.SortOrder++
				if _, err := planCandidatePlace.Update(ctx, tx, boil.Whitelist(entities.PlanCandidatePlaceColumns.SortOrder)); err != nil {
					return fmt.Errorf("failed to update plan candidate place: %w", err)
				}
			}
		}

		planCandidateEntity := factory.NewPlanCandidatePlaceEntityFromDomainModel(place, newOrder, planCandidateId, planId)
		if err := planCandidateEntity.Insert(ctx, tx, boil.Infer()); err != nil {
			return fmt.Errorf("failed to insert plan candidate place: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to run transaction: %w", err)
	}

	return nil
}

func (p PlanCandidateRepository) RemovePlaceFromPlan(ctx context.Context, planCandidateId string, planId string, placeId string) error {
	//TODO implement me
	panic("implement me")
}

func (p PlanCandidateRepository) UpdatePlacesOrder(ctx context.Context, planId string, planCandidate string, placeIdsOrdered []string) (*models.Plan, error) {
	//TODO implement me
	panic("implement me")
}

func (p PlanCandidateRepository) UpdatePlanCandidateMetaData(ctx context.Context, planCandidateId string, meta models.PlanCandidateMetaData) error {
	//TODO implement me
	panic("implement me")
}

func (p PlanCandidateRepository) ReplacePlace(ctx context.Context, planCandidateId string, planId string, placeIdToBeReplaced string, placeToReplace models.Place) error {
	//TODO implement me
	panic("implement me")
}

func (p PlanCandidateRepository) DeleteAll(ctx context.Context, planCandidateIds []string) error {
	//TODO implement me
	panic("implement me")
}

func (p PlanCandidateRepository) UpdateLikeToPlaceInPlanCandidate(ctx context.Context, planCandidateId string, placeId string, like bool) error {
	//TODO implement me
	panic("implement me")
}
