package cloudfunctions

import (
	"log"
	"net/http"
	"poroto.app/poroto/planner/internal/infrastructure/rdb"

	"poroto.app/poroto/planner/pkg/batch"
)

func DeleteExpiredPlanCandidates(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	db, err := rdb.InitDB(false)
	if err != nil {
		log.Printf("error while initializing db: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := batch.DeleteExpiredPlanCandidateSet(r.Context(), db); err != nil {
		log.Printf("error while deleting expired plan candidates: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
