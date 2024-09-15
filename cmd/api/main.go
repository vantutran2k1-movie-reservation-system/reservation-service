package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/routes"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/config"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	errors.RegisterCustomValidators()

	config.InitDB()
	config.InitRedis()

	router := routes.SetupRouter()
	if err := router.Run(":" + os.Getenv("APP_PORT")); err != nil {
		log.Fatal(err)
		return
	}
}
