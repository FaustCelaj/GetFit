package store

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Set struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	ExerciseID primitive.ObjectID `bson:"exercise_id" json:"exercise_id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	SetNumber  int16              `bson:"set_number" json:"set_number"`
	Weight     float32            `bson:"weight" json:"weight"`
	Reps       int16              `bson:"reps" json:"reps"`
	Date       time.Time          `bson:"date" json:"date"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

type SetStore struct {
	db *mongo.Database
}

const setCollection = "set"

// Insert SINGLE set
func (s *SetStore) AddSet(ctx context.Context, set *Set, exerciseID primitive.ObjectID, userID primitive.ObjectID) error {
	set.ID = primitive.NewObjectID()
	set.ExerciseID = exerciseID
	set.UserID = userID
	set.Date = time.Now()
	set.CreatedAt = time.Now()
	set.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := s.db.Collection(setCollection).InsertOne(ctx, set)
	if err != nil {
		return fmt.Errorf("failed to add set: %w", err)
	}

	return nil
}

// insert multiple sets
func (s *SetStore) AddMultipleSet(ctx context.Context, exerciseSets []Set, exerciseID primitive.ObjectID, userID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var documents []interface{}
	for _, set := range exerciseSets {
		set.ID = primitive.NewObjectID()
		set.ExerciseID = exerciseID
		set.UserID = userID
		set.CreatedAt = time.Now()
		set.UpdatedAt = time.Now()
		documents = append(documents, set)
	}

	_, err := s.db.Collection(setCollection).InsertMany(ctx, documents)
	if err != nil {
		return fmt.Errorf("failed to add sets: %w", err)
	}

	return nil
}
