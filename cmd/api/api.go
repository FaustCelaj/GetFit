package main

import (
	"time"

	"github.com/FaustCelaj/GetFit.git/internal/store"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	fiberApp.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173,https://your-deployed-frontend.vercel.app",
		AllowMethods: "GET,POST,PATCH,DELETE",
		AllowHeaders: "Content-Type,Authorization",
	}))

	// Routes
	api := fiberApp.Group("/v1")

	// Health Check
	api.Get("/", app.healthCheckHandler)

	user := api.Group("/user")
	user.Get("/:userID", app.getUserHandler)       // Get user details ✅
	user.Post("/", app.createUserHandler)          // Create a user ✅
	user.Patch("/:userID", app.patchUserHandler)   // Update user info ✅
	user.Delete("/:userID", app.deleteUserHandler) // Delete user ✅

	// Routine Routes
	routine := api.Group("/:userID/routine")
	routine.Post("/", app.createRoutineHandler)             // create a routine ✅
	routine.Get("/", app.getAllUserRoutinesIDHandler)       // Fetch all routines for a user ✅
	routine.Get("/:routineID", app.getRoutineByIDHandler)   // Get a signle routine ✅
	routine.Patch("/:routineID", app.patchRoutineHandler)   // Update a routine ✅
	routine.Delete("/:routineID", app.deleteRoutineHandler) // Delete a routine ✅

	// Exercise Routes
	exercise := api.Group("/:userID/exercise")
	exercise.Post("/", app.createExerciseHandler)              // Create a custom exercise ✅
	exercise.Get("/", app.getAllUserExercisesHandler)          // Fetch all exercises for a user ✅
	exercise.Get("/:exerciseID", app.getExerciseByIDHandler)   // Get a single exercise ✅
	exercise.Patch("/:exerciseID", app.updateExerciseHandler)  // Update an exercise ✅
	exercise.Delete("/:exerciseID", app.deleteExerciseHandler) // Delete an exercise ✅

	// Routes for Sets (Scoped by Exercise)
	set := api.Group("/:userID/exercise/:exerciseID/set")
	set.Post("/", app.addSetsHandler)           // Add a set to an exercise ✅
	set.Get("/", app.getAllExerciseSetsHandler) // Fetch all sets for an exercise ✅
	set.Get("/:setID", app.getSetByIDHandler)   // Fetch a single set ✅
	set.Patch("/:setID", app.updateSetHandler)  // Update a set
	set.Delete("/:setID", app.deleteSetHandler) // Delete a set ✅

	return fiberApp
}
