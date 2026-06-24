package models

import (
	"encoding/json"

	"openownership-workflow/backend/internal/workflow"
)

type AuditEvent struct {
	BaseModel
	SubmissionID    string           `json:"submissionId" gorm:"column:submission_id;type:uuid;not null;index:idx_audit_events_submission_created,priority:1"`
	SubmissionTitle string           `json:"submissionTitle" gorm:"-"`
	ActorID         string           `json:"actorId" gorm:"column:actor_id;type:uuid;not null"`
	ActorName       string           `json:"actorName" gorm:"-"`
	ActorRole       workflow.Role    `json:"actorRole" gorm:"-"`
	FromStatus      *workflow.Status `json:"fromStatus,omitempty" gorm:"column:from_status"`
	ToStatus        workflow.Status  `json:"toStatus" gorm:"column:to_status;not null"`
	Comment         string           `json:"comment" gorm:"not null;default:''"`
	Metadata        json.RawMessage  `json:"metadata" gorm:"type:jsonb;not null"`
	Submission      Submission       `json:"-" gorm:"foreignKey:SubmissionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Actor           User             `json:"-" gorm:"foreignKey:ActorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (AuditEvent) TableName() string {
	return "audit_events"
}

type SystemAuditEvent struct {
	BaseModel
	ActorID      string          `json:"actorId" gorm:"column:actor_id;type:uuid;not null;index"`
	ActorName    string          `json:"actorName" gorm:"-"`
	ActorRole    workflow.Role   `json:"actorRole" gorm:"-"`
	EventType    string          `json:"eventType" gorm:"column:event_type;not null;index"`
	ResourceType string          `json:"resourceType" gorm:"column:resource_type;not null;index"`
	ResourceID   string          `json:"resourceId" gorm:"column:resource_id;type:uuid;not null;index"`
	ResourceName string          `json:"resourceName" gorm:"column:resource_name;not null;default:''"`
	Summary      string          `json:"summary" gorm:"not null;default:''"`
	Metadata     json.RawMessage `json:"metadata" gorm:"type:jsonb;not null"`
	Actor        User            `json:"-" gorm:"foreignKey:ActorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (SystemAuditEvent) TableName() string {
	return "system_audit_events"
}

type SessionAuditEvent struct {
	BaseModel
	ActorID   *string         `json:"actorId,omitempty" gorm:"column:actor_id;type:uuid;index"`
	ActorName string          `json:"actorName" gorm:"-"`
	ActorRole workflow.Role   `json:"actorRole" gorm:"-"`
	Email     string          `json:"email" gorm:"not null;default:'';index"`
	EventType string          `json:"eventType" gorm:"column:event_type;not null;index"`
	Success   bool            `json:"success" gorm:"not null;index"`
	IPAddress string          `json:"ipAddress" gorm:"column:ip_address;not null;default:''"`
	UserAgent string          `json:"userAgent" gorm:"column:user_agent;not null;default:''"`
	Browser   string          `json:"browser" gorm:"not null;default:''"`
	Reason    string          `json:"reason" gorm:"not null;default:''"`
	Metadata  json.RawMessage `json:"metadata" gorm:"type:jsonb;not null"`
	Actor     *User           `json:"-" gorm:"foreignKey:ActorID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (SessionAuditEvent) TableName() string {
	return "session_audit_events"
}

type ActivityAuditEvent struct {
	BaseModel
	ActorID    string          `json:"actorId" gorm:"column:actor_id;type:uuid;not null;index"`
	ActorName  string          `json:"actorName" gorm:"-"`
	ActorRole  workflow.Role   `json:"actorRole" gorm:"-"`
	Method     string          `json:"method" gorm:"not null;index"`
	Path       string          `json:"path" gorm:"not null;index"`
	Query      string          `json:"query" gorm:"not null;default:''"`
	StatusCode int             `json:"statusCode" gorm:"column:status_code;not null;index"`
	Success    bool            `json:"success" gorm:"not null;index"`
	DurationMs int64           `json:"durationMs" gorm:"column:duration_ms;not null;default:0"`
	IPAddress  string          `json:"ipAddress" gorm:"column:ip_address;not null;default:''"`
	UserAgent  string          `json:"userAgent" gorm:"column:user_agent;not null;default:''"`
	Browser    string          `json:"browser" gorm:"not null;default:''"`
	Metadata   json.RawMessage `json:"metadata" gorm:"type:jsonb;not null"`
	Actor      User            `json:"-" gorm:"foreignKey:ActorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (ActivityAuditEvent) TableName() string {
	return "activity_audit_events"
}
