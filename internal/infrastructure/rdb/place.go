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
			googlePlaceEntity, err := entities.GooglePlaces(
				qm.Where(entities.GooglePlaceColumns.GooglePlaceID, googlePlace.PlaceId),
				qm.Load(entities.GooglePlaceRels.Place),
				qm.Load(entities.GooglePlaceRels.GooglePlaceTypes),
				qm.Load(entities.GooglePlaceRels.GooglePlacePhotoReferences),
				qm.Load(entities.GooglePlaceRels.GooglePlacePhotos),
				qm.Load(entities.GooglePlaceRels.GooglePlacePhotoAttributions),
				qm.Load(entities.GooglePlaceRels.GooglePlaceReviews),
				qm.Load(entities.GooglePlaceRels.GooglePlaceOpeningPeriods),
			).One(ctx, tx)

			if err != nil {
				return fmt.Errorf("failed to find google place: %w", err)
			}

			placeSaved, err := factory.NewPlaceFromGooglePlaceEntity(*googlePlaceEntity)
			if err != nil {
				return fmt.Errorf("failed to convert google place entity to place: %w", err)
			}
			if placeSaved == nil {
				return fmt.Errorf("failed to convert google place entity to place")
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

		// GooglePlacePhotoReference, GooglePlacePhoto, GooglePlacePhotoAttributionを保存
		if _, err := p.saveGooglePlacePhotoReferenceTx(ctx, tx, saveGooglePlacePhotoReferenceTxInput{
			GooglePlacePhotoReferences: googlePlace.PhotoReferences,
			GooglePlaceDetail:          googlePlace.PlaceDetail,
			GooglePlaceEntity:          googlePlaceEntity,
			GooglePlacePhotos:          googlePlace.Photos,
		}); err != nil {
			return fmt.Errorf("failed to insert google place photo reference: %w", err)
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
	googlePlaceEntity, err := entities.GooglePlaces(
		entities.GooglePlaceWhere.GooglePlaceID.EQ(googlePlaceID),
		qm.Load(entities.GooglePlaceRels.Place),
		qm.Load(entities.GooglePlaceRels.GooglePlaceTypes),
		qm.Load(entities.GooglePlaceRels.GooglePlacePhotoReferences),
		qm.Load(entities.GooglePlaceRels.GooglePlacePhotos),
		qm.Load(entities.GooglePlaceRels.GooglePlacePhotoAttributions),
		qm.Load(entities.GooglePlaceRels.GooglePlaceReviews),
		qm.Load(entities.GooglePlaceRels.GooglePlaceOpeningPeriods),
	).One(ctx, p.db)
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

func (p PlaceRepository) FindByPlanCandidateId(ctx context.Context, planCandidateId string) ([]models.Place, error) {
	//TODO implement me
	panic("implement me")
}

func (p PlaceRepository) SaveGooglePlacePhotos(ctx context.Context, googlePlaceId string, photos []models.GooglePlacePhoto) error {
	//TODO implement me
	panic("implement me")
}

func (p PlaceRepository) SaveGooglePlaceDetail(ctx context.Context, googlePlaceId string, googlePlaceDetail models.GooglePlaceDetail) error {
	if err := runTransaction(ctx, p, func(ctx context.Context, tx *sql.Tx) error {
		googlePlaceEntity, err := entities.GooglePlaces(
			qm.Where(entities.GooglePlaceColumns.GooglePlaceID, googlePlaceId),
			qm.Load(entities.GooglePlaceRels.GooglePlaceOpeningPeriods),
			qm.Load(entities.GooglePlaceRels.GooglePlaceReviews),
			qm.Load(
				entities.GooglePlaceRels.GooglePlacePhotoReferences,
				qm.Load(entities.GooglePlacePhotoReferenceRels.PhotoReferenceGooglePlacePhotos),
				qm.Load(entities.GooglePlacePhotoReferenceRels.PhotoReferenceGooglePlacePhotoAttributions),
			),
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
		if err := googlePlaceEntity.AddGooglePlacePhotoReferences(ctx, tx, true, googlePlacePhotoReferenceSlice...); err != nil {
			return fmt.Errorf("failed to insert google place photo reference: %w", err)
		}

		for _, googlePlacePhotoReference := range googlePlaceDetail.PhotoReferences {
			var googlePlacePhotoReferenceEntity *entities.GooglePlacePhotoReference
			for _, entity := range googlePlacePhotoReferenceSlice {
				if entity == nil {
					continue
				}

				if entity.PhotoReference == googlePlacePhotoReference.PhotoReference {
					googlePlacePhotoReferenceEntity = entity
					break
				}
			}

			if googlePlacePhotoReferenceEntity == nil {
				return fmt.Errorf("failed to find google place photo reference entity")
			}

			// HTMLAttributionを保存
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

type saveGooglePlacePhotoReferenceTxInput struct {
	GooglePlaceEntity          entities.GooglePlace
	GooglePlacePhotoReferences []models.GooglePlacePhotoReference
	GooglePlaceDetail          *models.GooglePlaceDetail
	GooglePlacePhotos          *[]models.GooglePlacePhoto
}

// saveGooglePlacePhotoReferenceTx google_place_photo_reference に google_place_photo を紐付ける
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

	photoReferenceToGooglePlacePhotoReferenceEntity := make(map[string]*entities.GooglePlacePhotoReference)
	for _, googlePlacePhotoReferenceEntity := range googlePlacePhotoReferenceEntities {
		photoReferenceToGooglePlacePhotoReferenceEntity[googlePlacePhotoReferenceEntity.PhotoReference] = googlePlacePhotoReferenceEntity
	}

	photoReferenceToGooglePlacePhoto := make(map[string]models.GooglePlacePhoto)
	if input.GooglePlacePhotos != nil {
		for _, googlePlacePhoto := range *input.GooglePlacePhotos {
			photoReferenceToGooglePlacePhoto[googlePlacePhoto.PhotoReference] = googlePlacePhoto
		}
	}

	for _, googlePlacePhotoReference := range googlePhotoReferences {
		googlePlacePhotoReferenceEntity, ok := photoReferenceToGooglePlacePhotoReferenceEntity[googlePlacePhotoReference.PhotoReference]
		if !ok {
			return nil, fmt.Errorf("failed to find google place photo reference entity")
		}
		if googlePlacePhotoReferenceEntity == nil {
			continue
		}

		// HTMLAttributionを保存
		googlePlacePhotoAttributionEntities := factory.NewGooglePlacePhotoAttributionSliceFromPhotoReference(googlePlacePhotoReference, input.GooglePlaceEntity.GooglePlaceID)
		if err := googlePlacePhotoReferenceEntity.AddPhotoReferenceGooglePlacePhotoAttributions(ctx, tx, true, googlePlacePhotoAttributionEntities...); err != nil {
			return nil, fmt.Errorf("failed to insert google place photo attribution: %w", err)
		}

		// Photoを保存
		googlePlacePhoto, ok := photoReferenceToGooglePlacePhoto[googlePlacePhotoReference.PhotoReference]
		if !ok {
			continue
		}
		if err := p.addGooglePlacePhotosTx(ctx, tx, googlePlacePhoto, *googlePlacePhotoReferenceEntity, input.GooglePlaceEntity.GooglePlaceID); err != nil {
			return nil, fmt.Errorf("failed to insert google place photo: %w", err)
		}
	}

	return &googlePlacePhotoReferenceEntities, nil
}

// addGooglePlacePhotosTx google_place_photo_reference に google_place_photo を紐付けする
func (p PlaceRepository) addGooglePlacePhotosTx(ctx context.Context, tx *sql.Tx, googlePlacePhoto models.GooglePlacePhoto, googlePlacePhotoReferenceEntity entities.GooglePlacePhotoReference, googlePlaceId string) error {
	if googlePlacePhotoReferenceEntity.R == nil {
		return fmt.Errorf("google place photo reference entity is nil")
	}

	// すでに紐付けがある場合は何もしない
	if len(googlePlacePhotoReferenceEntity.R.PhotoReferenceGooglePlacePhotos) > 0 {
		return nil
	}

	googlePlacePhotoEntities := factory.NewGooglePlacePhotoSliceFromDomainModel(googlePlacePhoto, googlePlaceId)
	if err := googlePlacePhotoReferenceEntity.AddPhotoReferenceGooglePlacePhotos(ctx, tx, true, googlePlacePhotoEntities...); err != nil {
		return fmt.Errorf("failed to insert google place photo: %w", err)
	}
	return nil
}
