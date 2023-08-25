package models

import (
	"googlemaps.github.io/maps"
	"poroto.app/poroto/planner/internal/domain/array"
)

// LocationCategory は場所の大まかなカテゴリを示す
// TODO: Photo は nilable にする
type LocationCategory struct {
	Name                  string
	DisplayName           string
	SubCategories         []string
	DefaultPhoto          string
	Photo                 string
	EstimatedStayDuration uint
}

var (
	// SEE: https://developers.google.com/maps/documentation/places/web-service/supported_types?hl=ja#table1
	CategoryAmusements = LocationCategory{
		Name:        "amusements",
		DisplayName: "遊び",
		SubCategories: []string{
			string(maps.PlaceTypeAmusementPark),
			string(maps.PlaceTypeBowlingAlley),
			string(maps.PlaceTypeMovieTheater),
			string(maps.PlaceTypeSpa),
			string(maps.PlaceTypeStadium),
		},
		DefaultPhoto:          "https://storage.googleapis.com/planner-public-asset-bucket/undraw_amusement_park_17oe.svg",
		EstimatedStayDuration: 90,
	}

	CategoryBookStore = LocationCategory{
		Name:        "book_store",
		DisplayName: "本屋",
		SubCategories: []string{
			string(maps.PlaceTypeBookStore),
		},
		DefaultPhoto:          "https://storage.googleapis.com/planner-public-asset-bucket/undraw_books_re_8gea.svg",
		EstimatedStayDuration: 30,
	}

	CategoryCafe = LocationCategory{
		Name:        "cafe",
		DisplayName: "カフェ",
		SubCategories: []string{
			string(maps.PlaceTypeCafe),
		},
		DefaultPhoto:          "https://storage.googleapis.com/planner-public-asset-bucket/undraw_coffee_re_x35h.svg",
		EstimatedStayDuration: 60,
	}

	CategoryCamp = LocationCategory{
		Name:        "camp",
		DisplayName: "キャンプ",
		SubCategories: []string{
			string(maps.PlaceTypeCampground),
			string(maps.PlaceTypeRvPark),
		},
		DefaultPhoto:          "https://storage.googleapis.com/planner-public-asset-bucket/undraw_camping_noc8.svg",
		EstimatedStayDuration: 300,
	}

	CategoryCulture = LocationCategory{
		Name:        "cultural_facility",
		DisplayName: "芸術や文化に触れる",
		SubCategories: []string{
			string(maps.PlaceTypeArtGallery),
			string(maps.PlaceTypeMuseum),
		},
		DefaultPhoto:          "https://storage.googleapis.com/planner-public-asset-bucket/undraw_art_lover_re_fn8g.svg",
		EstimatedStayDuration: 90,
	}

	CategoryNatural = LocationCategory{
		Name:        "natural_facility",
		DisplayName: "動物を見に行こう",
		SubCategories: []string{
			string(maps.PlaceTypeAquarium),
			string(maps.PlaceTypeZoo),
		},
		DefaultPhoto:          "https://storage.googleapis.com/planner-public-asset-bucket/undraw_fish_bowl_uu88.svg",
		EstimatedStayDuration: 120,
	}

	CategoryPark = LocationCategory{
		Name:        "park",
		DisplayName: "公園でゆったり",
		SubCategories: []string{
			string(maps.PlaceTypePark),
		},
		DefaultPhoto:          "https://storage.googleapis.com/planner-public-asset-bucket/undraw_a_day_at_the_park_re_9kxj.svg",
		EstimatedStayDuration: 30,
	}

	CategoryRestaurant = LocationCategory{
		Name:        "restaurant",
		DisplayName: "ご飯",
		SubCategories: []string{
			string(maps.PlaceDetailsFieldMaskServesVegetarianFood),
			string(maps.PlaceTypeBakery),
			string(maps.PlaceTypeBar),
			string(maps.PlaceTypeRestaurant),
		},
		DefaultPhoto:          "https://storage.googleapis.com/planner-public-asset-bucket/undraw_breakfast_psiw.svg",
		EstimatedStayDuration: 60,
	}

	CategoryLibrary = LocationCategory{
		Name:        "library",
		DisplayName: "図書館",
		SubCategories: []string{
			string(maps.PlaceTypeLibrary),
		},
		Photo:                 "https://placehold.jp/ff7070/ffffff/300x500.png?text=library",
		DefaultPhoto:          "https://storage.googleapis.com/planner-public-asset-bucket/undraw_book_reading_re_fu2c.svg",
		EstimatedStayDuration: 30,
	}

	CategoryMealTakeaway = LocationCategory{
		Name:        "meal_takeaway",
		DisplayName: "テイクアウト",
		SubCategories: []string{
			string(maps.PlaceTypeMealTakeaway),
		},
		DefaultPhoto:          "https://storage.googleapis.com/planner-public-asset-bucket/undraw_pizza_sharing_wxop.svg",
		EstimatedStayDuration: 30,
	}

	CategoryShopping = LocationCategory{
		Name:        "shopping",
		DisplayName: "ショッピング",
		SubCategories: []string{
			string(maps.PlaceTypeClothingStore),
			string(maps.PlaceTypeDepartmentStore),
			string(maps.PlaceTypeFurnitureStore),
			string(maps.PlaceTypeHardwareStore),
			string(maps.PlaceTypeHomeGoodsStore),
			string(maps.PlaceTypeMovieRental),
			string(maps.PlaceTypeShoeStore),
			string(maps.PlaceTypeStore),
		},
		DefaultPhoto:          "https://storage.googleapis.com/planner-public-asset-bucket/undraw_shopping_bags_o6w5.svg",
		EstimatedStayDuration: 60,
	}

	CategoryOther = LocationCategory{
		Name:                  "other",
		DisplayName:           "その他",
		SubCategories:         []string{},
		EstimatedStayDuration: 0,
	}

	CategoryIgnore = LocationCategory{
		Name: "ignore",
		SubCategories: []string{
			string(maps.PlaceTypeAccounting),

			string(maps.PlaceTypeAtm),

			string(maps.PlaceTypeBank),

			string(maps.PlaceTypeBeautySalon),
			string(maps.PlaceTypeBicycleStore),

			string(maps.PlaceTypeBusStation),

			string(maps.PlaceTypeCarDealer),
			string(maps.PlaceTypeCarRental),
			string(maps.PlaceTypeCarRepair),
			string(maps.PlaceTypeCarWash),
			string(maps.PlaceTypeCasino),
			string(maps.PlaceTypeCemetery),
			string(maps.PlaceTypeChurch),
			string(maps.PlaceTypeCityHall),

			string(maps.PlaceTypeConvenienceStore),
			string(maps.PlaceTypeCourthouse),
			string(maps.PlaceTypeDentist),
			string(maps.PlaceTypeDoctor),
			string(maps.PlaceTypeElectrician),
			string(maps.PlaceTypeEmbassy),
			string(maps.PlaceTypeFireStation),

			string(maps.PlaceTypeFuneralHome),
			string(maps.PlaceTypeGasStation),
			string(maps.PlaceTypeGym),
			string(maps.PlaceTypeHairCare),

			string(maps.PlaceTypeHinduTemple),

			string(maps.PlaceTypeHospital),
			string(maps.PlaceTypeInsuranceAgency),
			string(maps.PlaceTypeJewelryStore),
			string(maps.PlaceTypeLaundry),
			string(maps.PlaceTypeLawyer),

			string(maps.PlaceTypeLocalGovernmentOffice),
			string(maps.PlaceTypeLocksmith),
			string(maps.PlaceTypeLodging),

			string(maps.PlaceTypeMosque),

			string(maps.PlaceTypeMovingCompany),

			string(maps.PlaceTypeNightClub),
			string(maps.PlaceTypePainter),
			string(maps.PlaceTypePark),
			string(maps.PlaceTypeParking),
			string(maps.PlaceTypePharmacy),
			string(maps.PlaceTypePhysiotherapist),
			string(maps.PlaceTypePlumber),
			string(maps.PlaceTypePolice),
			string(maps.PlaceTypePostOffice),
			string(maps.PlaceTypeRealEstateAgency),

			string(maps.PlaceTypeRoofingContractor),
			string(maps.PlaceTypeRvPark),
			string(maps.PlaceTypeSchool),
			string(maps.PlaceTypeShoeStore),

			string(maps.PlaceTypeStorage),

			string(maps.PlaceTypeSynagogue),
			string(maps.PlaceTypeTaxiStand),

			string(maps.PlaceTypeTravelAgency),

			string(maps.PlaceTypeUniversity),
			string(maps.PlaceTypeVeterinaryCare),
		},
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
		CategoryMealTakeaway,
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
		CategoryLibrary,
		CategoryNatural,
		CategoryMealTakeaway,
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

// GetCategoriesFromSubCategories subCategories に対応する LocationCategory を重複が無いように返す
func GetCategoriesFromSubCategories(subCategories []string) []LocationCategory {
	categoryNames := make([]string, 0)
	categories := make([]LocationCategory, 0)
	for _, subCategory := range subCategories {
		category := CategoryOfSubCategory(subCategory)
		if category != nil && !array.IsContain(categoryNames, category.Name) {
			categories = append(categories, *category)
			categoryNames = append(categoryNames, category.Name)
		}
	}
	return categories
}
