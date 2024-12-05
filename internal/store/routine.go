package store

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Routine struct {
	ID          primitive.ObjectID   `bson:"_id" json:"id"`
	UserID      primitive.ObjectID   `bson:"user_id" json:"user_id"`
	Title       string               `bson:"title" json:"title"`
	Description *string              `bson:"description,omitempty" json:"description,omitempty"`
	ExerciseID  []primitive.ObjectID `bson:"exercise_id" json:"exercise_id"`
	Version     int16                `bson:"version" json:"version"`
	CreatedAt   time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time            `bson:"updated_at" json:"updated_at"`
}

type RoutineStore struct {
	db *mongo.Database
}

const routineCollection = "routine"

// Create a routine
func (s *RoutineStore) Create(ctx context.Context, routine *Routine, userID primitive.ObjectID) error {

	// checking for empty
	if len(routine.ExerciseID) == 0 {
		return fmt.Errorf("your routine needs at least 1 exercise")
	}

	if routine.Title == "" {
		return fmt.Errorf("title is required for a routine")
	}

	// assigning an ID
	routine.ID = primitive.NewObjectID()
	routine.CreatedAt = time.Now()
	routine.UpdatedAt = time.Now()

	// setting the version number
	if routine.Version == 0 {
		routine.Version = 1
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Insert the routine into the routine collection
	_, err := s.db.Collection(routineCollection).InsertOne(ctx, routine)
	if err != nil {
		return fmt.Errorf("failed to create routine: %w", err)
	}

	// Update the user document to include the routine ID
	_, err = s.db.Collection("user").UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{"$push": bson.M{"routines": routine.ID}},
	)
	if err != nil {
		return fmt.Errorf("failed to associate routine with user: %w", err)
	}

	return nil
}

// fetch all routines from a specified user
// returns an array of routines
func (s *RoutineStore) GetAllUserRoutines(ctx context.Context, userID primitive.ObjectID) ([]*Routine, error) {
	var routines []*Routine
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}

	cursor, err := s.db.Collection(routineCollection).Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch routines: %w", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &routines); err != nil {
		return nil, fmt.Errorf("failed to decode routines: %w", err)
	}

	return routines, nil
}

// Get a routine by routineID and userID
// returns a single routine to match specified userID and routineID
func (s *RoutineStore) GetByID(ctx context.Context, routineID, userID primitive.ObjectID) (*Routine, error) {
	routine := &Routine{}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Create a filter for both userID and routineID
	filter := bson.M{
		"_id":     routineID,
		"user_id": userID,
	}

	// Find the routine
	err := s.db.Collection(routineCollection).FindOne(ctx, filter).Decode(routine)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch routine: %w", err)
	}

	return routine, nil
}

// update a routine
func (s *RoutineStore) Update(ctx context.Context, routineID, userID primitive.ObjectID, updates map[string]interface{}, expectedVersion int16) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{
		"_id":     routineID,
		"user_id": userID,
		"version": expectedVersion,
	}

	updateFields := bson.M{}
	if exercises, ok := updates["exercise_id"]; ok {
		updateFields["exercise_id"] = exercises
	}
	if description, ok := updates["description"]; ok {
		updateFields["description"] = description
	}
	if title, ok := updates["title"]; ok {
		updateFields["title"] = title
	}
	updateFields["updated_at"] = time.Now()

	update := bson.M{
		"$set": updateFields,
		"$inc": bson.M{"version": 1},
	}

	result, err := s.db.Collection(routineCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update routine: %w", err)
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("no routine found with ID %s", routineID.Hex())
	}

	return nil
}

// Delete a routine
func (s *RoutineStore) Delete(ctx context.Context, routineID primitive.ObjectID, userID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": routineID, "user_id": userID}

	result, err := s.db.Collection(routineCollection).DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete routine: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no routine found for the given user and routine ID")
	}

	// Update the user's `routines` array to remove the deleted routine
	_, err = s.db.Collection("user").UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{"$pull": bson.M{"routines": routineID}},
	)
	if err != nil {
		return fmt.Errorf("failed to remove routine from user's routines array: %w", err)
	}

	return nil
}
