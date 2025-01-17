package main

import (
	"github.com/FaustCelaj/GetFit.git/internal/store"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Add sets to an exercise
func (app *application) addSetsHandler(c *fiber.Ctx) error {
	userIDStr := c.Params("userID")
	exerciseIDStr := c.Params("exerciseID")

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid userID format",
		})
	}

	exerciseID, err := primitive.ObjectIDFromHex(exerciseIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid exerciseID format",
		})
	}

	// Parse the request body into a SetWithMetadata struct
	var setMetadata store.SetWithMetadata
	if err := c.BodyParser(&setMetadata); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Ensure the set metadata is valid
	if len(setMetadata.Sets) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No sets provided in the request body",
		})
	}

	// Add the set(s)
	if len(setMetadata.Sets) == 1 {
		err = app.store.Set.Add(c.Context(), &setMetadata, exerciseID, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to create set",
				"details": err.Error(),
			})
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Set created successfully",
			"set":     setMetadata,
		})
	}

	// Add multiple sets
	err = app.store.Set.AddMultiple(c.Context(), setMetadata, exerciseID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create sets",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Sets created successfully",
		"sets":    setMetadata.Sets,
	})
}

// Fetch all sets for an exercise
func (app *application) getAllExerciseSetsHandler(c *fiber.Ctx) error {
	userIDStr := c.Params("userID")
	exerciseIDStr := c.Params("exerciseID")

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid userID format",
		})
	}

	exerciseID, err := primitive.ObjectIDFromHex(exerciseIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid exerciseID format",
		})
	}

	sets, err := app.store.Set.GetAll(c.Context(), exerciseID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to fetch sets",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Sets retrieved successfully",
		"sets":    sets,
	})
}

// Fetch a single set by ID
func (app *application) getSetByIDHandler(c *fiber.Ctx) error {
	setIDStr := c.Params("setID")

	setID, err := primitive.ObjectIDFromHex(setIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid setID format",
		})
	}

	set, err := app.store.Set.GetByID(c.Context(), setID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Set not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch set",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Set retrieved successfully",
		"set":     set,
	})
}

// Update a set by ID
func (app *application) updateSetHandler(c *fiber.Ctx) error {
	setIDStr := c.Params("setID")

	setID, err := primitive.ObjectIDFromHex(setIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid setID format",
		})
	}

	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := app.store.Set.Update(c.Context(), setID, updates); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update set",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Set updated successfully",
	})
}

// Delete a set by ID
func (app *application) deleteSetHandler(c *fiber.Ctx) error {
	setIDStr := c.Params("setID")

	setID, err := primitive.ObjectIDFromHex(setIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid setID format",
		})
	}

	if err := app.store.Set.Delete(c.Context(), setID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete set",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Set deleted successfully",
	})
}
