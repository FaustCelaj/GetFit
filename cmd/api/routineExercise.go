package main

import (
	"github.com/FaustCelaj/GetFit.git/internal/store"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddExerciseToRoutine godoc
//
//	@Summary		Add exercise to routine
//	@Description	Add an exercise with template sets to a workout routine
//	@Tags			routine-exercises
//	@Accept			json
//	@Produce		json
//	@Param			userID		path		string	true	"User ID"
//	@Param			routineID	path		string	true	"Routine ID"
//	@Param			exerciseID	path		string	true	"Exercise ID"
//	@Param			data		body		object	true	"Template sets information with version"
//	@Success		200			{object}	string	"Exercise added to routine successfully"
//	@Failure		400			{object}	error	"Invalid request body or IDs"
//	@Failure		500			{object}	error	"Failed to add exercise to routine"
//
// @Security		ApiKeyAuth
//
//	@Router			/users/{userID}/routine/{routineID}/exercise/{exerciseID} [post]
func (app *application) addExerciseToRoutineHandler(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	if userID == primitive.NilObjectID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "userID not found in context",
		})
	}

	routineID := c.Params("routineID")
	exerciseID := c.Params("exerciseID")

	if routineID == "" || exerciseID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "userID, routineID, and exerciseID are required",
		})
	}

	routineObjectID, err := primitive.ObjectIDFromHex(routineID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid routineID format",
		})
	}

	exerciseObjectID, err := primitive.ObjectIDFromHex(exerciseID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid exerciseID format",
		})
	}

	// Parse template sets from request body
	var payload struct {
		TemplateSets []store.TemplateSet `json:"template_sets"`
		Version      int16               `json:"expected_version"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if payload.Version == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "expected_version is required",
		})
	}

	// Add the exercise with template sets to the routine
	err = app.store.Routine.AddExerciseToRoutine(
		c.Context(),
		routineObjectID,
		userID,
		exerciseObjectID,
		payload.TemplateSets,
		payload.Version,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to add exercise to routine",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "exercise added to routine successfully",
	})
}

// UpdateExerciseInRoutine godoc
//
//	@Summary		Update exercise in routine
//	@Description	Update template sets for an exercise in a routine
//	@Tags			routine-exercises
//	@Accept			json
//	@Produce		json
//	@Param			userID		path		string	true	"User ID"
//	@Param			routineID	path		string	true	"Routine ID"
//	@Param			exerciseID	path		string	true	"Exercise ID"
//	@Param			data		body		object	true	"Updated template sets with version"
//	@Success		200			{object}	string	"Exercise template sets updated successfully"
//	@Failure		400			{object}	error	"Invalid request body or IDs"
//	@Failure		500			{object}	error	"Failed to update exercise in routine"
//
// @Security		ApiKeyAuth
//
//	@Router			/users/{userID}/routine/{routineID}/exercise/{exerciseID} [patch]
func (app *application) updateExerciseInRoutineHandler(c *fiber.Ctx) error {
	userID := c.Params("userID")
	routineID := c.Params("routineID")
	exerciseID := c.Params("exerciseID")

	if userID == "" || routineID == "" || exerciseID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "userID, routineID, and exerciseID are required",
		})
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid userID format",
		})
	}

	routineObjectID, err := primitive.ObjectIDFromHex(routineID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid routineID format",
		})
	}

	exerciseObjectID, err := primitive.ObjectIDFromHex(exerciseID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid exerciseID format",
		})
	}

	// Parse template sets from request body
	var payload struct {
		TemplateSets []store.TemplateSet `json:"template_sets"`
		Version      int16               `json:"expected_version"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if payload.Version == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "expected_version is required",
		})
	}

	// Update the exercise template sets in the routine
	err = app.store.Routine.UpdateExerciseInRoutine(
		c.Context(),
		routineObjectID,
		userObjectID,
		exerciseObjectID,
		payload.TemplateSets,
		payload.Version,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to update exercise in routine",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "exercise template sets updated successfully",
	})
}

// RemoveExerciseFromRoutine godoc
//
//	@Summary		Remove exercise from routine
//	@Description	Remove an exercise from a workout routine
//	@Tags			routine-exercises
//	@Accept			json
//	@Produce		json
//	@Param			userID		path		string	true	"User ID"
//	@Param			routineID	path		string	true	"Routine ID"
//	@Param			exerciseID	path		string	true	"Exercise ID"
//	@Param			version		body		object	true	"Expected version for optimistic concurrency"
//	@Success		200			{object}	string	"Exercise removed from routine successfully"
//	@Failure		400			{object}	error	"Invalid IDs or version"
//	@Failure		500			{object}	error	"Failed to remove exercise from routine"
//
// @Security		ApiKeyAuth
//
//	@Router			/users/{userID}/routine/{routineID}/exercise/{exerciseID} [delete]
func (app *application) removeExerciseFromRoutineHandler(c *fiber.Ctx) error {
	userID := c.Params("userID")
	routineID := c.Params("routineID")
	exerciseID := c.Params("exerciseID")

	if userID == "" || routineID == "" || exerciseID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "userID, routineID, and exerciseID are required",
		})
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid userID format",
		})
	}

	routineObjectID, err := primitive.ObjectIDFromHex(routineID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid routineID format",
		})
	}

	exerciseObjectID, err := primitive.ObjectIDFromHex(exerciseID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid exerciseID format",
		})
	}

	// Parse version from request body
	var payload struct {
		Version int16 `json:"expected_version"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if payload.Version == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "expected_version is required",
		})
	}

	// Remove the exercise from the routine
	err = app.store.Routine.RemoveExerciseFromRoutine(
		c.Context(),
		routineObjectID,
		userObjectID,
		exerciseObjectID,
		payload.Version,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to remove exercise from routine",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "exercise removed from routine successfully",
	})
}
