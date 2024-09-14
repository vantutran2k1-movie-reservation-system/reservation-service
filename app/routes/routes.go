package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/controllers"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/middlewares"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	apiV1 := router.Group("/api/v1")
	{
		user := apiV1.Group("/users")
		{
			user.POST("/", controllers.CreateUser)
			user.POST("/login", controllers.LoginUser)
			user.POST("/logout", middlewares.AuthMiddleware(), controllers.LogoutUser)
		}
	}

	return router
}
