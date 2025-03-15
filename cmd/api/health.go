package main

import "github.com/gofiber/fiber/v2"

// HealthCheck godoc
//
//	@Summary		Health check endpoint
//	@Description	Check if the API is up and running
//	@Tags			health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string	"API status and version"
//	@Router			/health [get]
func (app *application) healthCheckHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "available",
		"version": version,
	})
}
