package main

import (
	"time"

	"github.com/FaustCelaj/GetFit.git/internal/store"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// "/workout"
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

// "/from-routine/:routineID"
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

// "/:sessionID"
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

// "/exercise/:exerciseID/sets"
// adding a set
type addSetPayload struct {
	Weight    float32 `json:"weight"`
	Reps      int16   `json:"reps"`
	SetNumber int16   `json:"set_number"`
}

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
