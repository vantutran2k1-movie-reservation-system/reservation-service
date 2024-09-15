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
		users := apiV1.Group("/users")
		{
			users.POST("/", controllers.CreateUser)
			users.POST("/login", controllers.LoginUser)
			users.POST("/logout", middlewares.AuthMiddleware(), controllers.LogoutUser)
			users.PUT("/password", middlewares.AuthMiddleware(), controllers.UpdatePassword)
		}
	}

	return router
}
