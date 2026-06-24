package models

import (
	"encoding/json"

	"openownership-workflow/backend/internal/workflow"
)

type Submission struct {
	UpdatableModel
	Title      string          `json:"title" gorm:"not null"`
	Summary    string          `json:"summary" gorm:"not null"`
	Data       json.RawMessage `json:"data" gorm:"type:jsonb;not null"`
	Status     workflow.Status `json:"status" gorm:"not null;default:draft;index;check:status IN ('draft','submitted','changes_required','approved','rejected','withdrawn')"`
	OwnerID    string          `json:"ownerId" gorm:"column:owner_id;type:uuid;not null;index"`
	OwnerName  string          `json:"ownerName" gorm:"-"`
	ReviewerID *string         `json:"reviewerId,omitempty" gorm:"column:reviewer_id;type:uuid"`
	Version    int             `json:"version" gorm:"not null;default:1"`
	Owner      User            `json:"-" gorm:"foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Reviewer   *User           `json:"-" gorm:"foreignKey:ReviewerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (Submission) TableName() string {
	return "submissions"
}
