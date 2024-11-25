package store

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoutineExercise struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	ExerciseID primitive.ObjectID `bson:"exercise_id" json:"exercise_id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	Notes      *string            `bson:"notes,omitempty" json:"notes,omitempty"`
	Sets       []ExerciseSet      `bson:"sets" json:"sets"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

type RoutineExerciseStore struct {
	db *mongo.Database
}

const routineExerciseCollection = "routineExercise"

// Creating a routineExercise
func (s *RoutineExerciseStore) Create(ctx context.Context, routineExercise *RoutineExercise, userID primitive.ObjectID, exerciseID primitive.ObjectID) error {
	// Populate the fields
	routineExercise.ID = primitive.NewObjectID()
	routineExercise.UserID = userID
	routineExercise.ExerciseID = exerciseID
	routineExercise.CreatedAt = time.Now()
	routineExercise.UpdatedAt = time.Now()

	// Set a timeout for the database operation
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Check for duplicates
	filter := bson.M{"user_id": userID, "exercise_id": exerciseID}
	err := s.db.Collection(routineExerciseCollection).FindOne(ctx, filter).Err()
	if err == nil {
		// If no error, a duplicate exists
		return fmt.Errorf("duplicate routine exercise: userID %s, exerciseID %s already exists", userID.Hex(), exerciseID.Hex())
	} else if err != mongo.ErrNoDocuments {
		// If another error occurs, return it
		return fmt.Errorf("failed to check for duplicate: %w", err)
	}

	// If no duplicate exists, insert the new document
	_, err = s.db.Collection(routineExerciseCollection).InsertOne(ctx, routineExercise)
	if err != nil {
		return fmt.Errorf("failed to create routine exercise: %w", err)
	}

	return nil
}

// geting routineExercise by ID
func (s *RoutineExerciseStore) GetByID(ctx context.Context, routineExerciseID primitive.ObjectID) (*RoutineExercise, error) {
	routineExercise := &RoutineExercise{}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Fetch the routine exercise
	filter := bson.M{"_id": routineExerciseID}
	err := s.db.Collection(routineExerciseCollection).FindOne(ctx, filter).Decode(routineExercise)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("routine exercise not found: %w", err)
		}
		return nil, fmt.Errorf("failed to fetch routine exercise: %w", err)
	}

	// Fetch the associated sets
	var sets []ExerciseSet
	setFilter := bson.M{"routine_exercise_id": routineExerciseID}
	cursor, err := s.db.Collection("exerciseSet").Find(ctx, setFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch exercise sets: %w", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &sets); err != nil {
		return nil, fmt.Errorf("failed to decode exercise sets: %w", err)
	}

	// Populate the sets field
	routineExercise.Sets = sets

	return routineExercise, nil
}
