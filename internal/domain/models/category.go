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
		Photo: "https://placehold.jp/3d4070/ffffff/300x500.png?text=amusement",
		// TODO: implement me!
		EstimatedStayDuration: 30,
	}

	CategoryBook = LocationCategory{
		Name:        "book",
		DisplayName: "本",
		SubCategories: []string{
			"book_store",
			"library",
		},
		// TODO: implement me!
		Photo: "https://placehold.jp/80ddff/ffffff/300x500.png?text=book",
		// TODO: implement me!
		EstimatedStayDuration: 30,
	}

	CategoryCafe = LocationCategory{
		Name:        "cafe",
		DisplayName: "カフェ",
		SubCategories: []string{
			"cafe",
		},
		// TODO: implement me!
		Photo: "https://placehold.jp/ff9620/ffffff/300x500.png?text=cafe",
		// TODO: implement me!
		EstimatedStayDuration: 30,
	}

	CategoryCamp = LocationCategory{
		Name:        "camp",
		DisplayName: "キャンプ",
		SubCategories: []string{
			"campground",
			"rv_park",
		},
		// TODO: implement me!
		Photo: "https://placehold.jp/40ff20/ffffff/300x500.png?text=camp",
		// TODO: implement me!
		EstimatedStayDuration: 30,
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
		Photo: "https://placehold.jp/8f8f8f/ffffff/300x500.png?text=cultural%0Afacility",
		// TODO: implement me!
		EstimatedStayDuration: 30,
	}

	CategoryNatural = LocationCategory{
		Name:        "natural_facility",
		DisplayName: "動物を見に行こう",
		SubCategories: []string{
			"aquarium",
			"zoo",
		},
		// TODO: implement me!
		Photo: "https://placehold.jp/00ffbf/ffffff/300x500.png?text=natural%0Afacility",
		// TODO: implement me!
		EstimatedStayDuration: 30,
	}

	CategoryPark = LocationCategory{
		Name:        "park",
		DisplayName: "公園でゆったり",
		SubCategories: []string{
			"park",
		},
		// TODO: implement me!
		Photo: "https://placehold.jp/fbff00/ffffff/300x500.png?text=park",
		// TODO: implement me!
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
			"meal_delivery",
			"meal_takeaway",
		},
		// TODO: implement me!
		Photo: "https://placehold.jp/ff7070/ffffff/300x500.png?text=restaurant",
		// TODO: implement me!
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
		Photo: "https://placehold.jp/70dbff/ffffff/300x500.png?text=shopping",
		// TODO: implement me!
		EstimatedStayDuration: 30,
	}
)

func GetCategoryOfName(name string) LocationCategory {
	return map[string]LocationCategory{
		CategoryAmusements.Name: CategoryAmusements,
		CategoryBook.Name:       CategoryBook,
		CategoryCafe.Name:       CategoryCafe,
		CategoryCamp.Name:       CategoryCamp,
		CategoryCulture.Name:    CategoryCulture,
		CategoryNatural.Name:    CategoryNatural,
		CategoryPark.Name:       CategoryPark,
		CategoryRestaurant.Name: CategoryRestaurant,
		CategoryShopping.Name:   CategoryShopping,
	}[name]
}

// SubCategory がどの大カテゴリに所属するか
func CategoryOfSubCategory(subCategory string) *LocationCategory {
	var allCategory = []LocationCategory{
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

	for _, category := range allCategory {
		if array.IsContain(category.SubCategories, subCategory) {
			return &category
		}
	}

	return nil
}
