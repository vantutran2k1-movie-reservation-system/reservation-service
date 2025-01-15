package main

import (
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/config"
	"log"
)

func main() {
	a := app.InitApp()
	if err := a.Router.Run(":" + config.AppEnv.AppPort); err != nil {
		log.Fatal(err)
		return
	}
}
