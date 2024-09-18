package app

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/routes"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/config"
)

type App struct {
	Router *gin.Engine
}

func InitApp() *App {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	errors.RegisterCustomValidators()

	config.InitDB()
	config.InitRedis()

	router := routes.RegisterRoutes()
	return &App{
		Router: router,
	}
}
