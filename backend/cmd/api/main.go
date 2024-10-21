package main

import (
	"log"
	"os"
)

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
