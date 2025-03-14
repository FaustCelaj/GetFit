package main

import (
	"github.com/FaustCelaj/GetFit.git/internal/store"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (app *application) userContextMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// fetching the userID from URL
		userIDStr := c.Params("userID")
		if userIDStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "userID is required",
			})
		}

		// converting ID to objectID
		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid user ID format",
			})
		}

		// get user from the store
		user, err := app.store.Users.GetByID(c.Context(), userID)
		if err != nil {
			if err == store.ErrNotFound || err.Error() == "mongo: no documents in result" {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "User not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch user",
			})
		}

		// add user and userID to request context
		c.Locals("user", user)
		c.Locals("userID", userID)

		return c.Next()
	}
}

// getUserFromContext retrieves the user from the request context
func getUserFromContext(c *fiber.Ctx) *store.User {
	user, ok := c.Locals("user").(*store.User)
	if !ok {
		return nil
	}
	return user
}

// getUserIDFromContext retrieves the userID from the request context
func getUserIDFromContext(c *fiber.Ctx) primitive.ObjectID {
	userID, ok := c.Locals("userID").(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID
	}
	return userID
}
