package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"openownership-workflow/backend/internal/models"
	"openownership-workflow/backend/internal/repositories"
	"openownership-workflow/backend/internal/workflow"
)

var (
	ErrInvalidPermission = errors.New("permission name is required")
	ErrInvalidRoleName   = errors.New("role name is required")
)

type AccessService struct {
	repo *repositories.Repository
}

type CreatePermissionInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateRoleInput struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	PermissionIDs []string `json:"permissionIds"`
}

type UpdateRoleInput struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	PermissionIDs []string `json:"permissionIds"`
}

func NewAccessService(repo *repositories.Repository) *AccessService {
	return &AccessService{repo: repo}
}

func (s *AccessService) ListPermissions(ctx context.Context, actor models.User) ([]models.Permission, error) {
	if actor.Role != workflow.Admin {
		return nil, ErrAdminRequired
	}
	return s.repo.ListPermissions(ctx)
}

func (s *AccessService) CreatePermission(ctx context.Context, actor models.User, input CreatePermissionInput) (models.Permission, error) {
	if actor.Role != workflow.Admin {
		return models.Permission{}, ErrAdminRequired
	}
	name := normalizeAccessName(input.Name)
	if name == "" {
		return models.Permission{}, ErrInvalidPermission
	}
	permission, err := s.repo.CreatePermission(ctx, repositories.CreatePermissionParams{
		Name:        name,
		Description: strings.TrimSpace(input.Description),
	})
	if err != nil {
		return models.Permission{}, err
	}
	_ = s.recordAudit(ctx, actor, "permission.created", "permission", permission.ID, permission.Name, map[string]any{
		"name":        permission.Name,
		"description": permission.Description,
	})
	return permission, nil
}

func (s *AccessService) ListRoles(ctx context.Context, actor models.User) ([]models.AccessRole, error) {
	if actor.Role != workflow.Admin {
		return nil, ErrAdminRequired
	}
	return s.repo.ListRoles(ctx)
}

func (s *AccessService) CreateRole(ctx context.Context, actor models.User, input CreateRoleInput) (models.AccessRole, error) {
	if actor.Role != workflow.Admin {
		return models.AccessRole{}, ErrAdminRequired
	}
	name := normalizeAccessName(input.Name)
	if name == "" {
		return models.AccessRole{}, ErrInvalidRoleName
	}
	role, err := s.repo.CreateRole(ctx, repositories.CreateRoleParams{
		Name:          name,
		Description:   strings.TrimSpace(input.Description),
		PermissionIDs: input.PermissionIDs,
	})
	if err != nil {
		return models.AccessRole{}, err
	}
	_ = s.recordAudit(ctx, actor, "role.created", "role", role.ID, role.Name, roleAuditMetadata(role))
	return role, nil
}

func (s *AccessService) UpdateRole(ctx context.Context, actor models.User, id string, input UpdateRoleInput) (models.AccessRole, error) {
	if actor.Role != workflow.Admin {
		return models.AccessRole{}, ErrAdminRequired
	}
	name := normalizeAccessName(input.Name)
	if name == "" {
		return models.AccessRole{}, ErrInvalidRoleName
	}
	beforeRoles, err := s.repo.ListRoles(ctx)
	if err != nil {
		return models.AccessRole{}, err
	}
	var before models.AccessRole
	for _, role := range beforeRoles {
		if role.ID == id {
			before = role
			break
		}
	}
	role, err := s.repo.UpdateRole(ctx, id, repositories.UpdateRoleParams{
		Name:          name,
		Description:   strings.TrimSpace(input.Description),
		PermissionIDs: input.PermissionIDs,
	})
	if err != nil {
		return models.AccessRole{}, err
	}
	_ = s.recordAudit(ctx, actor, "role.updated", "role", role.ID, role.Name, map[string]any{
		"before": roleAuditMetadata(before),
		"after":  roleAuditMetadata(role),
	})
	return role, nil
}

func normalizeAccessName(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func (s *AccessService) recordAudit(ctx context.Context, actor models.User, eventType, resourceType, resourceID, resourceName string, metadata map[string]any) error {
	body, err := json.Marshal(metadata)
	if err != nil {
		return err
	}
	return s.repo.RecordSystemAudit(ctx, repositories.SystemAuditParams{
		ActorID:      actor.ID,
		EventType:    eventType,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		ResourceName: resourceName,
		Summary:      fmt.Sprintf("%s %s", eventType, resourceName),
		Metadata:     body,
	})
}

func roleAuditMetadata(role models.AccessRole) map[string]any {
	permissions := make([]string, 0, len(role.Permissions))
	for _, permission := range role.Permissions {
		permissions = append(permissions, permission.Name)
	}
	return map[string]any{
		"name":        role.Name,
		"description": role.Description,
		"permissions": permissions,
	}
}
