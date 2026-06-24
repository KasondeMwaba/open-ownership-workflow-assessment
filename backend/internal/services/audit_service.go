package services

import (
	"context"
	"encoding/json"

	"openownership-workflow/backend/internal/models"
	"openownership-workflow/backend/internal/repositories"
	"openownership-workflow/backend/internal/workflow"
)

type AuditService struct {
	repo *repositories.Repository
}

type AuditTrailResult struct {
	SubmissionEvents []models.AuditEvent         `json:"submissionEvents"`
	SystemEvents     []models.SystemAuditEvent   `json:"systemEvents"`
	SessionEvents    []models.SessionAuditEvent  `json:"sessionEvents"`
	ActivityEvents   []models.ActivityAuditEvent `json:"activityEvents"`
}

type SessionAuditInput struct {
	ActorID   *string
	Email     string
	EventType string
	Success   bool
	IPAddress string
	UserAgent string
	Browser   string
	Reason    string
	Metadata  map[string]any
}

type ActivityAuditInput struct {
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
	Metadata   map[string]any
}

func NewAuditService(repo *repositories.Repository) *AuditService {
	return &AuditService{repo: repo}
}

func (s *AuditService) ListSystemEvents(ctx context.Context, actor models.User) ([]models.SystemAuditEvent, error) {
	if actor.Role != workflow.Admin {
		return nil, ErrAdminRequired
	}
	return s.repo.ListSystemAuditEvents(ctx)
}

func (s *AuditService) RecordSessionEvent(ctx context.Context, input SessionAuditInput) {
	metadata := []byte(`{}`)
	if input.Metadata != nil {
		if payload, err := json.Marshal(input.Metadata); err == nil {
			metadata = payload
		}
	}
	_ = s.repo.RecordSessionAudit(ctx, repositories.SessionAuditParams{
		ActorID:   input.ActorID,
		Email:     input.Email,
		EventType: input.EventType,
		Success:   input.Success,
		IPAddress: input.IPAddress,
		UserAgent: input.UserAgent,
		Browser:   input.Browser,
		Reason:    input.Reason,
		Metadata:  metadata,
	})
}

func (s *AuditService) RecordActivityEvent(ctx context.Context, input ActivityAuditInput) {
	metadata := []byte(`{}`)
	if input.Metadata != nil {
		if payload, err := json.Marshal(input.Metadata); err == nil {
			metadata = payload
		}
	}
	_ = s.repo.RecordActivityAudit(ctx, repositories.ActivityAuditParams{
		ActorID:    input.ActorID,
		Method:     input.Method,
		Path:       input.Path,
		Query:      input.Query,
		StatusCode: input.StatusCode,
		Success:    input.Success,
		DurationMs: input.DurationMs,
		IPAddress:  input.IPAddress,
		UserAgent:  input.UserAgent,
		Browser:    input.Browser,
		Metadata:   metadata,
	})
}

func (s *AuditService) ListVisibleEvents(ctx context.Context, actor models.User) (AuditTrailResult, error) {
	submissionEvents, err := s.repo.ListVisibleSubmissionAuditEvents(ctx, actor)
	if err != nil {
		return AuditTrailResult{}, err
	}
	result := AuditTrailResult{SubmissionEvents: submissionEvents, SystemEvents: []models.SystemAuditEvent{}, SessionEvents: []models.SessionAuditEvent{}, ActivityEvents: []models.ActivityAuditEvent{}}
	if actor.Role == workflow.Admin {
		systemEvents, err := s.repo.ListSystemAuditEvents(ctx)
		if err != nil {
			return AuditTrailResult{}, err
		}
		sessionEvents, err := s.repo.ListSessionAuditEvents(ctx)
		if err != nil {
			return AuditTrailResult{}, err
		}
		activityEvents, err := s.repo.ListActivityAuditEvents(ctx)
		if err != nil {
			return AuditTrailResult{}, err
		}
		result.SystemEvents = systemEvents
		result.SessionEvents = sessionEvents
		result.ActivityEvents = activityEvents
	}
	return result, nil
}
