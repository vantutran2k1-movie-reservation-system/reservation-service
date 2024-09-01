package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/config"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	config.InitDB()

	appPort := os.Getenv("APP_PORT")
	log.Printf("Starting server on :%s", appPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", appPort), nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
