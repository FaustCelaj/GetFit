package main

import (
	"github.com/FaustCelaj/GetFit.git/internal/store"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (app *application) routineContextMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		routineIDStr := c.Params("routineID")
		if routineIDStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "routineID is required",
			})
		}

		userID := getUserIDFromContext(c)
		if userID == primitive.NilObjectID {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "userID not found in context",
			})
		}

		routineID, err := primitive.ObjectIDFromHex(routineIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid routine ID format",
			})
		}

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

		c.Locals("routine", routine)
		c.Locals("routineID", routineID)

		return c.Next()
	}
}

func getRoutineFromContext(c *fiber.Ctx) *store.Routine {
	routine, ok := c.Locals("routine").(*store.Routine)
	if !ok {
		return nil
	}
	return routine
}

func getRoutineIDFromContext(c *fiber.Ctx) primitive.ObjectID {
	routineID, ok := c.Locals("routineID").(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID
	}
	return routineID
}
