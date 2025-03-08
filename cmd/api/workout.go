package main

import (
	"time"

	"github.com/FaustCelaj/GetFit.git/internal/store"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Create a new blank workout session
func (app *application) createWorkoutSessionHandler(c *fiber.Ctx) error {
	userID := c.Params("userID")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "userID is required",
		})
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid userID format",
		})
	}

	var session store.WorkoutSession
	if err := c.BodyParser(&session); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Validate the session
	if session.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "workout title is required",
		})
	}

	// Set the userID for the session
	session.UserID = userObjectID

	// Call the Create method
	err = app.store.WorkoutSession.Create(c.Context(), &session, userObjectID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to create workout session",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "workout session created successfully",
		"session": session,
	})
}

// Create a workout session from a routine template
func (app *application) createWorkoutFromRoutineHandler(c *fiber.Ctx) error {
	userID := c.Params("userID")
	routineID := c.Params("routineID")

	if userID == "" || routineID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "userID and routineID are required",
		})
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid userID format",
		})
	}

	routineObjectID, err := primitive.ObjectIDFromHex(routineID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid routineID format",
		})
	}

	// Create a workout session from the routine
	session, err := app.store.WorkoutSession.CreateFromRoutine(c.Context(), routineObjectID, userObjectID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to create workout from routine",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "workout created from routine successfully",
		"session": session,
	})
}

// Get all workout sessions for a user
func (app *application) getAllWorkoutSessionsHandler(c *fiber.Ctx) error {
	userID := c.Params("userID")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "userID is required",
		})
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid userID format",
		})
	}

	sessions, err := app.store.WorkoutSession.GetAllUserSessions(c.Context(), userObjectID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to fetch workout sessions",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "workout sessions retrieved successfully",
		"sessions": sessions,
	})
}

// Get a specific workout session
func (app *application) getWorkoutSessionByIDHandler(c *fiber.Ctx) error {
	userID := c.Params("userID")
	sessionID := c.Params("sessionID")

	if userID == "" || sessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "userID and sessionID are required",
		})
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid userID format",
		})
	}

	sessionObjectID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid sessionID format",
		})
	}

	session, err := app.store.WorkoutSession.GetByID(c.Context(), sessionObjectID, userObjectID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "workout session not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch workout session",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "workout session retrieved successfully",
		"session": session,
	})
}

// Add a set to an exercise in a workout
type addSetPayload struct {
	Weight    float32 `json:"weight"`
	Reps      int16   `json:"reps"`
	SetNumber int16   `json:"set_number"`
}

func (app *application) addSetToWorkoutHandler(c *fiber.Ctx) error {
	userID := c.Params("userID")
	sessionID := c.Params("sessionID")
	exerciseID := c.Params("exerciseID")

	if userID == "" || sessionID == "" || exerciseID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "userID, sessionID, and exerciseID are required",
		})
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid userID format",
		})
	}

	sessionObjectID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid sessionID format",
		})
	}

	exerciseObjectID, err := primitive.ObjectIDFromHex(exerciseID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid exerciseID format",
		})
	}

	var payload addSetPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Create the session set
	set := store.SessionSet{
		Weight:      payload.Weight,
		Reps:        payload.Reps,
		SetNumber:   payload.SetNumber,
		CompletedAt: time.Now(),
	}

	// Add the set to the exercise in the workout
	err = app.store.WorkoutSession.AddSetToExercise(c.Context(), sessionObjectID, userObjectID, exerciseObjectID, set)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to add set to workout",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "set added to workout successfully",
		"set":     set,
	})
}

// Complete a workout session
func (app *application) completeWorkoutSessionHandler(c *fiber.Ctx) error {
	userID := c.Params("userID")
	sessionID := c.Params("sessionID")

	if userID == "" || sessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "userID and sessionID are required",
		})
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid userID format",
		})
	}

	sessionObjectID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid sessionID format",
		})
	}

	err = app.store.WorkoutSession.CompleteWorkout(c.Context(), sessionObjectID, userObjectID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to complete workout",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "workout completed successfully",
	})
}

// Delete a workout session
func (app *application) deleteWorkoutSessionHandler(c *fiber.Ctx) error {
	userID := c.Params("userID")
	sessionID := c.Params("sessionID")

	if userID == "" || sessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "userID and sessionID are required",
		})
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid userID format",
		})
	}

	sessionObjectID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid sessionID format",
		})
	}

	err = app.store.WorkoutSession.Delete(c.Context(), sessionObjectID, userObjectID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to delete workout session",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "workout session deleted successfully",
	})
}
