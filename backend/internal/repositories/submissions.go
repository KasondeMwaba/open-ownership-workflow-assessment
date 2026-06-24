package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"openownership-workflow/backend/internal/models"
	"openownership-workflow/backend/internal/workflow"
)

type CreateSubmissionParams struct {
	Title   string
	Summary string
	Data    json.RawMessage
	OwnerID string
}

type UpdateSubmissionParams struct {
	ID      string
	Title   string
	Summary string
	Data    json.RawMessage
	ActorID string
}

type TransitionParams struct {
	ID      string
	To      workflow.Status
	Comment string
	Actor   models.User
}

type SubmissionPage struct {
	Items    []models.Submission
	Total    int64
	Page     int
	PageSize int
}

func (r *Repository) ListSubmissions(ctx context.Context, user models.User, status string) ([]models.Submission, error) {
	page, err := r.ListSubmissionsPage(ctx, user, status, 1, 200)
	return page.Items, err
}

func (r *Repository) ListSubmissionsPage(ctx context.Context, user models.User, status string, page, pageSize int) (SubmissionPage, error) {
	query := r.db.WithContext(ctx).
		Table("submissions AS s").
		Select("s.id, s.created_at, s.updated_at, s.title, s.summary, s.data, s.status, s.owner_id, u.name AS owner_name, s.reviewer_id, s.version").
		Joins("JOIN users u ON u.id = s.owner_id")
	if shouldScopeToOwner(user) {
		query = query.Where("s.owner_id = ?", user.ID)
	}
	if status != "" {
		query = query.Where("s.status = ?", status)
	}

	var total int64
	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return SubmissionPage{}, err
	}
	submissions := []models.Submission{}
	err := query.Order("s.updated_at DESC").Limit(pageSize).Offset((page - 1) * pageSize).Scan(&submissions).Error
	return SubmissionPage{Items: submissions, Total: total, Page: page, PageSize: pageSize}, err
}

func (r *Repository) GetSubmission(ctx context.Context, id string, user models.User) (models.Submission, error) {
	var item models.Submission
	err := r.db.WithContext(ctx).
		Table("submissions AS s").
		Select("s.id, s.created_at, s.updated_at, s.title, s.summary, s.data, s.status, s.owner_id, u.name AS owner_name, s.reviewer_id, s.version").
		Joins("JOIN users u ON u.id = s.owner_id").
		Where("s.id = ?", id).
		Scan(&item).Error
	if err != nil {
		return item, err
	}
	if item.ID == "" {
		return item, gorm.ErrRecordNotFound
	}
	if shouldScopeToOwner(user) && item.OwnerID != user.ID {
		return item, gorm.ErrRecordNotFound
	}
	return item, nil
}

func (r *Repository) CreateSubmission(ctx context.Context, params CreateSubmissionParams) (models.Submission, error) {
	var item models.Submission
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		item = models.Submission{
			Title:   params.Title,
			Summary: params.Summary,
			Data:    params.Data,
			OwnerID: params.OwnerID,
			Status:  workflow.Draft,
			Version: 1,
		}
		if err := tx.Create(&item).Error; err != nil {
			return err
		}
		event := models.AuditEvent{
			SubmissionID: item.ID,
			ActorID:      params.OwnerID,
			ToStatus:     item.Status,
			Comment:      "Submission created",
			Metadata:     json.RawMessage(`{}`),
		}
		return tx.Create(&event).Error
	})
	if err != nil {
		return models.Submission{}, err
	}
	item.OwnerName = ""
	return item, nil
}

func (r *Repository) UpdateSubmission(ctx context.Context, params UpdateSubmissionParams) (models.Submission, error) {
	var item models.Submission
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&item, "id = ?", params.ID).Error; err != nil {
			return err
		}
		current := item.Status
		if item.OwnerID != params.ActorID {
			return errors.New("only the owner can edit a submission")
		}
		if current != workflow.Draft && current != workflow.ChangesRequired {
			return errors.New("only draft or changes-required submissions can be edited")
		}
		item.Title = params.Title
		item.Summary = params.Summary
		item.Data = params.Data
		item.Version++
		if err := tx.Save(&item).Error; err != nil {
			return err
		}
		event := models.AuditEvent{
			SubmissionID: params.ID,
			ActorID:      params.ActorID,
			FromStatus:   &current,
			ToStatus:     current,
			Comment:      "Submission edited",
			Metadata:     json.RawMessage(fmt.Sprintf(`{"version":%d}`, item.Version)),
		}
		return tx.Create(&event).Error
	})
	if err != nil {
		return models.Submission{}, err
	}
	return item, nil
}

func (r *Repository) TransitionSubmission(ctx context.Context, params TransitionParams) (models.Submission, error) {
	var item models.Submission
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&item, "id = ?", params.ID).Error; err != nil {
			return err
		}
		from := item.Status
		if shouldScopeToOwner(params.Actor) && item.OwnerID != params.Actor.ID {
			return errors.New("submission creators can only transition their own submissions")
		}
		if err := validateTransitionForUser(from, params.To, params.Actor); err != nil {
			return err
		}
		item.Status = params.To
		if params.Actor.HasPermission("submissions:review") || params.Actor.Role == workflow.Admin {
			item.ReviewerID = &params.Actor.ID
		}
		if err := tx.Save(&item).Error; err != nil {
			return err
		}
		event := models.AuditEvent{
			SubmissionID: params.ID,
			ActorID:      params.Actor.ID,
			FromStatus:   &from,
			ToStatus:     params.To,
			Comment:      params.Comment,
			Metadata:     json.RawMessage(`{}`),
		}
		return tx.Create(&event).Error
	})
	if err != nil {
		return models.Submission{}, err
	}
	return item, nil
}

func shouldScopeToOwner(user models.User) bool {
	return user.HasPermission("submissions:create") && !user.HasPermission("submissions:review") && user.Role != workflow.Admin
}

func validateTransitionForUser(from, to workflow.Status, actor models.User) error {
	if (from == workflow.Draft || from == workflow.ChangesRequired) && (to == workflow.Submitted || to == workflow.Withdrawn) {
		if actor.HasPermission("submissions:create") {
			return nil
		}
		return workflow.ErrTransitionNotAllowed
	}
	if from == workflow.Submitted && (to == workflow.ChangesRequired || to == workflow.Approved || to == workflow.Rejected) {
		if actor.HasPermission("submissions:review") {
			return nil
		}
		return workflow.ErrTransitionNotAllowed
	}
	return workflow.ErrTransitionNotAllowed
}
