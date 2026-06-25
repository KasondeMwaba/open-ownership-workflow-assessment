package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"openownership-workflow/backend/internal/dto"
	"openownership-workflow/backend/internal/models"
	"openownership-workflow/backend/internal/services"
)

func (api API) listSubmissions(c echo.Context) error {
	r := c.Request()
	query := r.URL.Query()
	if query.Get("page") != "" || query.Get("pageSize") != "" {
		page := parsePositiveInt(query.Get("page"), 1)
		pageSize := parsePositiveInt(query.Get("pageSize"), 10)
		result, err := api.submissions.ListPage(r.Context(), currentUser(c), query.Get("status"), page, pageSize)
		if err != nil {
			return writeError(c, http.StatusInternalServerError, err.Error())
		}
		return writeJSON(c, http.StatusOK, dto.PaginatedResponse[models.Submission]{
			Items:    result.Items,
			Total:    result.Total,
			Page:     result.Page,
			PageSize: result.PageSize,
		})
	}
	items, err := api.submissions.List(r.Context(), currentUser(c), query.Get("status"))
	if err != nil {
		return writeError(c, http.StatusInternalServerError, err.Error())
	}
	return writeJSON(c, http.StatusOK, items)
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

func (api API) getSubmission(c echo.Context) error {
	item, err := api.submissions.Get(c.Request().Context(), c.Param("id"), currentUser(c))
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = http.StatusNotFound
		}
		return writeError(c, status, "submission not found")
	}
	return writeJSON(c, http.StatusOK, item)
}

func (api API) createSubmission(c echo.Context) error {
	user := currentUser(c)
	if !user.HasPermission("submissions:create") {
		return writeError(c, http.StatusForbidden, "missing submissions:create permission")
	}
	payload, ok, err := api.readSubmissionPayload(c)
	if !ok {
		return err
	}
	item, err := api.submissions.Create(c.Request().Context(), user, payload)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, services.ErrInvalidSubmission) {
			status = http.StatusBadRequest
		}
		return writeError(c, status, err.Error())
	}
	return writeJSON(c, http.StatusCreated, item)
}

func (api API) updateSubmission(c echo.Context) error {
	payload, ok, err := api.readSubmissionPayload(c)
	if !ok {
		return err
	}
	item, err := api.submissions.Update(c.Request().Context(), c.Param("id"), currentUser(c), payload)
	if err != nil {
		return writeError(c, http.StatusBadRequest, err.Error())
	}
	return writeJSON(c, http.StatusOK, item)
}

func (api API) transitionSubmission(c echo.Context) error {
	var payload dto.TransitionSubmissionRequest
	if err := readJSON(c, &payload); err != nil {
		return writeError(c, http.StatusBadRequest, err.Error())
	}
	item, err := api.submissions.Transition(c.Request().Context(), c.Param("id"), currentUser(c), payload.Status, strings.TrimSpace(payload.Comment))
	if err != nil {
		return writeError(c, http.StatusBadRequest, err.Error())
	}
	return writeJSON(c, http.StatusOK, item)
}

func (api API) auditEvents(c echo.Context) error {
	events, err := api.submissions.AuditEvents(c.Request().Context(), c.Param("id"), currentUser(c))
	if err != nil {
		return writeError(c, http.StatusNotFound, "submission not found")
	}
	return writeJSON(c, http.StatusOK, events)
}

func (api API) readSubmissionPayload(c echo.Context) (services.SubmissionPayload, bool, error) {
	var payload dto.SubmissionRequest
	if err := readJSON(c, &payload); err != nil {
		return services.SubmissionPayload{}, false, writeError(c, http.StatusBadRequest, err.Error())
	}
	if strings.TrimSpace(payload.Title) == "" {
		return services.SubmissionPayload{}, false, writeError(c, http.StatusBadRequest, "title is required")
	}
	if strings.TrimSpace(payload.Summary) == "" {
		return services.SubmissionPayload{}, false, writeError(c, http.StatusBadRequest, "summary is required")
	}
	if len(payload.Data) == 0 || !json.Valid(payload.Data) {
		return services.SubmissionPayload{}, false, writeError(c, http.StatusBadRequest, "data must be valid JSON")
	}
	return services.SubmissionPayload{
		Title:   payload.Title,
		Summary: payload.Summary,
		Data:    payload.Data,
	}, true, nil
}
