package domain

import (
	"context"
	"time"
)

type ChatMessage struct {
	Role    string `bson:"role" json:"role"`
	Content string `bson:"content" json:"content"`
}

type Session struct {
	ID        string        `bson:"_id,omitempty" json:"id"`
	ProjectID string        `bson:"projectId" json:"projectId"`
	UserID    string        `bson:"userId" json:"userId"`
	Messages  []ChatMessage `bson:"messages" json:"messages"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt" json:"updatedAt"`
}

type SessionRepository interface {
	Create(ctx context.Context, session *Session) error
	FindByProjectID(ctx context.Context, projectID string, userID string) (*Session, error)
	Update(ctx context.Context, session *Session) error
}
