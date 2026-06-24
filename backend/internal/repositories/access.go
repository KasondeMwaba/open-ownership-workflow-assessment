package repositories

import (
	"context"

	"openownership-workflow/backend/internal/models"

	"gorm.io/gorm"
)

type CreatePermissionParams struct {
	Name        string
	Description string
}

type CreateRoleParams struct {
	Name          string
	Description   string
	PermissionIDs []string
}

type UpdateRoleParams struct {
	Name          string
	Description   string
	PermissionIDs []string
}

func (r *Repository) ListPermissions(ctx context.Context) ([]models.Permission, error) {
	permissions := []models.Permission{}
	err := r.db.WithContext(ctx).Order("name ASC").Find(&permissions).Error
	return permissions, err
}

func (r *Repository) CreatePermission(ctx context.Context, params CreatePermissionParams) (models.Permission, error) {
	permission := models.Permission{Name: params.Name, Description: params.Description}
	err := r.db.WithContext(ctx).Create(&permission).Error
	return permission, err
}

func (r *Repository) ListRoles(ctx context.Context) ([]models.AccessRole, error) {
	roles := []models.AccessRole{}
	err := r.db.WithContext(ctx).Preload("Permissions").Order("name ASC").Find(&roles).Error
	return roles, err
}

func (r *Repository) CreateRole(ctx context.Context, params CreateRoleParams) (models.AccessRole, error) {
	var role models.AccessRole
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		permissions, err := findPermissionsByID(tx, params.PermissionIDs)
		if err != nil {
			return err
		}
		role = models.AccessRole{Name: params.Name, Description: params.Description}
		if err := tx.Create(&role).Error; err != nil {
			return err
		}
		if err := tx.Model(&role).Association("Permissions").Replace(permissions); err != nil {
			return err
		}
		return tx.Preload("Permissions").First(&role, "id = ?", role.ID).Error
	})
	return role, err
}

func (r *Repository) UpdateRole(ctx context.Context, id string, params UpdateRoleParams) (models.AccessRole, error) {
	var role models.AccessRole
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&role, "id = ?", id).Error; err != nil {
			return err
		}
		permissions, err := findPermissionsByID(tx, params.PermissionIDs)
		if err != nil {
			return err
		}
		role.Name = params.Name
		role.Description = params.Description
		if err := tx.Save(&role).Error; err != nil {
			return err
		}
		if err := tx.Model(&role).Association("Permissions").Replace(permissions); err != nil {
			return err
		}
		return tx.Preload("Permissions").First(&role, "id = ?", role.ID).Error
	})
	return role, err
}

func (r *Repository) RoleExists(ctx context.Context, name string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.AccessRole{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

func (r *Repository) PermissionsForRole(ctx context.Context, roleName string) ([]string, error) {
	var role models.AccessRole
	if err := r.db.WithContext(ctx).Preload("Permissions").Where("name = ?", roleName).First(&role).Error; err != nil {
		return nil, err
	}
	permissions := make([]string, 0, len(role.Permissions))
	for _, permission := range role.Permissions {
		permissions = append(permissions, permission.Name)
	}
	return permissions, nil
}

func (r *Repository) AttachUserPermissions(ctx context.Context, user models.User) (models.User, error) {
	permissions, err := r.PermissionsForRole(ctx, string(user.Role))
	if err != nil {
		return user, err
	}
	user.Permissions = permissions
	return user, nil
}

func findPermissionsByID(tx *gorm.DB, ids []string) ([]models.Permission, error) {
	if len(ids) == 0 {
		return []models.Permission{}, nil
	}
	permissions := []models.Permission{}
	if err := tx.Where("id IN ?", ids).Find(&permissions).Error; err != nil {
		return nil, err
	}
	if len(permissions) != len(ids) {
		return nil, gorm.ErrRecordNotFound
	}
	return permissions, nil
}
