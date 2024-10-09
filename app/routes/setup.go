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
		UserRepository:               repositories.NewUserRepository(config.DB),
		LoginTokenRepository:         repositories.NewLoginTokenRepository(config.DB),
		UserSessionRepository:        repositories.NewUserSessionRepository(config.RedisClient),
		UserProfileRepository:        repositories.NewUserProfileRepository(config.DB),
		ProfilePictureRepository:     repositories.NewProfilePictureRepository(config.MinioClient),
		MovieRepository:              repositories.NewMovieRepository(config.DB),
		FeatureFlagRepository:        repositories.NewFeatureFlagRepository(config.ConfigcatClient),
		GenreRepository:              repositories.NewGenreRepository(config.DB),
		PasswordResetTokenRepository: repositories.NewPasswordResetTokenRepository(config.DB),
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
			repositories.PasswordResetTokenRepository,
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
