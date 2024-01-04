package rdb

import (
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"poroto.app/poroto/planner/internal/domain/array"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/entities"
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
		qm.Load(relation + "." + entities.PlaceRels.GooglePlaces),
		qm.Load(relation + "." + entities.PlaceRels.GooglePlaces + "." + entities.GooglePlaceRels.GooglePlaceTypes),
		qm.Load(relation + "." + entities.PlaceRels.GooglePlaces + "." + entities.GooglePlaceRels.GooglePlacePhotos),
		qm.Load(relation + "." + entities.PlaceRels.GooglePlaces + "." + entities.GooglePlaceRels.GooglePlacePhotoAttributions),
		qm.Load(relation + "." + entities.PlaceRels.GooglePlaces + "." + entities.GooglePlaceRels.GooglePlaceReviews),
		qm.Load(relation + "." + entities.PlaceRels.GooglePlaces + "." + entities.GooglePlaceRels.GooglePlaceOpeningPeriods),
	}
}
