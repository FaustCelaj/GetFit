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

	// Validate required fields
	if exercise.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "exercise name is required",
		})
	}
	if exercise.Category == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "exercise category is required",
		})
	}

	// Ensure nullable fields are explicitly set to nil if not provided
	if exercise.Force == nil {
		exercise.Force = nil // Explicitly set to nil
	}
	if exercise.Level == nil {
		exercise.Level = nil // Explicitly set to nil
	}
	if exercise.Mechanic == nil {
		exercise.Mechanic = nil // Explicitly set to nil
	}
	if exercise.Equipment == nil {
		exercise.Equipment = nil // Explicitly set to nil
	}
	if exercise.PrimaryMuscles == nil {
		exercise.PrimaryMuscles = &[]string{} // Default to empty array
	}
	if exercise.SecondaryMuscles == nil {
		exercise.SecondaryMuscles = &[]string{} // Default to empty array
	}
	if exercise.Instructions == nil {
		exercise.Instructions = &[]string{} // Default to empty array
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

	// Parse and validate payload
	var payload updateExercisePayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// // Validate that all fields are present in the payload
	// if payload.Name == nil || payload.Force == nil || payload.Level == nil ||
	// 	payload.Mechanic == nil || payload.Equipment == nil || payload.PrimaryMuscles == nil ||
	// 	payload.SecondaryMuscles == nil || payload.Instructions == nil || payload.Category == nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"error": "All fields are required, even if null",
	// 	})
	// }

	// Validate required fields (Name and Category cannot be nil or empty)
	if payload.Name == nil || *payload.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name cannot be empty",
		})
	}
	if payload.Category == nil || *payload.Category == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Category cannot be empty",
		})
	}

	if payload.ExpectedVersion == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Current version is required",
		})
	}

	// Build updates map
	updates := map[string]interface{}{
		"name":             *payload.Name,
		"force":            *payload.Force,
		"level":            *payload.Level,
		"mechanic":         *payload.Mechanic,
		"equipment":        *payload.Equipment,
		"primaryMuscles":   *payload.PrimaryMuscles,
		"secondaryMuscles": *payload.SecondaryMuscles,
		"instructions":     *payload.Instructions,
		"category":         *payload.Category,
	}

	// Perform the update in the database
	if err := app.store.Exercise.Update(c.Context(), exerciseID, userID, updates, payload.ExpectedVersion); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update exercise",
		})
	}

	// Return success response
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
