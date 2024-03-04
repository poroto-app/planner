package rdb

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/repository"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/entities"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/factory"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
)

const (
	// 半径2km圏内のプランを検索する
	defaultDistanceToSearchPlan = 2 * 1000
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

	if len(plan.Places) == 0 {
		return fmt.Errorf("plan places is empty")
	}

	startLocation := plan.Places[0].Location

	if err := runTransaction(ctx, p, func(ctx context.Context, tx *sql.Tx) error {
		planEntity := factory.NewPlanEntityFromDomainModel(*plan)
		if _, err := queries.Raw(
			fmt.Sprintf(
				"INSERT INTO %s (%s, %s, %s, %s) VALUES (?, ?, ?, POINT(?, ?))",
				generated.TableNames.Plans,
				generated.PlanColumns.ID,
				generated.PlanColumns.UserID,
				generated.PlanColumns.Name,
				generated.PlanColumns.Location,
			),
			planEntity.ID,
			planEntity.UserID,
			planEntity.Name,
			startLocation.Longitude,
			startLocation.Latitude,
		).ExecContext(ctx, tx); err != nil {
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
func (p PlanRepository) SortedByCreatedAt(ctx context.Context, queryCursor *repository.SortedByCreatedAtQueryCursor, limit int) (*[]models.Plan, *repository.SortedByCreatedAtQueryCursor, error) {
	planQueryMod := []qm.QueryMod{
		qm.Load(generated.PlanRels.PlanPlaces),
		qm.OrderBy(fmt.Sprintf("%s %s, %s %s", generated.PlanColumns.CreatedAt, "desc", generated.PlanColumns.ID, "desc")),
		qm.Limit(limit),
		qm.Load(generated.PlanRels.User),
		qm.Load(generated.PlanRels.PlanPlaces + "." + generated.PlanPlaceRels.Place + "." + generated.PlaceRels.GooglePlaces),
	}

	if queryCursor != nil {
		dateTime, err := parseSortByCreatedAtQueryCursor(*queryCursor)
		if err != nil {
			return nil, nil, err
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
		return nil, nil, fmt.Errorf("failed to find plans: %w", err)
	}

	if len(planEntities) == 0 {
		return &[]models.Plan{}, nil, nil
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
				array.Filter(planPlace.R.Place.R.PlacePhotos, func(placePhoto *generated.PlacePhoto) bool {
					return placePhoto.PlaceID == planPlace.PlaceID
				}),
			)
		})
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to map plan places: %w", err)
	}

	plans, err := array.MapWithErr(planEntities, func(planEntity *generated.Plan) (*models.Plan, error) {
		var author *models.User
		if planEntity.R.User != nil {
			author = factory.NewUserFromUserEntity(*planEntity.R.User)
		}

		return factory.NewPlanFromEntity(
			*planEntity,
			planEntity.R.PlanPlaces,
			array.Flatten(*places),
			author,
		)
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to map plans: %w", err)
	}

	var nextQueryCursor *repository.SortedByCreatedAtQueryCursor
	if len(*plans) == limit {
		qc := newSortByCreatedAtQueryCursor(planEntities[limit-1].CreatedAt)
		nextQueryCursor = &qc
	}

	return plans, nextQueryCursor, nil
}

func (p PlanRepository) Find(ctx context.Context, planId string) (*models.Plan, error) {
	planEntity, err := generated.Plans(concatQueryMod(
		[]qm.QueryMod{
			generated.PlanWhere.ID.EQ(planId),
			qm.Load(generated.PlanRels.PlanPlaces),
			qm.Load(generated.PlanRels.User),
		},
		placeQueryModes(generated.PlanRels.PlanPlaces, generated.PlanPlaceRels.Place),
		placeQueryModes(generated.PlanRels.PlanPlaces, generated.PlanPlaceRels.Place, generated.PlaceRels.PlacePhotos),
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

	placePhotoSlice := planEntity.R.PlanPlaces.GetLoadedPlaces().GetLoadedPlacePhotos()

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
			array.Filter(placePhotoSlice, func(placePhoto *generated.PlacePhoto) bool {
				return placePhoto.PlaceID == planPlace.PlaceID
			}),
		)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to map plan places: %w", err)
	}

	var author *models.User
	if planEntity.R.User != nil {
		author = factory.NewUserFromUserEntity(*planEntity.R.User)
	}

	plan, err := factory.NewPlanFromEntity(
		*planEntity,
		planEntity.R.PlanPlaces,
		*places,
		author,
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
			qm.Load(generated.PlanRels.User),
		},
		placeQueryModes(generated.PlanRels.PlanPlaces, generated.PlanPlaceRels.Place),
		placeQueryModes(generated.PlanRels.PlanPlaces, generated.PlanPlaceRels.Place, generated.PlaceRels.PlacePhotos),
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

	placePhotoSlice := planEntities.GetLoadedPlanPlaces().GetLoadedPlaces().GetLoadedPlacePhotos()
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
				array.Filter(placePhotoSlice, func(placePhoto *generated.PlacePhoto) bool {
					return placePhoto.PlaceID == planPlace.PlaceID
				}),
			)
		})
	})
	if err != nil {
		return nil, fmt.Errorf("failed to map plan places: %w", err)
	}

	plans, err := array.MapWithErr(planEntities, func(planEntity *generated.Plan) (*models.Plan, error) {
		var author *models.User
		if planEntity.R.User != nil {
			author = factory.NewUserFromUserEntity(*planEntity.R.User)
		}

		return factory.NewPlanFromEntity(
			*planEntity,
			planEntity.R.PlanPlaces,
			array.Flatten(*places),
			author,
		)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to map plans: %w", err)
	}

	return plans, nil
}

// TODO: ページングしない（範囲だけ指定させて、ソートも行わない）
func (p PlanRepository) SortedByLocation(ctx context.Context, location models.GeoLocation, queryCursor *string, limit int) (*[]models.Plan, *string, error) {
	minLocation, maxLocation := calculateMBR(location, defaultDistanceToSearchPlan)

	planEntities, err := generated.Plans(concatQueryMod(
		[]qm.QueryMod{
			qm.Where(
				fmt.Sprintf("MBRContains(LineString(Point(?, ?), Point(?, ?)), %s)", generated.PlanColumns.Location),
				minLocation.Longitude,
				minLocation.Latitude,
				maxLocation.Longitude,
				maxLocation.Latitude,
			),
			qm.OrderBy(fmt.Sprintf("%s %s", generated.PlanColumns.CreatedAt, "desc")),
			qm.Limit(limit),
			qm.Load(generated.PlanRels.PlanPlaces),
			qm.Load(generated.PlanRels.User),
		},
		placeQueryModes(generated.PlanRels.PlanPlaces, generated.PlanPlaceRels.Place),
		placeQueryModes(generated.PlanRels.PlanPlaces, generated.PlanPlaceRels.Place, generated.PlaceRels.PlacePhotos),
	)...).All(ctx, p.db)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find plans: %w", err)
	}

	if len(planEntities) == 0 {
		return &[]models.Plan{}, nil, nil
	}

	planCandidateSetPlaceLikeCounts, err := countPlaceLikeCounts(ctx, p.db, array.Map(planEntities, func(planEntity *generated.Plan) string {
		return planEntity.ID
	})...)
	if err != nil {
		// いいね数の取得に失敗してもエラーにしない
		p.logger.Warn("failed to count place like counts", zap.Error(err))
	}

	placePhotoSlice := planEntities.GetLoadedPlanPlaces().GetLoadedPlaces().GetLoadedPlacePhotos()
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
				array.Filter(placePhotoSlice, func(placePhoto *generated.PlacePhoto) bool {
					return placePhoto.PlaceID == planPlace.PlaceID
				}),
			)
		})
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to map plan places: %w", err)
	}

	plans, err := array.MapWithErr(planEntities, func(planEntity *generated.Plan) (*models.Plan, error) {
		var author *models.User
		if planEntity.R.User != nil {
			author = factory.NewUserFromUserEntity(*planEntity.R.User)
		}

		return factory.NewPlanFromEntity(
			*planEntity,
			planEntity.R.PlanPlaces,
			array.Flatten(*places),
			author,
		)
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to map plans: %w", err)
	}

	return plans, nil, nil
}

func newSortByCreatedAtQueryCursor(createdAt time.Time) repository.SortedByCreatedAtQueryCursor {
	return repository.SortedByCreatedAtQueryCursor(fmt.Sprintf("%d", createdAt.Unix()))
}

func parseSortByCreatedAtQueryCursor(queryCursor repository.SortedByCreatedAtQueryCursor) (*time.Time, error) {
	unixTime, err := strconv.ParseInt(string(queryCursor), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid query cursor: %s", queryCursor)
	}
	dateTime := time.Unix(unixTime, 0)
	return &dateTime, nil
}

// calculateMBR 特定の位置からの距離を元に、緯度の差分を計算する
func calculateMBR(location models.GeoLocation, distance float64) (minLocation models.GeoLocation, maxLocation models.GeoLocation) {
	// 地球の半径（メートル単位）
	const earthRadius = 6371e3

	// 1度あたりの距離（メートル単位）
	const metersPerDegree = earthRadius * math.Pi / 180

	// 緯度の増減値
	latDelta := distance / metersPerDegree

	// 経度の増減値（緯度に依存）
	lngDelta := distance / (metersPerDegree * math.Cos(math.Pi*location.Latitude/180))

	// 緯度経度の範囲を計算
	minLocation = models.GeoLocation{
		Latitude:  location.Latitude - latDelta*180/math.Pi,
		Longitude: location.Longitude - lngDelta*180/math.Pi,
	}

	maxLocation = models.GeoLocation{
		Latitude:  location.Latitude + latDelta*180/math.Pi,
		Longitude: location.Longitude + lngDelta*180/math.Pi,
	}

	return minLocation, maxLocation
}
