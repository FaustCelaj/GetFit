package main

import (
	"github.com/FaustCelaj/GetFit.git/internal/store"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateExercise godoc
//
//	@Summary		Create a new exercise
//	@Description	Create a custom exercise for a user
//	@Tags			exercises
//	@Accept			json
//	@Produce		json
//	@Param			userID		path		string			true	"User ID"
//	@Param			exercise	body		store.Exercise	true	"Exercise information"
//	@Success		201			{object}	string			"Exercise created successfully"
//	@Failure		400			{object}	error			"Invalid request body or missing required fields"
//	@Failure		500			{object}	error			"Failed to create exercise"
//
// @Security		ApiKeyAuth
//
//	@Router			/users/{userID}/exercise [post]
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

// GetAllUserExercises godoc
//
//	@Summary		Get all exercises for a user
//	@Description	Retrieve all exercises created by a specific user
//	@Tags			exercises
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		string			true	"User ID"
//	@Success		200		{array}		store.Exercise	"List of exercises"
//	@Failure		400		{object}	error			"Invalid user ID"
//	@Failure		500		{object}	error			"Failed to fetch exercises"
//
// @Security		ApiKeyAuth
//
//	@Router			/users/{userID}/exercise [get]
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

// GetExerciseByID godoc
//
//	@Summary		Get exercise by ID
//	@Description	Retrieve a specific exercise by its ID
//	@Tags			exercises
//	@Accept			json
//	@Produce		json
//	@Param			userID		path		string			true	"User ID"
//	@Param			exerciseID	path		string			true	"Exercise ID"
//	@Success		200			{object}	store.Exercise	"Exercise information"
//	@Failure		400			{object}	error			"Invalid ID format"
//	@Failure		404			{object}	error			"Exercise not found"
//	@Failure		500			{object}	error			"Failed to fetch exercise"
//
// @Security		ApiKeyAuth
//
//	@Router			/users/{userID}/exercise/{exerciseID} [get]
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

// UpdateExercise godoc
//
//	@Summary		Update an exercise
//	@Description	Update an existing exercise's details
//	@Tags			exercises
//	@Accept			json
//	@Produce		json
//	@Param			userID		path		string					true	"User ID"
//	@Param			exerciseID	path		string					true	"Exercise ID"
//	@Param			exercise	body		updateExercisePayload	true	"Updated exercise information"
//	@Success		200			{object}	string					"Exercise updated successfully"
//	@Failure		400			{object}	error					"Invalid request body or missing fields"
//	@Failure		404			{object}	error					"Exercise not found"
//	@Failure		500			{object}	error					"Failed to update exercise"
//
// @Security		ApiKeyAuth
//
//	@Router			/users/{userID}/exercise/{exerciseID} [patch]
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

// DeleteExercise godoc
//
//	@Summary		Delete an exercise
//	@Description	Remove an exercise from the system
//	@Tags			exercises
//	@Accept			json
//	@Produce		json
//	@Param			userID		path		string	true	"User ID"
//	@Param			exerciseID	path		string	true	"Exercise ID"
//	@Success		200			{object}	string	"Exercise successfully deleted"
//	@Failure		400			{object}	error	"Invalid ID format"
//	@Failure		500			{object}	error	"Failed to delete exercise"
//
// @Security		ApiKeyAuth
//
//	@Router			/users/{userID}/exercise/{exerciseID} [delete]
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

// SearchExerciseByID godoc
//
//	@Summary		Search for an exercise by ID
//	@Description	Find an exercise using its ID without user context
//	@Tags			exercises
//	@Accept			json
//	@Produce		json
//	@Param			exerciseID	path		string			true	"Exercise ID"
//	@Success		200			{object}	store.Exercise	"Exercise information"
//	@Failure		400			{object}	error			"Invalid exercise ID format"
//	@Failure		404			{object}	error			"Exercise not found"
//	@Failure		500			{object}	error			"Failed to fetch exercise"
//
// @Security		ApiKeyAuth
//
//	@Router			/search/{exerciseID} [get]
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
