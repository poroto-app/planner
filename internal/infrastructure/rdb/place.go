package rdb

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/entities"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/factory"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
)

const (
	defaultMaxDistance = 1000 * 5
)

type PlaceRepository struct {
	db     *sql.DB
	logger zap.Logger
}

func NewPlaceRepository(db *sql.DB) (*PlaceRepository, error) {
	logger, err := utils.NewLogger(utils.LoggerOption{
		Tag: "PlaceRepository",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	return &PlaceRepository{
		db:     db,
		logger: *logger,
	}, nil
}

func (p PlaceRepository) GetDB() *sql.DB {
	return p.db
}

func (p PlaceRepository) SavePlacesFromGooglePlace(ctx context.Context, googlePlace models.GooglePlace) (*models.Place, error) {
	var place *models.Place

	if err := runTransaction(ctx, p, func(ctx context.Context, tx *sql.Tx) error {
		// 同じGooglePlaceIDの場所があるか確認
		exists, err := generated.GooglePlaceExists(ctx, tx, googlePlace.PlaceId)
		if err != nil {
			return fmt.Errorf("failed to check google place exists: %w", err)
		}

		if exists {
			placeSaved, err := p.findByGooglePlaceId(ctx, tx, googlePlace.PlaceId)
			if err != nil {
				return err
			}

			place = placeSaved

			return nil
		}

		// Placeを保存
		placeEntity := factory.NewPlaceEntityFromGooglePlaceEntity(googlePlace)
		if err := placeEntity.Insert(ctx, tx, boil.Infer()); err != nil {
			return fmt.Errorf("failed to insert place: %w", err)
		}

		// GooglePlaceを保存
		// Point型を保存するのにカスタムクエリを使う必要がある
		googlePlaceEntity := factory.NewGooglePlaceEntityFromGooglePlace(googlePlace, placeEntity.ID)
		if _, err := queries.Raw(
			fmt.Sprintf(
				"INSERT INTO %s (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s) VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, POINT(?, ?) )",
				generated.TableNames.GooglePlaces,
				generated.GooglePlaceColumns.GooglePlaceID,
				generated.GooglePlaceColumns.PlaceID,
				generated.GooglePlaceColumns.Name,
				generated.GooglePlaceColumns.FormattedAddress,
				generated.GooglePlaceColumns.Vicinity,
				generated.GooglePlaceColumns.PriceLevel,
				generated.GooglePlaceColumns.Rating,
				generated.GooglePlaceColumns.UserRatingsTotal,
				generated.GooglePlaceColumns.Latitude,
				generated.GooglePlaceColumns.Longitude,
				generated.GooglePlaceColumns.Location,
			),
			googlePlaceEntity.GooglePlaceID,
			googlePlaceEntity.PlaceID,
			googlePlaceEntity.Name,
			googlePlaceEntity.FormattedAddress,
			googlePlaceEntity.Vicinity,
			googlePlaceEntity.PriceLevel,
			googlePlaceEntity.Rating,
			googlePlaceEntity.UserRatingsTotal,
			googlePlaceEntity.Latitude,
			googlePlaceEntity.Longitude,
			googlePlaceEntity.Longitude,
			googlePlaceEntity.Latitude,
		).Exec(tx); err != nil {
			return fmt.Errorf("failed to insert google place: %w", err)
		}

		// GooglePlacePhotoReference, GooglePlacePhotoAttributionを保存
		if _, err := p.saveGooglePlacePhotoReferenceTx(ctx, tx, saveGooglePlacePhotoReferenceTxInput{
			GooglePlaceEntity:          &googlePlaceEntity,
			GooglePlacePhotoReferences: googlePlace.PhotoReferences,
			GooglePlaceDetail:          googlePlace.PlaceDetail,
		}); err != nil {
			return fmt.Errorf("failed to insert google place photo reference: %w", err)
		}

		// GooglePlacePhotoを保存
		if googlePlace.Photos != nil {
			for _, photo := range *googlePlace.Photos {
				if err := p.addGooglePlacePhotosTx(ctx, tx, addGooglePlacePhotosTxInput{
					GooglePlaceId:                  googlePlace.PlaceId,
					GooglePlacePhotoReferenceSlice: googlePlaceEntity.R.GooglePlacePhotoReferences,
					GooglePlacePhoto:               photo,
					SavedGooglePlacePhotoSlice:     googlePlaceEntity.R.GooglePlacePhotos,
				}); err != nil {
					return fmt.Errorf("failed to insert google place photo: %w", err)
				}
			}
		}

		// GooglePlaceTypeを保存
		googlePlaceTypeEntities := factory.NewGooglePlaceTypeSliceFromGooglePlace(googlePlace)
		if err := googlePlaceEntity.AddGooglePlaceTypes(ctx, tx, true, googlePlaceTypeEntities...); err != nil {
			return fmt.Errorf("failed to insert google place type: %w", err)
		}

		// GooglePlaceReviewを保存
		googlePlaceReviewEntities := factory.NewGooglePlaceReviewSliceFromGooglePlace(googlePlace)
		if err := googlePlaceEntity.AddGooglePlaceReviews(ctx, tx, true, googlePlaceReviewEntities...); err != nil {
			return fmt.Errorf("failed to insert google place review: %w", err)
		}

		// GooglePlaceOpeningPeriodを保存
		googlePlaceOpeningPeriodEntities := factory.NewGooglePlaceOpeningPeriodSliceFromGooglePlace(googlePlace)
		if err := googlePlaceEntity.AddGooglePlaceOpeningPeriods(ctx, tx, true, googlePlaceOpeningPeriodEntities...); err != nil {
			return fmt.Errorf("failed to insert google place opening period: %w", err)
		}

		// 自前のクエリを用いてInsertしているため、関連付けを手動で行う
		if googlePlaceEntity.R != nil {
			googlePlaceEntity.R.Place = &placeEntity
		}

		placeSaved, err := factory.NewPlaceFromEntity(
			placeEntity,
			googlePlaceEntity,
			googlePlaceEntity.R.GooglePlaceTypes,
			googlePlaceEntity.R.GooglePlacePhotoReferences,
			googlePlaceEntity.R.GooglePlacePhotoAttributions,
			googlePlaceEntity.R.GooglePlacePhotos,
			googlePlaceEntity.R.GooglePlaceReviews,
			googlePlaceEntity.R.GooglePlaceOpeningPeriods,
			0,
		)

		if err != nil {
			return fmt.Errorf("failed to convert google place entity to place: %w", err)
		}
		if placeSaved == nil {
			return fmt.Errorf("failed to convert google place entity to place")
		}

		place = placeSaved

		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to run transaction: %w", err)
	}

	return place, nil
}

func (p PlaceRepository) FindByLocation(ctx context.Context, location models.GeoLocation) ([]models.Place, error) {
	googlePlaceEntities, err := generated.GooglePlaces(
		qm.Where("ST_Distance_Sphere(POINT(?, ?), location) < ?", location.Longitude, location.Latitude, defaultMaxDistance),
		qm.Load(generated.GooglePlaceRels.Place),
		qm.Load(generated.GooglePlaceRels.GooglePlaceTypes),
		qm.Load(generated.GooglePlaceRels.GooglePlacePhotoReferences),
		qm.Load(generated.GooglePlaceRels.GooglePlacePhotos),
		qm.Load(generated.GooglePlaceRels.GooglePlacePhotoAttributions),
		qm.Load(generated.GooglePlaceRels.GooglePlaceReviews),
		qm.Load(generated.GooglePlaceRels.GooglePlaceOpeningPeriods),
	).All(ctx, p.db)
	if err != nil {
		return nil, fmt.Errorf("failed to find google places: %w", err)
	}

	planCandidateSetLikePlaceCounts, err := countPlaceLikeCounts(ctx, p.db, array.MapAndFilter(googlePlaceEntities, func(googlePlaceEntity *generated.GooglePlace) (string, bool) {
		if googlePlaceEntity == nil {
			return "", false
		}
		return googlePlaceEntity.PlaceID, true
	})...)
	if err != nil {
		// いいね数の取得に失敗してもエラーにしない
		p.logger.Warn("failed to count place like counts", zap.Error(err))
	}

	var places []models.Place
	for _, googlePlaceEntity := range googlePlaceEntities {
		if googlePlaceEntity == nil || googlePlaceEntity.R.Place == nil {
			continue
		}

		place, err := factory.NewPlaceFromEntity(
			*googlePlaceEntity.R.Place,
			*googlePlaceEntity,
			googlePlaceEntity.R.GooglePlaceTypes,
			googlePlaceEntity.R.GooglePlacePhotoReferences,
			googlePlaceEntity.R.GooglePlacePhotoAttributions,
			googlePlaceEntity.R.GooglePlacePhotos,
			googlePlaceEntity.R.GooglePlaceReviews,
			googlePlaceEntity.R.GooglePlaceOpeningPeriods,
			entities.CountLikeOfPlace(planCandidateSetLikePlaceCounts, googlePlaceEntity.PlaceID),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to convert google place googlePlaceEntity to place: %w", err)
		}
		if place == nil {
			continue
		}

		places = append(places, *place)
	}

	return places, nil
}

func (p PlaceRepository) FindByGooglePlaceID(ctx context.Context, googlePlaceID string) (*models.Place, error) {
	return p.findByGooglePlaceId(ctx, p.db, googlePlaceID)
}

func (p PlaceRepository) FindByPlanCandidateId(ctx context.Context, planCandidateId string) ([]models.Place, error) {
	planCandidateSetSearchedPlaceSlice, err := generated.PlanCandidateSetSearchedPlaces(concatQueryMod(
		[]qm.QueryMod{generated.PlanCandidateSetSearchedPlaceWhere.PlanCandidateSetID.EQ(planCandidateId)},
		placeQueryModes(generated.PlanCandidateSetSearchedPlaceRels.Place),
	)...).All(ctx, p.db)
	if err != nil {
		return nil, fmt.Errorf("failed to find plan candidate set searched places: %w", err)
	}

	planCandidateSetPlaceLikeCounts, err := countPlaceLikeCounts(ctx, p.db, array.MapAndFilter(planCandidateSetSearchedPlaceSlice, func(planCandidateSetSearchedPlace *generated.PlanCandidateSetSearchedPlace) (string, bool) {
		if planCandidateSetSearchedPlace == nil {
			return "", false
		}
		return planCandidateSetSearchedPlace.PlaceID, true
	})...)
	if err != nil {
		// いいね数の取得に失敗してもエラーにしない
		p.logger.Warn("failed to count place like counts", zap.Error(err))
	}

	var places []models.Place
	for _, planCandidateSetSearchedPlace := range planCandidateSetSearchedPlaceSlice {
		if planCandidateSetSearchedPlace == nil {
			continue
		}

		if planCandidateSetSearchedPlace.R == nil {
			panic("planCandidateSetSearchedPlace.R is nil")
		}

		if planCandidateSetSearchedPlace.R.Place == nil {
			p.logger.Warn("planCandidateSetSearchedPlace.R.Place is nil", zap.String("plan_candidate_set_searched_place_id", planCandidateSetSearchedPlace.ID))
			continue
		}

		if planCandidateSetSearchedPlace.R.Place.R == nil {
			panic("planCandidateSetSearchedPlace.R.Place.R is nil")
		}

		if len(planCandidateSetSearchedPlace.R.Place.R.GooglePlaces) == 0 {
			p.logger.Warn("planCandidateSetSearchedPlace.R.Place.R.GooglePlaces is empty", zap.String("plan_candidate_set_searched_place_id", planCandidateSetSearchedPlace.ID))
			continue
		}

		place, err := factory.NewPlaceFromEntity(
			*planCandidateSetSearchedPlace.R.Place,
			*planCandidateSetSearchedPlace.R.Place.R.GooglePlaces[0],
			planCandidateSetSearchedPlace.R.Place.R.GooglePlaces[0].R.GooglePlaceTypes,
			planCandidateSetSearchedPlace.R.Place.R.GooglePlaces[0].R.GooglePlacePhotoReferences,
			planCandidateSetSearchedPlace.R.Place.R.GooglePlaces[0].R.GooglePlacePhotoAttributions,
			planCandidateSetSearchedPlace.R.Place.R.GooglePlaces[0].R.GooglePlacePhotos,
			planCandidateSetSearchedPlace.R.Place.R.GooglePlaces[0].R.GooglePlaceReviews,
			planCandidateSetSearchedPlace.R.Place.R.GooglePlaces[0].R.GooglePlaceOpeningPeriods,
			entities.CountLikeOfPlace(planCandidateSetPlaceLikeCounts, planCandidateSetSearchedPlace.PlaceID),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to convert google place googlePlaceEntity to place: %w", err)
		}

		places = append(places, *place)
	}

	return places, nil
}

func (p PlaceRepository) SaveGooglePlacePhotos(ctx context.Context, googlePlaceId string, photos []models.GooglePlacePhoto) error {
	if err := runTransaction(ctx, p, func(ctx context.Context, tx *sql.Tx) error {
		googlePlaceEntity, err := generated.GooglePlaces(
			generated.GooglePlaceWhere.GooglePlaceID.EQ(googlePlaceId),
			qm.Load(generated.GooglePlaceRels.GooglePlacePhotoReferences),
			qm.Load(generated.GooglePlaceRels.GooglePlacePhotos),
		).One(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to find google place: %w", err)
		}

		googlePlacePhotoReferenceEntities := googlePlaceEntity.R.GooglePlacePhotoReferences
		if len(googlePlacePhotoReferenceEntities) == 0 {
			return fmt.Errorf("google place photo reference is empty")
		}

		for _, photo := range photos {
			if err := p.addGooglePlacePhotosTx(ctx, tx, addGooglePlacePhotosTxInput{
				GooglePlaceId:                  googlePlaceId,
				GooglePlacePhotoReferenceSlice: googlePlacePhotoReferenceEntities,
				GooglePlacePhoto:               photo,
				SavedGooglePlacePhotoSlice:     googlePlaceEntity.R.GooglePlacePhotos,
			}); err != nil {
				return fmt.Errorf("failed to insert google place photo: %w", err)
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to run transaction: %w", err)
	}

	return nil
}

func (p PlaceRepository) SaveGooglePlaceDetail(ctx context.Context, googlePlaceId string, googlePlaceDetail models.GooglePlaceDetail) error {
	if err := runTransaction(ctx, p, func(ctx context.Context, tx *sql.Tx) error {
		googlePlaceEntity, err := generated.GooglePlaces(
			generated.GooglePlaceWhere.GooglePlaceID.EQ(googlePlaceId),
			qm.Load(generated.GooglePlaceRels.GooglePlaceOpeningPeriods),
			qm.Load(generated.GooglePlaceRels.GooglePlaceReviews),
			qm.Load(generated.GooglePlaceRels.GooglePlacePhotoReferences),
			qm.Load(generated.GooglePlaceRels.GooglePlacePhotos),
			qm.Load(generated.GooglePlaceRels.GooglePlacePhotoAttributions),
		).One(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to find google place: %w", err)
		}

		// GooglePlaceReviewを保存
		if len(googlePlaceEntity.R.GooglePlaceReviews) == 0 {
			googlePlaceOpeningPeriodEntities := factory.NewGooglePlaceReviewSliceFromGooglePlaceDetail(googlePlaceDetail)
			if err != nil {
				return fmt.Errorf("failed to convert google place opening period: %w", err)
			}
			if err := googlePlaceEntity.AddGooglePlaceReviews(ctx, tx, true, googlePlaceOpeningPeriodEntities...); err != nil {
				return fmt.Errorf("failed to insert google place opening period: %w", err)
			}
		}

		// GooglePlaceOpeningPeriodを保存
		if len(googlePlaceEntity.R.GooglePlaceOpeningPeriods) == 0 {
			googlePlaceOpeningPeriodEntities := factory.NewGooglePlaceOpeningPeriodSliceFromGooglePlaceDetail(googlePlaceDetail, googlePlaceId)
			if err := googlePlaceEntity.AddGooglePlaceOpeningPeriods(ctx, tx, true, googlePlaceOpeningPeriodEntities...); err != nil {
				return fmt.Errorf("failed to insert google place opening period: %w", err)
			}
		}

		// GooglePlacePhotoReferenceを保存
		googlePlacePhotoReferenceSlice := factory.NewGooglePlacePhotoReferenceSliceFromGooglePlacePhotoReferences(googlePlaceDetail.PhotoReferences, googlePlaceId)
		googlePlacePhotoReferenceSlice = array.Filter(googlePlacePhotoReferenceSlice, func(googlePlacePhotoReferenceEntity *generated.GooglePlacePhotoReference) bool {
			if googlePlacePhotoReferenceEntity == nil {
				return false
			}

			_, ok := array.Find(googlePlaceEntity.R.GooglePlacePhotoReferences, func(savedGooglePlacePhotoReferenceEntity *generated.GooglePlacePhotoReference) bool {
				if savedGooglePlacePhotoReferenceEntity == nil {
					return false
				}
				return savedGooglePlacePhotoReferenceEntity.PhotoReference == googlePlacePhotoReferenceEntity.PhotoReference
			})

			// すでに保存済みのものはスキップ
			return !ok
		})
		if err := googlePlaceEntity.AddGooglePlacePhotoReferences(ctx, tx, true, googlePlacePhotoReferenceSlice...); err != nil {
			return fmt.Errorf("failed to insert google place photo reference: %w", err)
		}

		// HTMLAttributionを保存
		for _, googlePlacePhotoReference := range googlePlaceDetail.PhotoReferences {
			googlePlacePhotoReferenceEntity, ok := array.Find(googlePlaceEntity.R.GooglePlacePhotoReferences, func(savedGooglePlacePhotoReferenceEntity *generated.GooglePlacePhotoReference) bool {
				if savedGooglePlacePhotoReferenceEntity == nil {
					return false
				}
				return savedGooglePlacePhotoReferenceEntity.PhotoReference == googlePlacePhotoReference.PhotoReference
			})
			if !ok || googlePlacePhotoReferenceEntity == nil {
				return fmt.Errorf("failed to find google place photo reference entity")
			}

			// すでに紐付けがある場合はスキップ
			if len(googlePlacePhotoReferenceEntity.R.PhotoReferenceGooglePlacePhotoAttributions) > 0 {
				continue
			}

			googlePlacePhotoAttributionEntities := factory.NewGooglePlacePhotoAttributionSliceFromPhotoReference(googlePlacePhotoReference, googlePlaceId)
			if err := googlePlacePhotoReferenceEntity.AddPhotoReferenceGooglePlacePhotoAttributions(ctx, tx, true, googlePlacePhotoAttributionEntities...); err != nil {
				return fmt.Errorf("failed to insert google place photo attribution: %w", err)
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to run transaction: %w", err)
	}

	return nil
}

func (p PlaceRepository) findByGooglePlaceId(ctx context.Context, exec boil.ContextExecutor, googlePlaceId string) (*models.Place, error) {
	googlePlaceEntity, err := generated.GooglePlaces(
		generated.GooglePlaceWhere.GooglePlaceID.EQ(googlePlaceId),
		qm.Load(generated.GooglePlaceRels.Place),
		qm.Load(generated.GooglePlaceRels.GooglePlaceTypes),
		qm.Load(generated.GooglePlaceRels.GooglePlacePhotoReferences),
		qm.Load(generated.GooglePlaceRels.GooglePlacePhotos),
		qm.Load(generated.GooglePlaceRels.GooglePlacePhotoAttributions),
		qm.Load(generated.GooglePlaceRels.GooglePlaceReviews),
		qm.Load(generated.GooglePlaceRels.GooglePlaceOpeningPeriods),
	).One(ctx, exec)
	if err != nil {
		return nil, fmt.Errorf("failed to find google place: %w", err)
	}

	if googlePlaceEntity == nil || googlePlaceEntity.R.Place == nil {
		return nil, nil
	}

	planCandidateSetPlaceLikeCounts, err := countPlaceLikeCounts(ctx, exec, googlePlaceEntity.PlaceID)
	if err != nil {
		// いいね数の取得に失敗してもエラーにしない
		p.logger.Warn("failed to count place like counts", zap.Error(err))
	}

	place, err := factory.NewPlaceFromEntity(
		*googlePlaceEntity.R.Place,
		*googlePlaceEntity,
		googlePlaceEntity.R.GooglePlaceTypes,
		googlePlaceEntity.R.GooglePlacePhotoReferences,
		googlePlaceEntity.R.GooglePlacePhotoAttributions,
		googlePlaceEntity.R.GooglePlacePhotos,
		googlePlaceEntity.R.GooglePlaceReviews,
		googlePlaceEntity.R.GooglePlaceOpeningPeriods,
		entities.CountLikeOfPlace(planCandidateSetPlaceLikeCounts, googlePlaceEntity.PlaceID),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to convert google place entity to place: %w", err)
	}

	return place, nil
}

func countPlaceLikeCounts(ctx context.Context, exec boil.ContextExecutor, placeIds ...string) (*[]entities.PlanCandidateSetPlaceLikeCount, error) {
	var planCandidateSetPlaceLikeCounts []entities.PlanCandidateSetPlaceLikeCount
	if err := generated.NewQuery(
		qm.Select(
			entities.PlanCandidateSetPlaceLikeCountColumns.Name,
			fmt.Sprintf("COUNT(*) as `%s`", entities.PlanCandidateSetPlaceLikeCountColumns.LikeCount),
		),
		qm.From(generated.TableNames.PlanCandidateSetLikePlaces),
		qm.WhereIn(
			fmt.Sprintf("%s IN ?", generated.PlanCandidateSetLikePlaceColumns.PlaceID),
			array.Map(placeIds, func(placeId string) interface{} { return placeId })...,
		),
		qm.GroupBy(generated.PlanCandidateSetLikePlaceTableColumns.PlaceID),
	).Bind(ctx, exec, &planCandidateSetPlaceLikeCounts); err != nil {
		return nil, fmt.Errorf("failed to count place like counts: %w", err)
	}
	return &planCandidateSetPlaceLikeCounts, nil
}

type saveGooglePlacePhotoReferenceTxInput struct {
	GooglePlaceEntity          *generated.GooglePlace
	GooglePlacePhotoReferences []models.GooglePlacePhotoReference
	GooglePlaceDetail          *models.GooglePlaceDetail
}

// saveGooglePlacePhotoReferenceTx google_place に google_place_photo_reference, google_place_photo_attributions を紐付ける
func (p PlaceRepository) saveGooglePlacePhotoReferenceTx(ctx context.Context, tx *sql.Tx, input saveGooglePlacePhotoReferenceTxInput) (*generated.GooglePlacePhotoReferenceSlice, error) {
	// NearbySearchで取得したものとPlaceDetailで取得したものをマージする
	var googlePhotoReferences []models.GooglePlacePhotoReference
	googlePhotoReferences = input.GooglePlacePhotoReferences
	if input.GooglePlaceDetail != nil {
		// TODO: 重複を削除する
		googlePhotoReferences = append(googlePhotoReferences, input.GooglePlaceDetail.PhotoReferences...)
	}

	// GooglePlacePhotoReferenceを保存
	googlePlacePhotoReferenceEntities := factory.NewGooglePlacePhotoReferenceSliceFromGooglePlacePhotoReferences(googlePhotoReferences, input.GooglePlaceEntity.GooglePlaceID)
	if err := input.GooglePlaceEntity.AddGooglePlacePhotoReferences(ctx, tx, true, googlePlacePhotoReferenceEntities...); err != nil {
		return nil, fmt.Errorf("failed to insert google place photo reference: %w", err)
	}

	for _, googlePlacePhotoReference := range googlePhotoReferences {
		googlePlacePhotoReferenceEntity, ok := array.Find(googlePlacePhotoReferenceEntities, func(googlePlacePhotoReferenceEntity *generated.GooglePlacePhotoReference) bool {
			if googlePlacePhotoReferenceEntity == nil {
				return false
			}
			return googlePlacePhotoReferenceEntity.PhotoReference == googlePlacePhotoReference.PhotoReference
		})
		if !ok || googlePlacePhotoReferenceEntity == nil {
			return nil, fmt.Errorf("failed to find google place photo reference entity")
		}

		// HTMLAttributionを保存
		googlePlacePhotoAttributionEntities := factory.NewGooglePlacePhotoAttributionSliceFromPhotoReference(googlePlacePhotoReference, input.GooglePlaceEntity.GooglePlaceID)
		if err := googlePlacePhotoReferenceEntity.AddPhotoReferenceGooglePlacePhotoAttributions(ctx, tx, true, googlePlacePhotoAttributionEntities...); err != nil {
			return nil, fmt.Errorf("failed to insert google place photo attribution: %w", err)
		}
	}

	return &googlePlacePhotoReferenceEntities, nil
}

type addGooglePlacePhotosTxInput struct {
	GooglePlaceId                  string
	GooglePlacePhotoReferenceSlice generated.GooglePlacePhotoReferenceSlice
	GooglePlacePhoto               models.GooglePlacePhoto
	SavedGooglePlacePhotoSlice     generated.GooglePlacePhotoSlice
}

// addGooglePlacePhotosTx google_place_photo_reference に google_place_photo を紐付けする
// TODO:　複数の写真を一気に保存できるようにする
func (p PlaceRepository) addGooglePlacePhotosTx(ctx context.Context, tx *sql.Tx, input addGooglePlacePhotosTxInput) error {
	googlePlacePhotoReferenceEntity, ok := array.Find(input.GooglePlacePhotoReferenceSlice, func(googlePlacePhotoReferenceEntity *generated.GooglePlacePhotoReference) bool {
		if googlePlacePhotoReferenceEntity == nil {
			return false
		}
		return googlePlacePhotoReferenceEntity.PhotoReference == input.GooglePlacePhoto.PhotoReference
	})

	if !ok || googlePlacePhotoReferenceEntity == nil {
		return fmt.Errorf("failed to find google place photo reference entity: %s", input.GooglePlacePhoto.PhotoReference)
	}

	googlePlacePhotoEntities := factory.NewGooglePlacePhotoSliceFromDomainModel(input.GooglePlacePhoto, input.GooglePlaceId)

	googlePlacePhotoEntities = array.Filter(googlePlacePhotoEntities, func(googlePlacePhotoEntity *generated.GooglePlacePhoto) bool {
		if googlePlacePhotoEntity == nil {
			return false
		}

		// すでに保存されている場合はスキップ
		if _, ok := array.Find(input.SavedGooglePlacePhotoSlice, func(savedGooglePlacePhotoEntity *generated.GooglePlacePhoto) bool {
			if savedGooglePlacePhotoEntity == nil {
				return false
			}
			return savedGooglePlacePhotoEntity.Width == googlePlacePhotoEntity.Width && savedGooglePlacePhotoEntity.Height == googlePlacePhotoEntity.Height
		}); ok {
			p.logger.Debug(
				"skip insert google place photo because already exists",
				zap.String("photo_reference", googlePlacePhotoEntity.PhotoReference),
				zap.Int("width", googlePlacePhotoEntity.Width),
				zap.Int("height", googlePlacePhotoEntity.Height),
			)
			return false
		}

		return true
	})

	if err := googlePlacePhotoReferenceEntity.AddPhotoReferenceGooglePlacePhotos(ctx, tx, true, googlePlacePhotoEntities...); err != nil {
		return fmt.Errorf("failed to insert google place photo: %w", err)
	}
	return nil
}
