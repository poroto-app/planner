package rdb

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/factory"
)

type PlanRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewPlanRepository(db *sql.DB) (*PlanRepository, error) {
	logger, err := utils.NewLogger(utils.LoggerOption{
		Tag: "PlanRepository",
	})
	if err != nil {
		return nil, err
	}

	return &PlanRepository{
		db:     db,
		logger: logger,
	}, nil
}

func (p PlanRepository) GetDB() *sql.DB {
	return p.db
}

func (p PlanRepository) Save(ctx context.Context, plan *models.Plan) error {
	// TODO: ポインタ型の引数にしない
	if plan == nil {
		return nil
	}

	if err := runTransaction(ctx, p, func(ctx context.Context, tx *sql.Tx) error {
		planEntity := factory.NewPlanEntityFromDomainModel(*plan)
		if err := planEntity.Insert(ctx, tx, boil.Infer()); err != nil {
			return fmt.Errorf("failed to insert plan: %w", err)
		}

		planPlaceSlice := factory.NewPlanPlaceSliceFromDomainMode(plan.Places, plan.Id)
		if _, err := planPlaceSlice.InsertAll(ctx, tx, boil.Infer()); err != nil {
			return fmt.Errorf("failed to insert plan places: %w", err)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (p PlanRepository) SortedByCreatedAt(ctx context.Context, queryCursor *string, limit int) (*[]models.Plan, error) {
	//TODO implement me
	panic("implement me")
}

func (p PlanRepository) Find(ctx context.Context, planId string) (*models.Plan, error) {
	//TODO implement me
	panic("implement me")
}

func (p PlanRepository) FindByAuthorId(ctx context.Context, authorId string) (*[]models.Plan, error) {
	//TODO implement me
	panic("implement me")
}

func (p PlanRepository) SortedByLocation(ctx context.Context, location models.GeoLocation, queryCursor *string, limit int) (*[]models.Plan, *string, error) {
	//TODO implement me
	panic("implement me")
}
