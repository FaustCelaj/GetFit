package main

import (
	"context"
	"log"

	"github.com/FaustCelaj/GetFit.git/internal/db"
	"github.com/FaustCelaj/GetFit.git/internal/env"
	"github.com/FaustCelaj/GetFit.git/internal/store"
)

const version = "0.0.1"

func main() {
	// Set up configuration
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "mongodb://localhost:27017"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDEL_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
	}

	client, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleTime)

	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())
	log.Println("MongoDB connection established")

	store := store.NewMongoDBStorage(client.Database("getfit"))

	// Create an application instance
	app := &application{
		config: cfg,
		store:  store,
	}

	// Mount routes
	fiberApp := app.mount()

	// Start the server
	log.Printf("Server is running on http://localhost%s", app.config.addr)
	if err := fiberApp.Listen(app.config.addr); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
