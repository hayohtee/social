package main

import (
	"log"

	"github.com/hayohtee/social/internal/database"
	"github.com/hayohtee/social/internal/env"
	"github.com/hayohtee/social/internal/repository"
	"github.com/joho/godotenv"
)

const version = "0.0.1"

//	@title			Social API
//	@version		0.0.1
//	@description	API documentation for Social, a social network for Gophers 
//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath					/v1
//
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		env:  env.GetString("ENV", "development"),
		apiUrl: env.GetString("EXTERNAL_URL", "localhost:8080"),
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

	app := &application{
		config:     cfg,
		repository: repository.NewRepository(db),
	}

	routes := app.routes()

	if err := app.serve(routes); err != nil {
		log.Fatal(err)
	}
}
