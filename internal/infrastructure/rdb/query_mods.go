package rdb

import (
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
	"strings"
)

func concatQueryMod(qms ...[]qm.QueryMod) []qm.QueryMod {
	return array.Flatten(qms)
}

// placeQueryModes models.Place を作成するのに必要な関連をロードするための qm.QueryMod を返す
// relations には X.X."Places" というように "Places" までの関連を指定する
func placeQueryModes(relationsToPlaces ...string) []qm.QueryMod {
	if len(relationsToPlaces) == 0 {
		panic("relationsToPlaces must be specified")
	}
	relation := strings.Join(relationsToPlaces, ".")
	return []qm.QueryMod{
		qm.Load(relation),
		qm.Load(relation + "." + generated.PlaceRels.GooglePlaces),
		qm.Load(relation + "." + generated.PlaceRels.GooglePlaces + "." + generated.GooglePlaceRels.GooglePlaceTypes),
		qm.Load(relation + "." + generated.PlaceRels.GooglePlaces + "." + generated.GooglePlaceRels.GooglePlacePhotoReferences),
		qm.Load(relation + "." + generated.PlaceRels.GooglePlaces + "." + generated.GooglePlaceRels.GooglePlacePhotos),
		qm.Load(relation + "." + generated.PlaceRels.GooglePlaces + "." + generated.GooglePlaceRels.GooglePlacePhotoAttributions),
		qm.Load(relation + "." + generated.PlaceRels.GooglePlaces + "." + generated.GooglePlaceRels.GooglePlaceReviews),
		qm.Load(relation + "." + generated.PlaceRels.GooglePlaces + "." + generated.GooglePlaceRels.GooglePlaceOpeningPeriods),
	}
}
