package service

import (
	"context"
	"errors"
	"strings"

	"github.com/dhruvsolanki0811/webgen/internal/domain"
)

type AuthService struct {
	users  domain.UserRepository
	tokens domain.TokenService
}

func NewAuthService(users domain.UserRepository, tokens domain.TokenService) *AuthService {
	return &AuthService{users: users, tokens: tokens}
}

type AuthResult struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (s *AuthService) Signup(ctx context.Context, email string, password string) (*AuthResult, error) {
	email = normalizeEmail(email)

	if err := validateSignupInput(email, password); err != nil {
		return nil, err
	}

	existing, err := s.users.FindByEmail(ctx, email)
	if err != nil && !isNotFound(err) {
		return nil, domain.ErrInternal(err)
	}
	if existing != nil {
		return nil, domain.ErrConflict("email already registered")
	}

	user := &domain.User{Email: email}
	if err := user.SetPassword(password); err != nil {
		return nil, domain.ErrInternal(err)
	}

	if err := s.users.Create(ctx, user); err != nil {
		return nil, domain.ErrInternal(err)
	}

	return s.generateTokens(user.ID)
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (*AuthResult, error) {
	email = normalizeEmail(email)

	if email == "" || password == "" {
		return nil, domain.ErrBadRequest("email and password required")
	}

	user, err := s.users.FindByEmail(ctx, email)
	if isNotFound(err) {
		return nil, domain.ErrUnauthorized("invalid credentials")
	}
	if err != nil {
		return nil, domain.ErrInternal(err)
	}

	if !user.CheckPassword(password) {
		return nil, domain.ErrUnauthorized("invalid credentials")
	}

	return s.generateTokens(user.ID)
}

func (s *AuthService) generateTokens(userID string) (*AuthResult, error) {
	access, err := s.tokens.GenerateAccess(userID)
	if err != nil {
		return nil, domain.ErrInternal(err)
	}

	refresh, err := s.tokens.GenerateRefresh(userID)
	if err != nil {
		return nil, domain.ErrInternal(err)
	}

	return &AuthResult{AccessToken: access, RefreshToken: refresh}, nil
}

func validateSignupInput(email string, password string) error {
	if email == "" {
		return domain.ErrBadRequest("email is required")
	}
	if !strings.Contains(email, "@") {
		return domain.ErrBadRequest("invalid email format")
	}
	if len(password) < domain.MinPasswordLen {
		return domain.ErrBadRequest("password must be at least 8 characters")
	}
	return nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func isNotFound(err error) bool {
	var appErr *domain.AppError
	return errors.As(err, &appErr) && appErr.Code == 404
}
