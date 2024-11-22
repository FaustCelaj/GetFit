package main

import (
	"log"
	"os"

	"github.com/FaustCelaj/GetFit.git/internal/env"
)

const version = "0.0.1"

func main() {
	// Set up configuration
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
	}

	// Create an application instance
	app := &application{
		config: cfg,
	}

	os.LookupEnv("PATH")

	// Mount routes
	fiberApp := app.mount()

	// Start the server
	log.Printf("Server is running on %s...", app.config.addr)
	if err := fiberApp.Listen(app.config.addr); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
