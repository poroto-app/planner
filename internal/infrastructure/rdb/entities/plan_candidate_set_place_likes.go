package entities

import "poroto.app/poroto/planner/internal/domain/array"

type PlanCandidateSetPlaceLikeCount struct {
	PlaceId    string `boil:"place_id"`
	LikesCount int    `boil:"likes_count"`
}

func CountLikeOfPlace(planCandidateSetPlaceLikeCounts *[]PlanCandidateSetPlaceLikeCount, placeId string) int {
	var likeCount int
	if planCandidateSetPlaceLikeCounts != nil {
		likeCountOfPlace, ok := array.Find(*planCandidateSetPlaceLikeCounts, func(planCandidateSetLikePlace PlanCandidateSetPlaceLikeCount) bool {
			return planCandidateSetLikePlace.PlaceId == placeId
		})
		if ok {
			likeCount = likeCountOfPlace.LikesCount
		}
	}
	return likeCount
}
