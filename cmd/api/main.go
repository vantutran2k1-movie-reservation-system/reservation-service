package main

import (
	"log"
	"os"

	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app"
)

func main() {
	app := app.InitApp()
	if err := app.Router.Run(":" + os.Getenv("APP_PORT")); err != nil {
		log.Fatal(err)
		return
	}
}
