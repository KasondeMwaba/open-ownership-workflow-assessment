package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (api API) listVisibleAudit(c echo.Context) error {
	events, err := api.audit.ListVisibleEvents(c.Request().Context(), currentUser(c))
	if err != nil {
		return writeError(c, http.StatusInternalServerError, err.Error())
	}
	return writeJSON(c, http.StatusOK, events)
}

func (api API) listSystemAudit(c echo.Context) error {
	events, err := api.audit.ListSystemEvents(c.Request().Context(), currentUser(c))
	if err != nil {
		return writeServiceError(c, err)
	}
	return writeJSON(c, http.StatusOK, events)
}
