package main

import (
	"time"

	"github.com/FaustCelaj/GetFit.git/internal/store"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type application struct {
	config config
	store  store.Storage
}

type config struct {
	addr string
	db   dbConfig
	env  string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleTime  string
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
	exercises := api.Group("/exercises")
	exercises.Post("/", app.createExerciseHandler)
	// exercises.Route("/:exerciseName", func(router fiber.Router) {
	// 	exercises.Get("/", app.getExerciseHandler)
	// })

	users := api.Group("/users")
	users.Post("/", app.createUserHandler)
	users.Delete("/:userID", app.deleteUserHandler)
	users.Patch("/:userID", app.patchUserHandler)

	return fiberApp
}
