package models

import (
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/infrastructure/api/google/places"
)

type PlacesFilter struct {
	placesToFilter []places.Place
}

func NewPlacesFilter(placesToFilter []places.Place) PlacesFilter {
	return PlacesFilter{
		placesToFilter: placesToFilter,
	}
}

func (f PlacesFilter) Places() []places.Place {
	return f.placesToFilter
}

func (f PlacesFilter) Find(findFunc func(place places.Place) bool) *places.Place {
	for _, place := range f.placesToFilter {
		if findFunc(place) {
			copy := place
			return &copy
		}
	}
	return nil
}

func (f PlacesFilter) filterPlaces(filterFunc func(place places.Place) bool) []places.Place {
	var values []places.Place
	for _, place := range f.placesToFilter {
		if filterFunc(place) {
			values = append(values, place)
		}
	}
	return values
}

func (f PlacesFilter) FilterByCategory(categories []LocationCategory) PlacesFilter {
	var subCategories []string
	for _, category := range categories {
		subCategories = append(subCategories, category.SubCategories...)
	}

	f.placesToFilter = f.filterPlaces(func(place places.Place) bool {
		for _, category := range place.Types {
			if array.IsContain(subCategories, category) {
				return true
			}
		}
		return false
	})

	return f
}

func (f PlacesFilter) FilterByOpeningNow() PlacesFilter {
	f.placesToFilter = f.filterPlaces(func(place places.Place) bool {
		return place.OpenNow
	})
	return f
}

func (f PlacesFilter) FilterWithinDistanceRange(
	currentLocation GeoLocation,
	startInMeter float64,
	endInMeter float64,
) PlacesFilter {
	f.placesToFilter = f.filterPlaces(func(place places.Place) bool {
		distance := currentLocation.DistanceInMeter(place.Location.ToGeoLocation())
		return startInMeter <= distance && distance < endInMeter
	})
	return f
}
