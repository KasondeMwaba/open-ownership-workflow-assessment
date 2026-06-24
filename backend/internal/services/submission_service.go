package services

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"openownership-workflow/backend/internal/models"
	"openownership-workflow/backend/internal/repositories"
	"openownership-workflow/backend/internal/workflow"
)

var ErrInvalidSubmission = errors.New("submission title, summary, company, jurisdiction, registration number, and valid beneficial owners are required")

type SubmissionService struct {
	repo      *repositories.Repository
	dashboard *DashboardService
}

type SubmissionPayload struct {
	Title   string          `json:"title"`
	Summary string          `json:"summary"`
	Data    json.RawMessage `json:"data"`
}

type submissionData struct {
	Company            string                `json:"company"`
	Jurisdiction       string                `json:"jurisdiction"`
	RegistrationNumber string                `json:"registrationNumber"`
	BeneficialOwners   []beneficialOwnerData `json:"beneficialOwners"`
}

type beneficialOwnerData struct {
	Name             string  `json:"name"`
	OwnershipPercent float64 `json:"ownershipPercent"`
	ControlType      string  `json:"controlType"`
}

func NewSubmissionService(repo *repositories.Repository, dashboard *DashboardService) *SubmissionService {
	return &SubmissionService{repo: repo, dashboard: dashboard}
}

func (s *SubmissionService) List(ctx context.Context, user models.User, status string) ([]models.Submission, error) {
	return s.repo.ListSubmissions(ctx, user, status)
}

func (s *SubmissionService) ListPage(ctx context.Context, user models.User, status string, page, pageSize int) (repositories.SubmissionPage, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return s.repo.ListSubmissionsPage(ctx, user, status, page, pageSize)
}

func (s *SubmissionService) Get(ctx context.Context, id string, user models.User) (models.Submission, error) {
	return s.repo.GetSubmission(ctx, id, user)
}

func (s *SubmissionService) Create(ctx context.Context, user models.User, payload SubmissionPayload) (models.Submission, error) {
	payload, err := validateSubmissionPayload(payload)
	if err != nil {
		return models.Submission{}, err
	}
	item, err := s.repo.CreateSubmission(ctx, repositories.CreateSubmissionParams{
		Title: payload.Title, Summary: payload.Summary, Data: payload.Data, OwnerID: user.ID,
	})
	if err == nil {
		s.dashboard.Invalidate(ctx)
	}
	return item, err
}

func (s *SubmissionService) Update(ctx context.Context, id string, user models.User, payload SubmissionPayload) (models.Submission, error) {
	payload, err := validateSubmissionPayload(payload)
	if err != nil {
		return models.Submission{}, err
	}
	item, err := s.repo.UpdateSubmission(ctx, repositories.UpdateSubmissionParams{
		ID: id, Title: payload.Title, Summary: payload.Summary, Data: payload.Data, ActorID: user.ID,
	})
	if err == nil {
		s.dashboard.Invalidate(ctx)
	}
	return item, err
}

func validateSubmissionPayload(payload SubmissionPayload) (SubmissionPayload, error) {
	payload.Title = strings.TrimSpace(payload.Title)
	payload.Summary = strings.TrimSpace(payload.Summary)
	if payload.Title == "" || payload.Summary == "" || len(payload.Data) == 0 || !json.Valid(payload.Data) {
		return payload, ErrInvalidSubmission
	}

	var data submissionData
	if err := json.Unmarshal(payload.Data, &data); err != nil {
		return payload, ErrInvalidSubmission
	}
	if strings.TrimSpace(data.Company) == "" || strings.TrimSpace(data.Jurisdiction) == "" || strings.TrimSpace(data.RegistrationNumber) == "" || len(data.BeneficialOwners) == 0 {
		return payload, ErrInvalidSubmission
	}
	for _, owner := range data.BeneficialOwners {
		if strings.TrimSpace(owner.Name) == "" || strings.TrimSpace(owner.ControlType) == "" || owner.OwnershipPercent < 0 || owner.OwnershipPercent > 100 {
			return payload, ErrInvalidSubmission
		}
	}
	return payload, nil
}

func (s *SubmissionService) Transition(ctx context.Context, id string, user models.User, status workflow.Status, comment string) (models.Submission, error) {
	item, err := s.repo.TransitionSubmission(ctx, repositories.TransitionParams{
		ID: id, To: status, Comment: comment, Actor: user,
	})
	if err == nil {
		s.dashboard.Invalidate(ctx)
	}
	return item, err
}

func (s *SubmissionService) AuditEvents(ctx context.Context, id string, user models.User) ([]models.AuditEvent, error) {
	if _, err := s.repo.GetSubmission(ctx, id, user); err != nil {
		return nil, err
	}
	return s.repo.ListAuditEvents(ctx, id)
}
