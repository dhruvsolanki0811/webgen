package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dhruvsolanki0811/webgen/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type userRepo struct {
	coll *mongo.Collection
}

func NewUserRepo(db *DB) domain.UserRepository {
	return &userRepo{coll: db.Collection(domain.CollectionUsers)}
}

func (r *userRepo) Create(ctx context.Context, user *domain.User) error {
	user.CreatedAt = time.Now()

	result, err := r.coll.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}

	if oid, ok := result.InsertedID.(bson.ObjectID); ok {
		user.ID = oid.Hex()
	}

	return nil
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, domain.ErrNotFound("user")
	}
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}

	return &user, nil
}

func (r *userRepo) FindByID(ctx context.Context, id string) (*domain.User, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrNotFound("user")
	}

	var user domain.User
	err = r.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, domain.ErrNotFound("user")
	}
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}

	return &user, nil
}
