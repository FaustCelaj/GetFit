package store

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Storage struct {
	Exercises interface {
		Create(context.Context, *Exercise) error
	}
	Users interface {
		Create(context.Context, *User) error
	}
}

func NewMongoDBStorage(db *mongo.Database) Storage {
	return Storage{
		Exercises: &ExerciseStore{db},
		Users:     &UserStore{db},
	}
}
