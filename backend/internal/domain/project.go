package domain

import (
	"context"
	"time"
)

type ProjectSpec struct {
	Description  string   `bson:"description" json:"description"`
	Features     []string `bson:"features" json:"features"`
	Pages        []string `bson:"pages" json:"pages"`
	Styling      string   `bson:"styling" json:"styling"`
	HasAuth      bool     `bson:"hasAuth" json:"hasAuth"`
	HasForms     bool     `bson:"hasForms" json:"hasForms"`
	HasAPIRoutes bool     `bson:"hasApiRoutes" json:"hasApiRoutes"`
}

type Project struct {
	ID        string      `bson:"_id,omitempty" json:"id"`
	UserID    string      `bson:"userId" json:"userId"`
	Name      string      `bson:"name" json:"name"`
	Spec      ProjectSpec `bson:"spec" json:"spec"`
	RepoURL   string      `bson:"repoUrl" json:"repoUrl"`
	DeployURL string      `bson:"deployUrl" json:"deployUrl"`
	Status    string      `bson:"status" json:"status"`
	Version   int         `bson:"version" json:"version"`
	CreatedAt time.Time   `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time   `bson:"updatedAt" json:"updatedAt"`
}

type ProjectRepository interface {
	Create(ctx context.Context, project *Project) error
	FindByID(ctx context.Context, id string, userID string) (*Project, error)
	FindByUserID(ctx context.Context, userID string) ([]Project, error)
	UpdateStatus(ctx context.Context, id string, userID string, status string) error
	UpdateSpec(ctx context.Context, id string, userID string, spec *ProjectSpec) error
	UpdateDeployment(ctx context.Context, id string, userID string, repoURL string, deployURL string) error
}
