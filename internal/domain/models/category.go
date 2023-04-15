package models

type LocationCategory struct {
	Name          string
	SubCategories []string
}

var (
	CategoryAmusements = LocationCategory{
		Name: "amusements",
		SubCategories: []string{
			"amusement_park",
			"aquarium",
			"art_gallery",
			"museum",
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
)
