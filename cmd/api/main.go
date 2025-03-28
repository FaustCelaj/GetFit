package main

import (
	"context"

	"github.com/FaustCelaj/GetFit.git/internal/db"
	"github.com/FaustCelaj/GetFit.git/internal/env"
	"github.com/FaustCelaj/GetFit.git/internal/store"
	"go.uber.org/zap"
)

const version = "0.0.1"

//	@title			GetFit API
//	@version		0.0.1
//	@description	API for tracking workouts, exercises, and fitness routines
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.yoursite.com/support
//	@contact.email	your.email@example.com

//	@license.name	MIT
//	@license.url	https://opensource.org/licenses/MIT

//	@host						localhost:8080
//	@BasePath					/api/v1
//
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description
//	@schemes	http

func main() {
	// Set up configuration
	cfg := config{
		addr:   env.GetString("ADDR", ":8080"),
		apiURL: env.GetString("EXTERNAL_URL", "localHost:8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "mongodb://localhost:27017"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDEL_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
	}

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	client, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleTime)

	if err != nil {
		logger.Fatal(err)
	}
	defer client.Disconnect(context.Background())
	logger.Info("MongoDB connection established")

	store := store.NewMongoDBStorage(client.Database("getfit"))

	// Create an application instance
	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
	}

	// Mount routes
	fiberApp := app.mount()

	// Start the server
	logger.Infof("Server is running on http://localhost%s", app.config.addr)
	if err := fiberApp.Listen(app.config.addr); err != nil {
		logger.Fatal(err)
	}
}
