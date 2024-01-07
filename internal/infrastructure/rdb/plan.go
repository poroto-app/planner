package rdb

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/entities"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/factory"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
	"strconv"
	"time"
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

// TODO: QueryCursor をこの関数内で生成する
func (p PlanRepository) SortedByCreatedAt(ctx context.Context, queryCursor *string, limit int) (*[]models.Plan, error) {
	planQueryMod := []qm.QueryMod{
		qm.Load(generated.PlanRels.PlanPlaces),
		qm.OrderBy(fmt.Sprintf("%s %s, %s %s", generated.PlanColumns.CreatedAt, "desc", generated.PlanColumns.ID, "desc")),
		qm.Limit(limit),
	}

	if queryCursor != nil {
		dateTime, err := parseSortByCreatedAtQueryCursor(*queryCursor)
		if err != nil {
			return nil, err
		}

		// WHERE (created_at) < (dateTime)
		planQueryMod = append(planQueryMod, qm.Where(
			fmt.Sprintf("%s < ?", generated.PlanColumns.CreatedAt),
			dateTime,
		))
	}

	planEntities, err := generated.Plans(concatQueryMod(
		planQueryMod,
		placeQueryModes(generated.PlanRels.PlanPlaces, generated.PlanPlaceRels.Place),
	)...).All(ctx, p.db)
	if err != nil {
		return nil, fmt.Errorf("failed to find plans: %w", err)
	}

	if len(planEntities) == 0 {
		return &[]models.Plan{}, nil
	}

	planCandidateSetPlaceLikeCounts, err := countPlaceLikeCounts(ctx, p.db, array.Map(planEntities, func(planEntity *generated.Plan) string {
		return planEntity.ID
	})...)
	if err != nil {
		// いいね数の取得に失敗してもエラーにしない
		p.logger.Warn("failed to count place like counts", zap.Error(err))
	}

	places, err := array.MapWithErr(planEntities, func(planEntity *generated.Plan) (*[]models.Place, error) {
		if planEntity.R == nil {
			return nil, fmt.Errorf("planEntity.R is nil")
		}

		return array.MapWithErr(planEntity.R.PlanPlaces, func(planPlace *generated.PlanPlace) (*models.Place, error) {
			if planPlace.R == nil {
				return nil, fmt.Errorf("planPlace.R is nil")
			}

			if planPlace.R.Place == nil {
				return nil, fmt.Errorf("planPlace.R.Place is nil")
			}

			if len(planPlace.R.Place.R.GooglePlaces) == 0 || planPlace.R.Place.R.GooglePlaces[0].R == nil {
				return nil, fmt.Errorf("planPlace.R.Place.R.GooglePlaces is nil")
			}

			return factory.NewPlaceFromEntity(
				*planPlace.R.Place,
				*planPlace.R.Place.R.GooglePlaces[0],
				planPlace.R.Place.R.GooglePlaces[0].R.GooglePlaceTypes,
				planPlace.R.Place.R.GooglePlaces[0].R.GooglePlacePhotoReferences,
				planPlace.R.Place.R.GooglePlaces[0].R.GooglePlacePhotoAttributions,
				planPlace.R.Place.R.GooglePlaces[0].R.GooglePlacePhotos,
				planPlace.R.Place.R.GooglePlaces[0].R.GooglePlaceReviews,
				planPlace.R.Place.R.GooglePlaces[0].R.GooglePlaceOpeningPeriods,
				entities.CountLikeOfPlace(planCandidateSetPlaceLikeCounts, planPlace.PlaceID),
			)
		})
	})
	if err != nil {
		return nil, fmt.Errorf("failed to map plan places: %w", err)
	}

	plans, err := array.MapWithErr(planEntities, func(planEntity *generated.Plan) (*models.Plan, error) {
		return factory.NewPlanFromEntity(
			*planEntity,
			planEntity.R.PlanPlaces,
			array.Flatten(*places),
		)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to map plans: %w", err)
	}

	return plans, nil
}

func (p PlanRepository) Find(ctx context.Context, planId string) (*models.Plan, error) {
	planEntity, err := generated.Plans(concatQueryMod(
		[]qm.QueryMod{
			generated.PlanWhere.ID.EQ(planId),
			qm.Load(generated.PlanRels.PlanPlaces),
		},
		placeQueryModes(generated.PlanRels.PlanPlaces, generated.PlanPlaceRels.Place),
	)...).One(ctx, p.db)
	if err != nil {
		return nil, fmt.Errorf("failed to find plan: %w", err)
	}

	if planEntity.R == nil {
		return nil, fmt.Errorf("planEntity.R is nil")
	}

	planCandidateSetPlaceLikeCounts, err := countPlaceLikeCounts(ctx, p.db, array.Map(planEntity.R.PlanPlaces, func(planPlace *generated.PlanPlace) string {
		return planPlace.PlaceID
	})...)
	if err != nil {
		// いいね数の取得に失敗してもエラーにしない
		p.logger.Warn("failed to count place like counts", zap.Error(err))
	}

	places, err := array.MapWithErr(planEntity.R.PlanPlaces, func(planPlace *generated.PlanPlace) (*models.Place, error) {
		if planPlace.R == nil {
			return nil, fmt.Errorf("planPlace.R is nil")
		}

		if planPlace.R.Place == nil {
			return nil, fmt.Errorf("planPlace.R.Place is nil")
		}

		if len(planPlace.R.Place.R.GooglePlaces) == 0 || planPlace.R.Place.R.GooglePlaces[0].R == nil {
			return nil, fmt.Errorf("planPlace.R.Place.R.GooglePlaces is nil")
		}

		return factory.NewPlaceFromEntity(
			*planPlace.R.Place,
			*planPlace.R.Place.R.GooglePlaces[0],
			planPlace.R.Place.R.GooglePlaces[0].R.GooglePlaceTypes,
			planPlace.R.Place.R.GooglePlaces[0].R.GooglePlacePhotoReferences,
			planPlace.R.Place.R.GooglePlaces[0].R.GooglePlacePhotoAttributions,
			planPlace.R.Place.R.GooglePlaces[0].R.GooglePlacePhotos,
			planPlace.R.Place.R.GooglePlaces[0].R.GooglePlaceReviews,
			planPlace.R.Place.R.GooglePlaces[0].R.GooglePlaceOpeningPeriods,
			entities.CountLikeOfPlace(planCandidateSetPlaceLikeCounts, planPlace.PlaceID),
		)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to map plan places: %w", err)
	}

	plan, err := factory.NewPlanFromEntity(
		*planEntity,
		planEntity.R.PlanPlaces,
		*places,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to map plan: %w", err)
	}

	return plan, nil
}

func (p PlanRepository) FindByAuthorId(ctx context.Context, authorId string) (*[]models.Plan, error) {
	planEntities, err := generated.Plans(concatQueryMod(
		[]qm.QueryMod{
			generated.PlanWhere.UserID.EQ(null.StringFrom(authorId)),
			qm.Load(generated.PlanRels.PlanPlaces),
			qm.OrderBy(fmt.Sprintf("%s %s", generated.PlanColumns.CreatedAt, "desc")),
		},
		placeQueryModes(generated.PlanRels.PlanPlaces, generated.PlanPlaceRels.Place),
	)...).All(ctx, p.db)
	if err != nil {
		return nil, fmt.Errorf("failed to find plans: %w", err)
	}

	if len(planEntities) == 0 {
		return &[]models.Plan{}, nil
	}

	planCandidateSetPlaceLikeCounts, err := countPlaceLikeCounts(ctx, p.db, array.Map(planEntities, func(planEntity *generated.Plan) string {
		return planEntity.ID
	})...)
	if err != nil {
		// いいね数の取得に失敗してもエラーにしない
		p.logger.Warn("failed to count place like counts", zap.Error(err))
	}

	places, err := array.MapWithErr(planEntities, func(planEntity *generated.Plan) (*[]models.Place, error) {
		if planEntity.R == nil {
			return nil, fmt.Errorf("planEntity.R is nil")
		}

		if planEntity.R.PlanPlaces == nil {
			return nil, fmt.Errorf("planEntity.R.PlanPlaces is nil")
		}

		return array.MapWithErr(planEntity.R.PlanPlaces, func(planPlace *generated.PlanPlace) (*models.Place, error) {
			if planPlace.R == nil {
				return nil, fmt.Errorf("planPlace.R is nil")
			}

			if planPlace.R.Place == nil {
				return nil, fmt.Errorf("planPlace.R.Place is nil")
			}

			if len(planPlace.R.Place.R.GooglePlaces) == 0 || planPlace.R.Place.R.GooglePlaces[0].R == nil {
				return nil, fmt.Errorf("planPlace.R.Place.R.GooglePlaces is nil")
			}

			return factory.NewPlaceFromEntity(
				*planPlace.R.Place,
				*planPlace.R.Place.R.GooglePlaces[0],
				planPlace.R.Place.R.GooglePlaces[0].R.GooglePlaceTypes,
				planPlace.R.Place.R.GooglePlaces[0].R.GooglePlacePhotoReferences,
				planPlace.R.Place.R.GooglePlaces[0].R.GooglePlacePhotoAttributions,
				planPlace.R.Place.R.GooglePlaces[0].R.GooglePlacePhotos,
				planPlace.R.Place.R.GooglePlaces[0].R.GooglePlaceReviews,
				planPlace.R.Place.R.GooglePlaces[0].R.GooglePlaceOpeningPeriods,
				entities.CountLikeOfPlace(planCandidateSetPlaceLikeCounts, planPlace.PlaceID),
			)
		})
	})
	if err != nil {
		return nil, fmt.Errorf("failed to map plan places: %w", err)
	}

	plans, err := array.MapWithErr(planEntities, func(planEntity *generated.Plan) (*models.Plan, error) {
		return factory.NewPlanFromEntity(
			*planEntity,
			planEntity.R.PlanPlaces,
			array.Flatten(*places),
		)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to map plans: %w", err)
	}

	return plans, nil
}

func (p PlanRepository) SortedByLocation(ctx context.Context, location models.GeoLocation, queryCursor *string, limit int) (*[]models.Plan, *string, error) {
	//TODO implement me
	panic("implement me")
}

func newSortByCreatedAtQueryCursor(createdAt time.Time) string {
	return fmt.Sprintf("%d", createdAt.Unix())
}

func parseSortByCreatedAtQueryCursor(queryCursor string) (*time.Time, error) {
	unixTime, err := strconv.ParseInt(queryCursor, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid query cursor: %s", queryCursor)
	}
	dateTime := time.Unix(unixTime, 0)
	return &dateTime, nil
}
