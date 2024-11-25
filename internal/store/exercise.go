package store

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// this shoudldnt have any access for users to manipulate
// TODO: consider this struct for custom exerscies that are unique to them?
type Exercise struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name             string             `bson:"name" json:"name"`
	Force            *string            `bson:"force,omitempty" json:"force,omitempty"`                       // Nullable, string or null
	Level            *string            `bson:"level,omitempty" json:"level,omitempty"`                       // Nullable, string or null
	Mechanic         *string            `bson:"mechanic,omitempty" json:"mechanic,omitempty"`                 // Nullable
	Equipment        *string            `bson:"equipment,omitempty" json:"equipment,omitempty"`               // Nullable
	PrimaryMuscles   *[]string          `bson:"primaryMuscles,omitempty" json:"primaryMuscles,omitempty"`     // Nullable array
	SecondaryMuscles *[]string          `bson:"secondaryMuscles,omitempty" json:"secondaryMuscles,omitempty"` // Nullable array
	Instructions     *[]string          `bson:"instructions,omitempty" json:"instructions,omitempty"`         // Nullable array
	Category         string             `bson:"category" json:"category"`                                     // Required
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`
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
