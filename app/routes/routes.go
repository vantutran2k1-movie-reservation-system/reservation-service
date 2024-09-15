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

		profiles := apiV1.Group("/profiles")
		{
			profiles.POST("/", middlewares.AuthMiddleware(), controllers.CreateUserProfile)
			profiles.PUT("/", middlewares.AuthMiddleware(), controllers.UpdateUserProfile)
		}
	}

	return router
}
