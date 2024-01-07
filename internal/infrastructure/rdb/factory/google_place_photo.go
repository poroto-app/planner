package factory

import (
	"github.com/google/uuid"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
	"sort"
)

func NewGooglePlacePhotoFromEntity(
	googlePlacePhotoReferenceEntity generated.GooglePlacePhotoReference,
	googlePlacePhotoSlice generated.GooglePlacePhotoSlice,
	googlePlacePhotoAttributionSlice generated.GooglePlacePhotoAttributionSlice,
) *models.GooglePlacePhoto {
	googlePlacePhotoEntitiesFiltered := array.MapAndFilter(googlePlacePhotoSlice, func(googlePlacePhotoEntity *generated.GooglePlacePhoto) (generated.GooglePlacePhoto, bool) {
		if googlePlacePhotoEntity == nil {
			return generated.GooglePlacePhoto{}, false
		}

		// PhotoReferenceが一致するものだけを抽出
		if googlePlacePhotoEntity.PhotoReference != googlePlacePhotoReferenceEntity.PhotoReference {
			return generated.GooglePlacePhoto{}, false
		}

		return *googlePlacePhotoEntity, true
	})

	if len(googlePlacePhotoEntitiesFiltered) == 0 {
		return nil
	}

	// googlePlacePhotoEntitiesを画像サイズの昇順にソート
	// 一番小さい画像をSmallに、一番大きい画像をLargeに設定する
	sort.Slice(googlePlacePhotoEntitiesFiltered, func(i, j int) bool {
		return googlePlacePhotoEntitiesFiltered[i].Width < googlePlacePhotoEntitiesFiltered[j].Width
	})

	googlePlacePhotoAttributions := array.MapAndFilter(googlePlacePhotoAttributionSlice, func(googlePlacePhotoAttributionEntity *generated.GooglePlacePhotoAttribution) (string, bool) {
		if googlePlacePhotoAttributionEntity == nil {
			return "", false
		}

		// PhotoReferenceが一致するものだけを抽出
		if googlePlacePhotoAttributionEntity.PhotoReference != googlePlacePhotoReferenceEntity.PhotoReference {
			return "", false
		}

		return googlePlacePhotoAttributionEntity.HTMLAttribution, true
	})

	return &models.GooglePlacePhoto{
		PhotoReference:   googlePlacePhotoReferenceEntity.PhotoReference,
		Width:            googlePlacePhotoReferenceEntity.Width,
		Height:           googlePlacePhotoReferenceEntity.Height,
		HTMLAttributions: googlePlacePhotoAttributions,
		Small:            utils.StrOmitEmpty(googlePlacePhotoEntitiesFiltered[0].URL),
		Large:            utils.StrOmitEmpty(googlePlacePhotoEntitiesFiltered[len(googlePlacePhotoEntitiesFiltered)-1].URL),
	}
}

func NewGooglePlacePhotoSliceFromDomainModel(googlePlacePhoto models.GooglePlacePhoto, googlePlaceId string) generated.GooglePlacePhotoSlice {
	var googlePlacePhotoEntities generated.GooglePlacePhotoSlice

	if googlePlacePhoto.Small != nil {
		googlePlacePhotoEntities = append(googlePlacePhotoEntities, &generated.GooglePlacePhoto{
			ID:             uuid.New().String(),
			PhotoReference: googlePlacePhoto.PhotoReference,
			GooglePlaceID:  googlePlaceId,
			Width:          googlePlacePhoto.Width,
			Height:         googlePlacePhoto.Height,
			URL:            *googlePlacePhoto.Small,
		})
	}

	if googlePlacePhoto.Large != nil {
		googlePlacePhotoEntities = append(googlePlacePhotoEntities, &generated.GooglePlacePhoto{
			ID:             uuid.New().String(),
			PhotoReference: googlePlacePhoto.PhotoReference,
			GooglePlaceID:  googlePlaceId,
			Width:          googlePlacePhoto.Width,
			Height:         googlePlacePhoto.Height,
			URL:            *googlePlacePhoto.Large,
		})
	}

	return googlePlacePhotoEntities
}
