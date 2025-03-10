package main

import (
	"github.com/FaustCelaj/GetFit.git/internal/store"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CREATE Handler
func (app *application) createUserHandler(c *fiber.Ctx) error {
	var user store.User

	// Parse the JSON body into the Exercise struct.
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := app.store.Users.Create(c.Context(), &user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id": user.ID.Hex(),
	})
}

// GET Handler
func (app *application) getUserHandler(c *fiber.Ctx) error {
	// Get the userID from the URL parameter
	userIDStr := c.Params("userID")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "userID is required",
		})
	}

	// Convert the userID to primitive.ObjectID
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user ID format",
		})
	}

	// Fetch the user from the database
	user, err := app.store.Users.GetByID(c.Context(), userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to fetch user",
			"details": err.Error(),
		})
	}

	// Return the user
	return c.Status(fiber.StatusOK).JSON(user)
}

type updateUserPayload struct {
	Username        *string `json:"username"`
	Email           *string `json:"email"`
	FirstName       *string `json:"first_name,omitempty"`
	Age             *int8   `json:"age,omitempty"`
	Title           *string `json:"title,omitempty"`
	Bio             *string `json:"bio,omitempty"`
	ExpectedVersion int16   `json:"expected_version"`
}

// PATCH Handler
func (app *application) patchUserHandler(c *fiber.Ctx) error {
	var payload updateUserPayload
	// Parse the JSON body into the payload struct
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Extract the userID from the URL
	userIDStr := c.Params("userID")
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	// validate version
	if payload.ExpectedVersion == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Current version is required",
		})
	}

	// Build the updates map dynamically
	updates := make(map[string]interface{})
	if payload.Username != nil {
		updates["username"] = *payload.Username
	}
	if payload.Email != nil {
		updates["email"] = *payload.Email
	}
	if payload.FirstName != nil {
		updates["first_name"] = *payload.FirstName
	}
	if payload.Age != nil {
		updates["age"] = *payload.Age
	}
	if payload.Title != nil {
		updates["title"] = *payload.Title
	}
	if payload.Bio != nil {
		updates["bio"] = *payload.Bio
	}

	// If no fields to update, return a bad request
	if len(updates) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No fields to update",
		})
	}

	// Perform the update in the database
	if err := app.store.Users.Update(c.Context(), userID, updates, payload.ExpectedVersion); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	// Return a success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User updated successfully",
	})
}

// DELETE Handler
func (app *application) deleteUserHandler(c *fiber.Ctx) error {
	// Get the userID from the URL parameter
	userIDStr := c.Params("userID")

	// Convert the userID to primitive.ObjectID
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	// Delete the user from the store
	if err := app.store.Users.Delete(c.Context(), userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	// Return a success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User was successfully deleted",
	})
}

// // middlewear to fetch userID
// func (app *application) userContextMiddleware(next fiber.Handler) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		// Get the userID from the URL parameter
// 		userIDStr := c.Params("userID")
// 		// Convert the userID to primitive.ObjectID
// 		userID, err := primitive.ObjectIDFromHex(userIDStr)
// 		if err != nil {
// 			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 				"error": "Invalid user ID format",
// 			})
// 		}
// 		// Fetch the user from the store
// 		user, err := app.store.Users.GetById(c.Context(), userID)
// 		if err != nil {
// 			// Handle case when user is not found
// 			if errors.Is(err, store.ErrNotFound) {
// 				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
// 					"error": "User not found",
// 				})
// 			}
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 				"error": "Failed to fetch user",
// 			})
// 		}
// 		// Add the user to the request context
// 		c.Locals("user", user)
// 		// Call the next handler
// 		return next(c)
// 	}
// }
// // helper func
// func getUserFromCtx(c *fiber.Ctx) *store.User {
// 	user, _ := c.Locals("user").(*store.User)
// 	return user
// }
