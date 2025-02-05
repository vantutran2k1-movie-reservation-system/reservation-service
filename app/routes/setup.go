package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/auth"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/constants"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/controllers"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/middlewares"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/repositories"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/services"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/transaction"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/config"
	"log"
	"time"
)

var r *Repositories
var s *Services
var c *Controllers
var m *Middlewares
var router *gin.Engine

var authenticator = auth.NewAuthenticator()
var transactionManager = transaction.NewTransactionManager()

type Repositories struct {
	UserRepository                  repositories.UserRepository
	UserRegistrationTokenRepository repositories.UserRegistrationTokenRepository
	LoginTokenRepository            repositories.LoginTokenRepository
	UserSessionRepository           repositories.UserSessionRepository
	UserProfileRepository           repositories.UserProfileRepository
	ProfilePictureRepository        repositories.ProfilePictureRepository
	MovieRepository                 repositories.MovieRepository
	FeatureFlagRepository           repositories.FeatureFlagRepository
	GenreRepository                 repositories.GenreRepository
	PasswordResetTokenRepository    repositories.PasswordResetTokenRepository
	MovieGenreRepository            repositories.MovieGenreRepository
	CountryRepository               repositories.CountryRepository
	StateRepository                 repositories.StateRepository
	CityRepository                  repositories.CityRepository
	TheaterRepository               repositories.TheaterRepository
	TheaterLocationRepository       repositories.TheaterLocationRepository
	SeatRepository                  repositories.SeatRepository
	ShowRepository                  repositories.ShowRepository
	NotificationRepository          repositories.NotificationRepository
}

type Services struct {
	UserService        services.UserService
	UserProfileService services.UserProfileService
	MovieService       services.MovieService
	GenreService       services.GenreService
	LocationService    services.LocationService
	TheaterService     services.TheaterService
	ShowService        services.ShowService
	RateLimiterService services.RateLimiterService
}

type Controllers struct {
	UserController        controllers.UserController
	UserProfileController controllers.UserProfileController
	MovieController       controllers.MovieController
	GenreController       controllers.GenreController
	LocationController    controllers.LocationController
	TheaterController     controllers.TheaterController
	ShowController        controllers.ShowController
}

type Middlewares struct {
	AuthMiddleware        middlewares.AuthMiddleware
	FilesUploadMiddleware middlewares.FilesUploadMiddleware
	ContextMiddleware     middlewares.ContextMiddleware
	RateLimitMiddleware   middlewares.RateLimitMiddleware
}

func setupRepositories() {
	r = &Repositories{
		UserRepository:                  repositories.NewUserRepository(config.DB),
		UserRegistrationTokenRepository: repositories.NewUserRegistrationTokenRepository(config.DB),
		LoginTokenRepository:            repositories.NewLoginTokenRepository(config.DB),
		UserSessionRepository:           repositories.NewUserSessionRepository(config.RedisClient),
		UserProfileRepository:           repositories.NewUserProfileRepository(config.DB),
		ProfilePictureRepository:        repositories.NewProfilePictureRepository(config.MinioClient),
		MovieRepository:                 repositories.NewMovieRepository(config.DB),
		FeatureFlagRepository:           repositories.NewFeatureFlagRepository(config.ConfigcatClient),
		GenreRepository:                 repositories.NewGenreRepository(config.DB),
		PasswordResetTokenRepository:    repositories.NewPasswordResetTokenRepository(config.DB),
		MovieGenreRepository:            repositories.NewMovieGenreRepository(config.DB),
		CountryRepository:               repositories.NewCountryRepository(config.DB),
		StateRepository:                 repositories.NewStateRepository(config.DB),
		CityRepository:                  repositories.NewCityRepository(config.DB),
		TheaterRepository:               repositories.NewTheaterRepository(config.DB),
		TheaterLocationRepository:       repositories.NewTheaterLocationRepository(config.DB),
		SeatRepository:                  repositories.NewSeatRepository(config.DB),
		ShowRepository:                  repositories.NewShowRepository(config.DB),
		NotificationRepository:          repositories.NewNotificationRepository(config.KafkaProducerClient),
	}
}

func setupServices(repositories *Repositories) {
	s = &Services{
		UserService: services.NewUserService(
			config.DB,
			config.RedisClient,
			authenticator,
			transactionManager,
			repositories.UserRepository,
			repositories.UserProfileRepository,
			repositories.LoginTokenRepository,
			repositories.UserSessionRepository,
			repositories.PasswordResetTokenRepository,
			repositories.UserRegistrationTokenRepository,
			repositories.NotificationRepository,
		),
		UserProfileService: services.NewUserProfileService(
			config.DB,
			transactionManager,
			repositories.UserRepository,
			repositories.UserProfileRepository,
			repositories.ProfilePictureRepository,
		),
		MovieService: services.NewMovieService(
			config.DB,
			transactionManager,
			repositories.MovieRepository,
			repositories.GenreRepository,
			repositories.MovieGenreRepository,
			repositories.FeatureFlagRepository,
		),
		GenreService: services.NewGenreService(
			config.DB,
			transactionManager,
			repositories.GenreRepository,
			repositories.MovieGenreRepository,
		),
		LocationService: services.NewLocationService(
			config.DB,
			transactionManager,
			repositories.CountryRepository,
			repositories.StateRepository,
			repositories.CityRepository,
		),
		TheaterService: services.NewTheaterService(
			config.DB,
			transactionManager,
			repositories.TheaterRepository,
			repositories.TheaterLocationRepository,
			repositories.SeatRepository,
			repositories.CityRepository,
			services.NewUserLocationService(config.AppEnv.UserLocationApiUrl, config.AppEnv.UserLocationApiTimeout),
		),
		ShowService: services.NewShowService(
			config.DB,
			transactionManager,
			repositories.ShowRepository,
			repositories.MovieRepository,
			repositories.TheaterRepository,
		),
		RateLimiterService: services.NewRateLimiterService(
			config.RedisClient,
			config.AppEnv.MaxRequestsPerMinute,
			time.Minute,
		),
	}
}

func setupControllers(services *Services) {
	c = &Controllers{
		UserController:        *controllers.NewUserController(&services.UserService),
		UserProfileController: *controllers.NewUserProfileController(&services.UserProfileService),
		MovieController:       *controllers.NewMovieController(&services.MovieService),
		GenreController:       *controllers.NewGenreController(&services.GenreService),
		LocationController:    *controllers.NewLocationController(&services.LocationService),
		TheaterController:     *controllers.NewTheaterController(&services.TheaterService),
		ShowController:        *controllers.NewShowController(&services.ShowService),
	}
}

func setupMiddlewares(repositories *Repositories) {
	m = &Middlewares{
		AuthMiddleware:        *middlewares.NewAuthMiddleware(repositories.UserSessionRepository, repositories.FeatureFlagRepository),
		FilesUploadMiddleware: *middlewares.NewFilesUploadMiddleware(),
		ContextMiddleware:     *middlewares.NewContextMiddleware(),
		RateLimitMiddleware:   *middlewares.NewRateLimitMiddleware(),
	}
}

func setupRouter() {
	if config.AppEnv.GinMode == constants.GinReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	router = gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{config.AppEnv.PlatformUiEndpoint},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", constants.ContentType, "Authorization", constants.UserVerificationToken},
		AllowCredentials: true,
	}))
}

func registerCronJobs() {
	_, err := config.CronJobManager.AddFunc("0 0 * * * *", func() {
		if err := s.ShowService.ScheduleUpdateShowStatus(); err != nil {
			log.Println(err)
		}
	})
	if err != nil {
		log.Fatal(err)
	}
}

func setupRoutes() {
	setupRepositories()
	setupServices(r)
	setupControllers(s)
	setupMiddlewares(r)
	registerCronJobs()
	setupRouter()
}
