package plangen

import (
	"fmt"
	"poroto.app/poroto/planner/internal/domain/models"
	"time"
)

// isOpeningWithIn は，指定された場所が指定された時間内に開いているかを判定する
func (s Service) isOpeningWithIn(place models.PlaceInPlanCandidate, startTime time.Time, duration time.Duration) (bool, error) {
	isOpeningAtStartTime, err := place.Google.IsOpening(startTime)
	if err != nil {
		return false, fmt.Errorf("error while checking opening hours: %v", err)
	}

	endTime := startTime.Add(duration)
	isOpeningAtEndTime, err := place.Google.IsOpening(endTime)
	if err != nil {
		return false, fmt.Errorf("error while checking opening hours: %v", err)
	}

	return isOpeningAtStartTime && isOpeningAtEndTime, nil
}
