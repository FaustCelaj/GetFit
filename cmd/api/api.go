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
	exercise := api.Group("/exercise")
	// creating a custom exercise
	exercise.Post("/", app.createExerciseHandler)

	user := api.Group("/user")
	user.Get("/:userID", app.getUserHandler)       //✅
	user.Post("/", app.createUserHandler)          //✅
	user.Patch("/:userID", app.patchUserHandler)   //✅
	user.Delete("/:userID", app.deleteUserHandler) //✅

	// Routine Routes
	routine := api.Group("/:userID/routine")
	routine.Post("/", app.createRoutineHandler) //✅
	// Fetch all routines by userID
	routine.Get("/", app.getRoutinesByUserIDHandler) //✅
	// returns a single routine
	routine.Get("/:routineID", app.getRoutineByIDHandler) //✅
	// Delete a single routine
	routine.Delete("/:routineID", app.deleteRoutineHandler)

	// Sets Routes
	// set := api.Group("/:routineID")
	// Add a set to a specific RoutineExercise
	// set.Post("/", app.addSetHandler)

	return fiberApp
}
