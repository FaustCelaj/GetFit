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
	api := fiberApp.Group("api/v1")

	// Health Check
	api.Get("/health", app.healthCheckHandler)

	user := api.Group("/user")
	user.Get("/:userID", app.getUserHandler)
	user.Post("/", app.createUserHandler)
	user.Patch("/:userID", app.patchUserHandler)
	user.Delete("/:userID", app.deleteUserHandler)

	search := api.Group("/search")
	search.Get("/:exerciseID", app.searchExerciseByIDHandler)

	userScoped := api.Group("/users/:userID")

	// Exercise Routes
	exercise := userScoped.Group("/exercise")
	exercise.Post("/", app.createExerciseHandler)
	exercise.Get("/", app.getAllUserExerciseHandler)
	exercise.Get("/:exerciseID", app.getExerciseByIDHandler)
	exercise.Patch("/:exerciseID", app.updateExerciseHandler)
	exercise.Delete("/:exerciseID", app.deleteExerciseHandler)

	// Routine Routes
	routine := userScoped.Group("/routine")
	routine.Post("/", app.createRoutineHandler)
	routine.Get("/", app.getAllUserRoutinesIDHandler)
	routine.Get("/:routineID", app.getRoutineByIDHandler)
	routine.Patch("/:routineID", app.patchRoutineHandler)
	routine.Delete("/:routineID", app.deleteRoutineHandler)

	// Editing exercises in routines
	routineExercise := routine.Group("/:routineID/exercise")
	routineExercise.Post("/:exerciseID", app.addExerciseToRoutineHandler)
	routineExercise.Patch("/:exerciseID", app.updateExerciseInRoutineHandler)
	routineExercise.Delete("/:exerciseID", app.removeExerciseFromRoutineHandler)

	// Workout Session Routes (Actual performed workouts)
	workouts := userScoped.Group("/workout")
	workouts.Post("/", app.createWorkoutSessionHandler)
	workouts.Post("/from-routine/:routineID", app.createWorkoutFromRoutineHandler)
	workouts.Get("/", app.getAllWorkoutSessionsHandler)
	workouts.Get("/:sessionID", app.getWorkoutSessionByIDHandler)
	workouts.Post("/:sessionID/complete", app.completeWorkoutSessionHandler)
	workouts.Delete("/:sessionID", app.deleteWorkoutSessionHandler)

	// Routes for adding sets to exercises within a workout session
	workoutSets := workouts.Group("/:sessionID/exercise/:exerciseID/sets")
	workoutSets.Post("/", app.addSetToWorkoutHandler)

	return fiberApp
}
