package place

import (
	"context"
	"poroto.app/poroto/planner/internal/domain/models"
)

// FetchPlacesRecommended は場所を指定してプランを作成するときに、おすすめの場所を取得する
func (s Service) FetchPlacesRecommended(ctx context.Context) (*[]models.Place, error) {
	return s.placeRepository.FindRecommendPlacesForCreatePlan(ctx)
}
