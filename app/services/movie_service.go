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
	CreateMovie(createdBy uuid.UUID, title string, description *string, releaseDate string, duration int, language *string, rating *float64) (*models.Movie, *errors.ApiError)
	UpdateMovie(id, updatedBy uuid.UUID, title string, description *string, releaseDate string, duration int, language *string, rating *float64) (*models.Movie, *errors.ApiError)
}

type movieService struct {
	db                 *gorm.DB
	transactionManager transaction.TransactionManager
	movieRepo          repositories.MovieRepository
}

func NewMovieService(
	db *gorm.DB,
	transactionManager transaction.TransactionManager,
	movieRepo repositories.MovieRepository,
) MovieService {
	return &movieService{
		db:                 db,
		transactionManager: transactionManager,
		movieRepo:          movieRepo,
	}
}

func (s *movieService) CreateMovie(createdBy uuid.UUID, title string, description *string, releaseDate string, duration int, language *string, rating *float64) (*models.Movie, *errors.ApiError) {
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
		CreatedBy:       createdBy,
		LastUpdatedBy:   createdBy,
	}
	if err := s.transactionManager.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		return s.movieRepo.CreateMovie(tx, &m)
	}); err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return &m, nil
}

func (s *movieService) UpdateMovie(id, updatedBy uuid.UUID, title string, description *string, releaseDate string, duration int, language *string, rating *float64) (*models.Movie, *errors.ApiError) {
	m, err := s.movieRepo.GetMovie(id)
	if err != nil {
		if errors.IsRecordNotFoundError(err) {
			return nil, errors.NotFoundError("Movie not found")
		}

		return nil, errors.InternalServerError(err.Error())
	}

	m.Title = title
	m.Description = description
	m.ReleaseDate = releaseDate
	m.DurationMinutes = duration
	m.Language = language
	m.Rating = rating
	m.UpdatedAt = time.Now().UTC()
	m.LastUpdatedBy = updatedBy
	if err := s.transactionManager.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		return s.movieRepo.UpdateMovie(tx, m)
	}); err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return m, nil
}
