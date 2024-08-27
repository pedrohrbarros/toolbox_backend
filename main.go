package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err:= godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	database_url := os.Getenv("DATABASE_URL")
	log.Default().Println(database_url)
}