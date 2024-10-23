package services

import (
	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/payloads"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/repositories"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/transaction"
	"gorm.io/gorm"
)

type GenreService interface {
	GetGenre(id uuid.UUID) (*models.Genre, *errors.ApiError)
	GetGenres() ([]*models.Genre, *errors.ApiError)
	CreateGenre(req payloads.CreateGenreRequest) (*models.Genre, *errors.ApiError)
	UpdateGenre(id uuid.UUID, req payloads.UpdateGenreRequest) (*models.Genre, *errors.ApiError)
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
		return nil, errors.InternalServerError(err.Error())
	}
	if g == nil {
		return nil, errors.NotFoundError("genre not found")
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

func (s *genreService) CreateGenre(req payloads.CreateGenreRequest) (*models.Genre, *errors.ApiError) {
	g, err := s.genreRepo.GetGenreByName(req.Name)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	if g != nil {
		return nil, errors.BadRequestError("duplicate genre name")
	}

	g = &models.Genre{
		ID:   uuid.New(),
		Name: req.Name,
	}
	if err := s.transactionManager.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		return s.genreRepo.CreateGenre(tx, g)
	}); err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return g, nil
}

func (s *genreService) UpdateGenre(id uuid.UUID, req payloads.UpdateGenreRequest) (*models.Genre, *errors.ApiError) {
	g, err := s.genreRepo.GetGenre(id)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	if g == nil {
		return nil, errors.NotFoundError("genre not found")
	}

	g, err = s.genreRepo.GetGenreByName(req.Name)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	if g != nil {
		return nil, errors.BadRequestError("duplicate genre name")
	}

	g = &models.Genre{
		ID:   id,
		Name: req.Name,
	}
	if err := s.transactionManager.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		return s.genreRepo.UpdateGenre(tx, g)
	}); err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return g, nil
}
