package store

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Exercise struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	Name         string             `bson:"name" json:"name"`
	ExerciseType string             `bson:"exercise_type" json:"exercise_type"`
	Muscle       string             `bson:"muscle" json:"muscle"`
	Equipment    string             `bson:"equipment" json:"equipment"`
	Difficulty   string             `bson:"difficulty" json:"difficulty"`
	Instructions string             `bson:"instructions" json:"instructions"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}

type ExerciseStore struct {
	db *mongo.Database
}

const exerciseCollection = "exercise"

func (s *ExerciseStore) Create(ctx context.Context, exercise *Exercise) error {
	// Generate a new ObjectID and assign a timestamp.
	exercise.ID = primitive.NewObjectID()
	exercise.CreatedAt = time.Now()

	// Set a timeout for the database operation.
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Insert the exercise into the MongoDB collection.
	_, err := s.db.Collection(exerciseCollection).InsertOne(ctx, exercise)
	if err != nil {
		return fmt.Errorf("failed to insert exercise: %w", err)
	}

	return nil
}
