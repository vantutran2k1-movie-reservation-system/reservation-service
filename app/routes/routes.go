package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/constants"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/middlewares"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/config"
)

func RegisterRoutes() *gin.Engine {
	setupRoutes()

	apiV1 := router.Group("/api/v1")
	apiV1.Use(m.ContextMiddleware.AddRequestContext())
	apiV1.Use(m.RateLimitMiddleware.NotExceedMaxRequests(s.RateLimiterService))
	{
		users := apiV1.Group("/users")
		{
			users.GET("/me", m.AuthMiddleware.RequireAuthMiddleware(), c.UserController.GetCurrentUser)
			users.GET(
				"/:userId",
				m.AuthMiddleware.RequireAuthMiddleware(),
				m.AuthMiddleware.RequireFeatureFlagMiddleware(constants.CanModifyUsers),
				c.UserController.GetUser,
			)
			users.GET("/exists", c.UserController.UserExistsByEmail)
			users.POST("/", c.UserController.CreateUser)

			users.POST("/login", c.UserController.LoginUser)
			users.POST("/logout", m.AuthMiddleware.RequireAuthMiddleware(), c.UserController.LogoutUser)

			users.POST("/verify", c.UserController.VerifyUser)

			users.PUT("/password", m.AuthMiddleware.RequireAuthMiddleware(), c.UserController.UpdateUserPassword)
			users.POST("/password-reset-token", c.UserController.CreatePasswordResetToken)
			users.POST("/password-reset", c.UserController.ResetPassword)
		}

		profiles := apiV1.Group("/profiles")
		profiles.Use(m.AuthMiddleware.RequireAuthMiddleware())
		{
			profiles.GET("/me", c.UserProfileController.GetProfileByUserID)
			profiles.PUT("/", c.UserProfileController.UpdateUserProfile)

			profiles.PUT(
				"/profile-picture",
				m.FilesUploadMiddleware.RequireNumberOfUploadedFilesMiddleware(constants.ProfilePictureRequestFormKey, 1),
				m.FilesUploadMiddleware.IsAllowedFileTypeMiddleware(constants.ProfilePictureRequestFormKey, middlewares.DefaultImageFileTypes),
				m.FilesUploadMiddleware.NotExceedMaxSizeLimitMiddleware(constants.ProfilePictureRequestFormKey, config.AppEnv.MaxProfilePictureFileSize),
				c.UserProfileController.UpdateProfilePicture,
			)
			profiles.DELETE("/profile-picture", c.UserProfileController.DeleteProfilePicture)
		}

		movies := apiV1.Group("/movies")
		{
			movies.GET("/:id", m.AuthMiddleware.OptionalAuthMiddleware(), c.MovieController.GetMovie)
			movies.GET("/", m.AuthMiddleware.OptionalAuthMiddleware(), c.MovieController.GetMovies)
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
			movies.DELETE(
				"/:id",
				m.AuthMiddleware.RequireAuthMiddleware(),
				m.AuthMiddleware.RequireFeatureFlagMiddleware(constants.CanModifyMovies),
				c.MovieController.DeleteMovie,
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
				c.GenreController.UpdateGenre,
			)
			genres.DELETE(
				"/:id",
				m.AuthMiddleware.RequireAuthMiddleware(),
				m.AuthMiddleware.RequireFeatureFlagMiddleware(constants.CanModifyGenres),
				c.GenreController.DeleteGenre,
			)
		}

		countries := apiV1.Group("/countries")
		{
			countries.GET("/", c.LocationController.GetCountries)
			countries.POST(
				"/",
				m.AuthMiddleware.RequireAuthMiddleware(),
				m.AuthMiddleware.RequireFeatureFlagMiddleware(constants.CanModifyLocations),
				c.LocationController.CreateCountry,
			)

			states := countries.Group("/:countryId/states")
			{
				states.GET("/", c.LocationController.GetStatesByCountry)
				states.POST(
					"/",
					m.AuthMiddleware.RequireAuthMiddleware(),
					m.AuthMiddleware.RequireFeatureFlagMiddleware(constants.CanModifyLocations),
					c.LocationController.CreateState,
				)

				cities := states.Group("/:stateId/cities")
				{
					cities.GET("/", c.LocationController.GetCitiesByState)
					cities.POST(
						"/",
						m.AuthMiddleware.RequireAuthMiddleware(),
						m.AuthMiddleware.RequireFeatureFlagMiddleware(constants.CanModifyLocations),
						c.LocationController.CreateCity,
					)
				}
			}
		}

		theaters := apiV1.Group("/theaters")
		{
			theaters.GET("/:theaterId", c.TheaterController.GetTheater)
			theaters.GET("/", c.TheaterController.GetTheaters)
			theaters.GET("/nearby", c.TheaterController.GetNearbyTheaters)
			theaters.POST(
				"/",
				m.AuthMiddleware.RequireAuthMiddleware(),
				m.AuthMiddleware.RequireFeatureFlagMiddleware(constants.CanModifyTheaters),
				c.TheaterController.CreateTheater,
			)

			theaterLocations := theaters.Group("/:theaterId/locations")
			{
				theaterLocations.POST(
					"/",
					m.AuthMiddleware.RequireAuthMiddleware(),
					m.AuthMiddleware.RequireFeatureFlagMiddleware(constants.CanModifyTheaters),
					c.TheaterController.CreateTheaterLocation,
				)
				theaterLocations.PUT(
					"/",
					m.AuthMiddleware.RequireAuthMiddleware(),
					m.AuthMiddleware.RequireFeatureFlagMiddleware(constants.CanModifyTheaters),
					c.TheaterController.UpdateTheaterLocation,
				)
			}

			seats := theaters.Group("/:theaterId/seats")
			{
				seats.POST(
					"/",
					m.AuthMiddleware.RequireAuthMiddleware(),
					m.AuthMiddleware.RequireFeatureFlagMiddleware(constants.CanModifyTheaters),
					c.TheaterController.CreateSeat,
				)
			}
		}

		shows := apiV1.Group("/shows")
		{
			shows.GET("/active", c.ShowController.GetActiveShows)
			shows.GET("/scheduled", c.ShowController.GetScheduledShows)
			shows.POST(
				"/",
				m.AuthMiddleware.RequireAuthMiddleware(),
				m.AuthMiddleware.RequireFeatureFlagMiddleware(constants.CanModifyShows),
				c.ShowController.CreateShow,
			)
		}
	}

	return router
}
