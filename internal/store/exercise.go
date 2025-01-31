package store

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Exercise struct {
	ID               primitive.ObjectID `bson:"_id" json:"id"`
	UserID           primitive.ObjectID `bson:"user_id" json:"user_id"`
	Name             string             `bson:"name" json:"name"`
	Force            *string            `bson:"force" json:"force"`
	Level            *string            `bson:"level" json:"level"`
	Mechanic         *string            `bson:"mechanic" json:"mechanic"`
	Equipment        *string            `bson:"equipment" json:"equipment"`
	PrimaryMuscles   *[]string          `bson:"primaryMuscles" json:"primaryMuscles"`
	SecondaryMuscles *[]string          `bson:"secondaryMuscles" json:"secondaryMuscles"`
	Instructions     *[]string          `bson:"instructions" json:"instructions"`
	Category         string             `bson:"category" json:"category"`
	IsCustom         bool               `bson:"is_custom" json:"is_custom"`
	Version          int16              `bson:"version" json:"version"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`
}

type ExerciseStore struct {
	db *mongo.Database
}

const exerciseCollection = "exercise"

// create a custom exercise
func (s *ExerciseStore) Create(ctx context.Context, exercise *Exercise, userID primitive.ObjectID) error {
	// Generate a new ObjectID and assign a timestamp.
	exercise.ID = primitive.NewObjectID()
	exercise.CreatedAt = time.Now()
	exercise.UpdatedAt = time.Now()
	exercise.IsCustom = true

	if exercise.Version == 0 {
		exercise.Version = 1
	}

	// Set a timeout for the database operation.
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Insert the exercise into the MongoDB collection.
	_, err := s.db.Collection(exerciseCollection).InsertOne(ctx, exercise)
	if err != nil {
		return fmt.Errorf("failed to insert exercise: %w", err)
	}

	// Update user to include custom exerciseID
	_, err = s.db.Collection("user").UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{"$push": bson.M{"custom_exercises": exercise.ID}},
	)
	if err != nil {
		return fmt.Errorf("failed to associate custom exercise with user: %w", err)
	}
	return nil
}

// return array of all user custom exercises
func (s *ExerciseStore) GetAllUserExercises(ctx context.Context, userID primitive.ObjectID) ([]*Exercise, error) {
	var exercises []*Exercise

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}

	cursor, err := s.db.Collection(exerciseCollection).Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch exercises: %w", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &exercises); err != nil {
		return nil, fmt.Errorf("failed to decode exercises: %w", err)
	}

	return exercises, nil
}

// return a single exercise by a exercsie ID and user id for validation
func (s *ExerciseStore) GetByID(ctx context.Context, exerciseID, userID primitive.ObjectID) (*Exercise, error) {
	exercise := &Exercise{}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{
		"_id":     exerciseID,
		"user_id": userID,
	}

	err := s.db.Collection(exerciseCollection).FindOne(ctx, filter).Decode(exercise)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch exercise: %w", err)
	}

	return exercise, nil
}

func (s *ExerciseStore) Update(ctx context.Context, exerciseID, userID primitive.ObjectID, updates map[string]interface{}, expectedVersion int16) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Build filter to ensure the exercise belongs to the user and matches the expected version
	filter := bson.M{
		"_id":     exerciseID,
		"user_id": userID,
		"version": expectedVersion,
	}

	// Prepare update fields
	updateFields := bson.M{}
	for key, value := range updates {
		updateFields[key] = value
	}

	// Add updated_at timestamp and increment version
	updateFields["updated_at"] = time.Now()
	updateFields["version"] = expectedVersion + 1

	// Build the update query
	update := bson.M{
		"$set": updateFields,
	}

	// Perform the update
	result, err := s.db.Collection(exerciseCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update exercise: %w", err)
	}

	// Check if the exercise was found and updated
	if result.ModifiedCount == 0 {
		return fmt.Errorf("no exercise found with ID %s or version mismatch", exerciseID.Hex())
	}

	return nil
}

// delete a custom exercise
func (s *ExerciseStore) Delete(ctx context.Context, exerciseID primitive.ObjectID, userID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": exerciseID, "user_id": userID}

	result, err := s.db.Collection(exerciseCollection).DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete exercise: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no exercise found for the given user and or exercise ID")
	}

	_, err = s.db.Collection("user").UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{"$pull": bson.M{"custom_exercises": exerciseID}},
	)
	if err != nil {
		return fmt.Errorf("failed to remove exercise from user's routines array: %w", err)
	}

	return nil
}
