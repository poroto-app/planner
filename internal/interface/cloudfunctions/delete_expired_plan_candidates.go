package cloudfunctions

import (
	"log"
	"net/http"

	"poroto.app/poroto/planner/pkg/batch"
)

func DeleteExpiredPlanCandidates(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if err := batch.DeleteExpiredPlanCandidate(r.Context()); err != nil {
		log.Printf("error while deleting expired plan candidates: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
