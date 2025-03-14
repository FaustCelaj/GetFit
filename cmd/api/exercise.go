package main

import (
	"github.com/FaustCelaj/GetFit.git/internal/store"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// "/exercise"
func (app *application) createExerciseHandler(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	if userID == primitive.NilObjectID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "userID not found in context",
		})
	}

	var exercise store.Exercise
	if err := c.BodyParser(&exercise); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"details": err.Error(),
			"error":   "invalid request body",
		})
	}

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

	if exercise.Force == nil {
		exercise.Force = nil
	}
	if exercise.Level == nil {
		exercise.Level = nil
	}
	if exercise.Mechanic == nil {
		exercise.Mechanic = nil
	}
	if exercise.Equipment == nil {
		exercise.Equipment = nil
	}
	if exercise.PrimaryMuscles == nil {
		exercise.PrimaryMuscles = &[]string{}
	}
	if exercise.SecondaryMuscles == nil {
		exercise.SecondaryMuscles = &[]string{}
	}
	if exercise.Instructions == nil {
		exercise.Instructions = &[]string{}
	}

	exercise.UserID = userID

	if err := app.store.Exercise.Create(c.Context(), &exercise, userID); err != nil {
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

func (app *application) getAllUserExerciseHandler(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	if userID == primitive.NilObjectID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "userID not found in context",
		})
	}

	exercises, err := app.store.Exercise.GetAllUserExercises(c.Context(), userID)
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

// "/:exerciseID"
// return custom exercise by user
func (app *application) getExerciseByIDHandler(c *fiber.Ctx) error {
	userID, exerciseID := getUserIDFromContext(c), getExerciseIDFromContext(c)
	if userID == primitive.NilObjectID || exerciseID == primitive.NilObjectID {
		missingID := "userID"
		if exerciseID == primitive.NilObjectID {
			missingID = "exerciseID"
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": missingID + " not found in context",
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
	userID, exerciseID := getUserIDFromContext(c), getExerciseIDFromContext(c)
	if userID == primitive.NilObjectID || exerciseID == primitive.NilObjectID {
		missingID := "userID"
		if exerciseID == primitive.NilObjectID {
			missingID = "exerciseID"
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": missingID + " not found in context",
		})
	}

	var payload updateExercisePayload
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

	if payload.Name != nil {
		updates["name"] = &payload.Name
	}
	if payload.Category != nil {
		updates["category"] = &payload.Category
	}

	if len(updates) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No fields to update",
		})
	}

	if err := app.store.Exercise.Update(c.Context(), exerciseID, userID, updates, payload.ExpectedVersion); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to update routine",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Exercise updated successfully",
	})
}

func (app *application) deleteExerciseHandler(c *fiber.Ctx) error {
	userID, exerciseID := getUserIDFromContext(c), getExerciseIDFromContext(c)
	if userID == primitive.NilObjectID || exerciseID == primitive.NilObjectID {
		missingID := "userID"
		if exerciseID == primitive.NilObjectID {
			missingID = "exerciseID"
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": missingID + " not found in context",
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

// Search routes
func (app *application) searchExerciseByIDHandler(c *fiber.Ctx) error {
	exerciseIDstr := c.Params("exerciseID")

	exerciseID, err := primitive.ObjectIDFromHex(exerciseIDstr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid routineID format",
		})
	}

	exercise, err := app.store.Exercise.SearchExerciseByID(c.Context(), exerciseID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Exercise not found",
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
