package main

import (
	"github.com/FaustCelaj/GetFit.git/internal/store"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// "/routine"
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

// "/:routineID"
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
