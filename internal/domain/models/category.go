package models

type LocationCategory struct {
	Name          string
	SubCategories []string
}

var (
	// SEE: https://developers.google.com/maps/documentation/places/web-service/supported_types?hl=ja#table1
	CategoryActivity = LocationCategory{
		Name: "activity",
		SubCategories: []string{
			"bowling_alley",
		},
	}

	CategoryAmusements = LocationCategory{
		Name: "amusements",
		SubCategories: []string{
			"amusement_park",
			"campground",
			"movie_rental",
			"movie_theater",
			"park",
			"rv_park",
			"spa",
			"stadium",
		},
	}

	CategoryBook = LocationCategory{
		Name: "book",
		SubCategories: []string{
			"book_store",
			"library",
		},
	}

	CategoryCulture = LocationCategory{
		Name: "cultural Facility",
		SubCategories: []string{
			"art_gallery",
			"museum",
			"tourist_attraction",
		},
	}

	CategoryNatural = LocationCategory{
		Name: "natural Facility",
		SubCategories: []string{
			"aquarium",
			"zoo",
		},
	}

	CategoryRestaurant = LocationCategory{
		Name: "restaurant",
		SubCategories: []string{
			"bakery",
			"bar",
			"cafe",
			"food",
			"restaurant",
		},
	}

	CategoryShopping = LocationCategory{
		Name: "shopping",
		SubCategories: []string{
			"clothing_store",
			"department_store",
			"furniture_store",
			"hardware_store",
			"home_goods_store",
			"shoe_store",
			"store",
		},
	}
)
