package main

import (
	"log"

	"github.com/hayohtee/social/internal/env"
	"github.com/hayohtee/social/internal/repository"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
	}

	app := &application{
		config: cfg,
		repository: repository.NewRepository(nil),
	}

	routes := app.routes()

	if err := app.serve(routes); err != nil {
		log.Fatal(err)
	}
}
