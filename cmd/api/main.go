package main

import (
	"time"

	"github.com/hayohtee/social/internal/database"
	"github.com/hayohtee/social/internal/env"
	"github.com/hayohtee/social/internal/repository"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

const version = "0.0.1"

func main() {
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	err := godotenv.Load()
	if err != nil {
		logger.Fatal(err)
	}

	cfg := config{
		addr: env.GetString("ADDR", ":4000"),
		env:  env.GetString("ENV", "development"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://user:password@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		mail: mailConfig{
			exp: time.Hour * 24 * 3, // 3 days
		},
	}

	db, err := database.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("database connection pool established")

	app := &application{
		config:     cfg,
		repository: repository.NewRepository(db),
		logger:     logger,
	}

	routes := app.routes()

	if err := app.serve(routes); err != nil {
		logger.Fatal(err)
	}
}
