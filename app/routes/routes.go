package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/controllers"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/repositories"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/services"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/config"
)

func RegisterRoutes() *gin.Engine {
	controllers := setupControllers()

	router := gin.Default()

	apiV1 := router.Group("/api/v1")
	{
		users := apiV1.Group("/users")
		{
			users.POST("/", controllers.UserController.CreateUser)
			users.POST("/login", controllers.UserController.LoginUser)
		}
	}

	return router
}

type Repositories struct {
	UserRepository        repositories.UserRepository
	LoginTokenRepository  repositories.LoginTokenRepository
	UserSessionRepository repositories.UserSessionRepository
}

type Services struct {
	UserService services.UserService
}

type Controllers struct {
	UserController controllers.UserController
}

func setupRepositories() *Repositories {
	return &Repositories{
		UserRepository:        repositories.NewUserRepository(config.DB),
		LoginTokenRepository:  repositories.NewLoginTokenRepository(config.DB),
		UserSessionRepository: repositories.NewUserSessionRepository(config.RedisClient),
	}
}

func setupServices() *Services {
	repositories := setupRepositories()

	return &Services{
		UserService: services.NewUserService(
			config.DB,
			config.RedisClient,
			repositories.UserRepository,
			repositories.LoginTokenRepository,
			repositories.UserSessionRepository,
		),
	}
}

func setupControllers() *Controllers {
	services := setupServices()

	return &Controllers{
		UserController: *controllers.NewUserController(&services.UserService),
	}
}
