package services

import (
	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/repositories"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/transaction"
	"gorm.io/gorm"
)

type CountryService interface {
	CreateCountry(name string, code string) (*models.Country, *errors.ApiError)
}

func NewCountryService(
	db *gorm.DB,
	transactionManager transaction.TransactionManager,
	countryRepo repositories.CountryRepository,
) CountryService {
	return &countryService{
		db:                 db,
		transactionManager: transactionManager,
		countryRepo:        countryRepo,
	}
}

type countryService struct {
	db                 *gorm.DB
	transactionManager transaction.TransactionManager
	countryRepo        repositories.CountryRepository
}

func (s *countryService) CreateCountry(name string, code string) (*models.Country, *errors.ApiError) {
	c, err := s.countryRepo.GetCountryByName(name)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	if c != nil {
		return nil, errors.BadRequestError("duplicate country name")
	}

	c, err = s.countryRepo.GetCountryByCode(code)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	if c != nil {
		return nil, errors.BadRequestError("duplicate country code")
	}

	c = &models.Country{
		ID:   uuid.New(),
		Name: name,
		Code: code,
	}
	if err := s.transactionManager.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		return s.countryRepo.CreateCountry(tx, c)
	}); err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return c, nil
}
