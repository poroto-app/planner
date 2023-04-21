package models

type LocationCategory struct {
	Name          string
	DisplayName   string
	SubCategories []string
}

var (
	// SEE: https://developers.google.com/maps/documentation/places/web-service/supported_types?hl=ja#table1
	CategoryAmusements = LocationCategory{
		Name:        "amusements",
		DisplayName: "遊び",
		SubCategories: []string{
			"amusement_park",
			"bowling_alley",
			"movie_theater",
			"spa",
			"stadium",
		},
	}

	CategoryBook = LocationCategory{
		Name:        "book",
		DisplayName: "本",
		SubCategories: []string{
			"book_store",
			"library",
		},
	}

	CategoryCamp = LocationCategory{
		Name:        "camp",
		DisplayName: "キャンプ",
		SubCategories: []string{
			"campground",
			"rv_park",
		},
	}

	CategoryCafe = LocationCategory{
		Name:        "cafe",
		DisplayName: "カフェ",
		SubCategories: []string{
			"cafe",
		},
	}

	CategoryCulture = LocationCategory{
		Name:        "cultural_facility",
		DisplayName: "芸術や文化に触れる",
		SubCategories: []string{
			"art_gallery",
			"museum",
			"tourist_attraction",
		},
	}

	CategoryNatural = LocationCategory{
		Name:        "natural_facility",
		DisplayName: "動物を見に行こう",
		SubCategories: []string{
			"aquarium",
			"zoo",
		},
	}

	CategoryPark = LocationCategory{
		Name:        "park",
		DisplayName: "公園でゆったり",
		SubCategories: []string{
			"park",
		},
	}

	CategoryRestaurant = LocationCategory{
		Name:        "restaurant",
		DisplayName: "ご飯",
		SubCategories: []string{
			"bakery",
			"bar",
			"food",
			"restaurant",
		},
	}

	CategoryShopping = LocationCategory{
		Name:        "shopping",
		DisplayName: "ショッピング",
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
