package repository

import (
	"context"
	"user-service/model"

	"go.mongodb.org/mongo-driver/mongo"
)

type RedisMongo struct {
	Collection *mongo.Collection
	JwtSecret  string
}

func NewUserRepo(db *mongo.Database) *RedisMongo {
	return &RedisMongo{
		Collection: db.Collection("users"),
	}
}

func (r *RedisMongo) CreateUser(ctx context.Context, user *model.User) error {
	_, err := r.Collection.InsertOne(ctx, user)
	return err
}

func (r *RedisMongo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.Collection.FindOne(ctx, map[string]interface{}{"email": email}).Decode(&user)

	if err != nil {
		return nil, err
	}
	return &user, nil
}
