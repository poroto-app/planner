package rdb

import (
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
	"strings"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"poroto.app/poroto/planner/internal/domain/array"
)

func toInterfaceArray[T any](arr []T) []interface{} {
	var result []interface{}
	for _, a := range arr {
		result = append(result, a)
	}
	return result
}

func concatQueryMod(qms ...[]qm.QueryMod) []qm.QueryMod {
	return array.Flatten(qms)
}

// placeQueryModes models.Place を作成するのに必要な関連をロードするための qm.QueryMod を返す
// relations には X.X."Places" というように "Places" までの関連を指定する
// places から直接参照する場合は 引数を指定しない
func placeQueryModes(relationsToPlaces ...string) []qm.QueryMod {
	var relation string
	if len(relationsToPlaces) == 0 {
		// Places がら直接参照する場合は PlacePhotos, GooglePlaces... となるようにする
		relation = ""
	} else {
		relation = strings.Join(relationsToPlaces, ".")
	}

	var queryMods []qm.QueryMod
	if len(relationsToPlaces) > 0 {
		// Places をロード
		queryMods = append(queryMods, qm.Load(relation))

		// XXX.Places.PlacePhotos, XXX.Places.GooglePlaces... となるようにする
		relation += "."
	}

	return concatQueryMod(
		queryMods,
		[]qm.QueryMod{
			qm.Load(relation + generated.PlaceRels.PlacePhotos),
			qm.Load(relation + generated.PlaceRels.GooglePlaces),
			qm.Load(relation + generated.PlaceRels.GooglePlaces + "." + generated.GooglePlaceRels.GooglePlaceTypes),
			qm.Load(relation + generated.PlaceRels.GooglePlaces + "." + generated.GooglePlaceRels.GooglePlacePhotoReferences),
			qm.Load(relation + generated.PlaceRels.GooglePlaces + "." + generated.GooglePlaceRels.GooglePlacePhotos),
			qm.Load(relation + generated.PlaceRels.GooglePlaces + "." + generated.GooglePlaceRels.GooglePlacePhotoAttributions),
			qm.Load(relation + generated.PlaceRels.GooglePlaces + "." + generated.GooglePlaceRels.GooglePlaceReviews),
			qm.Load(relation + generated.PlaceRels.GooglePlaces + "." + generated.GooglePlaceRels.GooglePlaceOpeningPeriods),
		},
	)
}
