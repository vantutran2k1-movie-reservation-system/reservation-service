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
		}
	}

	return router
}

type Repositories struct {
	UserRepository repositories.UserRepository
}

type Services struct {
	UserService services.UserService
}

type Controllers struct {
	UserController controllers.UserController
}

func setupRepositories() *Repositories {
	return &Repositories{
		UserRepository: repositories.NewUserRepository(config.DB),
	}
}

func setupServices() *Services {
	repositories := setupRepositories()

	return &Services{
		UserService: services.NewUserService(repositories.UserRepository, config.DB),
	}
}

func setupControllers() *Controllers {
	services := setupServices()

	return &Controllers{
		UserController: *controllers.NewUserController(&services.UserService),
	}
}
