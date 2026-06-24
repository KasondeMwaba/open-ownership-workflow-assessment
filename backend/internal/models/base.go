package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        string    `json:"id" gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
}

type UpdatableModel struct {
	BaseModel
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

func (base *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if base.ID == "" {
		base.ID = uuid.NewString()
	}
	return nil
}
