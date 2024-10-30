package main

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/hayohtee/social/internal/database"
	"github.com/hayohtee/social/internal/env"
	"github.com/hayohtee/social/internal/repository"
	"github.com/joho/godotenv"
)

const version = "0.0.1"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		env:  env.GetString("ENV", "development"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://user:password@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	db, err := database.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("database connection pool established")

	log.Println("configuring validator")
	validate := validator.New(validator.WithRequiredStructEnabled())

	app := &application{
		config:     cfg,
		repository: repository.NewRepository(db),
		validate:   validate,
	}

	routes := app.routes()

	if err := app.serve(routes); err != nil {
		log.Fatal(err)
	}
}
