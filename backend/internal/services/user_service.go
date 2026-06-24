package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"openownership-workflow/backend/internal/auth"
	"openownership-workflow/backend/internal/models"
	"openownership-workflow/backend/internal/repositories"
	"openownership-workflow/backend/internal/workflow"
)

var (
	ErrAdminRequired     = errors.New("admin privileges are required")
	ErrInvalidRole       = errors.New("role does not exist")
	ErrInvalidUserInput  = errors.New("name and valid email are required")
	ErrPasswordTooShort  = errors.New("password must be at least 8 characters")
	ErrCannotDisableSelf = errors.New("admins cannot disable their own account")
)

type UserService struct {
	repo *repositories.Repository
}

type CreateUserInput struct {
	Name     string        `json:"name"`
	Email    string        `json:"email"`
	Password string        `json:"password"`
	Role     workflow.Role `json:"role"`
	IsActive bool          `json:"isActive"`
}

type UpdateUserInput struct {
	Name     string        `json:"name"`
	Email    string        `json:"email"`
	Role     workflow.Role `json:"role"`
	IsActive bool          `json:"isActive"`
}

func NewUserService(repo *repositories.Repository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) List(ctx context.Context, actor models.User) ([]models.User, error) {
	if actor.Role != workflow.Admin {
		return nil, ErrAdminRequired
	}
	return s.repo.ListUsers(ctx)
}

func (s *UserService) Create(ctx context.Context, actor models.User, input CreateUserInput) (models.User, error) {
	if actor.Role != workflow.Admin {
		return models.User{}, ErrAdminRequired
	}
	name, email, err := validateUserIdentity(input.Name, input.Email)
	if err != nil {
		return models.User{}, err
	}
	role, err := s.validateRole(ctx, input.Role)
	if err != nil {
		return models.User{}, err
	}
	if len(input.Password) < 8 {
		return models.User{}, ErrPasswordTooShort
	}
	hash, err := auth.HashPassword(input.Password)
	if err != nil {
		return models.User{}, err
	}
	user, err := s.repo.CreateUser(ctx, repositories.CreateUserParams{
		Name:         name,
		Email:        email,
		PasswordHash: hash,
		Role:         role,
		IsActive:     input.IsActive,
	})
	if err != nil {
		return models.User{}, err
	}
	_ = s.recordAudit(ctx, actor, "user.created", user, map[string]any{
		"email":    user.Email,
		"role":     user.Role,
		"isActive": user.IsActive,
	})
	return user, nil
}

func (s *UserService) Update(ctx context.Context, actor models.User, id string, input UpdateUserInput) (models.User, error) {
	if actor.Role != workflow.Admin {
		return models.User{}, ErrAdminRequired
	}
	name, email, err := validateUserIdentity(input.Name, input.Email)
	if err != nil {
		return models.User{}, err
	}
	role, err := s.validateRole(ctx, input.Role)
	if err != nil {
		return models.User{}, err
	}
	if actor.ID == id && !input.IsActive {
		return models.User{}, ErrCannotDisableSelf
	}
	before, err := s.repo.FindUserByID(ctx, id)
	if err != nil {
		return models.User{}, err
	}
	user, err := s.repo.UpdateUser(ctx, id, repositories.UpdateUserParams{
		Name:     name,
		Email:    email,
		Role:     role,
		IsActive: input.IsActive,
	})
	if err != nil {
		return models.User{}, err
	}
	_ = s.recordAudit(ctx, actor, "user.updated", user, map[string]any{
		"before": map[string]any{"name": before.Name, "email": before.Email, "role": before.Role, "isActive": before.IsActive},
		"after":  map[string]any{"name": user.Name, "email": user.Email, "role": user.Role, "isActive": user.IsActive},
	})
	return user, nil
}

func validateUserIdentity(name, email string) (string, string, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(strings.ToLower(email))
	if name == "" || email == "" {
		return "", "", ErrInvalidUserInput
	}
	parsed, err := mail.ParseAddress(email)
	if err != nil || parsed.Address != email {
		return "", "", ErrInvalidUserInput
	}
	return name, email, nil
}

func (s *UserService) SetActive(ctx context.Context, actor models.User, id string, active bool) (models.User, error) {
	if actor.Role != workflow.Admin {
		return models.User{}, ErrAdminRequired
	}
	if actor.ID == id && !active {
		return models.User{}, ErrCannotDisableSelf
	}
	user, err := s.repo.SetUserActive(ctx, id, active)
	if err != nil {
		return models.User{}, err
	}
	eventType := "user.disabled"
	if active {
		eventType = "user.enabled"
	}
	_ = s.recordAudit(ctx, actor, eventType, user, map[string]any{"isActive": user.IsActive})
	return user, nil
}

func (s *UserService) validateRole(ctx context.Context, role workflow.Role) (workflow.Role, error) {
	role = workflow.Role(strings.TrimSpace(string(role)))
	if role == "" {
		return "", ErrInvalidRole
	}
	exists, err := s.repo.RoleExists(ctx, string(role))
	if err != nil {
		return "", err
	}
	if !exists {
		return "", ErrInvalidRole
	}
	return role, nil
}

func (s *UserService) recordAudit(ctx context.Context, actor models.User, eventType string, user models.User, metadata map[string]any) error {
	body, err := json.Marshal(metadata)
	if err != nil {
		return err
	}
	return s.repo.RecordSystemAudit(ctx, repositories.SystemAuditParams{
		ActorID:      actor.ID,
		EventType:    eventType,
		ResourceType: "user",
		ResourceID:   user.ID,
		ResourceName: user.Email,
		Summary:      fmt.Sprintf("%s %s", eventType, user.Email),
		Metadata:     body,
	})
}
