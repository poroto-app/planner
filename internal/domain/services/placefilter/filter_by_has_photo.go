package placefilter

import "poroto.app/poroto/planner/internal/domain/models"

// FilterByHasPhoto は写真を取得済み または　取得可能な場所のみを返す
func FilterByHasPhoto(placesToFilter []models.Place) []models.Place {
	return FilterPlaces(placesToFilter, func(place models.Place) bool {
		hasPhoto := place.Google.Photos != nil && len(*place.Google.Photos) > 0
		hasPhotoReference := place.Google.PhotoReferences != nil && len(place.Google.PhotoReferences) > 0
		hasPhotoReferenceInPlaceDetail := place.Google.PlaceDetail != nil && len(place.Google.PlaceDetail.PhotoReferences) > 0
		return hasPhoto || hasPhotoReference || hasPhotoReferenceInPlaceDetail
	})
}
