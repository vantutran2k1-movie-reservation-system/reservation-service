package services

import (
	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/repositories"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/transaction"
	"gorm.io/gorm"
)

type CityService interface {
	CreateCity(countryID, stateID uuid.UUID, name string) (*models.City, *errors.ApiError)
}

func NewCityService(
	db *gorm.DB,
	transactionManager transaction.TransactionManager,
	countryRepo repositories.CountryRepository,
	stateRepo repositories.StateRepository,
	cityRepo repositories.CityRepository,
) CityService {
	return &cityService{
		db:                 db,
		transactionManager: transactionManager,
		countryRepo:        countryRepo,
		stateRepo:          stateRepo,
		cityRepo:           cityRepo,
	}
}

type cityService struct {
	db                 *gorm.DB
	transactionManager transaction.TransactionManager
	countryRepo        repositories.CountryRepository
	stateRepo          repositories.StateRepository
	cityRepo           repositories.CityRepository
}

func (s *cityService) CreateCity(countryID, stateID uuid.UUID, name string) (*models.City, *errors.ApiError) {
	country, err := s.countryRepo.GetCountry(countryID)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	if country == nil {
		return nil, errors.NotFoundError("country does not exist")
	}

	state, err := s.stateRepo.GetState(stateID)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	if state == nil {
		return nil, errors.NotFoundError("state does not exist")
	}

	city, err := s.cityRepo.GetCityByName(stateID, name)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	if city != nil {
		return nil, errors.BadRequestError("duplicate city name for this state")
	}

	city = &models.City{
		ID:      uuid.New(),
		Name:    name,
		StateID: stateID,
	}
	if err := s.transactionManager.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		return s.cityRepo.CreateCity(tx, city)
	}); err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return city, nil
}
