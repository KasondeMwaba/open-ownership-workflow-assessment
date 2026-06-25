package dto

import "openownership-workflow/backend/internal/workflow"

type CreateUserRequest struct {
	Name     string        `json:"name"`
	Email    string        `json:"email"`
	Password string        `json:"password"`
	Role     workflow.Role `json:"role"`
	IsActive bool          `json:"isActive"`
}

type UpdateUserRequest struct {
	Name     string        `json:"name"`
	Email    string        `json:"email"`
	Role     workflow.Role `json:"role"`
	IsActive bool          `json:"isActive"`
}

type SetUserStatusRequest struct {
	IsActive bool `json:"isActive"`
}
