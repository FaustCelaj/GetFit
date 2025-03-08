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
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	Title       string             `bson:"title" json:"title"`
	Description *string            `bson:"description,omitempty" json:"description,omitempty"`
	Exercises   []RoutineExercise  `bson:"exercises" json:"exercises"`
	Version     int16              `bson:"version" json:"version"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type RoutineExercise struct {
	ExerciseID primitive.ObjectID `bson:"exercise_id" json:"exercise_id"`
	Order      int                `bson:"order" json:"order"`
	Sets       []TemplateSet      `bson:"template_sets" json:"template_sets"`
}

type TemplateSet struct {
	Weight    float32 `bson:"weight" json:"weight"`
	Reps      int16   `bson:"reps" json:"reps"`
	SetNumber int16   `bson:"set_number" json:"set_number"`
}

type RoutineStore struct {
	db *mongo.Database
}

const routineCollection = "routine"

// Create a routine
func (s *RoutineStore) Create(ctx context.Context, routine *Routine, userID primitive.ObjectID) error {

	// checking for empty
	if len(routine.Exercises) == 0 {
		return fmt.Errorf("your routine needs at least 1 exercise")
	}

	if routine.Title == "" {
		return fmt.Errorf("title is required for a routine")
	}

	// assigning an ID
	routine.ID = primitive.NewObjectID()
	routine.UserID = userID
	routine.CreatedAt = time.Now()
	routine.UpdatedAt = time.Now()

	// setting the version number
	if routine.Version == 0 {
		routine.Version = 1
	}

	// setting the order if not provided
	for i := range routine.Exercises {
		if routine.Exercises[i].Order == 0 {
			routine.Exercises[i].Order = i
		}

		// setting set numbers if not provided
		for j := range routine.Exercises[i].Sets {
			if routine.Exercises[i].Sets[j].SetNumber == 0 {
				routine.Exercises[i].Sets[j].SetNumber = int16(j + 1)
			}
		}
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
func (s *RoutineStore) GetByID(ctx context.Context, routineID, userID primitive.ObjectID) (*Routine, error) {
	routine := &Routine{}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

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
	for key, value := range updates {
		updateFields[key] = value
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

// add an exercise to a routine
func (s *RoutineStore) AddExerciseToRoutine(ctx context.Context, routineID, userID, exerciseID primitive.ObjectID, templateSets []TemplateSet, expectedVersion int16) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Get current routine to find the next order number
	routine, err := s.GetByID(ctx, routineID, userID)
	if err != nil {
		return fmt.Errorf("failed to fetch routine: %w", err)
	}

	// Find the next order value
	nextOrder := 0
	for _, ex := range routine.Exercises {
		if ex.Order >= nextOrder {
			nextOrder = ex.Order + 1
		}
	}

	// Set numbers for the sets
	for i := range templateSets {
		if templateSets[i].SetNumber == 0 {
			templateSets[i].SetNumber = int16(i + 1)
		}
	}

	// Create new routine exercise
	newExercise := RoutineExercise{
		ExerciseID: exerciseID,
		Order:      nextOrder,
		Sets:       templateSets,
	}

	filter := bson.M{
		"_id":     routineID,
		"user_id": userID,
		"version": expectedVersion,
	}

	update := bson.M{
		"$push": bson.M{"exercises": newExercise},
		"$set":  bson.M{"updated_at": time.Now()},
		"$inc":  bson.M{"version": 1},
	}

	result, err := s.db.Collection(routineCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to add exercise to routine: %w", err)
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("no routine found with ID %s or version mismatch", routineID.Hex())
	}

	return nil
}

// update an exercise in a routine
func (s *RoutineStore) UpdateExerciseInRoutine(ctx context.Context, routineID, userID, exerciseID primitive.ObjectID, templateSets []TemplateSet, expectedVersion int16) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Set numbers for the sets
	for i := range templateSets {
		if templateSets[i].SetNumber == 0 {
			templateSets[i].SetNumber = int16(i + 1)
		}
	}

	filter := bson.M{
		"_id":                   routineID,
		"user_id":               userID,
		"exercises.exercise_id": exerciseID,
		"version":               expectedVersion,
	}

	update := bson.M{
		"$set": bson.M{
			"exercises.$.template_sets": templateSets,
			"updated_at":                time.Now(),
		},
		"$inc": bson.M{"version": 1},
	}

	result, err := s.db.Collection(routineCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update exercise in routine: %w", err)
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("no routine or exercise found with the given IDs or version mismatch")
	}

	return nil
}

// remove exercise from a routine
func (s *RoutineStore) RemoveExerciseFromRoutine(ctx context.Context, routineID, userID, exerciseID primitive.ObjectID, expectedVersion int16) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{
		"_id":     routineID,
		"user_id": userID,
		"version": expectedVersion,
	}

	update := bson.M{
		"$pull": bson.M{"exercises": bson.M{"exercise_id": exerciseID}},
		"$set":  bson.M{"updated_at": time.Now()},
		"$inc":  bson.M{"version": 1},
	}

	result, err := s.db.Collection(routineCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to remove exercise from routine: %w", err)
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("no routine found with ID %s or version mismatch", routineID.Hex())
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
