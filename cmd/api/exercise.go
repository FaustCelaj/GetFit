package main

import (
	"github.com/FaustCelaj/GetFit.git/internal/store"
	"github.com/gofiber/fiber/v2"
)

func (app *application) createExerciseHandler(c *fiber.Ctx) error {
	var exercise store.Exercise

	// Parse the JSON body into the Exercise struct.
	if err := c.BodyParser(&exercise); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := app.store.Exercises.Create(c.Context(), &exercise); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create exercise",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id": exercise.ID.Hex(),
	})
}
