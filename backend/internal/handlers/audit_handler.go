package handlers

import "net/http"

func (api API) listVisibleAudit(w http.ResponseWriter, r *http.Request) {
	events, err := api.audit.ListVisibleEvents(r.Context(), currentUser(r))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, events)
}

func (api API) listSystemAudit(w http.ResponseWriter, r *http.Request) {
	events, err := api.audit.ListSystemEvents(r.Context(), currentUser(r))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, events)
}
