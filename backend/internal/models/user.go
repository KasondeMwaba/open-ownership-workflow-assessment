package models

import "openownership-workflow/backend/internal/workflow"

type User struct {
	BaseModel
	Name         string        `json:"name" gorm:"not null"`
	Email        string        `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string        `json:"-" gorm:"column:password_hash"`
	Role         workflow.Role `json:"role" gorm:"not null;index"`
	IsActive     bool          `json:"isActive" gorm:"column:is_active;not null;default:true;index"`
	Permissions  []string      `json:"permissions" gorm:"-"`
}

func (User) TableName() string {
	return "users"
}

func (user User) HasPermission(permission string) bool {
	if user.Role == workflow.Admin {
		return true
	}
	for _, item := range user.Permissions {
		if item == permission {
			return true
		}
	}
	return false
}
