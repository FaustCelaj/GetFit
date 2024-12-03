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
	if len(routine.ExerciseID) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "a routine must contain at least one exercise",
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

func (app *application) getRoutinesByUserIDHandler(c *fiber.Ctx) error {
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
	routines, err := app.store.Routine.GetByUserID(c.Context(), userObjectID)
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
	routine, err := app.store.Routine.GetByIDAndUserID(c.Context(), routineID, userID)
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
			"error": "Failed to delete routine",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Routine was successfully deleted",
	})
}
