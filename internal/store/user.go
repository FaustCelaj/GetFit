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
	FirstName string             `bson:"first_name" json:"first_name"`
	LastName  string             `bson:"last_name" json:"last_name"`
	Age       int8               `bson:"age" json:"age"`
	Title     string             `bson:"title" json:"title"`
	Bio       string             `bson:"bio" json:"bio"`
	Version   int16              `bson:"version" json:"version"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type UserStore struct {
	db *mongo.Database
}

const userCollection = "user"

// GET user by ID
func (s *UserStore) GetByID(ctx context.Context, userID primitive.ObjectID) (*User, error) {
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

	// setting the version number
	if user.Version == 0 {
		user.Version = 1
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := s.db.Collection(userCollection).InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

// UPDATING user
func (s *UserStore) Update(ctx context.Context, userID primitive.ObjectID, updates map[string]interface{}, expectedVersion int16) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": userID, "version": expectedVersion}

	updateFields := bson.M{}
	if username, ok := updates["username"]; ok {
		updateFields["username"] = username
	}
	if email, ok := updates["email"]; ok {
		updateFields["email"] = email
	}
	if firstName, ok := updates["first_name"]; ok {
		updateFields["first_name"] = firstName
	}
	if lastName, ok := updates["last_name"]; ok {
		updateFields["last_name"] = lastName
	}
	if age, ok := updates["age"]; ok {
		updateFields["age"] = age
	}
	if title, ok := updates["title"]; ok {
		updateFields["title"] = title
	}
	if bio, ok := updates["bio"]; ok {
		updateFields["bio"] = bio
	}

	updateFields["updated_at"] = time.Now()

	update := bson.M{"$set": updateFields, "$inc": bson.M{"version": 1}}

	// Perform the update
	result, err := s.db.Collection(userCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.MatchedCount == 0 {
		return ErrVersionMismatch
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
