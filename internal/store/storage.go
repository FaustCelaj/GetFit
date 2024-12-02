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
		GetById(context.Context, primitive.ObjectID) (*User, error)
		Update(context.Context, primitive.ObjectID, map[string]interface{}) error
		Delete(context.Context, primitive.ObjectID) error
	}
	RoutineExercise interface {
		Create(context.Context, *RoutineExercise, primitive.ObjectID, primitive.ObjectID) error
		GetByID(context.Context, primitive.ObjectID) (*RoutineExercise, error)
	}
	Exercises interface {
		Create(context.Context, *Exercise) error
	}
	ExerciseSet interface {
		AddSet(context.Context, *ExerciseSet, primitive.ObjectID) error
		AddMultipleSet(context.Context, []ExerciseSet, primitive.ObjectID) error
	}
}

func NewMongoDBStorage(db *mongo.Database) Storage {
	return Storage{
		Users:           &UserStore{db},
		RoutineExercise: &RoutineExerciseStore{db},
		Exercises:       &ExerciseStore{db},
		ExerciseSet:     &ExerciseSetStore{db},
	}
}
