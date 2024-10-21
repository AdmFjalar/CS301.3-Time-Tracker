package main

import (
	"log"
	"os"
)

// main is the entry point of the application. It initializes the configuration,
// creates an application instance, mounts the routes, and starts the server.
func main() {
	addr := os.Getenv("SERVER_ADDRESS")
	if addr == "" {
		addr = ":8080"
	}

	cfg := config{
		addr: addr,
	}
	app := &application{
		config: cfg,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
