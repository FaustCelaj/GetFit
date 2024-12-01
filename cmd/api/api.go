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
	// creating a custom exercise
	exercises.Post("/", app.createExerciseHandler)

	users := api.Group("/users")
	users.Post("/", app.createUserHandler)
	users.Patch("/:userID", app.patchUserHandler)
	users.Delete("/:userID", app.deleteUserHandler)

	// Routine Exercise Routes
	routineExercises := api.Group("/:userID/routineExercises")
	routineExercises.Post("/", app.createRoutineExerciseHandler)
	routineExercises.Get("/:routineExerciseID", app.getRoutineExerciseByIDHandler)
	routineExercises.Patch("/:routineExerciseID", app.updateRoutineExerciseHandler)

	// Sets Routes
	exerciseSets := api.Group("/:routineExerciseID")
	// Add a set to a specific RoutineExercise
	exerciseSets.Post("/", app.addExerciseSetHandler)

	return fiberApp
}
