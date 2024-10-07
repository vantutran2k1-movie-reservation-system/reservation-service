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

	apiV1 := router.Group("/api/v1")
	{
		users := apiV1.Group("/users")
		{
			users.GET("/me", m.AuthMiddleware.RequireAuthMiddleware(), c.UserController.GetUser)
			users.POST("/", c.UserController.CreateUser)

			users.POST("/login", c.UserController.LoginUser)
			users.POST("/logout", m.AuthMiddleware.RequireAuthMiddleware(), c.UserController.LogoutUser)

			users.PUT("/password", m.AuthMiddleware.RequireAuthMiddleware(), c.UserController.UpdateUserPassword)
		}

		profiles := apiV1.Group("/profiles")
		profiles.Use(m.AuthMiddleware.RequireAuthMiddleware())
		{
			profiles.GET("/me", c.UserProfileController.GetProfileByUserID)
			profiles.POST("/", c.UserProfileController.CreateUserProfile)
			profiles.PUT("/", c.UserProfileController.UpdateUserProfile)

			profiles.PUT(
				"/profile-picture",
				m.FilesUploadMiddleware.RequireNumberOfUploadedFilesMiddleware(constants.PROFILE_PICTURE_REQUEST_FORM_KEY, 1),
				m.FilesUploadMiddleware.IsAllowedFileTypeMiddleware(constants.PROFILE_PICTURE_REQUEST_FORM_KEY, middlewares.DEFAULT_IMAGE_FILE_TYPES),
				m.FilesUploadMiddleware.NotExceedMaxSizeLimitMiddleware(constants.PROFILE_PICTURE_REQUEST_FORM_KEY, middlewares.GetMaxProfilePictureFileSize()),
				c.UserProfileController.UpdateProfilePicture,
			)
			profiles.DELETE("/profile-picture", c.UserProfileController.DeleteProfilePicture)
		}

		movies := apiV1.Group("/movies")
		{
			movies.GET("/:id", c.MovieController.GetMovie)
			movies.GET("/", c.MovieController.GetMovies)
			movies.POST(
				"/",
				m.AuthMiddleware.RequireAuthMiddleware(),
				m.AuthMiddleware.RequireFeatureFlagMiddleware(constants.CAN_MODIFY_MOVIES),
				c.MovieController.CreateMovie,
			)
			movies.PUT(
				"/:id",
				m.AuthMiddleware.RequireAuthMiddleware(),
				m.AuthMiddleware.RequireFeatureFlagMiddleware(constants.CAN_MODIFY_MOVIES),
				c.MovieController.UpdateMovie,
			)
		}

		genres := apiV1.Group("/genres")
		{
			genres.POST(
				"/",
				m.AuthMiddleware.RequireAuthMiddleware(),
				m.AuthMiddleware.RequireFeatureFlagMiddleware(constants.CAN_MODIFY_GENRES),
				c.GenreController.CreateGenre,
			)
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
	FeatureFlagRepository    repositories.FeatureFlagRepository
	GenreRepository          repositories.GenreRepository
}

type Services struct {
	UserService        services.UserService
	UserProfileService services.UserProfileService
	MovieService       services.MovieService
	GenreService       services.GenreService
}

type Controllers struct {
	UserController        controllers.UserController
	UserProfileController controllers.UserProfileController
	MovieController       controllers.MovieController
	GenreController       controllers.GenreController
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
		FeatureFlagRepository:    repositories.NewFeatureFlagRepository(config.ConfigcatClient),
		GenreRepository:          repositories.NewGenreRepository(config.DB),
	}
}

func setupServices(repositories *Repositories) *Services {
	return &Services{
		UserService: services.NewUserService(
			config.DB,
			config.RedisClient,
			auth.NewAuthenticator(),
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
		GenreService: services.NewGenreService(
			config.DB,
			transaction.NewTransactionManager(),
			repositories.GenreRepository,
		),
	}
}

func setupControllers(services *Services) *Controllers {
	return &Controllers{
		UserController:        *controllers.NewUserController(&services.UserService),
		UserProfileController: *controllers.NewUserProfileController(&services.UserProfileService),
		MovieController:       *controllers.NewMovieController(&services.MovieService),
		GenreController:       *controllers.NewGenreController(&services.GenreService),
	}
}

func setupMiddlewares(repositories *Repositories) *Middlewares {
	return &Middlewares{
		AuthMiddleware:        *middlewares.NewAuthMiddleware(repositories.UserSessionRepository, repositories.FeatureFlagRepository),
		FilesUploadMiddleware: *middlewares.NewFilesUploadMiddleware(),
	}
}
