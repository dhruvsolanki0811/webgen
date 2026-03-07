package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/dhruvsolanki0811/webgen/internal/domain"
)

type claims struct {
	UserID string `json:"uid"`
	jwt.RegisteredClaims
}

type tokenService struct {
	secret []byte
}

func NewTokenService(secret string) domain.TokenService {
	return &tokenService{secret: []byte(secret)}
}

func (s *tokenService) GenerateAccess(userID string) (string, error) {
	return s.generate(userID, domain.TokenAccessTTL)
}

func (s *tokenService) GenerateRefresh(userID string) (string, error) {
	return s.generate(userID, domain.TokenRefreshTTL)
}

func (s *tokenService) ValidateToken(tokenStr string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &claims{}, s.keyFunc)
	if err != nil {
		return "", domain.ErrUnauthorized("invalid or expired token")
	}

	c, ok := token.Claims.(*claims)
	if !ok || !token.Valid {
		return "", domain.ErrUnauthorized("invalid token claims")
	}

	return c.UserID, nil
}

func (s *tokenService) generate(userID string, ttl time.Duration) (string, error) {
	now := time.Now()
	c := claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return signed, nil
}

func (s *tokenService) keyFunc(_ *jwt.Token) (any, error) {
	return s.secret, nil
}
