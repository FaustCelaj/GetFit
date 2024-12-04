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
		GetByUserID(context.Context, primitive.ObjectID) ([]*Routine, error)
		GetByIDAndUserID(context.Context, primitive.ObjectID, primitive.ObjectID) (*Routine, error)
		Update(context.Context, primitive.ObjectID, primitive.ObjectID, map[string]interface{}, int16) error
		Delete(context.Context, primitive.ObjectID, primitive.ObjectID) error
	}
	Exercises interface {
		Create(context.Context, *Exercise) error
	}
	Set interface {
		AddSet(context.Context, *Set, primitive.ObjectID, primitive.ObjectID) error
		AddMultipleSet(context.Context, []Set, primitive.ObjectID, primitive.ObjectID) error
	}
}

func NewMongoDBStorage(db *mongo.Database) Storage {
	return Storage{
		Users:     &UserStore{db},
		Routine:   &RoutineStore{db},
		Exercises: &ExerciseStore{db},
		Set:       &SetStore{db},
	}
}
