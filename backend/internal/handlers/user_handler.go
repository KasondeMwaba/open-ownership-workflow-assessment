package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"openownership-workflow/backend/internal/dto"
	"openownership-workflow/backend/internal/services"
)

func (api API) listUsers(c echo.Context) error {
	users, err := api.users.List(c.Request().Context(), currentUser(c))
	if err != nil {
		return writeServiceError(c, err)
	}
	return writeJSON(c, http.StatusOK, users)
}

func (api API) createUser(c echo.Context) error {
	var payload dto.CreateUserRequest
	if err := readJSON(c, &payload); err != nil {
		return writeError(c, http.StatusBadRequest, err.Error())
	}
	user, err := api.users.Create(c.Request().Context(), currentUser(c), services.CreateUserInput{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: payload.Password,
		Role:     payload.Role,
		IsActive: payload.IsActive,
	})
	if err != nil {
		return writeServiceError(c, err)
	}
	return writeJSON(c, http.StatusCreated, user)
}

func (api API) updateUser(c echo.Context) error {
	var payload dto.UpdateUserRequest
	if err := readJSON(c, &payload); err != nil {
		return writeError(c, http.StatusBadRequest, err.Error())
	}
	user, err := api.users.Update(c.Request().Context(), currentUser(c), c.Param("id"), services.UpdateUserInput{
		Name:     payload.Name,
		Email:    payload.Email,
		Role:     payload.Role,
		IsActive: payload.IsActive,
	})
	if err != nil {
		return writeServiceError(c, err)
	}
	return writeJSON(c, http.StatusOK, user)
}

func (api API) setUserStatus(c echo.Context) error {
	var payload dto.SetUserStatusRequest
	if err := readJSON(c, &payload); err != nil {
		return writeError(c, http.StatusBadRequest, err.Error())
	}
	user, err := api.users.SetActive(c.Request().Context(), currentUser(c), c.Param("id"), payload.IsActive)
	if err != nil {
		return writeServiceError(c, err)
	}
	return writeJSON(c, http.StatusOK, user)
}

func writeServiceError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, services.ErrAdminRequired):
		return writeError(c, http.StatusForbidden, err.Error())
	case errors.Is(err, services.ErrInvalidRole), errors.Is(err, services.ErrInvalidUserInput), errors.Is(err, services.ErrPasswordTooShort), errors.Is(err, services.ErrCannotDisableSelf), errors.Is(err, services.ErrInvalidPermission), errors.Is(err, services.ErrInvalidRoleName):
		return writeError(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, gorm.ErrRecordNotFound):
		return writeError(c, http.StatusNotFound, "user not found")
	default:
		return writeError(c, http.StatusInternalServerError, err.Error())
	}
}
