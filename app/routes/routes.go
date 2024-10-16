package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/constants"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/middlewares"
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
			users.POST("/password-reset-token", c.UserController.CreatePasswordResetToken)
			users.POST("/password-reset", c.UserController.ResetPassword)
		}

		profiles := apiV1.Group("/profiles")
		profiles.Use(m.AuthMiddleware.RequireAuthMiddleware())
		{
			profiles.GET("/me", c.UserProfileController.GetProfileByUserID)
			profiles.POST("/", c.UserProfileController.CreateUserProfile)
			profiles.PUT("/", c.UserProfileController.UpdateUserProfile)

			profiles.PUT(
				"/profile-picture",
				m.FilesUploadMiddleware.RequireNumberOfUploadedFilesMiddleware(constants.ProfilePictureRequestFormKey, 1),
				m.FilesUploadMiddleware.IsAllowedFileTypeMiddleware(constants.ProfilePictureRequestFormKey, middlewares.DEFAULT_IMAGE_FILE_TYPES),
				m.FilesUploadMiddleware.NotExceedMaxSizeLimitMiddleware(constants.ProfilePictureRequestFormKey, middlewares.GetMaxProfilePictureFileSize()),
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
				m.AuthMiddleware.RequireFeatureFlagMiddleware(constants.CanModifyMovies),
				c.MovieController.CreateMovie,
			)
			movies.PUT(
				"/:id",
				m.AuthMiddleware.RequireAuthMiddleware(),
				m.AuthMiddleware.RequireFeatureFlagMiddleware(constants.CanModifyMovies),
				c.MovieController.UpdateMovie,
			)
			movies.PUT(
				"/:id/genres",
				m.AuthMiddleware.RequireAuthMiddleware(),
				m.AuthMiddleware.RequireFeatureFlagMiddleware(constants.CanModifyMovies),
				c.MovieController.UpdateMovieGenres,
			)
		}

		genres := apiV1.Group("/genres")
		{
			genres.GET("/:id", c.GenreController.GetGenre)
			genres.GET("/", c.GenreController.GetGenres)
			genres.POST(
				"/",
				m.AuthMiddleware.RequireAuthMiddleware(),
				m.AuthMiddleware.RequireFeatureFlagMiddleware(constants.CanModifyGenres),
				c.GenreController.CreateGenre,
			)
			genres.PUT(
				"/:id",
				m.AuthMiddleware.RequireAuthMiddleware(),
				m.AuthMiddleware.RequireFeatureFlagMiddleware(constants.CanModifyGenres),
				c.GenreController.UpdateGenre)
		}

		countries := apiV1.Group("/countries")
		{
			countries.GET("/", c.CountryController.GetCountries)
			countries.POST(
				"/",
				m.AuthMiddleware.RequireAuthMiddleware(),
				m.AuthMiddleware.RequireFeatureFlagMiddleware(constants.CanModifyLocations),
				c.CountryController.CreateCountry,
			)
		}

		states := apiV1.Group("/countries/:id/states")
		{
			states.POST(
				"/",
				m.AuthMiddleware.RequireAuthMiddleware(),
				m.AuthMiddleware.RequireFeatureFlagMiddleware(constants.CanModifyLocations),
				c.StateController.CreateState,
			)
		}
	}

	return router
}
