package dto

type CreatePermissionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateRoleRequest struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	PermissionIDs []string `json:"permissionIds"`
}

type UpdateRoleRequest = CreateRoleRequest
