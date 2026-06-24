package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"openownership-workflow/backend/internal/services"
	"openownership-workflow/backend/internal/workflow"
)

func (api API) listSubmissions(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if query.Get("page") != "" || query.Get("pageSize") != "" {
		page := parsePositiveInt(query.Get("page"), 1)
		pageSize := parsePositiveInt(query.Get("pageSize"), 10)
		result, err := api.submissions.ListPage(r.Context(), currentUser(r), query.Get("status"), page, pageSize)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"items":    result.Items,
			"total":    result.Total,
			"page":     result.Page,
			"pageSize": result.PageSize,
		})
		return
	}
	items, err := api.submissions.List(r.Context(), currentUser(r), query.Get("status"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func parsePositiveInt(value string, fallback int) int {
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 1 {
		return fallback
	}
	return parsed
}

func (api API) getSubmission(w http.ResponseWriter, r *http.Request) {
	item, err := api.submissions.Get(r.Context(), chi.URLParam(r, "id"), currentUser(r))
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = http.StatusNotFound
		}
		writeError(w, status, "submission not found")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (api API) createSubmission(w http.ResponseWriter, r *http.Request) {
	user := currentUser(r)
	if !user.HasPermission("submissions:create") {
		writeError(w, http.StatusForbidden, "missing submissions:create permission")
		return
	}
	payload, ok := api.readSubmissionPayload(w, r)
	if !ok {
		return
	}
	item, err := api.submissions.Create(r.Context(), user, payload)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, services.ErrInvalidSubmission) {
			status = http.StatusBadRequest
		}
		writeError(w, status, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (api API) updateSubmission(w http.ResponseWriter, r *http.Request) {
	payload, ok := api.readSubmissionPayload(w, r)
	if !ok {
		return
	}
	item, err := api.submissions.Update(r.Context(), chi.URLParam(r, "id"), currentUser(r), payload)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (api API) transitionSubmission(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Status  workflow.Status `json:"status"`
		Comment string          `json:"comment"`
	}
	if err := readJSON(r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	item, err := api.submissions.Transition(r.Context(), chi.URLParam(r, "id"), currentUser(r), payload.Status, strings.TrimSpace(payload.Comment))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (api API) auditEvents(w http.ResponseWriter, r *http.Request) {
	events, err := api.submissions.AuditEvents(r.Context(), chi.URLParam(r, "id"), currentUser(r))
	if err != nil {
		writeError(w, http.StatusNotFound, "submission not found")
		return
	}
	writeJSON(w, http.StatusOK, events)
}

func (api API) readSubmissionPayload(w http.ResponseWriter, r *http.Request) (services.SubmissionPayload, bool) {
	var payload services.SubmissionPayload
	if err := readJSON(r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return payload, false
	}
	if strings.TrimSpace(payload.Title) == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return payload, false
	}
	if strings.TrimSpace(payload.Summary) == "" {
		writeError(w, http.StatusBadRequest, "summary is required")
		return payload, false
	}
	if len(payload.Data) == 0 || !json.Valid(payload.Data) {
		writeError(w, http.StatusBadRequest, "data must be valid JSON")
		return payload, false
	}
	return payload, true
}
