package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/dhruvsolanki0811/webgen/internal/domain"
	"github.com/dhruvsolanki0811/webgen/internal/handler"
)

type AuthMiddleware struct {
	tokens domain.TokenService
}

func NewAuthMiddleware(tokens domain.TokenService) *AuthMiddleware {
	return &AuthMiddleware{tokens: tokens}
}

func (m *AuthMiddleware) Required(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := m.extractUserID(r)
		if err != nil {
			handler.RespondError(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), domain.ContextKeyUserID, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) extractUserID(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", domain.ErrUnauthorized("missing authorization header")
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", domain.ErrUnauthorized("invalid authorization format")
	}

	userID, err := m.tokens.ValidateToken(parts[1])
	if err != nil {
		return "", err
	}

	return userID, nil
}

func UserIDFromContext(ctx context.Context) string {
	id, _ := ctx.Value(domain.ContextKeyUserID).(string)
	return id
}
