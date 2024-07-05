package rdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/entities"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/factory"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
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

// Create プラン候補を作成する
// TODO: PlanCandidateSet のすべての値を保存できるようにする
func (p PlanCandidateRepository) Create(cxt context.Context, planCandidateSetId string, expiresAt time.Time) error {
	if err := runTransaction(cxt, p, func(ctx context.Context, tx *sql.Tx) error {
		planCandidateSetEntity := generated.PlanCandidateSet{ID: planCandidateSetId, ExpiresAt: expiresAt}
		if err := planCandidateSetEntity.Insert(ctx, tx, boil.Infer()); err != nil {
			return fmt.Errorf("failed to insert plan candidate: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to run transaction: %w", err)
	}
	return nil
}

func (p PlanCandidateRepository) Find(ctx context.Context, planCandidateSetId string, now time.Time) (*models.PlanCandidateSet, error) {
	planCandidateSetEntity, err := generated.PlanCandidateSets(concatQueryMod(
		[]qm.QueryMod{
			generated.PlanCandidateSetWhere.ID.EQ(planCandidateSetId),
			generated.PlanCandidateSetWhere.ExpiresAt.GT(now),
			qm.Load(generated.PlanCandidateSetRels.PlanCandidates),
			qm.Load(generated.PlanCandidateSetRels.PlanCandidateSetMetaData),
			qm.Load(generated.PlanCandidateSetRels.PlanCandidateSetMetaDataCategories),
			qm.Load(generated.PlanCandidateSetRels.PlanCandidateSetLikePlaces),
			qm.Load(generated.PlanCandidateSetRels.PlanCandidateSetMetaDataCreateByCategories),
		},
		placeQueryModes(generated.PlanCandidateSetRels.PlanCandidatePlaces, generated.PlanCandidatePlaceRels.Place),
	)...).One(ctx, p.db)
	if err != nil {
		// TODO: エラーにする
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find plan candidate: %w", err)
	}

	if planCandidateSetEntity.R == nil {
		panic("planCandidateSetEntity.R is nil")
	}

	planCandidateSetPlaceLikeCounts, err := countPlaceLikeCounts(ctx, p.db, array.Map(planCandidateSetEntity.R.PlanCandidatePlaces, func(planCandidatePlace *generated.PlanCandidatePlace) string {
		return planCandidatePlace.PlaceID
	})...)
	if err != nil {
		// いいね数の取得に失敗してもエラーにしない
		p.logger.Warn("failed to count place like counts", zap.Error(err))
	}

	var places []models.Place
	for _, planCandidatePlace := range planCandidateSetEntity.R.PlanCandidatePlaces {
		if planCandidatePlace.R.Place == nil {
			p.logger.Warn("planCandidatePlace.R.Place is nil", zap.String("planCandidatePlaceId", planCandidatePlace.ID))
			continue
		}

		if planCandidatePlace.R.Place.R == nil {
			panic("planCandidatePlace.R.Place.R is nil")
		}

		if len(planCandidatePlace.R.Place.R.GooglePlaces) == 0 || planCandidatePlace.R.Place.R.GooglePlaces[0] == nil {
			p.logger.Warn("planCandidatePlace.R.Place.R.GooglePlaces is empty", zap.String("planCandidatePlaceId", planCandidatePlace.ID))
			continue
		}

		place, err := factory.NewPlaceFromEntity(
			*planCandidatePlace.R.Place,
			planCandidatePlace.R.Place.R.PlacePhotos,
			*planCandidatePlace.R.Place.R.GooglePlaces[0],
			planCandidatePlace.R.Place.R.GooglePlaces[0].R.GooglePlaceTypes,
			planCandidatePlace.R.Place.R.GooglePlaces[0].R.GooglePlacePhotoReferences,
			planCandidatePlace.R.Place.R.GooglePlaces[0].R.GooglePlacePhotoAttributions,
			planCandidatePlace.R.Place.R.GooglePlaces[0].R.GooglePlacePhotos,
			planCandidatePlace.R.Place.R.GooglePlaces[0].R.GooglePlaceReviews,
			planCandidatePlace.R.Place.R.GooglePlaces[0].R.GooglePlaceOpeningPeriods,
			entities.CountLikeOfPlace(planCandidateSetPlaceLikeCounts, planCandidatePlace.PlaceID),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create place: %w", err)
		}

		places = append(places, *place)
	}

	planCandidateSet, err := factory.NewPlanCandidateSetFromEntity(
		*planCandidateSetEntity,
		planCandidateSetEntity.R.PlanCandidates,
		planCandidateSetEntity.R.PlanCandidateSetMetaData,
		planCandidateSetEntity.R.PlanCandidateSetMetaDataCategories,
		planCandidateSetEntity.R.PlanCandidateSetMetaDataCreateByCategories,
		planCandidateSetEntity.R.PlanCandidatePlaces,
		planCandidateSetEntity.R.PlanCandidateSetLikePlaces,
		places,
		nil, // TODO: ユーザー情報を取得できるようにする
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create plan candidate: %w", err)
	}

	return planCandidateSet, nil
}

func (p PlanCandidateRepository) FindPlan(ctx context.Context, planCandidateId string, planId string) (*models.Plan, error) {
	planCandidate, err := generated.PlanCandidates(concatQueryMod(
		[]qm.QueryMod{generated.PlanCandidateWhere.ID.EQ(planId)},
		placeQueryModes(generated.PlanCandidateRels.PlanCandidatePlaces, generated.PlanCandidatePlaceRels.Place),
	)...).One(ctx, p.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find plan candidate: %w", err)
	}

	if planCandidate.R == nil {
		panic("planCandidate.R is nil")
	}

	planCandidateSetPlaceLikeCounts, err := countPlaceLikeCounts(
		ctx,
		p.db,
		array.Map(planCandidate.R.PlanCandidatePlaces, func(planCandidatePlace *generated.PlanCandidatePlace) string {
			return planCandidatePlace.PlaceID
		})...)
	if err != nil {
		// いいね数の取得に失敗してもエラーにしない
		p.logger.Warn("failed to count place like counts", zap.Error(err))
	}

	var places []models.Place
	for _, planCandidatePlace := range planCandidate.R.PlanCandidatePlaces {
		if planCandidatePlace.R.Place == nil {
			p.logger.Warn("planCandidatePlace.R.Place is nil", zap.String("planCandidatePlaceId", planCandidatePlace.ID))
			continue
		}

		if planCandidatePlace.R.Place.R == nil {
			panic("planCandidatePlace.R.Place.R is nil")
		}

		if len(planCandidatePlace.R.Place.R.GooglePlaces) == 0 || planCandidatePlace.R.Place.R.GooglePlaces[0] == nil {
			p.logger.Warn("planCandidatePlace.R.Place.R.GooglePlaces is empty", zap.String("planCandidatePlaceId", planCandidatePlace.ID))
			continue
		}

		place, err := factory.NewPlaceFromEntity(
			*planCandidatePlace.R.Place,
			planCandidatePlace.R.Place.R.PlacePhotos,
			*planCandidatePlace.R.Place.R.GooglePlaces[0],
			planCandidatePlace.R.Place.R.GooglePlaces[0].R.GooglePlaceTypes,
			planCandidatePlace.R.Place.R.GooglePlaces[0].R.GooglePlacePhotoReferences,
			planCandidatePlace.R.Place.R.GooglePlaces[0].R.GooglePlacePhotoAttributions,
			planCandidatePlace.R.Place.R.GooglePlaces[0].R.GooglePlacePhotos,
			planCandidatePlace.R.Place.R.GooglePlaces[0].R.GooglePlaceReviews,
			planCandidatePlace.R.Place.R.GooglePlaces[0].R.GooglePlaceOpeningPeriods,
			entities.CountLikeOfPlace(planCandidateSetPlaceLikeCounts, planCandidatePlace.PlaceID),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create place: %w", err)
		}

		places = append(places, *place)
	}

	// TODO: ユーザー情報を取得できるようにする
	plan, err := factory.NewPlanCandidateFromEntity(*planCandidate, planCandidate.R.PlanCandidatePlaces, places, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create plan candidate: %w", err)
	}

	return plan, nil
}

func (p PlanCandidateRepository) AddPlan(ctx context.Context, planCandidateId string, plans ...models.Plan) error {
	if err := runTransaction(ctx, p, func(ctx context.Context, tx *sql.Tx) error {
		planCountInPlanCandidateSet, err := generated.PlanCandidates(generated.PlanCandidateWhere.PlanCandidateSetID.EQ(planCandidateId)).Count(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to count plan candidate: %w", err)
		}

		var planCandidateSlice generated.PlanCandidateSlice
		var planCandidatePlaceSlice generated.PlanCandidatePlaceSlice
		for iPlan, plan := range plans {
			planSortOrder := iPlan + int(planCountInPlanCandidateSet)
			planCandidateEntity := factory.PlanCandidateEntityFromDomainModel(plan, planCandidateId, planSortOrder)
			planCandidateSlice = append(planCandidateSlice, &planCandidateEntity)

			planCandidatePlaceSlice = append(
				planCandidatePlaceSlice,
				factory.NewPlanCandidatePlaceSliceFromDomainModel(plan.Places, planCandidateId, plan.Id)...,
			)
		}

		if _, err := planCandidateSlice.InsertAll(ctx, tx, boil.Infer()); err != nil {
			return fmt.Errorf("failed to insert plan candidate: %w", err)
		}

		if _, err := planCandidatePlaceSlice.InsertAll(ctx, tx, boil.Infer()); err != nil {
			return fmt.Errorf("failed to insert plan candidate place: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to run transaction: %w", err)
	}

	return nil
}

func (p PlanCandidateRepository) AddPlaceToPlan(ctx context.Context, planCandidateId string, planId string, previousPlaceId string, place models.Place) error {
	if err := runTransaction(ctx, p, func(ctx context.Context, tx *sql.Tx) error {
		planCandidatePlaceSlice, err := generated.
			PlanCandidatePlaces(generated.PlanCandidatePlaceWhere.PlanCandidateSetID.EQ(planCandidateId)).
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
				if _, err := planCandidatePlace.Update(ctx, tx, boil.Whitelist(generated.PlanCandidatePlaceColumns.SortOrder)); err != nil {
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
	if err := runTransaction(ctx, p, func(ctx context.Context, tx *sql.Tx) error {
		planCandidateEntity, err := generated.PlanCandidates(
			generated.PlanCandidateWhere.ID.EQ(planId),
			generated.PlanCandidateWhere.PlanCandidateSetID.EQ(planCandidateId),
			qm.Load(generated.PlanCandidateRels.PlanCandidatePlaces),
		).One(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to get plan candidate: %w", err)
		}

		if planCandidateEntity.R == nil {
			panic("planCandidateEntity.R is nil")
		}

		planCandidatePlaceSlice := planCandidateEntity.R.PlanCandidatePlaces
		planCandidatePlaceToDelete, ok := array.Find(planCandidatePlaceSlice, func(planCandidatePlace *generated.PlanCandidatePlace) bool {
			if planCandidatePlace == nil {
				return false
			}
			return planCandidatePlace.PlaceID == placeId
		})
		if !ok {
			// もともと存在しない場所を削除しようとした場合は何もしない
			return nil
		}

		if _, err := planCandidatePlaceToDelete.Delete(ctx, tx); err != nil {
			return fmt.Errorf("failed to delete plan candidate place: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to run transaction: %w", err)
	}

	return nil
}

func (p PlanCandidateRepository) UpdatePlacesOrder(ctx context.Context, planId string, planCandidate string, placeIdsOrdered []string) error {
	if err := runTransaction(ctx, p, func(ctx context.Context, tx *sql.Tx) error {
		planCandidateEntity, err := generated.PlanCandidates(
			generated.PlanCandidateWhere.ID.EQ(planId),
			generated.PlanCandidateWhere.PlanCandidateSetID.EQ(planCandidate),
			qm.Load(generated.PlanCandidateRels.PlanCandidatePlaces),
		).One(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to get plan candidate: %w", err)
		}

		if planCandidateEntity.R == nil {
			panic("planCandidateEntity.R is nil")
		}

		planCandidatePlaceSlice := planCandidateEntity.R.PlanCandidatePlaces

		// 場所のID一覧に過不足がないかを確認
		if len(placeIdsOrdered) != len(planCandidatePlaceSlice) {
			return fmt.Errorf("invalid placeIdsOrdered length")
		}

		// すべての場所のIDが存在するかを確認
		for _, placeId := range placeIdsOrdered {
			if _, ok := array.Find(planCandidatePlaceSlice, func(planCandidatePlace *generated.PlanCandidatePlace) bool {
				if planCandidatePlace == nil {
					return false
				}
				return planCandidatePlace.PlaceID == placeId
			}); !ok {
				return fmt.Errorf("invalid placeId %s", placeId)
			}
		}

		// 場所の順序を更新
		for i, placeId := range placeIdsOrdered {
			planCandidatePlace, ok := array.Find(planCandidatePlaceSlice, func(planCandidatePlace *generated.PlanCandidatePlace) bool {
				if planCandidatePlace == nil {
					return false
				}
				return planCandidatePlace.PlaceID == placeId
			})
			if !ok {
				return fmt.Errorf("invalid placeId %s", placeId)
			}

			planCandidatePlace.SortOrder = i
			if _, err := planCandidatePlace.Update(ctx, tx, boil.Whitelist(generated.PlanCandidatePlaceColumns.SortOrder)); err != nil {
				return fmt.Errorf("failed to update plan candidate place: %w", err)
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to run transaction: %w", err)
	}

	return nil
}

func (p PlanCandidateRepository) UpdatePlanCandidateMetaData(ctx context.Context, planCandidateId string, meta models.PlanCandidateMetaData) error {
	if err := runTransaction(ctx, p, func(ctx context.Context, tx *sql.Tx) error {
		savedPlanCandidateSetMetaDataEntity, err := generated.PlanCandidateSetMetaData(generated.PlanCandidateSetMetaDatumWhere.PlanCandidateSetID.EQ(planCandidateId)).One(ctx, tx)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("failed to get plan candidate set meta data: %w", err)
			}
		}

		if savedPlanCandidateSetMetaDataEntity == nil {
			// 保存されていない場合は新規作成
			planCandidateSetMetaDataEntity := factory.NewPlanCandidateMetaDataFromDomainModel(meta, planCandidateId)
			if err := planCandidateSetMetaDataEntity.Insert(ctx, tx, boil.Infer()); err != nil {
				return fmt.Errorf("failed to insert plan candidate set meta data: %w", err)
			}
		} else {
			// 保存されている場合は更新
			planCandidateMetaDataEntity := factory.NewPlanCandidateMetaDataFromDomainModel(meta, planCandidateId)
			planCandidateMetaDataEntity.ID = savedPlanCandidateSetMetaDataEntity.ID
			if _, err := planCandidateMetaDataEntity.Update(ctx, tx, boil.Infer()); err != nil {
				return fmt.Errorf("failed to upsert plan candidate set meta data: %w", err)
			}
		}

		// カテゴリを更新
		if meta.CategoriesRejected != nil || meta.CategoriesPreferred != nil {
			// すでに登録されているカテゴリを削除
			if _, err := generated.PlanCandidateSetMetaDataCategories(generated.PlanCandidateSetMetaDataCategoryWhere.PlanCandidateSetID.EQ(planCandidateId)).DeleteAll(ctx, tx); err != nil {
				return fmt.Errorf("failed to delete plan candidate set categories: %w", err)
			}

			// カテゴリを登録
			planCandidateSetMetaDataCategorySlice := factory.NewPlanCandidateSetMetaDataCategorySliceFromDomainModel(meta.CategoriesPreferred, meta.CategoriesRejected, planCandidateId)
			if _, err := planCandidateSetMetaDataCategorySlice.InsertAll(ctx, tx, boil.Infer()); err != nil {
				return fmt.Errorf("failed to insert plan candidate set meta data category: %w", err)
			}
		}

		// カテゴリからプランを作成した場合の情報を更新
		if meta.CreateByCategoryMetaData != nil {
			// すでに保存されている場合は削除
			if _, err := generated.PlanCandidateSetMetaDataCreateByCategories(generated.PlanCandidateSetMetaDataCreateByCategoryWhere.PlanCandidateSetID.EQ(planCandidateId)).DeleteAll(ctx, tx); err != nil {
				return fmt.Errorf("failed to get create by category meta data: %w", err)
			}

			planCandidateSetMetaDataCreateByCategoryEntity := factory.NewPlanCandidateMetaDataCreateByCategoryFromDomainModel(planCandidateId, *meta.CreateByCategoryMetaData)
			if err := planCandidateSetMetaDataCreateByCategoryEntity.Insert(ctx, tx, boil.Infer()); err != nil {
				return fmt.Errorf("failed to insert plan candidate set meta data create by category: %w", err)
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to run transaction: %w", err)
	}

	return nil
}

func (p PlanCandidateRepository) UpdateIsPlaceSearched(ctx context.Context, planCandidateId string, isPlaceSearched bool) error {
	if err := runTransaction(ctx, p, func(ctx context.Context, tx *sql.Tx) error {
		planCandidateSetEntity, err := generated.PlanCandidateSets(
			generated.PlanCandidateSetWhere.ID.EQ(planCandidateId),
		).One(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to get plan candidate set: %w", err)
		}

		planCandidateSetEntity.IsPlaceSearched = isPlaceSearched
		if _, err := planCandidateSetEntity.Update(ctx, tx, boil.Whitelist(generated.PlanCandidateSetColumns.IsPlaceSearched)); err != nil {
			return fmt.Errorf("failed to update plan candidate set: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to run transaction: %w", err)
	}

	return nil
}

func (p PlanCandidateRepository) ReplacePlace(ctx context.Context, planCandidateId string, planId string, placeIdToBeReplaced string, placeToReplace models.Place) error {
	if err := runTransaction(ctx, p, func(ctx context.Context, tx *sql.Tx) error {
		planCandidatePlaceEntity, err := generated.PlanCandidatePlaces(
			generated.PlanCandidatePlaceWhere.PlanCandidateSetID.EQ(planCandidateId),
			generated.PlanCandidatePlaceWhere.PlanCandidateID.EQ(planId),
			generated.PlanCandidatePlaceWhere.PlaceID.EQ(placeIdToBeReplaced),
		).One(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to get plan candidate place: %w", err)
		}

		planCandidatePlaceEntity.PlaceID = placeToReplace.Id

		if _, err := planCandidatePlaceEntity.Update(ctx, tx, boil.Whitelist(generated.PlanCandidatePlaceColumns.PlaceID)); err != nil {
			return fmt.Errorf("failed to update plan candidate place: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to run transaction: %w", err)
	}

	return nil
}

func (p PlanCandidateRepository) UpdateLikeToPlaceInPlanCandidateSet(ctx context.Context, planCandidateId string, placeId string, like bool) error {
	if err := runTransaction(ctx, p, func(ctx context.Context, tx *sql.Tx) error {
		if like {
			isPlanCandidateSetLikePlaceEntityExist, err := generated.PlanCandidateSetLikePlaces(
				generated.PlanCandidateSetLikePlaceWhere.PlanCandidateSetID.EQ(planCandidateId),
				generated.PlanCandidateSetLikePlaceWhere.PlaceID.EQ(placeId),
			).Exists(ctx, tx)
			if err != nil {
				return fmt.Errorf("failed to check plan candidate set like place existence: %w", err)
			}

			if isPlanCandidateSetLikePlaceEntityExist {
				// すでに存在する場合は何もしない
				return nil
			}

			planCandidateSetPlaceEntity := generated.PlanCandidateSetLikePlace{
				ID:                 uuid.New().String(),
				PlanCandidateSetID: planCandidateId,
				PlaceID:            placeId,
			}
			if err := planCandidateSetPlaceEntity.Insert(ctx, tx, boil.Infer()); err != nil {
				return fmt.Errorf("failed to insert plan candidate set like place: %w", err)
			}
		} else {
			if _, err := generated.PlanCandidateSetLikePlaces(
				generated.PlanCandidateSetLikePlaceWhere.PlanCandidateSetID.EQ(planCandidateId),
				generated.PlanCandidateSetLikePlaceWhere.PlaceID.EQ(placeId),
			).DeleteAll(ctx, tx); err != nil {
				return fmt.Errorf("failed to delete plan candidate set like place: %w", err)
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to run transaction: %w", err)
	}

	return nil
}
