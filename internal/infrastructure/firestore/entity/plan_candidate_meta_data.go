package entity

import "poroto.app/poroto/planner/internal/domain/models"

type PlanCandidateMetaDataV1Entity struct {
	CreatedBasedOnCurrentLocation bool               `firestore:"created_based_on_current_location"`
	CategoriesPreferred           *[]string          `firestore:"categories_preferred,omitempty"`
	CategoriesRejected            *[]string          `firestore:"categories_rejected,omitempty"`
	LocationStart                 *GeoLocationEntity `firestore:"location_start,omitempty"`
	FreeTime                      *int               `firestore:"free_time,omitempty"`
}

func ToPlanCandidateMetaDataV1Entity(data models.PlanCandidateMetaData) PlanCandidateMetaDataV1Entity {
	var categoriesPreferred, categoriesRejected *[]string
	if data.CategoriesPreferred != nil {
		var categoryNames []string
		for _, category := range *data.CategoriesPreferred {
			categoryNames = append(categoryNames, category.Name)
		}
		categoriesPreferred = &categoryNames
	}

	if data.CategoriesRejected != nil {
		var categoryNames []string
		for _, category := range *data.CategoriesRejected {
			categoryNames = append(categoryNames, category.Name)
		}
		categoriesRejected = &categoryNames
	}

	var locationStart *GeoLocationEntity
	if data.LocationStart != nil {
		location := ToGeoLocationEntity(*data.LocationStart)
		locationStart = &location
	}

	return PlanCandidateMetaDataV1Entity{
		CreatedBasedOnCurrentLocation: data.CreatedBasedOnCurrentLocation,
		CategoriesPreferred:           categoriesPreferred,
		CategoriesRejected:            categoriesRejected,
		LocationStart:                 locationStart,
		FreeTime:                      data.FreeTime,
	}
}

func FromPlanCandidateMetaDataV1Entity(entity PlanCandidateMetaDataV1Entity) models.PlanCandidateMetaData {
	var categoriesPreferred, categoriesRejected *[]models.LocationCategory
	if entity.CategoriesPreferred != nil {
		var categories []models.LocationCategory
		for _, categoryName := range *entity.CategoriesPreferred {
			category := models.GetCategoryOfName(categoryName)
			if category != nil {
				categories = append(categories, *category)
			}
		}
		categoriesPreferred = &categories
	}

	if entity.CategoriesRejected != nil {
		var categories []models.LocationCategory
		for _, categoryName := range *entity.CategoriesRejected {
			category := models.GetCategoryOfName(categoryName)
			if category != nil {
				categories = append(categories, *category)
			}
		}
		categoriesRejected = &categories
	}

	var locationStart *models.GeoLocation
	if entity.LocationStart != nil {
		location := FromGeoLocationEntity(*entity.LocationStart)
		locationStart = &location
	}

	return models.PlanCandidateMetaData{
		CreatedBasedOnCurrentLocation: entity.CreatedBasedOnCurrentLocation,
		CategoriesPreferred:           categoriesPreferred,
		CategoriesRejected:            categoriesRejected,
		LocationStart:                 locationStart,
		FreeTime:                      entity.FreeTime,
	}
}
