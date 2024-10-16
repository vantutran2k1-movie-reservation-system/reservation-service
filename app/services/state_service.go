package services

import (
	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/repositories"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/transaction"
	"gorm.io/gorm"
)

type StateService interface {
	CreateState(countryID uuid.UUID, name string, code *string) (*models.State, *errors.ApiError)
}

func NewStateService(
	db *gorm.DB,
	transactionManager transaction.TransactionManager,
	countryRepo repositories.CountryRepository,
	stateRepo repositories.StateRepository,
) StateService {
	return &stateService{
		db:                 db,
		transactionManager: transactionManager,
		countryRepo:        countryRepo,
		stateRepo:          stateRepo,
	}
}

type stateService struct {
	db                 *gorm.DB
	transactionManager transaction.TransactionManager
	countryRepo        repositories.CountryRepository
	stateRepo          repositories.StateRepository
}

func (s *stateService) CreateState(countryID uuid.UUID, name string, code *string) (*models.State, *errors.ApiError) {
	country, err := s.countryRepo.GetCountry(countryID)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	if country == nil {
		return nil, errors.NotFoundError("country does not exist")
	}

	state, err := s.stateRepo.GetStateByName(countryID, name)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	if state != nil {
		return nil, errors.BadRequestError("duplicate state name for this country")
	}

	state = &models.State{
		ID:        uuid.New(),
		Name:      name,
		Code:      code,
		CountryID: countryID,
	}
	if err := s.transactionManager.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		return s.stateRepo.CreateState(tx, state)
	}); err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return state, nil
}
