package handlers

import "net/http"

func (api API) dashboardStats(w http.ResponseWriter, r *http.Request) {
	stats, err := api.dashboard.Stats(r.Context(), currentUser(r))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, stats)
}
