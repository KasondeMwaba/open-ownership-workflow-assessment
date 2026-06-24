package repositories

import (
	"context"
	"encoding/json"
	"time"

	"openownership-workflow/backend/internal/models"
	"openownership-workflow/backend/internal/workflow"
)

type SystemAuditParams struct {
	ActorID      string
	EventType    string
	ResourceType string
	ResourceID   string
	ResourceName string
	Summary      string
	Metadata     []byte
}

type SessionAuditParams struct {
	ActorID   *string
	Email     string
	EventType string
	Success   bool
	IPAddress string
	UserAgent string
	Browser   string
	Reason    string
	Metadata  []byte
}

type ActivityAuditParams struct {
	ActorID    string
	Method     string
	Path       string
	Query      string
	StatusCode int
	Success    bool
	DurationMs int64
	IPAddress  string
	UserAgent  string
	Browser    string
	Metadata   []byte
}

type auditEventRow struct {
	ID              string
	CreatedAt       time.Time
	SubmissionID    string
	SubmissionTitle string
	ActorID         string
	ActorName       string
	ActorRole       workflow.Role
	FromStatus      *workflow.Status
	ToStatus        workflow.Status
	Comment         string
	Metadata        json.RawMessage
}

func (r *Repository) ListAuditEvents(ctx context.Context, submissionID string) ([]models.AuditEvent, error) {
	rows := []auditEventRow{}
	err := r.db.WithContext(ctx).
		Table("audit_events AS e").
		Select("e.id, e.created_at, e.submission_id, s.title AS submission_title, e.actor_id, u.name AS actor_name, u.role AS actor_role, e.from_status, e.to_status, e.comment, e.metadata").
		Joins("JOIN users u ON u.id = e.actor_id").
		Joins("JOIN submissions s ON s.id = e.submission_id").
		Where("e.submission_id = ?", submissionID).
		Order("e.created_at DESC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return auditRowsToEvents(rows), nil
}

func (r *Repository) ListVisibleSubmissionAuditEvents(ctx context.Context, user models.User) ([]models.AuditEvent, error) {
	query := r.db.WithContext(ctx).
		Table("audit_events AS e").
		Select("e.id, e.created_at, e.submission_id, s.title AS submission_title, e.actor_id, u.name AS actor_name, u.role AS actor_role, e.from_status, e.to_status, e.comment, e.metadata").
		Joins("JOIN users u ON u.id = e.actor_id").
		Joins("JOIN submissions s ON s.id = e.submission_id")
	if shouldScopeToOwner(user) {
		query = query.Where("s.owner_id = ?", user.ID)
	}
	rows := []auditEventRow{}
	err := query.Order("e.created_at DESC").Limit(200).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return auditRowsToEvents(rows), nil
}

func auditRowsToEvents(rows []auditEventRow) []models.AuditEvent {
	events := make([]models.AuditEvent, 0, len(rows))
	for _, row := range rows {
		events = append(events, models.AuditEvent{
			BaseModel:       models.BaseModel{ID: row.ID, CreatedAt: row.CreatedAt},
			SubmissionID:    row.SubmissionID,
			SubmissionTitle: row.SubmissionTitle,
			ActorID:         row.ActorID,
			ActorName:       row.ActorName,
			ActorRole:       row.ActorRole,
			FromStatus:      row.FromStatus,
			ToStatus:        row.ToStatus,
			Comment:         row.Comment,
			Metadata:        row.Metadata,
		})
	}
	return events
}

func (r *Repository) RecordSystemAudit(ctx context.Context, params SystemAuditParams) error {
	event := models.SystemAuditEvent{
		ActorID:      params.ActorID,
		EventType:    params.EventType,
		ResourceType: params.ResourceType,
		ResourceID:   params.ResourceID,
		ResourceName: params.ResourceName,
		Summary:      params.Summary,
		Metadata:     params.Metadata,
	}
	return r.db.WithContext(ctx).Create(&event).Error
}

func (r *Repository) ListSystemAuditEvents(ctx context.Context) ([]models.SystemAuditEvent, error) {
	type systemAuditRow struct {
		ID           string
		CreatedAt    time.Time
		ActorID      string
		ActorName    string
		ActorRole    workflow.Role
		EventType    string
		ResourceType string
		ResourceID   string
		ResourceName string
		Summary      string
		Metadata     json.RawMessage
	}
	rows := []systemAuditRow{}
	err := r.db.WithContext(ctx).
		Table("system_audit_events AS e").
		Select("e.id, e.created_at, e.actor_id, u.name AS actor_name, u.role AS actor_role, e.event_type, e.resource_type, e.resource_id, e.resource_name, e.summary, e.metadata").
		Joins("JOIN users u ON u.id = e.actor_id").
		Order("e.created_at DESC").
		Limit(200).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	events := make([]models.SystemAuditEvent, 0, len(rows))
	for _, row := range rows {
		events = append(events, models.SystemAuditEvent{
			BaseModel:    models.BaseModel{ID: row.ID, CreatedAt: row.CreatedAt},
			ActorID:      row.ActorID,
			ActorName:    row.ActorName,
			ActorRole:    row.ActorRole,
			EventType:    row.EventType,
			ResourceType: row.ResourceType,
			ResourceID:   row.ResourceID,
			ResourceName: row.ResourceName,
			Summary:      row.Summary,
			Metadata:     row.Metadata,
		})
	}
	return events, nil
}

func (r *Repository) RecordSessionAudit(ctx context.Context, params SessionAuditParams) error {
	event := models.SessionAuditEvent{
		ActorID:   params.ActorID,
		Email:     params.Email,
		EventType: params.EventType,
		Success:   params.Success,
		IPAddress: params.IPAddress,
		UserAgent: params.UserAgent,
		Browser:   params.Browser,
		Reason:    params.Reason,
		Metadata:  params.Metadata,
	}
	return r.db.WithContext(ctx).Create(&event).Error
}

func (r *Repository) ListSessionAuditEvents(ctx context.Context) ([]models.SessionAuditEvent, error) {
	type sessionAuditRow struct {
		ID        string
		CreatedAt time.Time
		ActorID   *string
		ActorName string
		ActorRole workflow.Role
		Email     string
		EventType string
		Success   bool
		IPAddress string
		UserAgent string
		Browser   string
		Reason    string
		Metadata  json.RawMessage
	}
	rows := []sessionAuditRow{}
	err := r.db.WithContext(ctx).
		Table("session_audit_events AS e").
		Select("e.id, e.created_at, e.actor_id, u.name AS actor_name, u.role AS actor_role, e.email, e.event_type, e.success, e.ip_address, e.user_agent, e.browser, e.reason, e.metadata").
		Joins("LEFT JOIN users u ON u.id = e.actor_id").
		Order("e.created_at DESC").
		Limit(200).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	events := make([]models.SessionAuditEvent, 0, len(rows))
	for _, row := range rows {
		events = append(events, models.SessionAuditEvent{
			BaseModel: models.BaseModel{ID: row.ID, CreatedAt: row.CreatedAt},
			ActorID:   row.ActorID,
			ActorName: row.ActorName,
			ActorRole: row.ActorRole,
			Email:     row.Email,
			EventType: row.EventType,
			Success:   row.Success,
			IPAddress: row.IPAddress,
			UserAgent: row.UserAgent,
			Browser:   row.Browser,
			Reason:    row.Reason,
			Metadata:  row.Metadata,
		})
	}
	return events, nil
}

func (r *Repository) RecordActivityAudit(ctx context.Context, params ActivityAuditParams) error {
	event := models.ActivityAuditEvent{
		ActorID:    params.ActorID,
		Method:     params.Method,
		Path:       params.Path,
		Query:      params.Query,
		StatusCode: params.StatusCode,
		Success:    params.Success,
		DurationMs: params.DurationMs,
		IPAddress:  params.IPAddress,
		UserAgent:  params.UserAgent,
		Browser:    params.Browser,
		Metadata:   params.Metadata,
	}
	return r.db.WithContext(ctx).Create(&event).Error
}

func (r *Repository) ListActivityAuditEvents(ctx context.Context) ([]models.ActivityAuditEvent, error) {
	type activityAuditRow struct {
		ID         string
		CreatedAt  time.Time
		ActorID    string
		ActorName  string
		ActorRole  workflow.Role
		Method     string
		Path       string
		Query      string
		StatusCode int
		Success    bool
		DurationMs int64
		IPAddress  string
		UserAgent  string
		Browser    string
		Metadata   json.RawMessage
	}
	rows := []activityAuditRow{}
	err := r.db.WithContext(ctx).
		Table("activity_audit_events AS e").
		Select("e.id, e.created_at, e.actor_id, u.name AS actor_name, u.role AS actor_role, e.method, e.path, e.query, e.status_code, e.success, e.duration_ms, e.ip_address, e.user_agent, e.browser, e.metadata").
		Joins("JOIN users u ON u.id = e.actor_id").
		Order("e.created_at DESC").
		Limit(300).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	events := make([]models.ActivityAuditEvent, 0, len(rows))
	for _, row := range rows {
		events = append(events, models.ActivityAuditEvent{
			BaseModel:  models.BaseModel{ID: row.ID, CreatedAt: row.CreatedAt},
			ActorID:    row.ActorID,
			ActorName:  row.ActorName,
			ActorRole:  row.ActorRole,
			Method:     row.Method,
			Path:       row.Path,
			Query:      row.Query,
			StatusCode: row.StatusCode,
			Success:    row.Success,
			DurationMs: row.DurationMs,
			IPAddress:  row.IPAddress,
			UserAgent:  row.UserAgent,
			Browser:    row.Browser,
			Metadata:   row.Metadata,
		})
	}
	return events, nil
}
