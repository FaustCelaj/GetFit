package store

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Username  string             `bson:"username" json:"username"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password_hash" json:"password_hash"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type UserStore struct {
	db *mongo.Database
}

const userCollection = "user"

// GET user by ID
func (s *UserStore) GetById(ctx context.Context, userID primitive.ObjectID) (*User, error) {
	user := &User{}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": userID}

	err := s.db.Collection(userCollection).FindOne(ctx, filter).Decode(user)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// CREATING user
func (s *UserStore) Create(ctx context.Context, user *User) error {
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := s.db.Collection(userCollection).InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

// UPDATING user
func (s *UserStore) Update(ctx context.Context, userID primitive.ObjectID, updates map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": userID}
	update := bson.M{"$set": updates}

	// Perform the update
	_, err := s.db.Collection(userCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// DELETING user
func (s *UserStore) Delete(ctx context.Context, userID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": userID}

	_, err := s.db.Collection(userCollection).DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
