package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"openownership-workflow/backend/internal/services"
)

func (api API) listPermissions(w http.ResponseWriter, r *http.Request) {
	permissions, err := api.access.ListPermissions(r.Context(), currentUser(r))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, permissions)
}

func (api API) createPermission(w http.ResponseWriter, r *http.Request) {
	var payload services.CreatePermissionInput
	if err := readJSON(r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	permission, err := api.access.CreatePermission(r.Context(), currentUser(r), payload)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, permission)
}

func (api API) listRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := api.access.ListRoles(r.Context(), currentUser(r))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, roles)
}

func (api API) createRole(w http.ResponseWriter, r *http.Request) {
	var payload services.CreateRoleInput
	if err := readJSON(r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	role, err := api.access.CreateRole(r.Context(), currentUser(r), payload)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, role)
}

func (api API) updateRole(w http.ResponseWriter, r *http.Request) {
	var payload services.UpdateRoleInput
	if err := readJSON(r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	role, err := api.access.UpdateRole(r.Context(), currentUser(r), chi.URLParam(r, "id"), payload)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, role)
}
