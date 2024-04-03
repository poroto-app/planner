package rdb

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/friendsofgo/errors"
	"github.com/google/uuid"
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

func (p PlaceRepository) SavePlacesFromGooglePlaces(ctx context.Context, googlePlace ...models.GooglePlace) (*[]models.Place, error) {
	googlePlaceIds := array.Map(googlePlace, func(googlePlace models.GooglePlace) string {
		return googlePlace.PlaceId
	})

	if err := runTransaction(ctx, p, func(ctx context.Context, tx *sql.Tx) error {
		googlePlaceSliceSaved, err := generated.GooglePlaces(generated.GooglePlaceWhere.GooglePlaceID.IN(googlePlaceIds)).All(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to find google places: %w", err)
		}

		googlePlacesNotSaved := array.Filter(googlePlace, func(googlePlace models.GooglePlace) bool {
			_, found := array.Find(googlePlaceSliceSaved, func(savedGooglePlace *generated.GooglePlace) bool {
				if savedGooglePlace == nil {
					return false
				}
				return savedGooglePlace.GooglePlaceID == googlePlace.PlaceId
			})
			return !found
		})

		if len(googlePlacesNotSaved) == 0 {
			return nil
		}

		// Google Place から対応する Place を生成
		var googlePlaceEntities generated.GooglePlaceSlice
		var placeEntities generated.PlaceSlice
		for _, googlePlace := range googlePlacesNotSaved {
			placeEntity := factory.NewPlaceEntityFromGooglePlaceEntity(googlePlace)
			placeEntities = append(placeEntities, &placeEntity)

			googlePlaceEntity := factory.NewGooglePlaceEntityFromGooglePlace(googlePlace, placeEntity.ID)
			googlePlaceEntities = append(googlePlaceEntities, &googlePlaceEntity)
		}

		// Placeを保存
		if _, err := placeEntities.InsertAll(ctx, tx, boil.Infer()); err != nil {
			return fmt.Errorf("failed to insert place: %w", err)
		}

		// GooglePlaceを保存
		if _, err := googlePlaceEntities.InsertAll(ctx, tx, boil.Infer()); err != nil {
			return fmt.Errorf("failed to insert google place: %w", err)
		}

		// GooglePlacePhotoReference を保存
		var googlePlacePhotoReferenceSliceNearbySearch generated.GooglePlacePhotoReferenceSlice = array.FlatMap(googlePlacesNotSaved, func(googlePlace models.GooglePlace) []*generated.GooglePlacePhotoReference {
			return factory.NewGooglePlacePhotoReferenceSliceFromGooglePlacePhotoReferences(googlePlace.PhotoReferences, googlePlace.PlaceId)
		})
		var googlePlacePhotoReferenceSlicePlaceDetail generated.GooglePlacePhotoReferenceSlice = array.FlatMap(googlePlacesNotSaved, func(googlePlace models.GooglePlace) []*generated.GooglePlacePhotoReference {
			if googlePlace.PlaceDetail == nil {
				return nil
			}
			return factory.NewGooglePlacePhotoReferenceSliceFromGooglePlacePhotoReferences(googlePlace.PlaceDetail.PhotoReferences, googlePlace.PlaceId)
		})
		googlePlacePhotoReferenceSlice := append(googlePlacePhotoReferenceSliceNearbySearch, googlePlacePhotoReferenceSlicePlaceDetail...)
		googlePlacePhotoReferenceSlice = array.DistinctBy(googlePlacePhotoReferenceSlice, func(googlePlacePhotoReference *generated.GooglePlacePhotoReference) string {
			if googlePlacePhotoReference == nil {
				return ""
			}
			return googlePlacePhotoReference.PhotoReference
		})
		if _, err := googlePlacePhotoReferenceSlice.InsertAll(ctx, tx, boil.Infer()); err != nil {
			return fmt.Errorf("failed to insert google place photo reference: %w", err)
		}

		// GooglePlacePhotoAttributionを保存
		var googlePlacePhotoAttributionSliceNearbySearch generated.GooglePlacePhotoAttributionSlice = array.FlatMap(googlePlacesNotSaved, func(googlePlace models.GooglePlace) []*generated.GooglePlacePhotoAttribution {
			return array.FlatMap(googlePlace.PhotoReferences, func(googlePlacePhotoReference models.GooglePlacePhotoReference) []*generated.GooglePlacePhotoAttribution {
				return factory.NewGooglePlacePhotoAttributionSliceFromPhotoReference(googlePlacePhotoReference, googlePlace.PlaceId)
			})
		})
		var googlePlacePhotoAttributionSlicePlaceDetail generated.GooglePlacePhotoAttributionSlice = array.FlatMap(googlePlacesNotSaved, func(googlePlace models.GooglePlace) []*generated.GooglePlacePhotoAttribution {
			if googlePlace.PlaceDetail == nil {
				return nil
			}
			return array.FlatMap(googlePlace.PlaceDetail.PhotoReferences, func(googlePlacePhotoReference models.GooglePlacePhotoReference) []*generated.GooglePlacePhotoAttribution {
				return factory.NewGooglePlacePhotoAttributionSliceFromPhotoReference(googlePlacePhotoReference, googlePlace.PlaceId)
			})
		})
		googlePlacePhotoAttributionSlice := append(googlePlacePhotoAttributionSliceNearbySearch, googlePlacePhotoAttributionSlicePlaceDetail...)
		googlePlacePhotoAttributionSlice = array.DistinctBy(googlePlacePhotoAttributionSlice, func(googlePlacePhotoAttribution *generated.GooglePlacePhotoAttribution) string {
			if googlePlacePhotoAttribution == nil {
				return ""
			}
			return googlePlacePhotoAttribution.HTMLAttribution
		})
		if _, err := googlePlacePhotoAttributionSlice.InsertAll(ctx, tx, boil.Infer()); err != nil {
			return fmt.Errorf("failed to insert google place photo attribution: %w", err)
		}

		// GooglePlacePhotoを保存
		var googlePhotoSlice generated.GooglePlacePhotoSlice
		for _, googlePlace := range googlePlacesNotSaved {
			if googlePlace.Photos == nil {
				continue
			}

			for _, googlePhoto := range *googlePlace.Photos {
				googlePhotoSlice = append(googlePhotoSlice, factory.NewGooglePlacePhotoSliceFromDomainModel(googlePhoto, googlePlace.PlaceId)...)
			}
		}
		if _, err := googlePhotoSlice.InsertAll(ctx, tx, boil.Infer()); err != nil {
			return fmt.Errorf("failed to insert google place photo: %w", err)
		}

		// GooglePlaceTypeを保存
		var googlePlaceTypeSlice generated.GooglePlaceTypeSlice = array.FlatMap(googlePlacesNotSaved, func(googlePlace models.GooglePlace) []*generated.GooglePlaceType {
			return factory.NewGooglePlaceTypeSliceFromGooglePlace(googlePlace)
		})
		if _, err := googlePlaceTypeSlice.InsertAll(ctx, tx, boil.Infer()); err != nil {
			return fmt.Errorf("failed to insert google place type: %w", err)
		}

		// GooglePlaceReviewを保存
		var googlePlaceReviewSlice generated.GooglePlaceReviewSlice = array.FlatMap(googlePlacesNotSaved, func(googlePlace models.GooglePlace) []*generated.GooglePlaceReview {
			return factory.NewGooglePlaceReviewSliceFromGooglePlace(googlePlace)
		})
		if _, err := googlePlaceReviewSlice.InsertAll(ctx, tx, boil.Infer()); err != nil {
			return fmt.Errorf("failed to insert google place review: %w", err)
		}

		// GooglePlaceOpeningPeriodを保存
		var googlePlaceOpeningPeriodSlice generated.GooglePlaceOpeningPeriodSlice = array.FlatMap(googlePlacesNotSaved, func(googlePlace models.GooglePlace) []*generated.GooglePlaceOpeningPeriod {
			return factory.NewGooglePlaceOpeningPeriodSliceFromGooglePlace(googlePlace)
		})
		if _, err := googlePlaceOpeningPeriodSlice.InsertAll(ctx, tx, boil.Infer()); err != nil {
			return fmt.Errorf("failed to insert google place opening period: %w", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to run transaction: %w", err)
	}

	// 保存したPlaceを取得
	places, err := p.findAllByGooglePlaceId(ctx, p.GetDB(), googlePlaceIds)
	if err != nil {
		return nil, fmt.Errorf("failed to find places: %w", err)
	}

	return &places, nil
}

func (p PlaceRepository) FindByLocation(ctx context.Context, location models.GeoLocation, radius float64) ([]models.Place, error) {
	minLocation, maxLocation := location.CalculateMBR(radius)

	googlePlaceEntities, err := generated.GooglePlaces(
		generated.GooglePlaceWhere.Latitude.GT(minLocation.Latitude),
		generated.GooglePlaceWhere.Latitude.LT(maxLocation.Latitude),
		generated.GooglePlaceWhere.Longitude.GT(minLocation.Longitude),
		generated.GooglePlaceWhere.Longitude.LT(maxLocation.Longitude),
		qm.Load(generated.GooglePlaceRels.Place),
		qm.Load(generated.GooglePlaceRels.Place+"."+generated.PlaceRels.PlacePhotos),
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
			googlePlaceEntity.R.Place.R.PlacePhotos,
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

func (p PlaceRepository) FindByGooglePlaceType(ctx context.Context, googlePlaceType string, baseLocation models.GeoLocation, radius float64) (*[]models.Place, error) {
	minLocation, maxLocation := baseLocation.CalculateMBR(radius)
	googlePlaceEntities, err := generated.GooglePlaces(
		qm.InnerJoin(fmt.Sprintf(
			"%s on %s.%s = %s.%s",
			generated.TableNames.GooglePlaceTypes,
			generated.TableNames.GooglePlaceTypes,
			generated.GooglePlaceTypeColumns.GooglePlaceID,
			generated.TableNames.GooglePlaces,
			generated.GooglePlaceColumns.GooglePlaceID,
		)),
		qm.Where(fmt.Sprintf("%s.%s = ?", generated.TableNames.GooglePlaceTypes, generated.GooglePlaceTypeColumns.Type), googlePlaceType),
		generated.GooglePlaceWhere.Latitude.GT(minLocation.Latitude),
		generated.GooglePlaceWhere.Latitude.LT(maxLocation.Latitude),
		generated.GooglePlaceWhere.Longitude.GT(minLocation.Longitude),
		generated.GooglePlaceWhere.Longitude.LT(maxLocation.Longitude),
		qm.Load(generated.GooglePlaceRels.Place),
		qm.Load(generated.GooglePlaceRels.Place+"."+generated.PlaceRels.PlacePhotos),
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

	if len(googlePlaceEntities) == 0 {
		return &[]models.Place{}, nil
	}

	// 一対多の関係になるため、重複を排除する
	googlePlaceEntities = array.DistinctBy(googlePlaceEntities, func(googlePlaceEntity *generated.GooglePlace) string {
		if googlePlaceEntity == nil {
			return ""
		}
		return googlePlaceEntity.PlaceID
	})

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
			googlePlaceEntity.R.Place.R.PlacePhotos,
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

	return &places, nil
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
			planCandidateSetSearchedPlace.R.Place.R.PlacePhotos,
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

func (p PlaceRepository) FindLikePlacesByUserId(ctx context.Context, userId string) (*[]models.Place, error) {
	userLikePlaces, err := generated.UserLikePlaces(concatQueryMod(
		[]qm.QueryMod{generated.UserLikePlaceWhere.UserID.EQ(userId)},
		placeQueryModes(generated.PlanCandidateSetSearchedPlaceRels.Place),
	)...).All(ctx, p.db)
	if err != nil {
		return nil, fmt.Errorf("failed to find user like places: %w", err)
	}

	placeLikeCounts, err := countPlaceLikeCounts(ctx, p.db, array.MapAndFilter(userLikePlaces, func(userLikePlace *generated.UserLikePlace) (string, bool) {
		if userLikePlace == nil {
			return "", false
		}
		return userLikePlace.PlaceID, true
	})...)
	if err != nil {
		// いいね数の取得に失敗してもエラーにしない
		p.logger.Warn("failed to count place like counts", zap.Error(err))
	}

	places := make([]models.Place, 0, len(userLikePlaces))
	for _, userLikePlace := range userLikePlaces {
		if userLikePlace == nil {
			continue
		}

		if userLikePlace.R == nil {
			panic("userLikePlace.R is nil")
		}

		if userLikePlace.R.Place == nil {
			p.logger.Warn("userLikePlace.R.Place is nil", zap.String("user_like_place_id", userLikePlace.ID))
			continue
		}

		if userLikePlace.R.Place.R == nil {
			panic("userLikePlace.R.Place.R is nil")
		}

		if len(userLikePlace.R.Place.R.GooglePlaces) == 0 {
			p.logger.Warn("userLikePlace.R.Place.R.GooglePlaces is empty", zap.String("user_like_place_id", userLikePlace.ID))
			continue
		}

		place, err := factory.NewPlaceFromEntity(
			*userLikePlace.R.Place,
			userLikePlace.R.Place.R.PlacePhotos,
			*userLikePlace.R.Place.R.GooglePlaces[0],
			userLikePlace.R.Place.R.GooglePlaces[0].R.GooglePlaceTypes,
			userLikePlace.R.Place.R.GooglePlaces[0].R.GooglePlacePhotoReferences,
			userLikePlace.R.Place.R.GooglePlaces[0].R.GooglePlacePhotoAttributions,
			userLikePlace.R.Place.R.GooglePlaces[0].R.GooglePlacePhotos,
			userLikePlace.R.Place.R.GooglePlaces[0].R.GooglePlaceReviews,
			userLikePlace.R.Place.R.GooglePlaces[0].R.GooglePlaceOpeningPeriods,
			entities.CountLikeOfPlace(placeLikeCounts, userLikePlace.PlaceID),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to convert google place googlePlaceEntity to place: %w", err)
		}

		places = append(places, *place)
	}

	return &places, nil
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

		googlePlacePhotoSlicesNotSaved := array.MapAndFilter(photos, func(googlePlacePhoto models.GooglePlacePhoto) (*generated.GooglePlacePhotoSlice, bool) {
			// すでに保存されている場合はスキップ
			if _, found := array.Find(googlePlaceEntity.R.GooglePlacePhotos, func(savedGooglePlacePhotoEntity *generated.GooglePlacePhoto) bool {
				if savedGooglePlacePhotoEntity == nil {
					return false
				}
				alreadySaved := savedGooglePlacePhotoEntity.Width == googlePlacePhoto.Width && savedGooglePlacePhotoEntity.Height == googlePlacePhoto.Height
				return alreadySaved
			}); found {
				p.logger.Debug(
					"skip google place photo because already exists",
					zap.String("photo_reference", googlePlacePhoto.PhotoReference),
					zap.Int("width", googlePlacePhoto.Width),
					zap.Int("height", googlePlacePhoto.Height),
				)
				return nil, false
			}

			googlePlacePhotoEntity := factory.NewGooglePlacePhotoSliceFromDomainModel(googlePlacePhoto, googlePlaceId)
			return &googlePlacePhotoEntity, true
		})

		var googlePhotoSlice generated.GooglePlacePhotoSlice = array.FlatMap(googlePlacePhotoSlicesNotSaved, func(googlePlacePhotoSlice *generated.GooglePlacePhotoSlice) []*generated.GooglePlacePhoto {
			if googlePlacePhotoSlice == nil {
				return nil
			}
			return *googlePlacePhotoSlice
		})
		if _, err := googlePhotoSlice.InsertAll(ctx, tx, boil.Infer()); err != nil {
			return fmt.Errorf("failed to insert google place photo: %w", err)
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
			googlePlaceOpeningPeriodEntities := factory.NewGooglePlaceReviewSliceFromGooglePlaceDetail(googlePlaceDetail, googlePlaceId)
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

			_, found := array.Find(googlePlaceEntity.R.GooglePlacePhotoReferences, func(savedGooglePlacePhotoReferenceEntity *generated.GooglePlacePhotoReference) bool {
				if savedGooglePlacePhotoReferenceEntity == nil {
					return false
				}
				return savedGooglePlacePhotoReferenceEntity.PhotoReference == googlePlacePhotoReferenceEntity.PhotoReference
			})

			// すでに保存済みのものはスキップ
			return !found
		})
		if err := googlePlaceEntity.AddGooglePlacePhotoReferences(ctx, tx, true, googlePlacePhotoReferenceSlice...); err != nil {
			return fmt.Errorf("failed to insert google place photo reference: %w", err)
		}

		// HTMLAttributionを保存
		for _, googlePlacePhotoReference := range googlePlaceDetail.PhotoReferences {
			googlePlacePhotoReferenceEntity, found := array.Find(googlePlaceEntity.R.GooglePlacePhotoReferences, func(savedGooglePlacePhotoReferenceEntity *generated.GooglePlacePhotoReference) bool {
				if savedGooglePlacePhotoReferenceEntity == nil {
					return false
				}
				return savedGooglePlacePhotoReferenceEntity.PhotoReference == googlePlacePhotoReference.PhotoReference
			})
			if !found || googlePlacePhotoReferenceEntity == nil {
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

func (p PlaceRepository) SavePlacePhotos(ctx context.Context, photos []models.PlacePhoto) error {
	if err := runTransaction(ctx, p, func(ctx context.Context, tx *sql.Tx) error {
		placePhotoUrls := array.Map(photos, func(photo models.PlacePhoto) string {
			return photo.PhotoUrl
		})

		placePhotoSliceAlreadySaved, err := generated.PlacePhotos(
			generated.PlacePhotoWhere.PhotoURL.IN(placePhotoUrls),
		).All(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to find place photos: %w", err)
		}

		placePhotoUrlsAlreadySaved := array.Map(placePhotoSliceAlreadySaved, func(photoEntity *generated.PlacePhoto) string {
			return photoEntity.PhotoURL
		})

		if len(placePhotoUrlsAlreadySaved) > 0 {
			p.logger.Debug("skipped to save because place photo urls already saved", zap.Strings("place_photo_urls", placePhotoUrlsAlreadySaved))
		}

		placePhotosToSave := array.Filter(photos, func(photo models.PlacePhoto) bool {
			return !array.IsContain(placePhotoUrlsAlreadySaved, photo.PhotoUrl)
		})

		placePhotoSliceToSave := factory.NewPlacePhotoSliceFromDomainModel(placePhotosToSave)
		if placePhotoSliceToSave == nil {
			p.logger.Debug("no place photo to save")
			return nil
		}

		if _, err := placePhotoSliceToSave.InsertAll(ctx, tx, boil.Infer()); err != nil {
			return fmt.Errorf("failed to insert place photo slice: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to save place photos: %w", err)
	}
	return nil
}

func (p PlaceRepository) UpdateLikeByUserId(ctx context.Context, userId string, placeId string, like bool) error {
	if err := runTransaction(ctx, p, func(ctx context.Context, tx *sql.Tx) error {
		if !like {
			// いいねを取り消す
			if _, err := generated.UserLikePlaces(
				generated.UserLikePlaceWhere.UserID.EQ(userId),
				generated.UserLikePlaceWhere.PlaceID.EQ(placeId),
			).DeleteAll(ctx, tx); err != nil {
				return fmt.Errorf("failed to delete place like: %w", err)
			}
			return nil
		}

		userLikePlaceEntity := generated.UserLikePlace{
			ID:      uuid.New().String(),
			UserID:  userId,
			PlaceID: placeId,
		}
		if err := userLikePlaceEntity.Insert(ctx, tx, boil.Infer()); err != nil {
			return fmt.Errorf("failed to insert place like: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to run transaction: %w", err)
	}
	return nil
}

func (p PlaceRepository) UpdateLikeByPlanCandidateSetToUser(ctx context.Context, userId string, planCandidateSetIds []string) error {
	if err := runTransaction(ctx, p, func(ctx context.Context, tx *sql.Tx) error {
		planCandidateSetEntities, err := generated.PlanCandidateSets(
			generated.PlanCandidateSetWhere.ID.IN(planCandidateSetIds),
			qm.Load(generated.PlanCandidateSetRels.PlanCandidateSetLikePlaces),
		).All(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to find plan candidate sets: %w", err)
		}

		var likePlacesByPlanCandidates generated.PlanCandidateSetLikePlaceSlice
		likePlacesByPlanCandidates = array.FlatMap(planCandidateSetEntities, func(planCandidateSetEntity *generated.PlanCandidateSet) []*generated.PlanCandidateSetLikePlace {
			likePlaces := planCandidateSetEntity.R.GetPlanCandidateSetLikePlaces()
			if likePlaces == nil {
				return []*generated.PlanCandidateSetLikePlace{}
			}
			return likePlaces
		})

		likePlacesByPlanCandidates = array.DistinctBy(likePlacesByPlanCandidates, func(likePlace *generated.PlanCandidateSetLikePlace) string {
			if likePlace == nil {
				return ""
			}
			return likePlace.PlaceID
		})

		// ユーザーとしていいねを登録
		var userLikePlaceEntities generated.UserLikePlaceSlice = array.Map(likePlacesByPlanCandidates, func(likePlace *generated.PlanCandidateSetLikePlace) *generated.UserLikePlace {
			return &generated.UserLikePlace{
				ID:      uuid.New().String(),
				UserID:  userId,
				PlaceID: likePlace.PlaceID,
			}
		})
		if _, err := userLikePlaceEntities.InsertAll(ctx, tx, boil.Infer()); err != nil {
			return fmt.Errorf("failed to insert user like places: %w", err)
		}

		// プラン候補セットとしていいねした記録を削除
		if _, err := likePlacesByPlanCandidates.DeleteAll(ctx, tx); err != nil {
			return fmt.Errorf("failed to delete plan candidate set like places: %w", err)
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
		qm.Load(generated.GooglePlaceRels.Place+"."+generated.PlaceRels.PlacePhotos),
		qm.Load(generated.GooglePlaceRels.GooglePlaceTypes),
		qm.Load(generated.GooglePlaceRels.GooglePlacePhotoReferences),
		qm.Load(generated.GooglePlaceRels.GooglePlacePhotos),
		qm.Load(generated.GooglePlaceRels.GooglePlacePhotoAttributions),
		qm.Load(generated.GooglePlaceRels.GooglePlaceReviews),
		qm.Load(generated.GooglePlaceRels.GooglePlaceOpeningPeriods),
	).One(ctx, exec)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
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
		googlePlaceEntity.R.Place.R.PlacePhotos,
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

func (p PlaceRepository) findAllByGooglePlaceId(ctx context.Context, exec boil.ContextExecutor, googlePlaceIds []string) ([]models.Place, error) {
	googlePlaceEntities, err := generated.GooglePlaces(
		generated.GooglePlaceWhere.GooglePlaceID.IN(googlePlaceIds),
		qm.Load(generated.GooglePlaceRels.Place),
		qm.Load(generated.GooglePlaceRels.Place+"."+generated.PlaceRels.PlacePhotos),
		qm.Load(generated.GooglePlaceRels.GooglePlaceTypes),
		qm.Load(generated.GooglePlaceRels.GooglePlacePhotoReferences),
		qm.Load(generated.GooglePlaceRels.GooglePlacePhotos),
		qm.Load(generated.GooglePlaceRels.GooglePlacePhotoAttributions),
		qm.Load(generated.GooglePlaceRels.GooglePlaceReviews),
		qm.Load(generated.GooglePlaceRels.GooglePlaceOpeningPeriods),
	).All(ctx, exec)
	if err != nil {
		return nil, fmt.Errorf("failed to find google places: %w", err)
	}

	planCandidateSetPlaceLikeCounts, err := countPlaceLikeCounts(ctx, exec, googlePlaceIds...)
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
			googlePlaceEntity.R.Place.R.PlacePhotos,
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

		places = append(places, *place)
	}

	return places, nil
}

// countPlaceLikeCounts は場所ごとのいいね数をカウントする
// いいねはPlanCandidateSetとUserによって行われるが、その両方を考慮し、総数をカウントする
func countPlaceLikeCounts(ctx context.Context, exec boil.ContextExecutor, placeIds ...string) (*[]entities.PlanCandidateSetPlaceLikeCount, error) {
	var placeIdPlaceHolder string
	if len(placeIds) == 0 {
		return nil, nil
	} else if len(placeIds) == 1 {
		placeIdPlaceHolder = "?"
	} else {
		placeIdPlaceHolder = strings.Repeat("?,", len(placeIds)-1) + "?"
	}

	query := fmt.Sprintf(
		`SELECT place_likes.%s, COUNT(*) AS %s
FROM (
	SELECT %s, %s AS %s
	FROM %s
	UNION ALL
	SELECT %s, %s AS %s
	FROM %s
) AS place_likes	
WHERE %s IN (%s)
GROUP BY place_likes.place_id`,
		entities.PlanCandidateSetPlaceLikeCountColumns.PlaceId,
		entities.PlanCandidateSetPlaceLikeCountColumns.LikeCount,

		generated.PlanCandidateSetLikePlaceColumns.ID,
		generated.PlanCandidateSetLikePlaceColumns.PlaceID,
		entities.PlanCandidateSetPlaceLikeCountColumns.PlaceId,

		generated.TableNames.PlanCandidateSetLikePlaces,

		generated.UserLikePlaceColumns.ID,
		generated.UserLikePlaceColumns.PlaceID,
		entities.PlanCandidateSetPlaceLikeCountColumns.PlaceId,

		generated.TableNames.UserLikePlaces,

		entities.PlanCandidateSetPlaceLikeCountColumns.PlaceId,

		placeIdPlaceHolder,
	)

	var planCandidateSetPlaceLikeCounts []entities.PlanCandidateSetPlaceLikeCount
	if err := queries.
		Raw(query, toInterfaceArray(placeIds)...).
		Bind(ctx, exec, &planCandidateSetPlaceLikeCounts); err != nil {
		return nil, fmt.Errorf("failed to count place like counts: %w", err)
	}
	return &planCandidateSetPlaceLikeCounts, nil
}
