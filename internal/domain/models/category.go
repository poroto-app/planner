package models

type LocationCategory struct {
	Name          string
	SubCategories []string
}

var (
	// SEE: https://developers.google.com/maps/documentation/places/web-service/supported_types?hl=ja#table1
	CategoryAmusements = LocationCategory{
		Name: "amusements",
		SubCategories: []string{
			"amusement_park",
			"bowling_alley",
			"movie_theater",
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

	CategoryCafe = LocationCategory{
		Name: "cafe",
		SubCategories: []string{
			"cafe",
		},
	}

	CategoryCamp = LocationCategory{
		Name: "camp",
		SubCategories: []string{
			"campground",
			"rv_park",
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

	CategoryPark = LocationCategory{
		Name: "park",
		SubCategories: []string{
			"park",
		},
	}

	CategoryRestaurant = LocationCategory{
		Name: "restaurant",
		SubCategories: []string{
			"bakery",
			"bar",
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
			"movie_rental",
			"shoe_store",
			"store",
		},
	}
)

// 大カテゴリの集合
var AllCategory = []LocationCategory{
	CategoryAmusements,
	CategoryBook,
	CategoryCafe,
	CategoryCamp,
	CategoryCulture,
	CategoryNatural,
	CategoryPark,
	CategoryRestaurant,
	CategoryShopping,
}
