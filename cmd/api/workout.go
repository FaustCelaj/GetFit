package main

import (
	"time"

	"github.com/FaustCelaj/GetFit.git/internal/store"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateWorkoutSession godoc
//
//	@Summary		Create a new workout session
//	@Description	Start a new workout session from scratch
//	@Tags			workouts
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		string					true	"User ID"
//	@Param			session	body		store.WorkoutSession	true	"Workout session information"
//	@Success		201		{object}	string					"Workout session created successfully"
//	@Failure		400		{object}	error					"Invalid request body or missing fields"
//	@Failure		500		{object}	error					"Failed to create workout session"
//
// @Security		ApiKeyAuth
//
//	@Router			/users/{userID}/workout [post]
func (app *application) createWorkoutSessionHandler(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	if userID == primitive.NilObjectID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "userID not found in context",
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
	session.UserID = userID

	// Call the Create method
	err := app.store.WorkoutSession.Create(c.Context(), &session, userID)
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

// GetAllWorkoutSessions godoc
//
//	@Summary		Get all workout sessions
//	@Description	Retrieve all workout sessions for a user
//	@Tags			workouts
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		string					true	"User ID"
//	@Success		200		{array}		store.WorkoutSession	"List of workout sessions"
//	@Failure		400		{object}	error					"Invalid user ID"
//	@Failure		500		{object}	error					"Failed to fetch workout sessions"
//
// @Security		ApiKeyAuth
//
//	@Router			/users/{userID}/workout [get]
func (app *application) getAllWorkoutSessionsHandler(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	if userID == primitive.NilObjectID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "userID not found in context",
		})
	}

	sessions, err := app.store.WorkoutSession.GetAllUserSessions(c.Context(), userID)
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

// CreateWorkoutFromRoutine godoc
//
//	@Summary		Create workout from routine
//	@Description	Create a new workout session based on an existing routine
//	@Tags			workouts
//	@Accept			json
//	@Produce		json
//	@Param			userID		path		string	true	"User ID"
//	@Param			routineID	path		string	true	"Routine ID"
//	@Success		201			{object}	string	"Workout created from routine successfully"
//	@Failure		400			{object}	error	"Invalid IDs"
//	@Failure		500			{object}	error	"Failed to create workout from routine"
//
// @Security		ApiKeyAuth
//
//	@Router			/users/{userID}/workout/from-routine/{routineID} [post]
func (app *application) createWorkoutFromRoutineHandler(c *fiber.Ctx) error {
	userID, routineID := getUserIDFromContext(c), getRoutineIDFromContext(c)
	if userID == primitive.NilObjectID || routineID == primitive.NilObjectID {
		missingID := "userID"
		if routineID == primitive.NilObjectID {
			missingID = "routineID"
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": missingID + " not found in context",
		})
	}

	// Create a workout session from the routine
	session, err := app.store.WorkoutSession.CreateFromRoutine(c.Context(), routineID, userID)
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

// GetWorkoutSessionByID godoc
//
//	@Summary		Get workout session by ID
//	@Description	Retrieve a specific workout session by its ID
//	@Tags			workouts
//	@Accept			json
//	@Produce		json
//	@Param			userID		path		string					true	"User ID"
//	@Param			sessionID	path		string					true	"Session ID"
//	@Success		200			{object}	store.WorkoutSession	"Workout session information"
//	@Failure		400			{object}	error					"Invalid ID format"
//	@Failure		404			{object}	error					"Workout session not found"
//	@Failure		500			{object}	error					"Failed to fetch workout session"
//
// @Security		ApiKeyAuth
//
//	@Router			/users/{userID}/workout/{sessionID} [get]
func (app *application) getWorkoutSessionByIDHandler(c *fiber.Ctx) error {
	userID, sessionID := getUserIDFromContext(c), getSessionIDFromContext(c)
	if userID == primitive.NilObjectID || sessionID == primitive.NilObjectID {
		missingID := "userID"
		if sessionID == primitive.NilObjectID {
			missingID = "sessionID"
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": missingID + " not found in context",
		})
	}

	session, err := app.store.WorkoutSession.GetByID(c.Context(), sessionID, userID)
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

// CompleteWorkoutSession godoc
//
//	@Summary		Complete a workout session
//	@Description	Mark a workout session as completed and calculate metrics
//	@Tags			workouts
//	@Accept			json
//	@Produce		json
//	@Param			userID		path		string	true	"User ID"
//	@Param			sessionID	path		string	true	"Session ID"
//	@Success		200			{object}	string	"Workout completed successfully"
//	@Failure		400			{object}	error	"Invalid ID format"
//	@Failure		500			{object}	error	"Failed to complete workout"
//
// @Security		ApiKeyAuth
//
//	@Router			/users/{userID}/workout/{sessionID}/complete [post]
func (app *application) completeWorkoutSessionHandler(c *fiber.Ctx) error {
	userID, sessionID := getUserIDFromContext(c), getSessionIDFromContext(c)
	if userID == primitive.NilObjectID || sessionID == primitive.NilObjectID {
		missingID := "userID"
		if sessionID == primitive.NilObjectID {
			missingID = "sessionID"
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": missingID + " not found in context",
		})
	}

	err := app.store.WorkoutSession.CompleteWorkout(c.Context(), sessionID, userID)
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

// DeleteWorkoutSession godoc
//
//	@Summary		Delete a workout session
//	@Description	Remove a workout session from the system
//	@Tags			workouts
//	@Accept			json
//	@Produce		json
//	@Param			userID		path		string	true	"User ID"
//	@Param			sessionID	path		string	true	"Session ID"
//	@Success		200			{object}	string	"Workout session deleted successfully"
//	@Failure		400			{object}	error	"Invalid ID format"
//	@Failure		500			{object}	error	"Failed to delete workout session"
//
// @Security		ApiKeyAuth
//
//	@Router			/users/{userID}/workout/{sessionID} [delete]
func (app *application) deleteWorkoutSessionHandler(c *fiber.Ctx) error {
	userID, sessionID := getUserIDFromContext(c), getSessionIDFromContext(c)
	if userID == primitive.NilObjectID || sessionID == primitive.NilObjectID {
		missingID := "userID"
		if sessionID == primitive.NilObjectID {
			missingID = "sessionID"
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": missingID + " not found in context",
		})
	}

	err := app.store.WorkoutSession.Delete(c.Context(), sessionID, userID)
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

type addSetPayload struct {
	Weight    float32 `json:"weight"`
	Reps      int16   `json:"reps"`
	SetNumber int16   `json:"set_number"`
}

// AddSetToWorkout godoc
//
//	@Summary		Add a set to a workout exercise
//	@Description	Record a completed set for an exercise in a workout
//	@Tags			workout-sets
//	@Accept			json
//	@Produce		json
//	@Param			userID		path		string			true	"User ID"
//	@Param			sessionID	path		string			true	"Session ID"
//	@Param			exerciseID	path		string			true	"Exercise ID"
//	@Param			set			body		addSetPayload	true	"Set information"
//	@Success		200			{object}	string			"Set added to workout successfully"
//	@Failure		400			{object}	error			"Invalid request body or IDs"
//	@Failure		500			{object}	error			"Failed to add set to workout"
//
// @Security		ApiKeyAuth
//
//	@Router			/users/{userID}/workout/{sessionID}/exercise/{exerciseID}/sets [post]
func (app *application) addSetToWorkoutHandler(c *fiber.Ctx) error {
	userID, sessionID, exerciseID := getUserIDFromContext(c), getSessionIDFromContext(c), getExerciseIDFromContext(c)
	if userID == primitive.NilObjectID || sessionID == primitive.NilObjectID || exerciseID == primitive.NilObjectID {
		missingID := "userID"
		if sessionID == primitive.NilObjectID {
			missingID = "sessionID"
		}
		if exerciseID == primitive.NilObjectID {
			missingID = "exerciseID"
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": missingID + " not found in context",
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
	err := app.store.WorkoutSession.AddSetToExercise(c.Context(), sessionID, userID, exerciseID, set)
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
