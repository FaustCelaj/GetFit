package main

import (
	"github.com/FaustCelaj/GetFit.git/internal/store"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateRoutine godoc
//
//	@Summary		Create a new workout routine
//	@Description	Create a structured workout routine with exercise templates
//	@Tags			routines
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		string			true	"User ID"
//	@Param			routine	body		store.Routine	true	"Routine information"
//	@Success		201		{object}	string			"Routine created successfully"
//	@Failure		400		{object}	error			"Invalid request body or missing fields"
//	@Failure		500		{object}	error			"Failed to create routine"
//
// @Security		ApiKeyAuth
//
//	@Router			/users/{userID}/routine [post]
func (app *application) createRoutineHandler(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	if userID == primitive.NilObjectID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "userID not found in context",
		})
	}

	var routine store.Routine
	if err := c.BodyParser(&routine); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Validate the routine
	if routine.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "routine title is required",
		})
	}

	// Set the userID for the routine
	routine.UserID = userID

	// Call the Create method in RoutineStore
	err := app.store.Routine.Create(c.Context(), &routine, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to create routine",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "routine created successfully",
		"routine": routine,
	})
}

// GetAllRoutines godoc
//
//	@Summary		Get all user routines
//	@Description	Retrieve all workout routines created by a user
//	@Tags			routines
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		string			true	"User ID"
//	@Success		200		{array}		store.Routine	"List of routines"
//	@Failure		400		{object}	error			"Invalid user ID"
//	@Failure		500		{object}	error			"Failed to fetch routines"
//
// @Security		ApiKeyAuth
//
//	@Router			/users/{userID}/routine [get]
func (app *application) getAllUserRoutinesIDHandler(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	if userID == primitive.NilObjectID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "userID not found in context",
		})
	}

	routines, err := app.store.Routine.GetAllUserRoutines(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to fetch routines",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "routines retrieved successfully",
		"routines": routines,
	})
}

// GetRoutineByID godoc
//
//	@Summary		Get routine by ID
//	@Description	Retrieve a specific workout routine by its ID
//	@Tags			routines
//	@Accept			json
//	@Produce		json
//	@Param			userID		path		string			true	"User ID"
//	@Param			routineID	path		string			true	"Routine ID"
//	@Success		200			{object}	store.Routine	"Routine information"
//	@Failure		400			{object}	error			"Invalid ID format"
//	@Failure		404			{object}	error			"Routine not found"
//	@Failure		500			{object}	error			"Failed to fetch routine"
//
// @Security		ApiKeyAuth
//
//	@Router			/users/{userID}/routine/{routineID} [get]
func (app *application) getRoutineByIDHandler(c *fiber.Ctx) error {
	userID, routineID := getUserIDFromContext(c), getRoutineIDFromContext(c)
	if userID == primitive.NilObjectID || routineID == primitive.NilObjectID {
		missingID := "userID"
		if routineID == primitive.NilObjectID {
			missingID = "sessionID"
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": missingID + " not found in context",
		})
	}

	routine, err := app.store.Routine.GetByID(c.Context(), routineID, userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Routine not found or does not belong to the user",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch routine",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Routine retrieved successfully",
		"routine": routine,
	})
}

type updateRoutinePayload struct {
	Title           *string   `json:"title"`
	Description     *string   `json:"description"`
	ExerciseID      *[]string `json:"exercise_id"`
	ExpectedVersion int16     `json:"expected_version"`
}

// UpdateRoutine godoc
//
//	@Summary		Update a routine
//	@Description	Update details of an existing workout routine
//	@Tags			routines
//	@Accept			json
//	@Produce		json
//	@Param			userID		path		string					true	"User ID"
//	@Param			routineID	path		string					true	"Routine ID"
//	@Param			routine		body		updateRoutinePayload	true	"Updated routine information"
//	@Success		200			{object}	string					"Routine updated successfully"
//	@Failure		400			{object}	error					"Invalid request body or missing fields"
//	@Failure		404			{object}	error					"Routine not found"
//	@Failure		500			{object}	error					"Failed to update routine"
//
// @Security		ApiKeyAuth
//
//	@Router			/users/{userID}/routine/{routineID} [patch]
func (app *application) patchRoutineHandler(c *fiber.Ctx) error {
	userID, routineID := getUserIDFromContext(c), getRoutineIDFromContext(c)
	if userID == primitive.NilObjectID || routineID == primitive.NilObjectID {
		missingID := "userID"
		if routineID == primitive.NilObjectID {
			missingID = "sessionID"
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": missingID + " not found in context",
		})
	}

	var payload updateRoutinePayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if payload.ExpectedVersion == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Current version is required",
		})
	}

	updates := make(map[string]interface{})

	if payload.Title != nil {
		updates["title"] = &payload.Title
	}
	if payload.Description != nil {
		updates["description"] = &payload.Description
	}

	if requiredField, ok := updates["requiredField"]; !ok || requiredField == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "RequiredField cannot be empty",
		})
	}

	if len(updates) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No fields to update",
		})
	}

	// Perform the update in the database
	if err := app.store.Routine.Update(c.Context(), routineID, userID, updates, payload.ExpectedVersion); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to update routine",
			"details": err.Error(),
		})
	}

	// Return a success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Routine updated successfully",
	})
}

// DeleteRoutine godoc
//
//	@Summary		Delete a routine
//	@Description	Remove a workout routine from the system
//	@Tags			routines
//	@Accept			json
//	@Produce		json
//	@Param			userID		path		string	true	"User ID"
//	@Param			routineID	path		string	true	"Routine ID"
//	@Success		200			{object}	string	"Routine successfully deleted"
//	@Failure		400			{object}	error	"Invalid ID format"
//	@Failure		500			{object}	error	"Failed to delete routine"
//
// @Security		ApiKeyAuth
//
//	@Router			/users/{userID}/routine/{routineID} [delete]
func (app *application) deleteRoutineHandler(c *fiber.Ctx) error {
	userID, routineID := getUserIDFromContext(c), getRoutineIDFromContext(c)
	if userID == primitive.NilObjectID || routineID == primitive.NilObjectID {
		missingID := "userID"
		if routineID == primitive.NilObjectID {
			missingID = "sessionID"
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": missingID + " not found in context",
		})
	}

	if err := app.store.Routine.Delete(c.Context(), routineID, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to delete routine",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Routine was successfully deleted",
	})
}
