package main

import (
	"log"
	"os"

	"github.com/hayohtee/social/internal/database"
	"github.com/hayohtee/social/internal/repository"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.New(os.Getenv("DEV_DB_ADDR"), 5, 5, "15m")
	if err != nil {
		log.Fatal("error creating database conn:", err)
	}
	defer db.Close()

	repo := repository.NewRepository(db)
	database.Seed(repo)
}