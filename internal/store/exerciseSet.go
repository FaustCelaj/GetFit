package store

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// this contains multiple sets, not just one
// this represents the sets you do in a day
type ExerciseSet struct {
	ID                primitive.ObjectID `bson:"_id" json:"id"`
	RoutineExerciseID primitive.ObjectID `bson:"routine_exercise_id" json:"routine_exercise_id"`
	// UserID            primitive.ObjectID `bson:"user_id" json:"user_id"`
	SetNumber int16 `bson:"set_number" json:"set_number"`
	// IsWarmup          bool               `bson:"is_warmup" json:"is_warmup"`
	Weight    int16     `bson:"weight" json:"weight"`
	Reps      int16     `bson:"reps" json:"reps"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type ExerciseSetStore struct {
	db *mongo.Database
}

const exerciseSetCollection = "exerciseSet"

// Insert SINGLE set
func (s *ExerciseSetStore) AddSet(ctx context.Context, exerciseSet *ExerciseSet, routineExerciseID primitive.ObjectID) error {
	exerciseSet.ID = primitive.NewObjectID()
	exerciseSet.RoutineExerciseID = routineExerciseID
	exerciseSet.CreatedAt = time.Now()
	exerciseSet.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := s.db.Collection(exerciseSetCollection).InsertOne(ctx, exerciseSet)
	if err != nil {
		return fmt.Errorf("failed to add sets: %w", err)
	}

	return nil
}

func (s *ExerciseSetStore) AddMultipleSet(ctx context.Context, exerciseSets []ExerciseSet, routineExerciseID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// preparing the documents for a bulk insert
	var documents []interface{}
	for _, set := range exerciseSets {
		set.ID = primitive.NewObjectID()
		set.RoutineExerciseID = routineExerciseID
		set.CreatedAt = time.Now()
		set.UpdatedAt = time.Now()
		documents = append(documents, set)
	}

	_, err := s.db.Collection(exerciseSetCollection).InsertMany(ctx, documents)
	if err != nil {
		return fmt.Errorf("failed to add sets: %w", err)
	}

	return nil
}
