package rdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/array"
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

// Create プラン候補を作成する
// TODO: PlanCandidateSet のすべての値を保存できるようにする
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

func (p PlanCandidateRepository) Find(ctx context.Context, planCandidateId string, now time.Time) (*models.PlanCandidate, error) {
	planCandidateSetEntity, err := entities.PlanCandidateSets(concatQueryMod(
		[]qm.QueryMod{
			entities.PlanCandidateSetWhere.ID.EQ(planCandidateId),
			entities.PlanCandidateSetWhere.ExpiresAt.GT(now),
			qm.Load(entities.PlanCandidateSetRels.PlanCandidates),
			qm.Load(entities.PlanCandidateSetRels.PlanCandidateSetMetaData),
			qm.Load(entities.PlanCandidateSetRels.PlanCandidateSetMetaDataCategories),
		},
		placeQueryModes(entities.PlanCandidateSetRels.PlanCandidatePlaces, entities.PlanCandidatePlaceRels.Place),
	)...).One(ctx, p.db)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find plan candidate: %w", err)
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
			*planCandidatePlace.R.Place.R.GooglePlaces[0],
			planCandidatePlace.R.Place.R.GooglePlaces[0].R.GooglePlaceTypes,
			planCandidatePlace.R.Place.R.GooglePlaces[0].R.GooglePlacePhotoReferences,
			planCandidatePlace.R.Place.R.GooglePlaces[0].R.GooglePlacePhotoAttributions,
			planCandidatePlace.R.Place.R.GooglePlaces[0].R.GooglePlacePhotos,
			planCandidatePlace.R.Place.R.GooglePlaces[0].R.GooglePlaceReviews,
			planCandidatePlace.R.Place.R.GooglePlaces[0].R.GooglePlaceOpeningPeriods,
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
		planCandidateSetEntity.R.PlanCandidatePlaces,
		places,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create plan candidate: %w", err)
	}

	return planCandidateSet, nil
}

func (p PlanCandidateRepository) FindPlan(ctx context.Context, planCandidateId string, planId string) (*models.Plan, error) {
	planCandidate, err := entities.PlanCandidates(concatQueryMod(
		[]qm.QueryMod{
			entities.PlanCandidateWhere.ID.EQ(planId),
			entities.PlanCandidateWhere.PlanCandidateSetID.EQ(planCandidateId),
		},
		placeQueryModes(entities.PlanCandidateRels.PlanCandidatePlaces, entities.PlanCandidatePlaceRels.Place),
	)...).One(ctx, p.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find plan candidate: %w", err)
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
			*planCandidatePlace.R.Place.R.GooglePlaces[0],
			planCandidatePlace.R.Place.R.GooglePlaces[0].R.GooglePlaceTypes,
			planCandidatePlace.R.Place.R.GooglePlaces[0].R.GooglePlacePhotoReferences,
			planCandidatePlace.R.Place.R.GooglePlaces[0].R.GooglePlacePhotoAttributions,
			planCandidatePlace.R.Place.R.GooglePlaces[0].R.GooglePlacePhotos,
			planCandidatePlace.R.Place.R.GooglePlaces[0].R.GooglePlaceReviews,
			planCandidatePlace.R.Place.R.GooglePlaces[0].R.GooglePlaceOpeningPeriods,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create place: %w", err)
		}

		places = append(places, *place)
	}

	plan, err := factory.NewPlanCandidateFromEntity(*planCandidate, planCandidate.R.PlanCandidatePlaces, places)
	if err != nil {
		return nil, fmt.Errorf("failed to create plan candidate: %w", err)
	}

	return plan, nil
}

func (p PlanCandidateRepository) FindExpiredBefore(ctx context.Context, expiresAt time.Time) (*[]string, error) {
	planCandidateSetSlice, err := entities.PlanCandidateSets(entities.PlanCandidateSetWhere.ExpiresAt.LT(expiresAt)).All(ctx, p.db)
	if err != nil {
		return nil, fmt.Errorf("failed to find expired plan candidate sets: %w", err)
	}

	var planCandidateIds []string
	for _, planCandidateSet := range planCandidateSetSlice {
		planCandidateIds = append(planCandidateIds, planCandidateSet.ID)
	}

	return &planCandidateIds, nil
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
	if err := runTransaction(ctx, p, func(ctx context.Context, tx *sql.Tx) error {
		planCandidateEntity, err := entities.PlanCandidates(
			entities.PlanCandidateWhere.ID.EQ(planId),
			entities.PlanCandidateWhere.PlanCandidateSetID.EQ(planCandidateId),
			qm.Load(entities.PlanCandidateRels.PlanCandidatePlaces),
		).One(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to get plan candidate: %w", err)
		}

		if planCandidateEntity.R == nil {
			panic("planCandidateEntity.R is nil")
		}

		planCandidatePlaceSlice := planCandidateEntity.R.PlanCandidatePlaces
		planCandidatePlaceToDelete, ok := array.Find(planCandidatePlaceSlice, func(planCandidatePlace *entities.PlanCandidatePlace) bool {
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
		planCandidateEntity, err := entities.PlanCandidates(
			entities.PlanCandidateWhere.ID.EQ(planId),
			entities.PlanCandidateWhere.PlanCandidateSetID.EQ(planCandidate),
			qm.Load(entities.PlanCandidateRels.PlanCandidatePlaces),
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
			if _, ok := array.Find(planCandidatePlaceSlice, func(planCandidatePlace *entities.PlanCandidatePlace) bool {
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
			planCandidatePlace, ok := array.Find(planCandidatePlaceSlice, func(planCandidatePlace *entities.PlanCandidatePlace) bool {
				if planCandidatePlace == nil {
					return false
				}
				return planCandidatePlace.PlaceID == placeId
			})
			if !ok {
				return fmt.Errorf("invalid placeId %s", placeId)
			}

			planCandidatePlace.SortOrder = i
			if _, err := planCandidatePlace.Update(ctx, tx, boil.Whitelist(entities.PlanCandidatePlaceColumns.SortOrder)); err != nil {
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
		savedPlanCandidateSetMetaDataEntity, err := entities.PlanCandidateSetMetaData(entities.PlanCandidateSetMetaDatumWhere.PlanCandidateSetID.EQ(planCandidateId)).One(ctx, tx)
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
			if _, err := entities.PlanCandidateSetMetaDataCategories(entities.PlanCandidateSetMetaDataCategoryWhere.PlanCandidateSetID.EQ(planCandidateId)).DeleteAll(ctx, tx); err != nil {
				return fmt.Errorf("failed to delete plan candidate set categories: %w", err)
			}

			// カテゴリを登録
			// TODO: BatchInsertする
			planCandidateSetMetaDataCategorySlice := factory.NewPlanCandidateSetMetaDataCategorySliceFromDomainModel(meta.CategoriesPreferred, meta.CategoriesRejected, planCandidateId)
			for _, planCandidateSetMetaDataCategory := range planCandidateSetMetaDataCategorySlice {
				if err := planCandidateSetMetaDataCategory.Insert(ctx, tx, boil.Infer()); err != nil {
					return fmt.Errorf("failed to insert plan candidate set meta data category: %w", err)
				}
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to run transaction: %w", err)
	}

	return nil
}

func (p PlanCandidateRepository) ReplacePlace(ctx context.Context, planCandidateId string, planId string, placeIdToBeReplaced string, placeToReplace models.Place) error {
	if err := runTransaction(ctx, p, func(ctx context.Context, tx *sql.Tx) error {
		planCandidatePlaceEntity, err := entities.PlanCandidatePlaces(
			entities.PlanCandidatePlaceWhere.PlanCandidateSetID.EQ(planCandidateId),
			entities.PlanCandidatePlaceWhere.PlanCandidateID.EQ(planId),
			entities.PlanCandidatePlaceWhere.PlaceID.EQ(placeIdToBeReplaced),
		).One(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to get plan candidate place: %w", err)
		}

		planCandidatePlaceEntity.PlaceID = placeToReplace.Id

		if _, err := planCandidatePlaceEntity.Update(ctx, tx, boil.Whitelist(entities.PlanCandidatePlaceColumns.PlaceID)); err != nil {
			return fmt.Errorf("failed to update plan candidate place: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to run transaction: %w", err)
	}

	return nil
}

func (p PlanCandidateRepository) DeleteAll(ctx context.Context, planCandidateIds []string) error {
	if err := runTransaction(ctx, p, func(ctx context.Context, tx *sql.Tx) error {
		// プラン候補場所を削除
		if _, err := entities.PlanCandidatePlaces(entities.PlanCandidatePlaceWhere.PlanCandidateSetID.IN(planCandidateIds)).DeleteAll(ctx, tx); err != nil {
			return fmt.Errorf("failed to delete plan candidate places: %w", err)
		}

		// 検索履歴を削除
		if _, err := entities.PlanCandidateSetSearchedPlaces(entities.PlanCandidateSetSearchedPlaceWhere.PlanCandidateSetID.IN(planCandidateIds)).DeleteAll(ctx, tx); err != nil {
			return fmt.Errorf("failed to delete plan candidate set searched places: %w", err)
		}

		// プラン候補を削除
		if _, err := entities.PlanCandidates(entities.PlanCandidateWhere.PlanCandidateSetID.IN(planCandidateIds)).DeleteAll(ctx, tx); err != nil {
			return fmt.Errorf("failed to delete plan candidates: %w", err)
		}

		// プラン候補メタデータを削除
		if _, err := entities.PlanCandidateSetMetaData(entities.PlanCandidateSetMetaDatumWhere.PlanCandidateSetID.IN(planCandidateIds)).DeleteAll(ctx, tx); err != nil {
			return fmt.Errorf("failed to delete plan candidate set meta data: %w", err)
		}
		if _, err := entities.PlanCandidateSetMetaDataCategories(entities.PlanCandidateSetMetaDataCategoryWhere.PlanCandidateSetID.IN(planCandidateIds)).DeleteAll(ctx, tx); err != nil {
			return fmt.Errorf("failed to delete plan candidate set meta data categories: %w", err)
		}

		// プラン候補一覧を削除
		if _, err := entities.PlanCandidateSets(entities.PlanCandidateSetWhere.ID.IN(planCandidateIds)).DeleteAll(ctx, tx); err != nil {
			return fmt.Errorf("failed to delete plan candidate sets: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to run transaction: %w", err)
	}

	return nil
}

func (p PlanCandidateRepository) UpdateLikeToPlaceInPlanCandidate(ctx context.Context, planCandidateId string, placeId string, like bool) error {
	//TODO implement me
	panic("implement me")
}
