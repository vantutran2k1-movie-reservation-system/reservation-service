package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/auth"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/constants"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/controllers"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/middlewares"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/repositories"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/services"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/transaction"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/config"
)

func RegisterRoutes() *gin.Engine {
	r := setupRepositories()
	s := setupServices(r)
	c := setupControllers(s)
	m := setupMiddlewares(r)

	router := gin.Default()

	authMiddleware := m.AuthMiddleware.RequireAuthMiddleware()
	filesUploadMiddleware := m.FilesUploadMiddleware

	apiV1 := router.Group("/api/v1")
	{
		users := apiV1.Group("/users")
		{
			users.GET("/", authMiddleware, c.UserController.GetUser)
			users.POST("/", c.UserController.CreateUser)

			users.POST("/login", c.UserController.LoginUser)
			users.POST("/logout", authMiddleware, c.UserController.LogoutUser)

			users.PUT("/password", authMiddleware, c.UserController.UpdateUserPassword)
		}

		profiles := apiV1.Group("/profiles")
		profiles.Use(authMiddleware)
		{
			profiles.GET("/", c.UserProfileController.GetProfileByUserID)
			profiles.POST("/", c.UserProfileController.CreateUserProfile)
			profiles.PUT("/", c.UserProfileController.UpdateUserProfile)

			profiles.PUT(
				"/profile-picture",
				filesUploadMiddleware.RequireNumberOfUploadedFilesMiddleware(constants.PROFILE_PICTURE_REQUEST_FORM_KEY, 1),
				filesUploadMiddleware.IsAllowedFileTypeMiddleware(constants.PROFILE_PICTURE_REQUEST_FORM_KEY, middlewares.DEFAULT_IMAGE_FILE_TYPES),
				filesUploadMiddleware.NotExceedMaxSizeLimitMiddleware(constants.PROFILE_PICTURE_REQUEST_FORM_KEY, middlewares.GetMaxProfilePictureFileSize()),
				c.UserProfileController.UpdateProfilePicture,
			)
			profiles.DELETE("/profile-picture", c.UserProfileController.DeleteProfilePicture)
		}

		movies := apiV1.Group("/movies")
		{
			movies.POST("/", authMiddleware, c.MovieController.CreateMovie)
		}
	}

	return router
}

type Repositories struct {
	UserRepository           repositories.UserRepository
	LoginTokenRepository     repositories.LoginTokenRepository
	UserSessionRepository    repositories.UserSessionRepository
	UserProfileRepository    repositories.UserProfileRepository
	ProfilePictureRepository repositories.ProfilePictureRepository
	MovieRepository          repositories.MovieRepository
}

type Services struct {
	UserService        services.UserService
	UserProfileService services.UserProfileService
	MovieService       services.MovieService
}

type Controllers struct {
	UserController        controllers.UserController
	UserProfileController controllers.UserProfileController
	MovieController       controllers.MovieController
}

type Middlewares struct {
	AuthMiddleware        middlewares.AuthMiddleware
	FilesUploadMiddleware middlewares.FilesUploadMiddleware
}

func setupRepositories() *Repositories {
	return &Repositories{
		UserRepository:           repositories.NewUserRepository(config.DB),
		LoginTokenRepository:     repositories.NewLoginTokenRepository(config.DB),
		UserSessionRepository:    repositories.NewUserSessionRepository(config.RedisClient),
		UserProfileRepository:    repositories.NewUserProfileRepository(config.DB),
		ProfilePictureRepository: repositories.NewProfilePictureRepository(config.MinioClient),
		MovieRepository:          repositories.NewMovieRepository(config.DB),
	}
}

func setupServices(repositories *Repositories) *Services {
	return &Services{
		UserService: services.NewUserService(
			config.DB,
			config.RedisClient,
			auth.NewAuthenticator(),
			auth.NewTokenGenerator(),
			transaction.NewTransactionManager(),
			repositories.UserRepository,
			repositories.LoginTokenRepository,
			repositories.UserSessionRepository,
		),
		UserProfileService: services.NewUserProfileService(
			config.DB,
			transaction.NewTransactionManager(),
			repositories.UserProfileRepository,
			repositories.ProfilePictureRepository,
		),
		MovieService: services.NewMovieService(
			config.DB,
			transaction.NewTransactionManager(),
			repositories.MovieRepository,
		),
	}
}

func setupControllers(services *Services) *Controllers {
	return &Controllers{
		UserController:        *controllers.NewUserController(&services.UserService),
		UserProfileController: *controllers.NewUserProfileController(&services.UserProfileService),
		MovieController:       *controllers.NewMovieController(&services.MovieService),
	}
}

func setupMiddlewares(repositories *Repositories) *Middlewares {
	return &Middlewares{
		AuthMiddleware:        *middlewares.NewAuthMiddleware(repositories.UserSessionRepository),
		FilesUploadMiddleware: *middlewares.NewFilesUploadMiddleware(),
	}
}
