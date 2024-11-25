package main

import (
	"fmt"

	"github.com/FaustCelaj/GetFit.git/internal/store"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// get a routine exercise by ID
func (app *application) getRoutineExerciseByIDHandler(c *fiber.Ctx) error {
	// Extract routineExerciseID from the URL parameter
	routineExerciseIDStr := c.Params("routineExerciseID")
	routineExerciseID, err := primitive.ObjectIDFromHex(routineExerciseIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid routineExercise ID format",
		})
	}

	// Call the storage layer to retrieve the RoutineExercise
	routineExercise, err := app.store.RoutineExercise.GetByID(c.Context(), routineExerciseID)
	if err != nil {
		if err.Error() == mongo.ErrNoDocuments.Error() {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "RoutineExercise not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch routineExercise",
		})
	}

	// Return the RoutineExercise in the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":          "RoutineExercise retrieved successfully",
		"routine_exercise": routineExercise,
	})
}

// create a routineExercise
func (app *application) createRoutineExerciseHandler(c *fiber.Ctx) error {
	// Extract userID from the URL
	userIDStr := c.Params("userID")
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	// Parse request body to get the exerciseID and notes
	type requestPayload struct {
		ExerciseID string  `json:"exercise_id"`
		Notes      *string `json:"notes,omitempty"`
	}

	var payload requestPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Convert exerciseID to ObjectID
	exerciseID, err := primitive.ObjectIDFromHex(payload.ExerciseID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid exercise ID format",
		})
	}

	// Create a new RoutineExercise
	routineExercise := &store.RoutineExercise{
		Notes: payload.Notes,
	}

	err = app.store.RoutineExercise.Create(c.Context(), routineExercise, userID, exerciseID)
	if err != nil {
		// Handle specific errors
		if err.Error() == fmt.Sprintf("duplicate routine exercise: userID %s, exerciseID %s already exists", userID.Hex(), exerciseID.Hex()) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// For any other error, return a generic 500 response
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create routine exercise",
		})
	}

	// Return a success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Routine exercise created successfully",
		"id":      routineExercise.ID.Hex(),
	})
}

// add a set/sets
func (app *application) addExerciseSetHandler(c *fiber.Ctx) error {
	// Extract routineExerciseID from URL parameters
	routineExerciseIDStr := c.Params("routineExerciseID")
	routineExerciseID, err := primitive.ObjectIDFromHex(routineExerciseIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid routineExerciseID format",
		})
	}

	// Parse the request body to support both single and multiple sets
	var sets []store.ExerciseSet
	if err := c.BodyParser(&sets); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Validate that at least one set is provided
	if len(sets) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "At least one set must be provided",
		})
	}

	// Handle single set addition
	if len(sets) == 1 {
		err := app.store.ExerciseSet.AddSet(c.Context(), &sets[0], routineExerciseID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to add set: %v", err),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Set added successfully",
			"set_id":  sets[0].ID.Hex(),
		})
	}

	// Handle multiple set additions
	err = app.store.ExerciseSet.AddMultipleSet(c.Context(), sets, routineExerciseID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to add sets: %v", err),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": fmt.Sprintf("%d sets added successfully", len(sets)),
	})
}

// update a routine exercise
func (app *application) updateRoutineExerciseHandler(c *fiber.Ctx) error {
	return nil
}
