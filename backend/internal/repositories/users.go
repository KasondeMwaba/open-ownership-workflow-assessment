package repositories

import (
	"context"

	"openownership-workflow/backend/internal/models"
	"openownership-workflow/backend/internal/workflow"

	"gorm.io/gorm"
)

type CreateUserParams struct {
	Name         string
	Email        string
	PasswordHash string
	Role         workflow.Role
	IsActive     bool
}

type UpdateUserParams struct {
	Name     string
	Email    string
	Role     workflow.Role
	IsActive bool
}

func (r *Repository) FindUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Where("lower(email) = lower(?)", email).
		First(&user).Error
	return user, err
}

func (r *Repository) ListUsers(ctx context.Context) ([]models.User, error) {
	users := []models.User{}
	err := r.db.WithContext(ctx).Order("created_at DESC").Find(&users).Error
	return users, err
}

func (r *Repository) CreateUser(ctx context.Context, params CreateUserParams) (models.User, error) {
	user := models.User{
		Name:         params.Name,
		Email:        params.Email,
		PasswordHash: params.PasswordHash,
		Role:         params.Role,
		IsActive:     params.IsActive,
	}
	err := r.db.WithContext(ctx).Create(&user).Error
	return user, err
}

func (r *Repository) UpdateUser(ctx context.Context, id string, params UpdateUserParams) (models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&user, "id = ?", id).Error; err != nil {
			return err
		}
		user.Name = params.Name
		user.Email = params.Email
		user.Role = params.Role
		user.IsActive = params.IsActive
		return tx.Save(&user).Error
	})
	return user, err
}

func (r *Repository) SetUserActive(ctx context.Context, id string, active bool) (models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&user, "id = ?", id).Error; err != nil {
			return err
		}
		user.IsActive = active
		return tx.Save(&user).Error
	})
	return user, err
}

func (r *Repository) FindUserByID(ctx context.Context, id string) (models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
	return user, err
}
