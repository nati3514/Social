package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/nati3514/Social/internal/env"
)

func main() {
	godotenv.Load()
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
	}

	app := &application{
		config: cfg,
	}

	mux := app.mount()
	log.Println(app.run(mux))
}
