package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"openownership-workflow/backend/internal/services"
)

func (api API) listUsers(w http.ResponseWriter, r *http.Request) {
	users, err := api.users.List(r.Context(), currentUser(r))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, users)
}

func (api API) createUser(w http.ResponseWriter, r *http.Request) {
	var payload services.CreateUserInput
	if err := readJSON(r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	user, err := api.users.Create(r.Context(), currentUser(r), payload)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, user)
}

func (api API) updateUser(w http.ResponseWriter, r *http.Request) {
	var payload services.UpdateUserInput
	if err := readJSON(r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	user, err := api.users.Update(r.Context(), currentUser(r), chi.URLParam(r, "id"), payload)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (api API) setUserStatus(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		IsActive bool `json:"isActive"`
	}
	if err := readJSON(r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	user, err := api.users.SetActive(r.Context(), currentUser(r), chi.URLParam(r, "id"), payload.IsActive)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, services.ErrAdminRequired):
		writeError(w, http.StatusForbidden, err.Error())
	case errors.Is(err, services.ErrInvalidRole), errors.Is(err, services.ErrInvalidUserInput), errors.Is(err, services.ErrPasswordTooShort), errors.Is(err, services.ErrCannotDisableSelf), errors.Is(err, services.ErrInvalidPermission), errors.Is(err, services.ErrInvalidRoleName):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, gorm.ErrRecordNotFound):
		writeError(w, http.StatusNotFound, "user not found")
	default:
		writeError(w, http.StatusInternalServerError, err.Error())
	}
}
