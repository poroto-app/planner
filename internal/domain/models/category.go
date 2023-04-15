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
			"aquarium",
			"art_gallery",
			"bowling_alley",
			"campground",
			"movie_rental",
			"movie_theater",
			"museum",
			"park",
			"rv_park",
			"spa",
			"stadium",
			"tourist_attraction",
		},
	}

	CategoryRestaurant = LocationCategory{
		Name: "Restaurant",
		SubCategories: []string{
			"bakery",
			"bar",
			"cafe",
			"food",
			"restaurant",
		},
	}

	CategoryBook = LocationCategory{
		Name: "Book",
		SubCategories: []string{
			"book_store",
			"library",
		},
	}

	CategoryShopping = LocationCategory{
		Name: "Shopping",
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

	CategoryOutdoor = LocationCategory{}
)
