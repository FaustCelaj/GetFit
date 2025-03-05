package store

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WorkoutSession struct {
	ID          primitive.ObjectID     `bson:"_id" json:"id"`
	UserID      primitive.ObjectID     `bson:"user_id" json:"user_id"`
	RoutineID   *primitive.ObjectID    `bson:"routine_id,omitempty" json:"routine_id,omitempty"` // Optional: may be a routine-based or freestyle workout
	Title       string                 `bson:"title" json:"title"`
	Description *string                `bson:"description,omitempty" json:"description,omitempty"`
	Status      string                 `bson:"status" json:"status"` // "in_progress", "completed"
	StartTime   time.Time              `bson:"start_time" json:"start_time"`
	EndTime     *time.Time             `bson:"end_time,omitempty" json:"end_time,omitempty"`
	Exercises   []SessionExercise      `bson:"exercises" json:"exercises"`
	Notes       *string                `bson:"notes,omitempty" json:"notes,omitempty"`
	Metrics     map[string]interface{} `bson:"metrics,omitempty" json:"metrics,omitempty"` // For calculated values like total weight lifted
	Version     int16                  `bson:"version" json:"version"`
	CreatedAt   time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time              `bson:"updated_at" json:"updated_at"`
}

type SessionExercise struct {
	ExerciseID    primitive.ObjectID `bson:"exercise_id" json:"exercise_id"`
	Order         int                `bson:"order" json:"order"` // Position in the workout
	CompletedSets []SessionSet       `bson:"completed_sets" json:"completed_sets"`
}

type SessionSet struct {
	Weight      float32   `bson:"weight" json:"weight"`
	Reps        int16     `bson:"reps" json:"reps"`
	SetNumber   int16     `bson:"set_number" json:"set_number"`
	CompletedAt time.Time `bson:"completed_at" json:"completed_at"`
}

type WorkoutSessionStore struct {
	db *mongo.Database
}

const workoutSessionCollection = "workout_session"

// starting a workout session from scratch (adding as we go)
func (s *WorkoutSessionStore) Create(ctx context.Context, session *WorkoutSession, userID primitive.ObjectID) error {

	session.ID = primitive.NewObjectID()
	session.UserID = userID
	session.CreatedAt = time.Now()
	session.UpdatedAt = time.Now()
	session.StartTime = time.Now()
	session.Status = "in_progress"

	if session.Version == 0 {
		session.Version = 1
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := s.db.Collection(workoutSessionCollection).InsertOne(ctx, session)
	if err != nil {
		return fmt.Errorf("failed to create workout session: %w", err)
	}

	return nil
}

// create a workout session from a routine (template)
func (s *WorkoutSessionStore) CreateFromRoutine(ctx context.Context, routineID, userID primitive.ObjectID) (*WorkoutSession, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	routineStore := &RoutineStore{db: s.db}
	routine, err := routineStore.GetByID(ctx, routineID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch routine: %w", err)
	}

	session := &WorkoutSession{
		ID:          primitive.NewObjectID(),
		UserID:      userID,
		RoutineID:   &routineID,
		Title:       routine.Title,
		Description: routine.Description,
		StartTime:   time.Now(),
		Status:      "in_progress",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Exercises:   []SessionExercise{},
	}

	for i, exerciseID := range routine.Exercises {
		sessionExercise := SessionExercise{
			ExerciseID:    exerciseID,
			Order:         i,
			CompletedSets: []SessionSet{},
		}
		session.Exercises = append(session.Exercises, sessionExercise)
	}

	_, err = s.db.Collection(workoutSessionCollection).InsertOne(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to create workout session from routine: %w", err)
	}

	return session, nil
}

// get all workouts for a user
func (s *WorkoutSessionStore) GetAllUserSessions(ctx context.Context, userID primitive.ObjectID) ([]*WorkoutSession, error) {
	var sessions []*WorkoutSession
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}

	// Sort by start_time in descending order (newest first)
	opts := options.Find().SetSort(bson.M{"start_time": -1})

	cursor, err := s.db.Collection(workoutSessionCollection).Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch workout sessions: %w", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &sessions); err != nil {
		return nil, fmt.Errorf("failed to decode workout sessions: %w", err)
	}

	return sessions, nil
}

// Get a single workout for a user
func (s *WorkoutSessionStore) GetByID(ctx context.Context, sessionID, userID primitive.ObjectID) (*WorkoutSession, error) {
	session := &WorkoutSession{}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{
		"_id":     sessionID,
		"user_id": userID,
	}

	err := s.db.Collection(workoutSessionCollection).FindOne(ctx, filter).Decode(session)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch workout session: %w", err)
	}

	return session, nil
}

// Update workout session (add or update sets)
func (s *WorkoutSessionStore) AddSetToExercise(ctx context.Context, sessionID, userID, exerciseID primitive.ObjectID, set SessionSet, expectedVersion int16) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Set the completion time if not already set
	if set.CompletedAt.IsZero() {
		set.CompletedAt = time.Now()
	}

	// Find the workout session and update the specific exercise's sets
	filter := bson.M{
		"_id":                   sessionID,
		"user_id":               userID,
		"exercises.exercise_id": exerciseID,
		"version":               expectedVersion,
	}

	update := bson.M{
		"$push": bson.M{
			"exercises.$.completed_sets": set,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	result, err := s.db.Collection(workoutSessionCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to add set to exercise: %w", err)
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("no workout session found with ID %s or exercise ID %s", sessionID.Hex(), exerciseID.Hex())
	}

	return nil
}

// Complete a workout session
func (s *WorkoutSessionStore) CompleteWorkout(ctx context.Context, sessionID, userID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Get the session to calculate metrics
	session, err := s.GetByID(ctx, sessionID, userID)
	if err != nil {
		return fmt.Errorf("failed to fetch workout session: %w", err)
	}

	// Calculate metrics
	totalWeight := float32(0)
	totalReps := int16(0)
	totalSets := int16(0)

	for _, exercise := range session.Exercises {
		for _, set := range exercise.CompletedSets {
			totalWeight += set.Weight * float32(set.Reps)
			totalReps += set.Reps
			totalSets++
		}
	}

	metrics := map[string]interface{}{
		"total_weight": totalWeight,
		"total_reps":   totalReps,
		"total_sets":   totalSets,
		"duration":     time.Now().Sub(session.StartTime).Minutes(),
	}

	endTime := time.Now()

	// Update the session status and metrics
	filter := bson.M{
		"_id":     sessionID,
		"user_id": userID,
	}

	update := bson.M{
		"$set": bson.M{
			"status":     "completed",
			"end_time":   endTime,
			"metrics":    metrics,
			"updated_at": time.Now(),
		},
	}

	result, err := s.db.Collection(workoutSessionCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to complete workout session: %w", err)
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("no workout session found with ID %s", sessionID.Hex())
	}

	return nil
}

// Delete a workout session
func (s *WorkoutSessionStore) Delete(ctx context.Context, sessionID, userID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{
		"_id":     sessionID,
		"user_id": userID,
	}

	result, err := s.db.Collection(workoutSessionCollection).DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete workout session: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no workout session found with ID %s", sessionID.Hex())
	}

	return nil
}
