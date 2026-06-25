package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (api API) dashboardStats(c echo.Context) error {
	stats, err := api.dashboard.Stats(c.Request().Context(), currentUser(c))
	if err != nil {
		return writeError(c, http.StatusInternalServerError, err.Error())
	}
	return writeJSON(c, http.StatusOK, stats)
}
