package main

import (
	"github.com/FaustCelaj/GetFit.git/internal/store"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

type registerUserPayload struct {
	Username *string `json:"username" validate:"required,max=50"`
	Email    *string `json:"email" validate:"required,email,max=50"`
	Password *string `json:"password" validate:"required,min=8,max=72"`
}

// RegisterUser godoc
//
//	@Summary		Register a new user
//
//	@Description	Register a new user in the system with username, email, and password
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		store.User	true	"User information including username, email, and password"
//
//	@Success		201		{object}	error		"Returns the created user ID"
//	@Failure		400		{object}	error		"Invalid request body"
//	@Failure		500		{object}	error		"Failed to create user"
//	@Security		ApiKeyAuth
//
//	@Router			/user [post]
func (app *application) registerUserHandler(c *fiber.Ctx) error {
	var payload registerUserPayload

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	if err := validate.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		app.logger.Infof("Validation failed for user registration: %v", validationErrors)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": formatValidationErrors(validationErrors),
		})
	}

	user := &store.User{
		username: payload.Username,
		email:    payload.Email,
	}

	if err := user.Password.set(payload.Password); err != nil {
		return nil
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User registered successfully",
	})
}

// helper function for error validation
func formatValidationErrors(errors validator.ValidationErrors) map[string]string {
	formattedErrors := make(map[string]string)

	for _, err := range errors {
		field := err.Field()
		switch err.Tag() {
		case "required":
			formattedErrors[field] = field + " is required"
		case "email":
			formattedErrors[field] = field + " must be a valid email address"
		case "min":
			formattedErrors[field] = field + " must be at least " + err.Param() + " characters long"
		case "max":
			formattedErrors[field] = field + " must be at most " + err.Param() + " characters long"
		default:
			formattedErrors[field] = field + " is invalid"
		}
	}

	return formattedErrors
}
