package routes

import (
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/auth"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/controllers"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/middlewares"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/repositories"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/services"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/transaction"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/config"
)

var authenticator = auth.NewAuthenticator()
var transactionManager = transaction.NewTransactionManager()

type Repositories struct {
	UserRepository               repositories.UserRepository
	LoginTokenRepository         repositories.LoginTokenRepository
	UserSessionRepository        repositories.UserSessionRepository
	UserProfileRepository        repositories.UserProfileRepository
	ProfilePictureRepository     repositories.ProfilePictureRepository
	MovieRepository              repositories.MovieRepository
	FeatureFlagRepository        repositories.FeatureFlagRepository
	GenreRepository              repositories.GenreRepository
	PasswordResetTokenRepository repositories.PasswordResetTokenRepository
	MovieGenreRepository         repositories.MovieGenreRepository
	CountryRepository            repositories.CountryRepository
	StateRepository              repositories.StateRepository
	CityRepository               repositories.CityRepository
	TheaterRepository            repositories.TheaterRepository
}

type Services struct {
	UserService        services.UserService
	UserProfileService services.UserProfileService
	MovieService       services.MovieService
	GenreService       services.GenreService
	CountryService     services.CountryService
	StateService       services.StateService
	CityService        services.CityService
	TheaterService     services.TheaterService
}

type Controllers struct {
	UserController        controllers.UserController
	UserProfileController controllers.UserProfileController
	MovieController       controllers.MovieController
	GenreController       controllers.GenreController
	CountryController     controllers.CountryController
	StateController       controllers.StateController
	CityController        controllers.CityController
	TheaterController     controllers.TheaterController
}

type Middlewares struct {
	AuthMiddleware        middlewares.AuthMiddleware
	FilesUploadMiddleware middlewares.FilesUploadMiddleware
}

func setupRepositories() *Repositories {
	return &Repositories{
		UserRepository:               repositories.NewUserRepository(config.DB),
		LoginTokenRepository:         repositories.NewLoginTokenRepository(config.DB),
		UserSessionRepository:        repositories.NewUserSessionRepository(config.RedisClient),
		UserProfileRepository:        repositories.NewUserProfileRepository(config.DB),
		ProfilePictureRepository:     repositories.NewProfilePictureRepository(config.MinioClient),
		MovieRepository:              repositories.NewMovieRepository(config.DB),
		FeatureFlagRepository:        repositories.NewFeatureFlagRepository(config.ConfigcatClient),
		GenreRepository:              repositories.NewGenreRepository(config.DB),
		PasswordResetTokenRepository: repositories.NewPasswordResetTokenRepository(config.DB),
		MovieGenreRepository:         repositories.NewMovieGenreRepository(config.DB),
		CountryRepository:            repositories.NewCountryRepository(config.DB),
		StateRepository:              repositories.NewStateRepository(config.DB),
		CityRepository:               repositories.NewCityRepository(config.DB),
		TheaterRepository:            repositories.NewTheaterRepository(config.DB),
	}
}

func setupServices(repositories *Repositories) *Services {
	return &Services{
		UserService: services.NewUserService(
			config.DB,
			config.RedisClient,
			authenticator,
			transactionManager,
			repositories.UserRepository,
			repositories.LoginTokenRepository,
			repositories.UserSessionRepository,
			repositories.PasswordResetTokenRepository,
		),
		UserProfileService: services.NewUserProfileService(
			config.DB,
			transactionManager,
			repositories.UserProfileRepository,
			repositories.ProfilePictureRepository,
		),
		MovieService: services.NewMovieService(
			config.DB,
			transactionManager,
			repositories.MovieRepository,
			repositories.GenreRepository,
			repositories.MovieGenreRepository,
		),
		GenreService: services.NewGenreService(
			config.DB,
			transactionManager,
			repositories.GenreRepository,
		),
		CountryService: services.NewCountryService(
			config.DB,
			transactionManager,
			repositories.CountryRepository,
		),
		StateService: services.NewStateService(
			config.DB,
			transactionManager,
			repositories.CountryRepository,
			repositories.StateRepository,
		),
		CityService: services.NewCityService(
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
		),
	}
}

func setupControllers(services *Services) *Controllers {
	return &Controllers{
		UserController:        *controllers.NewUserController(&services.UserService),
		UserProfileController: *controllers.NewUserProfileController(&services.UserProfileService),
		MovieController:       *controllers.NewMovieController(&services.MovieService),
		GenreController:       *controllers.NewGenreController(&services.GenreService),
		CountryController:     *controllers.NewCountryController(&services.CountryService),
		StateController:       *controllers.NewStateController(&services.StateService),
		CityController:        *controllers.NewCityController(&services.CityService),
		TheaterController:     *controllers.NewTheaterController(&services.TheaterService),
	}
}

func setupMiddlewares(repositories *Repositories) *Middlewares {
	return &Middlewares{
		AuthMiddleware:        *middlewares.NewAuthMiddleware(repositories.UserSessionRepository, repositories.FeatureFlagRepository),
		FilesUploadMiddleware: *middlewares.NewFilesUploadMiddleware(),
	}
}
