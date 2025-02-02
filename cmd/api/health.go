package main

import "github.com/gofiber/fiber/v2"

func (app *application) healthCheckHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "available",
		"version": version,
	})
}
