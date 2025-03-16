package main

import (
	"fmt"
	"time"

	"github.com/FaustCelaj/GetFit.git/docs" // required to generate swagger docs
	"github.com/FaustCelaj/GetFit.git/internal/store"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	swagger "github.com/gofiber/swagger"
	"go.uber.org/zap"
)

type application struct {
	config config
	store  store.Storage
	logger *zap.SugaredLogger
}

type config struct {
	addr   string
	db     dbConfig
	env    string
	apiURL string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleTime  string
}

func (app *application) mount() *fiber.App {
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.apiURL
	docs.SwaggerInfo.BasePath = "/api/v1"

	// Create a new Fiber instance
	fiberApp := fiber.New(fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	})

	// Add middlewares
	// Logs all requests in HTTP level vs zap logs application level
	fiberApp.Use(logger.New(logger.Config{
		Format: "${time} ${status} - ${method} ${path}\n",
	}))
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
	docsURL := fmt.Sprintf("http://%s/api/v1/swagger/doc.json", app.config.apiURL)
	api.Get("/swagger/*", swagger.New(swagger.Config{
		URL: docsURL,
	}))

	// Public Routes
	auth := api.Group("/authentication")
	auth.Post("/register", app.registerUserHandler)

	user := api.Group("/user")
	user.Get("/:userID", app.getUserHandler)
	// user.Post("/", app.createUserHandler)
	user.Patch("/:userID", app.patchUserHandler)
	user.Delete("/:userID", app.deleteUserHandler)

	userScoped := api.Group("/users/:userID", app.userContextMiddleware())

	// Exercise Routes
	exercise := userScoped.Group("/exercise")
	exercise.Post("/", app.createExerciseHandler)
	exercise.Get("/", app.getAllUserExerciseHandler)

	exerciseWithID := exercise.Group("/:exerciseID", app.exerciseContextMiddleware())
	exerciseWithID.Get("/", app.getExerciseByIDHandler)
	exerciseWithID.Patch("/", app.updateExerciseHandler)
	exerciseWithID.Delete("/", app.deleteExerciseHandler)

	// Routine Routes
	routine := userScoped.Group("/routine")
	routine.Post("/", app.createRoutineHandler)
	routine.Get("/", app.getAllUserRoutinesIDHandler)

	routineWithID := routine.Group("/:routineID", app.routineContextMiddleware())
	routineWithID.Get("/", app.getRoutineByIDHandler)
	routineWithID.Patch("/", app.patchRoutineHandler)
	routineWithID.Delete("/", app.deleteRoutineHandler)

	// Editing exercises in routines
	routineExercise := routineWithID.Group("/exercise/:exerciseID", app.exerciseContextMiddleware())
	routineExercise.Post("/", app.addExerciseToRoutineHandler)
	routineExercise.Patch("/", app.updateExerciseInRoutineHandler)
	routineExercise.Delete("/", app.removeExerciseFromRoutineHandler)

	// Workout Session Routes (Actual performed workouts)
	workouts := userScoped.Group("/workout")
	workouts.Post("/", app.createWorkoutSessionHandler)
	workouts.Get("/", app.getAllWorkoutSessionsHandler)

	workoutsFromRoutine := workouts.Group("/from-routine/:routineID", app.routineContextMiddleware())
	workoutsFromRoutine.Post("/", app.createWorkoutFromRoutineHandler)

	workoutSession := workouts.Group("/:sessionID", app.workoutContextMiddleware())
	workoutSession.Get("/", app.getWorkoutSessionByIDHandler)
	workoutSession.Post("/complete", app.completeWorkoutSessionHandler)
	workoutSession.Delete("/", app.deleteWorkoutSessionHandler)

	workoutSets := workoutSession.Group("/exercise/:exerciseID/sets", app.exerciseContextMiddleware())
	workoutSets.Post("/", app.addSetToWorkoutHandler)

	// search := api.Group("/search")
	// search.Get("/:exerciseID", app.searchExerciseByIDHandler)
	return fiberApp
}
