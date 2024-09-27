package services

import (
	"time"

	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/repositories"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/transaction"
	"gorm.io/gorm"
)

type MovieService interface {
	CreateMovie(title string, description *string, releaseDate string, duration int, language *string, rating *float64) (*models.Movie, *errors.ApiError)
}

type movieService struct {
	db        *gorm.DB
	movieRepo repositories.MovieRepository
}

func NewMovieService(db *gorm.DB, movieRepo repositories.MovieRepository) MovieService {
	return &movieService{db: db, movieRepo: movieRepo}
}

func (s *movieService) CreateMovie(title string, description *string, releaseDate string, duration int, language *string, rating *float64) (*models.Movie, *errors.ApiError) {
	m := models.Movie{
		ID:              uuid.New(),
		Title:           title,
		Description:     description,
		ReleaseDate:     releaseDate,
		DurationMinutes: duration,
		Language:        language,
		Rating:          rating,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}
	if err := transaction.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		return s.movieRepo.CreateMovie(tx, &m)
	}); err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return &m, nil
}
