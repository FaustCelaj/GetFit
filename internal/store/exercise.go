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
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID           primitive.ObjectID `bson:"user_id" json:"user_id"`
	Name             string             `bson:"name" json:"name"`
	Force            *string            `bson:"force,omitempty" json:"force,omitempty"`                       // Nullable, string or null
	Level            *string            `bson:"level,omitempty" json:"level,omitempty"`                       // Nullable, string or null
	Mechanic         *string            `bson:"mechanic,omitempty" json:"mechanic,omitempty"`                 // Nullable
	Equipment        *string            `bson:"equipment,omitempty" json:"equipment,omitempty"`               // Nullable
	PrimaryMuscles   *[]string          `bson:"primaryMuscles,omitempty" json:"primaryMuscles,omitempty"`     // Nullable array
	SecondaryMuscles *[]string          `bson:"secondaryMuscles,omitempty" json:"secondaryMuscles,omitempty"` // Nullable array
	Instructions     *[]string          `bson:"instructions,omitempty" json:"instructions,omitempty"`         // Nullable array
	Category         string             `bson:"category" json:"category"`                                     // Required
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

	// updating user to include custom exerciseID
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

	filter := bson.M{
		"_id":     exerciseID,
		"user_id": userID,
		"version": expectedVersion,
	}

	updateFields := bson.M{}
	if name, ok := updates["name"]; ok {
		updateFields["name"] = name
	}
	if force, ok := updates["force"]; ok {
		updateFields["force"] = force
	}
	if level, ok := updates["level"]; ok {
		updateFields["level"] = level
	}
	if mechanic, ok := updates["mechanic"]; ok {
		updateFields["mechanic"] = mechanic
	}
	if equipment, ok := updates["equipment"]; ok {
		updateFields["equipment"] = equipment
	}
	if primaryMuscles, ok := updates["primaryMuscles"]; ok {
		updateFields["primaryMuscles"] = primaryMuscles
	}
	if secondaryMuscles, ok := updates["secondaryMuscles"]; ok {
		updateFields["secondaryMuscles"] = secondaryMuscles
	}
	if instructions, ok := updates["instructions"]; ok {
		updateFields["instructions"] = instructions
	}
	if category, ok := updates["category"]; ok {
		updateFields["category"] = category
	}

	updateFields["updated_at"] = time.Now()

	update := bson.M{
		"$set": updateFields,
		"$inc": bson.M{"version": 1},
	}

	result, err := s.db.Collection(exerciseCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update exercise: %w", err)
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("no exercise found with ID %s", exerciseID.Hex())
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
