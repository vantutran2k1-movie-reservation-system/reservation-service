package services

import (
	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/repositories"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/transaction"
	"gorm.io/gorm"
)

type GenreService interface {
	GetGenre(id uuid.UUID) (*models.Genre, *errors.ApiError)
	GetGenres() ([]*models.Genre, *errors.ApiError)
	CreateGenre(name string) (*models.Genre, *errors.ApiError)
}

func NewGenreService(db *gorm.DB, transactionManager transaction.TransactionManager, genreRepo repositories.GenreRepository) GenreService {
	return &genreService{db: db, transactionManager: transactionManager, genreRepo: genreRepo}
}

type genreService struct {
	db                 *gorm.DB
	transactionManager transaction.TransactionManager
	genreRepo          repositories.GenreRepository
}

func (s *genreService) GetGenre(id uuid.UUID) (*models.Genre, *errors.ApiError) {
	g, err := s.genreRepo.GetGenre(id)
	if err != nil {
		if errors.IsRecordNotFoundError(err) {
			return nil, errors.NotFoundError("genre not found")
		}

		return nil, errors.InternalServerError(err.Error())
	}

	return g, nil
}

func (s *genreService) GetGenres() ([]*models.Genre, *errors.ApiError) {
	genres, err := s.genreRepo.GetGenres()
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return genres, nil
}

func (s *genreService) CreateGenre(name string) (*models.Genre, *errors.ApiError) {
	_, err := s.genreRepo.GetGenreByName(name)
	if err == nil {
		return nil, errors.BadRequestError("Duplicate genre name")
	}
	if !errors.IsRecordNotFoundError(err) {
		return nil, errors.InternalServerError(err.Error())
	}

	g := models.Genre{
		ID:   uuid.New(),
		Name: name,
	}
	if err := s.transactionManager.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		return s.genreRepo.CreateGenre(tx, &g)
	}); err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return &g, nil
}
