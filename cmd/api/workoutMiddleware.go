package main

import (
	"github.com/FaustCelaj/GetFit.git/internal/store"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (app *application) workoutContextMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {

		sessionIDStr := c.Params("sessionID")
		if sessionIDStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "sessionID is required",
			})
		}

		userID := getUserIDFromContext(c)
		if userID == primitive.NilObjectID {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "userID not found in context",
			})
		}

		sessionID, err := primitive.ObjectIDFromHex(sessionIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid sessionID format",
			})
		}

		session, err := app.store.WorkoutSession.GetByID(c.Context(), sessionID, userID)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Workout session not found or does not belong to the user",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch workout session",
			})
		}

		c.Locals("session", session)
		c.Locals("sessionID", sessionID)

		return c.Next()
	}
}

func getSessionFromContext(c *fiber.Ctx) *store.WorkoutSession {
	session, ok := c.Locals("session").(*store.WorkoutSession)
	if !ok {
		return nil
	}
	return session
}

func getSessionIDFromContext(c *fiber.Ctx) primitive.ObjectID {
	sessionID, ok := c.Locals("sessionID").(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID
	}
	return sessionID
}
