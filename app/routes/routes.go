package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/controllers"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/api/users", controllers.CreateUser)
	router.POST("/api/users/login", controllers.LoginUser)

	return router
}
