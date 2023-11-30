package models

type LocationCategoryWithPlaces struct {
	Category LocationCategory
	Places   []Place
}

func NewLocationCategoryWithPlaces(category LocationCategory, places []Place) LocationCategoryWithPlaces {
	var placesToAdd []Place
	for _, place := range places {
		// 画像がない場合は追加しない
		if place.Google.Photos == nil || len(*place.Google.Photos) == 0 {
			continue
		}

		placesToAdd = append(placesToAdd, place)
	}

	return LocationCategoryWithPlaces{
		Category: category,
		Places:   placesToAdd,
	}
}
