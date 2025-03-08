package main

import (
	"github.com/FaustCelaj/GetFit.git/internal/store"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (app *application) createRoutineHandler(c *fiber.Ctx) error {
	userID := c.Params("userID")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "userID is required",
		})
	}

	// Convert userID to a primitive.ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid userID format",
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
	routine.UserID = userObjectID

	// Call the Create method in RoutineStore
	err = app.store.Routine.Create(c.Context(), &routine, userObjectID)
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

// fetch all routines from a specified user
// returns an array of routines
func (app *application) getAllUserRoutinesIDHandler(c *fiber.Ctx) error {
	userID := c.Params("userID")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "userID is required",
		})
	}

	// Convert userID to a primitive.ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid userID format",
		})
	}

	// Call the storage layer to retrieve routines
	routines, err := app.store.Routine.GetAllUserRoutines(c.Context(), userObjectID)
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

// fetch routine by routine ID (returns a single routine)
func (app *application) getRoutineByIDHandler(c *fiber.Ctx) error {
	userIDStr := c.Params("userID")
	routineIDStr := c.Params("routineID")

	// Validate and convert userID
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid userID format",
		})
	}

	// Validate and convert routineID
	routineID, err := primitive.ObjectIDFromHex(routineIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid routineID format",
		})
	}

	// Retrieve routine using the storage layer
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

	// Return the routine in the response
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

// update a specific routine by ID and user ID for validation
func (app *application) patchRoutineHandler(c *fiber.Ctx) error {
	userIDStr := c.Params("userID")
	routineIDStr := c.Params("routineID")

	// Validate and convert userID
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid userID format",
		})
	}

	// Validate and convert routineID
	routineID, err := primitive.ObjectIDFromHex(routineIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid routineID format",
		})
	}

	var payload updateRoutinePayload
	// validate payload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	// validate version
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
	userIDStr := c.Params("userID")
	routineIDStr := c.Params("routineID")

	// Validate and convert userID
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid userID format",
		})
	}

	// Validate and convert routineID
	routineID, err := primitive.ObjectIDFromHex(routineIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid routineID format",
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
