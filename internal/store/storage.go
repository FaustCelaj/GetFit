package store

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrConflict          = errors.New("resource already exists")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Users interface {
		Create(context.Context, *User) error
		GetByID(context.Context, primitive.ObjectID) (*User, error)
		Update(context.Context, primitive.ObjectID, map[string]interface{}, int16) error
		Delete(context.Context, primitive.ObjectID) error
	}
	Routine interface {
		Create(context.Context, *Routine, primitive.ObjectID) error
		GetAllUserRoutines(context.Context, primitive.ObjectID) ([]*Routine, error)
		GetByID(context.Context, primitive.ObjectID, primitive.ObjectID) (*Routine, error)
		Update(context.Context, primitive.ObjectID, primitive.ObjectID, map[string]interface{}, int16) error
		AddExerciseToRoutine(context.Context, primitive.ObjectID, primitive.ObjectID, primitive.ObjectID, []TemplateSet, int16) error
		UpdateExerciseInRoutine(context.Context, primitive.ObjectID, primitive.ObjectID, primitive.ObjectID, []TemplateSet, int16) error
		RemoveExerciseFromRoutine(context.Context, primitive.ObjectID, primitive.ObjectID, primitive.ObjectID, int16) error
		Delete(context.Context, primitive.ObjectID, primitive.ObjectID) error
	}
	Exercise interface {
		Create(context.Context, *Exercise, primitive.ObjectID) error
		GetAllUserExercises(context.Context, primitive.ObjectID) ([]*Exercise, error)
		GetByID(context.Context, primitive.ObjectID, primitive.ObjectID) (*Exercise, error)
		SearchExerciseByID(context.Context, primitive.ObjectID) (*Exercise, error)
		Update(context.Context, primitive.ObjectID, primitive.ObjectID, map[string]interface{}, int16) error
		Delete(context.Context, primitive.ObjectID, primitive.ObjectID) error
	}
	WorkoutSession interface {
		Create(context.Context, *WorkoutSession, primitive.ObjectID) error
		CreateFromRoutine(context.Context, primitive.ObjectID, primitive.ObjectID) (*WorkoutSession, error)
		GetAllUserSessions(context.Context, primitive.ObjectID) ([]*WorkoutSession, error)
		GetByID(context.Context, primitive.ObjectID, primitive.ObjectID) (*WorkoutSession, error)
		AddSetToExercise(context.Context, primitive.ObjectID, primitive.ObjectID, primitive.ObjectID, SessionSet) error
		CompleteWorkout(context.Context, primitive.ObjectID, primitive.ObjectID) error
		Delete(context.Context, primitive.ObjectID, primitive.ObjectID) error
	}
}

func NewMongoDBStorage(db *mongo.Database) Storage {
	return Storage{
		Users:          &UserStore{db},
		Routine:        &RoutineStore{db},
		Exercise:       &ExerciseStore{db},
		WorkoutSession: &WorkoutSessionStore{db},
	}
}
