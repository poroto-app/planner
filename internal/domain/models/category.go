package models

import "poroto.app/poroto/planner/internal/domain/array"

type LocationCategory struct {
	Name          string
	SubCategories []string
	Photo         string
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
		// TODO: implement me!
		Photo: "https://placehold.jp/3d4070/ffffff/300x500.png?text=amusement",
	}

	CategoryBook = LocationCategory{
		Name: "book",
		SubCategories: []string{
			"book_store",
			"library",
		},
		// TODO: implement me!
		Photo: "https://placehold.jp/80ddff/ffffff/300x500.png?text=book",
	}

	CategoryCafe = LocationCategory{
		Name: "cafe",
		SubCategories: []string{
			"cafe",
		},
		// TODO: implement me!
		Photo: "https://placehold.jp/ff9620/ffffff/300x500.png?text=cafe",
	}

	CategoryCamp = LocationCategory{
		Name: "camp",
		SubCategories: []string{
			"campground",
			"rv_park",
		},
		// TODO: implement me!
		Photo: "https://placehold.jp/40ff20/ffffff/300x500.png?text=camp",
	}

	CategoryCulture = LocationCategory{
		Name: "cultural Facility",
		SubCategories: []string{
			"art_gallery",
			"museum",
			"tourist_attraction",
		},
		// TODO: implement me!
		Photo: "https://placehold.jp/8f8f8f/ffffff/300x500.png?text=cultural%0Afacility",
	}

	CategoryNatural = LocationCategory{
		Name: "natural Facility",
		SubCategories: []string{
			"aquarium",
			"zoo",
		},
		// TODO: implement me!
		Photo: "https://placehold.jp/00ffbf/ffffff/300x500.png?text=natural%0Afacility",
	}

	CategoryPark = LocationCategory{
		Name: "park",
		SubCategories: []string{
			"park",
		},
		// TODO: implement me!
		Photo: "https://placehold.jp/fbff00/ffffff/300x500.png?text=park",
	}

	CategoryRestaurant = LocationCategory{
		Name: "restaurant",
		SubCategories: []string{
			"bakery",
			"bar",
			"food",
			"restaurant",
		},
		// TODO: implement me!
		Photo: "https://placehold.jp/ff7070/ffffff/300x500.png?text=restaurant",
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
		// TODO: implement me!
		Photo: "https://placehold.jp/70dbff/ffffff/300x500.png?text=shopping",
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
