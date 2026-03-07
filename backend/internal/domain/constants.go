package domain

import "time"

// Collection names
const (
	CollectionUsers    = "users"
	CollectionProjects = "projects"
	CollectionSessions = "sessions"
)

// Project statuses
const (
	StatusDraft      = "draft"
	StatusClarifying = "clarifying"
	StatusSpecReady  = "spec_ready"
	StatusGenerating = "generating"
	StatusValidated  = "validated"
	StatusDeploying  = "deploying"
	StatusDeployed   = "deployed"
	StatusFailed     = "failed"
)

// Validation limits
const (
	MaxIdeaLength    = 2000
	MaxMessageLength = 500
	MinPasswordLen   = 8
	MaxNameLength    = 100
)

// Auth
const (
	TokenAccessTTL  = 15 * time.Minute
	TokenRefreshTTL = 7 * 24 * time.Hour
)

// Rate limits
const (
	RateLimitRequestsPerMinute = 100
	RateLimitBurst             = 10
	RateLimitGenerationsPerDay = 5
	RateLimitCleanupInterval   = 5 * time.Minute
	RateLimitExpiry            = 10 * time.Minute
)

// Generation
const (
	MaxGenerationRetries   = 2
	ClaudeModel            = "claude-sonnet-4-5-20250929"
	ClarificationMaxTokens = 1024
	GenerationMaxTokens    = 64000
)

// Context keys
type contextKey string

const ContextKeyUserID contextKey = "userId"
