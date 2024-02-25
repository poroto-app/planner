package place

import (
	"context"
	"fmt"
	"log"

	"go.uber.org/zap"
	"poroto.app/poroto/planner/internal/domain/utils"
)

type UploadPlacePhotoInPlanInput struct {
	UserId   string
	PlaceId  string
	PhotoUrl string
	Width    int
	Height   int
}

func (s Service) UploadPlacePhotoInPlan(
	ctx context.Context,
	input UploadPlacePhotoInPlanInput,
) error {
	logger, err := utils.NewLogger(utils.LoggerOption{Tag: "GraphQL"})
	if err != nil {
		log.Println("error while initializing logger: ", err)
		return fmt.Errorf("internal server error")
	}

	err = s.placeRepository.SavePlacePhotos(ctx, input.UserId, input.PlaceId, input.PhotoUrl, input.Width, input.Height)
	if err != nil {
		logger.Error("error while saving place photos", zap.Error(err))
		return fmt.Errorf("internal server error")
	}
	return nil
}
