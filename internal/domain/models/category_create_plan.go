package models

import "googlemaps.github.io/maps"

type LocationCategorySetCreatePlan struct {
	Name             string
	DisplayNameJa    string
	DisplayNameEn    string
	Categories       []LocationCategoryCreatePlan
	GooglePlaceTypes []string
}

type LocationCategoryCreatePlan struct {
	Name                string
	DisplayNameJa       string
	DisplayNameEn       string
	GooglePlaceTypes    []string
	SearchRadiusMinInKm int // 検索半径（最小）
	Image               string
}

var (
	CreatePlanCategoryAmusements = LocationCategorySetCreatePlan{
		Name:          "amusements",
		DisplayNameJa: "遊び",
		DisplayNameEn: "Amusements",
		Categories: []LocationCategoryCreatePlan{
			{
				Name:             "amusement_park",
				DisplayNameJa:    "遊園地",
				DisplayNameEn:    "Amusement Park",
				GooglePlaceTypes: []string{string(maps.PlaceTypeAmusementPark)},
			},
			{
				Name:             "bowling_alley",
				DisplayNameJa:    "ボウリング場",
				DisplayNameEn:    "Bowling Alley",
				GooglePlaceTypes: []string{string(maps.PlaceTypeBowlingAlley)},
			},
			{
				Name:             "movie_theater",
				DisplayNameJa:    "映画館",
				DisplayNameEn:    "Movie Theater",
				GooglePlaceTypes: []string{string(maps.PlaceTypeMovieTheater)},
			},
		},
		GooglePlaceTypes: CategoryAmusements.SubCategories,
	}

	CreatePlaceCategoryCulture = LocationCategorySetCreatePlan{
		Name:          "cultural_facility",
		DisplayNameJa: "芸術・動物",
		DisplayNameEn: "Culture",
		Categories: []LocationCategoryCreatePlan{
			{
				Name:             "art_gallery",
				DisplayNameJa:    "美術館",
				DisplayNameEn:    "Art Gallery",
				GooglePlaceTypes: []string{string(maps.PlaceTypeArtGallery)},
			},
			{
				Name:             "museum",
				DisplayNameJa:    "博物館",
				DisplayNameEn:    "Museum",
				GooglePlaceTypes: []string{string(maps.PlaceTypeMuseum)},
			},
			{
				Name:                "aquarium",
				DisplayNameJa:       "水族館",
				DisplayNameEn:       "Aquarium",
				GooglePlaceTypes:    []string{string(maps.PlaceTypeAquarium)},
				SearchRadiusMinInKm: 30,
			},
			{
				Name:             "zoo",
				DisplayNameJa:    "動物園",
				DisplayNameEn:    "Zoo",
				GooglePlaceTypes: []string{string(maps.PlaceTypeZoo)},
			},
		},
	}

	CreatePlanCategorySetRelaxation = LocationCategorySetCreatePlan{
		Name:          "relaxation",
		DisplayNameJa: "リラックス",
		DisplayNameEn: "Relaxation",
		Categories: []LocationCategoryCreatePlan{
			{
				Name:             "spa",
				DisplayNameJa:    "温泉",
				DisplayNameEn:    "Spa",
				GooglePlaceTypes: []string{string(maps.PlaceTypeSpa)},
			},
			{
				Name:             "park",
				DisplayNameJa:    "公園",
				DisplayNameEn:    "Park",
				GooglePlaceTypes: []string{string(maps.PlaceTypePark)},
			},
		},
	}

	CreatePlanCategorySetShopping = LocationCategorySetCreatePlan{
		Name:          "shopping",
		DisplayNameJa: "ショッピング",
		DisplayNameEn: "Shopping",
		Categories: []LocationCategoryCreatePlan{
			{
				Name:             "shopping_mall",
				DisplayNameJa:    "ショッピングモール",
				DisplayNameEn:    "Shopping Mall",
				GooglePlaceTypes: []string{string(maps.PlaceTypeShoppingMall)},
			},
			{
				Name:             "本屋",
				DisplayNameJa:    "本屋",
				DisplayNameEn:    "Bookstore",
				GooglePlaceTypes: []string{string(maps.PlaceTypeBookStore)},
			},
		},
	}

	CreatePlanCategorySetEat = LocationCategorySetCreatePlan{
		Name:          "eat",
		DisplayNameJa: "食事",
		DisplayNameEn: "Eat",
		Categories: []LocationCategoryCreatePlan{
			{
				Name:             "restaurant",
				DisplayNameJa:    "レストラン",
				DisplayNameEn:    "Restaurant",
				GooglePlaceTypes: []string{string(maps.PlaceTypeRestaurant)},
			},
			{
				Name:             "cafe",
				DisplayNameJa:    "カフェ",
				DisplayNameEn:    "Cafe",
				GooglePlaceTypes: []string{string(maps.PlaceTypeCafe)},
			},
			{
				Name:             "bakery",
				DisplayNameJa:    "パン屋",
				DisplayNameEn:    "Bakery",
				GooglePlaceTypes: []string{string(maps.PlaceTypeBakery)},
			},
		},
	}

	CreatePlanCategorySetAttractions = LocationCategorySetCreatePlan{
		Name:          "attractions",
		DisplayNameJa: "観光",
		DisplayNameEn: "Attractions",
		Categories: []LocationCategoryCreatePlan{
			{
				Name:             "観光スポット",
				DisplayNameJa:    "観光スポット",
				DisplayNameEn:    "Sightseeing",
				GooglePlaceTypes: []string{string(maps.PlaceTypeTouristAttraction)},
			},
			{
				Name:             "寺・神社",
				DisplayNameJa:    "寺・神社",
				DisplayNameEn:    "Temples & Shrines",
				GooglePlaceTypes: []string{"place_of_worship"},
			},
		},
	}
)
