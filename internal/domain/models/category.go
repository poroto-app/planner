package models

import "poroto.app/poroto/planner/internal/domain/array"

type LocationCategory struct {
	Name                  string
	DisplayName           string
	SubCategories         []string
	Photo                 string
	EstimatedStayDuration uint
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
		// TODO: implement me!
		Photo:                 "https://placehold.jp/3d4070/ffffff/300x500.png?text=amusement",
		EstimatedStayDuration: 90,
	}

	CategoryBookStore = LocationCategory{
		Name:        "book_store",
		DisplayName: "本屋",
		SubCategories: []string{
			"book_store",
		},
		// TODO: implement me!
		Photo:                 "https://placehold.jp/80ddff/ffffff/300x500.png?text=book",
		EstimatedStayDuration: 30,
	}

	CategoryCafe = LocationCategory{
		Name:        "cafe",
		DisplayName: "カフェ",
		SubCategories: []string{
			"cafe",
		},
		// TODO: implement me!
		Photo:                 "https://placehold.jp/ff9620/ffffff/300x500.png?text=cafe",
		EstimatedStayDuration: 60,
	}

	CategoryCamp = LocationCategory{
		Name:        "camp",
		DisplayName: "キャンプ",
		SubCategories: []string{
			"campground",
			"rv_park",
		},
		// TODO: implement me!
		Photo:                 "https://placehold.jp/40ff20/ffffff/300x500.png?text=camp",
		EstimatedStayDuration: 300,
	}

	CategoryCulture = LocationCategory{
		Name:        "cultural_facility",
		DisplayName: "芸術や文化に触れる",
		SubCategories: []string{
			"art_gallery",
			"museum",
			"tourist_attraction",
		},
		// TODO: implement me!
		Photo:                 "https://placehold.jp/8f8f8f/ffffff/300x500.png?text=cultural%0Afacility",
		EstimatedStayDuration: 90,
	}

	CategoryNatural = LocationCategory{
		Name:        "natural_facility",
		DisplayName: "動物を見に行こう",
		SubCategories: []string{
			"aquarium",
			"zoo",
		},
		// TODO: implement me!
		Photo:                 "https://placehold.jp/00ffbf/ffffff/300x500.png?text=natural%0Afacility",
		EstimatedStayDuration: 120,
	}

	CategoryPark = LocationCategory{
		Name:        "park",
		DisplayName: "公園でゆったり",
		SubCategories: []string{
			"park",
		},
		// TODO: implement me!
		Photo:                 "https://placehold.jp/fbff00/ffffff/300x500.png?text=park",
		EstimatedStayDuration: 30,
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
		// TODO: implement me!
		Photo:                 "https://placehold.jp/ff7070/ffffff/300x500.png?text=restaurant",
		EstimatedStayDuration: 60,
	}

	CategoryLibrary = LocationCategory{
		Name:                  "library",
		DisplayName:           "図書館",
		SubCategories:         []string{"library"},
		Photo:                 "https://placehold.jp/ff7070/ffffff/300x500.png?text=library",
		EstimatedStayDuration: 30,
	}

	CategoryMealTakeaway = LocationCategory{
		Name:        "meal_takeaway",
		DisplayName: "テイクアウト",
		SubCategories: []string{
			"meal_takeaway",
		},
		// TODO: implement me!
		Photo:                 "https://placehold.jp/1d7187/ffffff/300x500.png?text=quick%0Aservice%0Arestaurant",
		EstimatedStayDuration: 30,
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
		// TODO: implement me!
		Photo:                 "https://placehold.jp/70dbff/ffffff/300x500.png?text=shopping",
		EstimatedStayDuration: 60,
	}
)

func GetCategoryToFilter() []LocationCategory {
	return []LocationCategory{
		CategoryAmusements,
		CategoryBookStore,
		CategoryCamp,
		CategoryCafe,
		CategoryCulture,
		CategoryLibrary,
		CategoryNatural,
		CategoryPark,
		CategoryRestaurant,
		CategoryShopping,
	}
}

func getAllCategories() []LocationCategory {
	return []LocationCategory{
		CategoryAmusements,
		CategoryBookStore,
		CategoryCafe,
		CategoryCamp,
		CategoryCulture,
		CategoryNatural,
		CategoryPark,
		CategoryRestaurant,
		CategoryShopping,
	}
}

// GetCategoryOfName name に対応する LocationCategory を返す
// name が見つからない場合は nil を返す
// NOTE: category の値が上書きされないようにコピーを渡している
func GetCategoryOfName(name string) *LocationCategory {
	for _, category := range getAllCategories() {
		if category.Name == name {
			c := category
			return &c
		}
	}
	return nil
}

// CategoryOfSubCategory SubCategory がどの大カテゴリに所属するか
func CategoryOfSubCategory(subCategory string) *LocationCategory {
	for _, category := range getAllCategories() {
		if array.IsContain(category.SubCategories, subCategory) {
			return &category
		}
	}

	return nil
}
