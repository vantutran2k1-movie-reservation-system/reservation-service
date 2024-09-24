package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/controllers"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/middlewares"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/repositories"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/services"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/config"
)

func RegisterRoutes() *gin.Engine {
	repositories := setupRepositories()
	services := setupServices(repositories)
	controllers := setupControllers(services)

	authMiddleware := middlewares.NewAuthMiddleware(repositories.UserSessionRepository).RequireBasicAuthMiddleware()

	router := gin.Default()

	apiV1 := router.Group("/api/v1")
	{
		users := apiV1.Group("/users")
		{
			users.GET("/", authMiddleware, controllers.UserController.GetUser)
			users.POST("/", controllers.UserController.CreateUser)
			users.POST("/login", controllers.UserController.LoginUser)
			users.POST("/logout", authMiddleware, controllers.UserController.LogoutUser)
			users.PUT("/password", authMiddleware, controllers.UserController.UpdateUserPassword)
		}

		profiles := apiV1.Group("/profiles")
		{
			profiles.GET("/", authMiddleware, controllers.UserProfileController.GetProfileByUserID)
			profiles.POST("/", authMiddleware, controllers.UserProfileController.CreateUserProfile)
			profiles.PUT("/", authMiddleware, controllers.UserProfileController.UpdateUserProfile)
		}
	}

	return router
}

type Repositories struct {
	UserRepository        repositories.UserRepository
	LoginTokenRepository  repositories.LoginTokenRepository
	UserSessionRepository repositories.UserSessionRepository
	UserProfileRepository repositories.UserProfileRepository
}

type Services struct {
	UserService        services.UserService
	UserProfileService services.UserProfileService
}

type Controllers struct {
	UserController        controllers.UserController
	UserProfileController controllers.UserProfileController
}

func setupRepositories() *Repositories {
	return &Repositories{
		UserRepository:        repositories.NewUserRepository(config.DB),
		LoginTokenRepository:  repositories.NewLoginTokenRepository(config.DB),
		UserSessionRepository: repositories.NewUserSessionRepository(config.RedisClient),
		UserProfileRepository: repositories.NewUserProfileRepository(config.DB),
	}
}

func setupServices(repositories *Repositories) *Services {
	return &Services{
		UserService: services.NewUserService(
			config.DB,
			config.RedisClient,
			repositories.UserRepository,
			repositories.LoginTokenRepository,
			repositories.UserSessionRepository,
		),
		UserProfileService: services.NewUserProfileService(
			config.DB,
			repositories.UserProfileRepository,
		),
	}
}

func setupControllers(services *Services) *Controllers {
	return &Controllers{
		UserController:        *controllers.NewUserController(&services.UserService),
		UserProfileController: *controllers.NewUserProfileController(&services.UserProfileService),
	}
}
