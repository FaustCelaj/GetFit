package main

import (
	"github.com/FaustCelaj/GetFit.git/internal/store"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (app *application) createExerciseHandler(c *fiber.Ctx) error {
	userID := c.Params("userID")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "userID is required",
		})
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"details": err.Error(),
			"error":   "invalid userID format",
		})
	}

	var exercise store.Exercise
	if err := c.BodyParser(&exercise); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"details": err.Error(),
			"error":   "invalid request body",
		})
	}

	exercise.UserID = userObjectID

	if err := app.store.Exercise.Create(c.Context(), &exercise, userObjectID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"details": err.Error(),
			"error":   "Failed to create exercise",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":       exercise.ID.Hex(),
		"message":  "exercise created successfully",
		"exercise": exercise,
	})
}

func (app *application) getAllUserExercisesHandler(c *fiber.Ctx) error {
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

	exercises, err := app.store.Exercise.GetAllUserExercises(c.Context(), userObjectID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to fetch exercises",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "exercises retrieved successfully",
		"exercises": exercises,
	})
}

func (app *application) getExerciseByIDHandler(c *fiber.Ctx) error {
	userIDstr := c.Params("userID")
	exerciseIDstr := c.Params("exerciseID")

	userID, err := primitive.ObjectIDFromHex(userIDstr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid userID format",
		})
	}

	// Validate and convert routineID
	exerciseID, err := primitive.ObjectIDFromHex(exerciseIDstr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid routineID format",
		})
	}

	exercise, err := app.store.Exercise.GetByID(c.Context(), exerciseID, userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Exercise not found or does not belong to the user",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch exercise",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "Routine retrieved successfully",
		"exercise": exercise,
	})
}

type updateExercisePayload struct {
	Name             *string   `json:"name"`
	Force            *string   `json:"force"`
	Level            *string   `json:"level"`
	Mechanic         *string   `json:"mechanic"`
	Equipment        *string   `json:"equipment"`
	PrimaryMuscles   *[]string `json:"primaryMuscles"`
	SecondaryMuscles *[]string `json:"secondaryMuscles"`
	Instructions     *[]string `json:"instructions"`
	Category         *string   `json:"category"`
	ExpectedVersion  int16     `json:"expected_version"`
}

func (app *application) updateExerciseHandler(c *fiber.Ctx) error {
	userIDStr := c.Params("userID")
	exerciseIDStr := c.Params("exerciseID")

	// Validate and convert userID
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid userID format",
		})
	}

	// Validate and convert exerciseID
	exerciseID, err := primitive.ObjectIDFromHex(exerciseIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid exerciseID format",
		})
	}

	var payload updateExercisePayload
	// Validate payload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate version
	if payload.ExpectedVersion == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Current version is required",
		})
	}

	updates := make(map[string]interface{})
	if payload.Name != nil {
		if *payload.Name == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Name cannot be empty",
			})
		}
		updates["name"] = *payload.Name
	}
	if payload.Force != nil {
		updates["force"] = *payload.Force
	}
	if payload.Level != nil {
		updates["level"] = *payload.Level
	}
	if payload.Mechanic != nil {
		updates["mechanic"] = *payload.Mechanic
	}
	if payload.Equipment != nil {
		updates["equipment"] = *payload.Equipment
	}
	if payload.PrimaryMuscles != nil {
		updates["primaryMuscles"] = *payload.PrimaryMuscles
	}
	if payload.SecondaryMuscles != nil {
		updates["secondaryMuscles"] = *payload.SecondaryMuscles
	}
	if payload.Instructions != nil {
		updates["instructions"] = *payload.Instructions
	}
	if payload.Category != nil {
		if *payload.Category == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Category cannot be empty",
			})
		}
		updates["category"] = *payload.Category
	}

	if len(updates) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No fields to update",
		})
	}

	// Perform the update in the database
	if err := app.store.Exercise.Update(c.Context(), exerciseID, userID, updates, payload.ExpectedVersion); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update exercise",
		})
	}

	// Return a success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Exercise updated successfully",
	})
}

func (app *application) deleteExerciseHandler(c *fiber.Ctx) error {
	userIDStr := c.Params("userID")
	exerciseIDStr := c.Params("exerciseID")

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"details": err.Error(),
			"error":   "Invalid userID format",
		})
	}

	exerciseID, err := primitive.ObjectIDFromHex(exerciseIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"details": err.Error(),
			"error":   "Invalid exerciseID format",
		})
	}

	if err := app.store.Exercise.Delete(c.Context(), exerciseID, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"details": err.Error(),
			"error":   "Failed to delete exercise",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Exercise was successfully deleted",
	})
}
