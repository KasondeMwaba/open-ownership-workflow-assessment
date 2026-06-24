package services

import (
	"context"
	"errors"

	"openownership-workflow/backend/internal/auth"
	"openownership-workflow/backend/internal/models"
	"openownership-workflow/backend/internal/repositories"
)

var ErrInvalidCredentials = errors.New("invalid email or password")
var ErrUserDisabled = errors.New("user account is disabled")

type AuthService struct {
	repo      *repositories.Repository
	jwtSecret string
}

type LoginResult struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func NewAuthService(repo *repositories.Repository, jwtSecret string) *AuthService {
	return &AuthService{repo: repo, jwtSecret: jwtSecret}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (LoginResult, error) {
	user, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil || !auth.VerifyPassword(user.PasswordHash, password) {
		return LoginResult{}, ErrInvalidCredentials
	}
	if !user.IsActive {
		return LoginResult{}, ErrUserDisabled
	}
	user, err = s.repo.AttachUserPermissions(ctx, user)
	if err != nil {
		return LoginResult{}, err
	}
	token, err := auth.IssueToken(s.jwtSecret, user)
	if err != nil {
		return LoginResult{}, err
	}
	return LoginResult{Token: token, User: user}, nil
}

func (s *AuthService) FindUserByID(ctx context.Context, id string) (models.User, error) {
	user, err := s.repo.FindUserByID(ctx, id)
	if err != nil {
		return user, err
	}
	return s.repo.AttachUserPermissions(ctx, user)
}

func (s *AuthService) FindUserByEmail(ctx context.Context, email string) (models.User, error) {
	return s.repo.FindUserByEmail(ctx, email)
}
