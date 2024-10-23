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

type TheaterService interface {
	CreateTheater(req payloads.CreateTheaterRequest) (*models.Theater, *errors.ApiError)
}

func NewTheaterService(
	db *gorm.DB,
	transactionManager transaction.TransactionManager,
	theaterRepo repositories.TheaterRepository,
) TheaterService {
	return &theaterService{
		db:                 db,
		transactionManager: transactionManager,
		theaterRepo:        theaterRepo,
	}
}

type theaterService struct {
	db                 *gorm.DB
	transactionManager transaction.TransactionManager
	theaterRepo        repositories.TheaterRepository
}

func (s *theaterService) CreateTheater(req payloads.CreateTheaterRequest) (*models.Theater, *errors.ApiError) {
	t, err := s.theaterRepo.GetTheaterByName(req.Name)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	if t != nil {
		return nil, errors.BadRequestError("duplicate theater name")
	}

	t = &models.Theater{
		ID:   uuid.New(),
		Name: req.Name,
	}
	if err := s.transactionManager.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		return s.theaterRepo.CreateTheater(tx, t)
	}); err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return t, nil
}
