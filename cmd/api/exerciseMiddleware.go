package main

import (
	"github.com/FaustCelaj/GetFit.git/internal/store"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (app *application) exerciseContextMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		exerciseIDStr := c.Params("exerciseID")
		if exerciseIDStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "exerciseID is required",
			})
		}

		userID := getUserIDFromContext(c)
		if userID == primitive.NilObjectID {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "userID not found in context",
			})
		}

		exerciseID, err := primitive.ObjectIDFromHex(exerciseIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid exerciseID format",
			})
		}

		exercise, err := app.store.Exercise.SearchExerciseByID(c.Context(), exerciseID)
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

		c.Locals("exercise", exercise)
		c.Locals("exerciseID", exerciseID)

		return c.Next()
	}
}

func getExerciseFromContext(c *fiber.Ctx) *store.Exercise {
	exercise, ok := c.Locals("exercise").(*store.Exercise)
	if !ok {
		return nil
	}
	return exercise
}

func getExerciseIDFromContext(c *fiber.Ctx) primitive.ObjectID {
	exerciseID, ok := c.Locals("exerciseID").(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID
	}
	return exerciseID
}
