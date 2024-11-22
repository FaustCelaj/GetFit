package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type application struct {
	config config
	addr   string
}

type config struct {
	addr string
}

func (app *application) mount() *fiber.App {
	// Create a new Fiber instance
	fiberApp := fiber.New(fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	})

	// Add middlewares
	fiberApp.Use(logger.New())  // Logs all requests
	fiberApp.Use(recover.New()) // Recovers from panics

	// Routes
	api := fiberApp.Group("/v1")

	// Health Check
	api.Get("/", app.healthCheckHandler)

	// Exercise Routes
	// exercises := api.Group("/exercises")
	// exercises.Exercise("/", app.createExerciseHandler)

	return fiberApp
}
