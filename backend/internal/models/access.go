package models

type Permission struct {
	BaseModel
	Name        string `json:"name" gorm:"uniqueIndex;not null"`
	Description string `json:"description" gorm:"not null;default:''"`
}

func (Permission) TableName() string {
	return "permissions"
}

type AccessRole struct {
	BaseModel
	Name        string       `json:"name" gorm:"uniqueIndex;not null"`
	Description string       `json:"description" gorm:"not null;default:''"`
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;constraint:OnDelete:CASCADE;"`
}

func (AccessRole) TableName() string {
	return "roles"
}
