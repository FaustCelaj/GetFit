package store

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// SetWithMetadata holds the metadata for sets, such as the user, exercise, and the list of sets
type SetWithMetadata struct {
	SetID      primitive.ObjectID `bson:"_id" json:"set_id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	ExerciseID primitive.ObjectID `bson:"exercise_id" json:"exercise_id"`
	Sets       []SetDetails       `bson:"sets" json:"sets"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

// SetDetails holds the individual set details such as weight, reps, and set number
type SetDetails struct {
	Weight    float32 `bson:"weight" json:"weight"`
	Reps      int16   `bson:"reps" json:"reps"`
	SetNumber int16   `bson:"set_number" json:"set_number"`
}

type SetStore struct {
	db *mongo.Database
}

const setCollection = "set"

// Add inserts a single set into the database
func (s *SetStore) Add(ctx context.Context, set *SetWithMetadata, exerciseID, userID primitive.ObjectID) error {
	// Set the ID, exerciseID, userID, timestamps, and set number
	set.SetID = primitive.NewObjectID()
	set.ExerciseID = exerciseID
	set.UserID = userID
	set.CreatedAt = time.Now()
	set.UpdatedAt = time.Now()

	// If SetNumber isn't assigned for any set in the list, assign it here (defaults to 1)
	if len(set.Sets) > 0 && set.Sets[0].SetNumber == 0 {
		for i := range set.Sets {
			set.Sets[i].SetNumber = int16(i + 1)
		}
	}

	// Create a context with a timeout for the database operation
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Insert the single set into the database
	_, err := s.db.Collection(setCollection).InsertOne(ctx, set)
	if err != nil {
		return fmt.Errorf("failed to add set: %w", err)
	}

	return nil
}

// AddMultiple inserts multiple sets into one document in the database
func (s *SetStore) AddMultiple(ctx context.Context, setWithMetadata SetWithMetadata, exerciseID, userID primitive.ObjectID) error {
	// Create a context with a timeout for the database operation
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Assign metadata fields for the document
	setWithMetadata.SetID = primitive.NewObjectID()
	setWithMetadata.ExerciseID = exerciseID
	setWithMetadata.UserID = userID
	setWithMetadata.CreatedAt = time.Now()
	setWithMetadata.UpdatedAt = time.Now()

	// Assign incremental set numbers
	for i := range setWithMetadata.Sets {
		setWithMetadata.Sets[i].SetNumber = int16(i + 1)
	}

	// Insert the single document into the database
	_, err := s.db.Collection(setCollection).InsertOne(ctx, setWithMetadata)
	if err != nil {
		return fmt.Errorf("failed to add sets: %w", err)
	}

	return nil
}

// GetAll fetches all sets for a specific exercise and user
func (s *SetStore) GetAll(ctx context.Context, exerciseID, userID primitive.ObjectID) ([]SetWithMetadata, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"exercise_id": exerciseID, "user_id": userID}
	cursor, err := s.db.Collection(setCollection).Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sets: %w", err)
	}
	defer cursor.Close(ctx)

	var sets []SetWithMetadata
	if err := cursor.All(ctx, &sets); err != nil {
		return nil, fmt.Errorf("failed to decode sets: %w", err)
	}

	return sets, nil
}

// GetByID fetches a single set by its ID
func (s *SetStore) GetByID(ctx context.Context, setID primitive.ObjectID) (*SetWithMetadata, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": setID}
	var set SetWithMetadata
	err := s.db.Collection(setCollection).FindOne(ctx, filter).Decode(&set)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("set not found")
		}
		return nil, fmt.Errorf("failed to fetch set: %w", err)
	}

	return &set, nil
}

// Update modifies a specific set's details
func (s *SetStore) Update(ctx context.Context, setID primitive.ObjectID, updates map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	updates["updated_at"] = time.Now()

	filter := bson.M{"_id": setID}
	update := bson.M{"$set": updates}

	result, err := s.db.Collection(setCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update set: %w", err)
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("no set found with ID %s", setID.Hex())
	}

	return nil
}

// UpdateMultiple modifies multiple sets
func (s *SetStore) UpdateMultiple(ctx context.Context, setID, userID primitive.ObjectID, sets []SetDetails) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var bulkOps []mongo.WriteModel

	// Loop through each set and prepare the bulk operations
	for _, set := range sets {
		updateDoc := bson.M{
			"$set": bson.M{
				"weight":     set.Weight,
				"reps":       set.Reps,
				"set_number": set.SetNumber,
				"updated_at": time.Now(),
			},
		}

		// Generate the filter for the specific set based on set_id, user_id, and exercise_id
		filter := bson.M{
			"_id":     setID,
			"user_id": userID,
		}

		// Create an update operation
		bulkOps = append(bulkOps, mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(updateDoc))
	}

	// Execute the bulk write operation
	if _, err := s.db.Collection(setCollection).BulkWrite(ctx, bulkOps); err != nil {
		return fmt.Errorf("failed to update multiple sets: %w", err)
	}

	return nil
}

// Delete removes a specific set by its ID
func (s *SetStore) Delete(ctx context.Context, setID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": setID}
	result, err := s.db.Collection(setCollection).DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete set: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no set found with ID %s", setID.Hex())
	}

	return nil
}
