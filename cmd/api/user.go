package main

import (
	"errors"

	"github.com/FaustCelaj/GetFit.git/internal/store"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CREATE Handler
func (app *application) createUserHandler(c *fiber.Ctx) error {
	var user store.User

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

	userIDStr := c.Params("userID")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "userID is required",
		})
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user ID format",
		})
	}

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

	return c.Status(fiber.StatusOK).JSON(user)
}

type updateUserPayload struct {
	Username        *string `json:"username"`
	Email           *string `json:"email"`
	FirstName       *string `json:"first_name,omitempty"`
	LastName        *string `json:"last_name,omitempty"`
	Age             *int8   `json:"age,omitempty"`
	Title           *string `json:"title,omitempty"`
	Bio             *string `json:"bio,omitempty"`
	ExpectedVersion int16   `json:"expected_version"`
}

// PATCH Handler
func (app *application) patchUserHandler(c *fiber.Ctx) error {
	userIDStr := c.Params("userID")
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	var payload updateUserPayload
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

	if payload.Username != nil {
		if *payload.Username == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Username cannot be empty",
			})
		}
		updates["username"] = *payload.Username
	}

	if payload.Email != nil {
		if *payload.Email == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Email cannot be empty",
			})
		}
		updates["email"] = *payload.Email
	}
	if payload.FirstName != nil {
		updates["first_name"] = *payload.FirstName
	}
	if payload.LastName != nil {
		updates["last_name"] = *payload.LastName
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

	if len(updates) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No fields to update",
		})
	}

	if err := app.store.Users.Update(c.Context(), userID, updates, payload.ExpectedVersion); err != nil {
		if errors.Is(err, store.ErrVersionMismatch) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "This record has been modified since you last viewed it. Please refresh and try again.",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User updated successfully",
	})
}

// DELETE Handler
func (app *application) deleteUserHandler(c *fiber.Ctx) error {
	userIDStr := c.Params("userID")

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	if err := app.store.Users.Delete(c.Context(), userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User was successfully deleted",
	})
}
