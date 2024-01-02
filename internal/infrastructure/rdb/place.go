package rdb

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/entities"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/factory"
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
		exists, err := entities.GooglePlaceExists(ctx, tx, googlePlace.PlaceId)
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
		googlePlaceEntity := factory.NewGooglePlaceEntityFromGooglePlace(googlePlace, placeEntity.ID)
		if err := placeEntity.AddGooglePlaces(ctx, tx, true, &googlePlaceEntity); err != nil {
			return fmt.Errorf("failed to insert google place: %w", err)
		}

		// GooglePlacePhotoReference, GooglePlacePhotoAttributionを保存
		if _, err := p.saveGooglePlacePhotoReferenceTx(ctx, tx, saveGooglePlacePhotoReferenceTxInput{
			GooglePlacePhotoReferences: googlePlace.PhotoReferences,
			GooglePlaceDetail:          googlePlace.PlaceDetail,
			GooglePlaceEntity:          googlePlaceEntity,
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
		googlePlaceOpeningPeriodEntities, err := factory.NewGooglePlaceOpeningPeriodSliceFromGooglePlace(googlePlace)
		if err != nil {
			return fmt.Errorf("failed to convert google place opening period: %w", err)
		}
		if err := googlePlaceEntity.AddGooglePlaceOpeningPeriods(ctx, tx, true, googlePlaceOpeningPeriodEntities...); err != nil {
			return fmt.Errorf("failed to insert google place opening period: %w", err)
		}

		placeSaved, err := factory.NewPlaceFromGooglePlaceEntity(googlePlaceEntity)
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
	entities, err := entities.GooglePlaces(
		qm.Where("ST_Distance_Sphere(POINT(?, ?), location) < ?", location.Longitude, location.Latitude, defaultMaxDistance),
		qm.Load(entities.GooglePlaceRels.Place),
		qm.Load(entities.GooglePlaceRels.GooglePlaceTypes),
		qm.Load(entities.GooglePlaceRels.GooglePlacePhotoReferences),
		qm.Load(entities.GooglePlaceRels.GooglePlacePhotos),
		qm.Load(entities.GooglePlaceRels.GooglePlacePhotoAttributions),
		qm.Load(entities.GooglePlaceRels.GooglePlaceReviews),
		qm.Load(entities.GooglePlaceRels.GooglePlaceOpeningPeriods),
	).All(ctx, p.db)
	if err != nil {
		return nil, fmt.Errorf("failed to find google places: %w", err)
	}

	var places []models.Place
	for _, entity := range entities {
		if entity == nil {
			continue
		}

		place, err := factory.NewPlaceFromGooglePlaceEntity(*entity)
		if err != nil {
			return nil, fmt.Errorf("failed to convert google place entity to place: %w", err)
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
	//TODO implement me
	panic("implement me")
}

func (p PlaceRepository) SaveGooglePlacePhotos(ctx context.Context, googlePlaceId string, photos []models.GooglePlacePhoto) error {
	if err := runTransaction(ctx, p, func(ctx context.Context, tx *sql.Tx) error {
		googlePlaceEntity, err := entities.GooglePlaces(
			entities.GooglePlaceWhere.GooglePlaceID.EQ(googlePlaceId),
			qm.Load(entities.GooglePlaceRels.GooglePlacePhotoReferences),
			qm.Load(entities.GooglePlaceRels.GooglePlacePhotos),
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
		googlePlaceEntity, err := entities.GooglePlaces(
			entities.GooglePlaceWhere.GooglePlaceID.EQ(googlePlaceId),
			qm.Load(entities.GooglePlaceRels.GooglePlaceOpeningPeriods),
			qm.Load(entities.GooglePlaceRels.GooglePlaceReviews),
			qm.Load(entities.GooglePlaceRels.GooglePlacePhotoReferences),
			qm.Load(entities.GooglePlaceRels.GooglePlacePhotos),
			qm.Load(entities.GooglePlaceRels.GooglePlacePhotoAttributions),
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
			googlePlaceOpeningPeriodEntities, err := factory.NewGooglePlaceOpeningPeriodSliceFromGooglePlaceDetail(googlePlaceDetail, googlePlaceId)
			if err != nil {
				return fmt.Errorf("failed to convert google place opening period: %w", err)
			}
			if err := googlePlaceEntity.AddGooglePlaceOpeningPeriods(ctx, tx, true, googlePlaceOpeningPeriodEntities...); err != nil {
				return fmt.Errorf("failed to insert google place opening period: %w", err)
			}
		}

		// GooglePlacePhotoReferenceを保存
		googlePlacePhotoReferenceSlice := factory.NewGooglePlacePhotoReferenceSliceFromGooglePlacePhotoReferences(googlePlaceDetail.PhotoReferences)
		googlePlacePhotoReferenceSlice = array.Filter(googlePlacePhotoReferenceSlice, func(googlePlacePhotoReferenceEntity *entities.GooglePlacePhotoReference) bool {
			if googlePlacePhotoReferenceEntity == nil {
				return false
			}

			_, ok := array.Find(googlePlaceEntity.R.GooglePlacePhotoReferences, func(savedGooglePlacePhotoReferenceEntity *entities.GooglePlacePhotoReference) bool {
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
			googlePlacePhotoReferenceEntity, ok := array.Find(googlePlaceEntity.R.GooglePlacePhotoReferences, func(savedGooglePlacePhotoReferenceEntity *entities.GooglePlacePhotoReference) bool {
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
	googlePlaceEntity, err := entities.GooglePlaces(
		entities.GooglePlaceWhere.GooglePlaceID.EQ(googlePlaceId),
		qm.Load(entities.GooglePlaceRels.Place),
		qm.Load(entities.GooglePlaceRels.GooglePlaceTypes),
		qm.Load(entities.GooglePlaceRels.GooglePlacePhotoReferences),
		qm.Load(entities.GooglePlaceRels.GooglePlacePhotos),
		qm.Load(entities.GooglePlaceRels.GooglePlacePhotoAttributions),
		qm.Load(entities.GooglePlaceRels.GooglePlaceReviews),
		qm.Load(entities.GooglePlaceRels.GooglePlaceOpeningPeriods),
	).One(ctx, exec)
	if err != nil {
		return nil, fmt.Errorf("failed to find google place: %w", err)
	}

	if googlePlaceEntity == nil {
		return nil, nil
	}

	place, err := factory.NewPlaceFromGooglePlaceEntity(*googlePlaceEntity)
	if err != nil {
		return nil, fmt.Errorf("failed to convert google place entity to place: %w", err)
	}

	return place, nil
}

type saveGooglePlacePhotoReferenceTxInput struct {
	GooglePlaceEntity          entities.GooglePlace
	GooglePlacePhotoReferences []models.GooglePlacePhotoReference
	GooglePlaceDetail          *models.GooglePlaceDetail
}

// saveGooglePlacePhotoReferenceTx google_place に google_place_photo_reference, google_place_photo_attributions を紐付ける
func (p PlaceRepository) saveGooglePlacePhotoReferenceTx(ctx context.Context, tx *sql.Tx, input saveGooglePlacePhotoReferenceTxInput) (*entities.GooglePlacePhotoReferenceSlice, error) {
	// NearbySearchで取得したものとPlaceDetailで取得したものをマージする
	var googlePhotoReferences []models.GooglePlacePhotoReference
	googlePhotoReferences = input.GooglePlacePhotoReferences
	if input.GooglePlaceDetail != nil {
		// TODO: 重複を削除する
		googlePhotoReferences = append(googlePhotoReferences, input.GooglePlaceDetail.PhotoReferences...)
	}

	// GooglePlacePhotoReferenceを保存
	googlePlacePhotoReferenceEntities := factory.NewGooglePlacePhotoReferenceSliceFromGooglePlacePhotoReferences(googlePhotoReferences)
	if err := input.GooglePlaceEntity.AddGooglePlacePhotoReferences(ctx, tx, true, googlePlacePhotoReferenceEntities...); err != nil {
		return nil, fmt.Errorf("failed to insert google place photo reference: %w", err)
	}

	for _, googlePlacePhotoReference := range googlePhotoReferences {
		googlePlacePhotoReferenceEntity, ok := array.Find(googlePlacePhotoReferenceEntities, func(googlePlacePhotoReferenceEntity *entities.GooglePlacePhotoReference) bool {
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
	GooglePlacePhotoReferenceSlice entities.GooglePlacePhotoReferenceSlice
	GooglePlacePhoto               models.GooglePlacePhoto
	SavedGooglePlacePhotoSlice     entities.GooglePlacePhotoSlice
}

// addGooglePlacePhotosTx google_place_photo_reference に google_place_photo を紐付けする
// TODO:　複数の写真を一気に保存できるようにする
func (p PlaceRepository) addGooglePlacePhotosTx(ctx context.Context, tx *sql.Tx, input addGooglePlacePhotosTxInput) error {
	googlePlacePhotoReferenceEntity, ok := array.Find(input.GooglePlacePhotoReferenceSlice, func(googlePlacePhotoReferenceEntity *entities.GooglePlacePhotoReference) bool {
		if googlePlacePhotoReferenceEntity == nil {
			return false
		}
		return googlePlacePhotoReferenceEntity.PhotoReference == input.GooglePlacePhoto.PhotoReference
	})

	if !ok || googlePlacePhotoReferenceEntity == nil {
		return fmt.Errorf("failed to find google place photo reference entity")
	}

	googlePlacePhotoEntities := factory.NewGooglePlacePhotoSliceFromDomainModel(input.GooglePlacePhoto, input.GooglePlaceId)

	googlePlacePhotoEntities = array.Filter(googlePlacePhotoEntities, func(googlePlacePhotoEntity *entities.GooglePlacePhoto) bool {
		if googlePlacePhotoEntity == nil {
			return false
		}

		// すでに保存されている場合はスキップ
		if _, ok := array.Find(input.SavedGooglePlacePhotoSlice, func(savedGooglePlacePhotoEntity *entities.GooglePlacePhoto) bool {
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
