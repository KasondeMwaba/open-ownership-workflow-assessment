package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"openownership-workflow/backend/internal/dto"
	"openownership-workflow/backend/internal/services"
)

func (api API) listPermissions(c echo.Context) error {
	permissions, err := api.access.ListPermissions(c.Request().Context(), currentUser(c))
	if err != nil {
		return writeServiceError(c, err)
	}
	return writeJSON(c, http.StatusOK, permissions)
}

func (api API) createPermission(c echo.Context) error {
	var payload dto.CreatePermissionRequest
	if err := readJSON(c, &payload); err != nil {
		return writeError(c, http.StatusBadRequest, err.Error())
	}
	permission, err := api.access.CreatePermission(c.Request().Context(), currentUser(c), services.CreatePermissionInput{
		Name:        payload.Name,
		Description: payload.Description,
	})
	if err != nil {
		return writeServiceError(c, err)
	}
	return writeJSON(c, http.StatusCreated, permission)
}

func (api API) listRoles(c echo.Context) error {
	roles, err := api.access.ListRoles(c.Request().Context(), currentUser(c))
	if err != nil {
		return writeServiceError(c, err)
	}
	return writeJSON(c, http.StatusOK, roles)
}

func (api API) createRole(c echo.Context) error {
	var payload dto.CreateRoleRequest
	if err := readJSON(c, &payload); err != nil {
		return writeError(c, http.StatusBadRequest, err.Error())
	}
	role, err := api.access.CreateRole(c.Request().Context(), currentUser(c), services.CreateRoleInput{
		Name:          payload.Name,
		Description:   payload.Description,
		PermissionIDs: payload.PermissionIDs,
	})
	if err != nil {
		return writeServiceError(c, err)
	}
	return writeJSON(c, http.StatusCreated, role)
}

func (api API) updateRole(c echo.Context) error {
	var payload dto.UpdateRoleRequest
	if err := readJSON(c, &payload); err != nil {
		return writeError(c, http.StatusBadRequest, err.Error())
	}
	role, err := api.access.UpdateRole(c.Request().Context(), currentUser(c), c.Param("id"), services.UpdateRoleInput{
		Name:          payload.Name,
		Description:   payload.Description,
		PermissionIDs: payload.PermissionIDs,
	})
	if err != nil {
		return writeServiceError(c, err)
	}
	return writeJSON(c, http.StatusOK, role)
}
