package factory

import (
	"github.com/google/uuid"
	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/domain/utils"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
	"sort"
)

func NewGooglePlacePhotoFromEntity(
	googlePlacePhotoReferenceEntity generated.GooglePlacePhotoReference,
	googlePlacePhotoEntities generated.GooglePlacePhotoSlice,
	googlePlacePhotoAttributionEntities generated.GooglePlacePhotoAttributionSlice,
) *models.GooglePlacePhoto {
	var googlePlacePhotoEntitiesFiltered []generated.GooglePlacePhoto
	for _, googlePlacePhotoEntity := range googlePlacePhotoEntities {
		if googlePlacePhotoEntity == nil {
			continue
		}

		// PhotoReferenceが一致するものだけを抽出
		if googlePlacePhotoEntity.PhotoReference != googlePlacePhotoReferenceEntity.PhotoReference {
			continue
		}

		googlePlacePhotoEntitiesFiltered = append(googlePlacePhotoEntitiesFiltered, *googlePlacePhotoEntity)
	}

	if len(googlePlacePhotoEntitiesFiltered) == 0 {
		return nil
	}

	// googlePlacePhotoEntitiesを画像サイズの昇順にソート
	// 一番小さい画像をSmallに、一番大きい画像をLargeに設定する
	sort.Slice(googlePlacePhotoEntitiesFiltered, func(i, j int) bool {
		return googlePlacePhotoEntitiesFiltered[i].Width < googlePlacePhotoEntitiesFiltered[j].Width
	})

	var googlePlacePhotoAttributions []string
	for _, googlePlacePhotoAttribution := range googlePlacePhotoAttributionEntities {
		if googlePlacePhotoAttribution == nil {
			continue
		}

		// PhotoReferenceが一致するものだけを抽出
		if googlePlacePhotoAttribution.PhotoReference != googlePlacePhotoReferenceEntity.PhotoReference {
			continue
		}

		googlePlacePhotoAttributions = append(googlePlacePhotoAttributions, googlePlacePhotoAttribution.HTMLAttribution)
	}

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
