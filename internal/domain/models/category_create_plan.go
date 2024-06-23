package models

import (
	"fmt"
	"googlemaps.github.io/maps"
	"os"
	"poroto.app/poroto/planner/internal/env"
)

type LocationCategorySetCreatePlan struct {
	Name             string
	DisplayNameJa    string
	DisplayNameEn    string
	Categories       []LocationCategoryCreatePlan
	GooglePlaceTypes []string
}

type LocationCategoryCreatePlan struct {
	Id                  string
	DisplayNameJa       string
	DisplayNameEn       string
	GooglePlaceTypes    []string
	SearchRadiusMinInKm float64 // 検索半径（最小）
	Image               string
}

func getCloudStorageImageUrl(filename string) string {
	env.LoadEnv(env.WithSkipErrors())
	return fmt.Sprintf("https://storage.googleapis.com/%s/public/images/create_plan_categories/%s", os.Getenv("CLOUD_STORAGE_IMAGE_BUCKET"), filename)
}

var (
	LocationCategorySetCreatePlanAmusements = LocationCategorySetCreatePlan{
		Name:          "amusements",
		DisplayNameJa: "遊び",
		DisplayNameEn: "Amusements",
		Categories: []LocationCategoryCreatePlan{
			{
				Id:               "amusement_park",
				DisplayNameJa:    "遊園地",
				DisplayNameEn:    "Amusement Park",
				GooglePlaceTypes: []string{string(maps.PlaceTypeAmusementPark)},
				Image:            getCloudStorageImageUrl("amusement_park.jpg"),
			},
			{
				Id:               "bowling_alley",
				DisplayNameJa:    "ボウリング場",
				DisplayNameEn:    "Bowling Alley",
				GooglePlaceTypes: []string{string(maps.PlaceTypeBowlingAlley)},
				Image:            getCloudStorageImageUrl("bowling.jpg"),
			},
			{
				Id:               "movie_theater",
				DisplayNameJa:    "映画館",
				DisplayNameEn:    "Movie Theater",
				GooglePlaceTypes: []string{string(maps.PlaceTypeMovieTheater)},
				Image:            getCloudStorageImageUrl("movie.jpg"),
			},
		},
		GooglePlaceTypes: CategoryAmusements.SubCategories,
	}

	LocationCategorySetCreatePlanAttractions = LocationCategorySetCreatePlan{
		Name:          "attractions",
		DisplayNameJa: "観光",
		DisplayNameEn: "Attractions",
		Categories: []LocationCategoryCreatePlan{
			{
				Id:               "観光スポット",
				DisplayNameJa:    "観光スポット",
				DisplayNameEn:    "Sightseeing",
				GooglePlaceTypes: []string{string(maps.PlaceTypeTouristAttraction)},
				Image:            getCloudStorageImageUrl("tourist_attraction.jpg"),
			},
			{
				Id:               "寺・神社",
				DisplayNameJa:    "寺・神社",
				DisplayNameEn:    "Temples & Shrines",
				GooglePlaceTypes: []string{"place_of_worship"},
				Image:            getCloudStorageImageUrl("temple.jpg"),
			},
		},
	}

	LocationCategorySetCreatePlanCulture = LocationCategorySetCreatePlan{
		Name:          "cultural_facility",
		DisplayNameJa: "芸術・動物",
		DisplayNameEn: "Culture",
		Categories: []LocationCategoryCreatePlan{
			{
				Id:               "art_gallery",
				DisplayNameJa:    "美術館",
				DisplayNameEn:    "Art Gallery",
				GooglePlaceTypes: []string{string(maps.PlaceTypeArtGallery)},
				Image:            getCloudStorageImageUrl("art_gallery.jpg"),
			},
			{
				Id:               "museum",
				DisplayNameJa:    "博物館",
				DisplayNameEn:    "Museum",
				GooglePlaceTypes: []string{string(maps.PlaceTypeMuseum)},
				Image:            getCloudStorageImageUrl("museum.jpg"),
			},
			{
				Id:                  "aquarium",
				DisplayNameJa:       "水族館",
				DisplayNameEn:       "Aquarium",
				GooglePlaceTypes:    []string{string(maps.PlaceTypeAquarium)},
				SearchRadiusMinInKm: 30,
				Image:               getCloudStorageImageUrl("aquarium.jpg"),
			},
			{
				Id:               "zoo",
				DisplayNameJa:    "動物園",
				DisplayNameEn:    "Zoo",
				GooglePlaceTypes: []string{string(maps.PlaceTypeZoo)},
				Image:            getCloudStorageImageUrl("zoo.jpg"),
			},
		},
	}

	LocationCategorySetCreatePlanEat = LocationCategorySetCreatePlan{
		Name:          "eat",
		DisplayNameJa: "食事",
		DisplayNameEn: "Eat",
		Categories: []LocationCategoryCreatePlan{
			{
				Id:               "restaurant",
				DisplayNameJa:    "レストラン",
				DisplayNameEn:    "Restaurant",
				GooglePlaceTypes: []string{string(maps.PlaceTypeRestaurant)},
				Image:            getCloudStorageImageUrl("restaurant.jpg"),
			},
			{
				Id:               "cafe",
				DisplayNameJa:    "カフェ",
				DisplayNameEn:    "Cafe",
				GooglePlaceTypes: []string{string(maps.PlaceTypeCafe)},
				Image:            getCloudStorageImageUrl("cafe.jpg"),
			},
			{
				Id:               "bakery",
				DisplayNameJa:    "パン屋",
				DisplayNameEn:    "Bakery",
				GooglePlaceTypes: []string{string(maps.PlaceTypeBakery)},
				Image:            getCloudStorageImageUrl("bakery.jpg"),
			},
		},
	}

	LocationCategorySetCreatePlanRelaxation = LocationCategorySetCreatePlan{
		Name:          "relaxation",
		DisplayNameJa: "リラックス",
		DisplayNameEn: "Relaxation",
		Categories: []LocationCategoryCreatePlan{
			{
				Id:               "spa",
				DisplayNameJa:    "温泉",
				DisplayNameEn:    "Spa",
				GooglePlaceTypes: []string{string(maps.PlaceTypeSpa)},
				Image:            getCloudStorageImageUrl("spa.jpg"),
			},
			{
				Id:               "park",
				DisplayNameJa:    "公園",
				DisplayNameEn:    "Park",
				GooglePlaceTypes: []string{string(maps.PlaceTypePark)},
				Image:            getCloudStorageImageUrl("park.jpg"),
			},
		},
	}

	LocationCategorySetCreatePlanShopping = LocationCategorySetCreatePlan{
		Name:          "shopping",
		DisplayNameJa: "ショッピング",
		DisplayNameEn: "Shopping",
		Categories: []LocationCategoryCreatePlan{
			{
				Id:               "shopping_mall",
				DisplayNameJa:    "ショッピングモール",
				DisplayNameEn:    "Shopping Mall",
				GooglePlaceTypes: []string{string(maps.PlaceTypeShoppingMall)},
				Image:            getCloudStorageImageUrl("shopping_mall.jpg"),
			},
			{
				Id:               "本屋",
				DisplayNameJa:    "本屋",
				DisplayNameEn:    "Bookstore",
				GooglePlaceTypes: []string{string(maps.PlaceTypeBookStore)},
				Image:            getCloudStorageImageUrl("bookstore.jpg"),
			},
		},
	}
)

func GetAllLocationCategorySetCreatePlan() []LocationCategorySetCreatePlan {
	return []LocationCategorySetCreatePlan{
		LocationCategorySetCreatePlanEat,
		LocationCategorySetCreatePlanRelaxation,
		LocationCategorySetCreatePlanAttractions,
		LocationCategorySetCreatePlanShopping,
		LocationCategorySetCreatePlanAmusements,
		LocationCategorySetCreatePlanCulture,
	}
}
