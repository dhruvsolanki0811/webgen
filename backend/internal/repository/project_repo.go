package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dhruvsolanki0811/webgen/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type projectRepo struct {
	coll *mongo.Collection
}

func NewProjectRepo(db *DB) domain.ProjectRepository {
	return &projectRepo{coll: db.Collection(domain.CollectionProjects)}
}

func (r *projectRepo) Create(ctx context.Context, project *domain.Project) error {
	now := time.Now()
	project.CreatedAt = now
	project.UpdatedAt = now
	project.Version = 1
	project.Status = domain.StatusDraft

	result, err := r.coll.InsertOne(ctx, project)
	if err != nil {
		return fmt.Errorf("insert project: %w", err)
	}

	if oid, ok := result.InsertedID.(bson.ObjectID); ok {
		project.ID = oid.Hex()
	}

	return nil
}

func (r *projectRepo) FindByID(ctx context.Context, id string, userID string) (*domain.Project, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrNotFound("project")
	}

	filter := bson.M{"_id": oid, "userId": userID}

	var project domain.Project
	err = r.coll.FindOne(ctx, filter).Decode(&project)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, domain.ErrNotFound("project")
	}
	if err != nil {
		return nil, fmt.Errorf("find project: %w", err)
	}

	return &project, nil
}

func (r *projectRepo) FindByUserID(ctx context.Context, userID string) ([]domain.Project, error) {
	filter := bson.M{"userId": userID}
	opts := options.Find().SetSort(bson.M{"createdAt": -1})

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("find projects: %w", err)
	}
	defer cursor.Close(ctx)

	var projects []domain.Project
	if err := cursor.All(ctx, &projects); err != nil {
		return nil, fmt.Errorf("decode projects: %w", err)
	}

	if projects == nil {
		projects = []domain.Project{}
	}

	return projects, nil
}

func (r *projectRepo) UpdateStatus(ctx context.Context, id string, userID string, status string) error {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrNotFound("project")
	}

	filter := bson.M{"_id": oid, "userId": userID}
	update := bson.M{"$set": bson.M{
		"status":    status,
		"updatedAt": time.Now(),
	}}

	result, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("update project status: %w", err)
	}
	if result.MatchedCount == 0 {
		return domain.ErrNotFound("project")
	}

	return nil
}

func (r *projectRepo) UpdateSpec(ctx context.Context, id string, userID string, spec *domain.ProjectSpec) error {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrNotFound("project")
	}

	filter := bson.M{"_id": oid, "userId": userID}
	update := bson.M{"$set": bson.M{
		"spec":      spec,
		"status":    domain.StatusGenerating,
		"updatedAt": time.Now(),
	}}

	result, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("update project spec: %w", err)
	}
	if result.MatchedCount == 0 {
		return domain.ErrNotFound("project")
	}

	return nil
}

func (r *projectRepo) UpdateDeployment(ctx context.Context, id string, userID string, repoURL string, deployURL string) error {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrNotFound("project")
	}

	filter := bson.M{"_id": oid, "userId": userID}
	update := bson.M{"$set": bson.M{
		"repoUrl":   repoURL,
		"deployUrl": deployURL,
		"status":    domain.StatusDeployed,
		"updatedAt": time.Now(),
	}}

	result, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("update project deployment: %w", err)
	}
	if result.MatchedCount == 0 {
		return domain.ErrNotFound("project")
	}

	return nil
}
