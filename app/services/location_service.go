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

type LocationService interface {
	GetCountries() ([]*models.Country, *errors.ApiError)
	CreateCountry(req payloads.CreateCountryRequest) (*models.Country, *errors.ApiError)
	GetStatesByCountry(countryID uuid.UUID) ([]*models.State, *errors.ApiError)
	CreateState(countryID uuid.UUID, req payloads.CreateStateRequest) (*models.State, *errors.ApiError)
	CreateCity(countryID, stateID uuid.UUID, req payloads.CreateCityRequest) (*models.City, *errors.ApiError)
}

func NewLocationService(
	db *gorm.DB,
	transactionManager transaction.TransactionManager,
	countryRepo repositories.CountryRepository,
	stateRepo repositories.StateRepository,
	cityRepo repositories.CityRepository,
) LocationService {
	return &locationService{
		db:                 db,
		transactionManager: transactionManager,
		countryRepo:        countryRepo,
		stateRepo:          stateRepo,
		cityRepo:           cityRepo,
	}
}

type locationService struct {
	db                 *gorm.DB
	transactionManager transaction.TransactionManager
	countryRepo        repositories.CountryRepository
	stateRepo          repositories.StateRepository
	cityRepo           repositories.CityRepository
}

func (s *locationService) GetCountries() ([]*models.Country, *errors.ApiError) {
	countries, err := s.countryRepo.GetCountries()
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return countries, nil
}

func (s *locationService) CreateCountry(req payloads.CreateCountryRequest) (*models.Country, *errors.ApiError) {
	c, err := s.countryRepo.GetCountryByName(req.Name)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	if c != nil {
		return nil, errors.BadRequestError("duplicate country name")
	}

	c, err = s.countryRepo.GetCountryByCode(req.Code)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	if c != nil {
		return nil, errors.BadRequestError("duplicate country code")
	}

	c = &models.Country{
		ID:   uuid.New(),
		Name: req.Name,
		Code: req.Code,
	}
	if err := s.transactionManager.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		return s.countryRepo.CreateCountry(tx, c)
	}); err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return c, nil
}

func (s *locationService) GetStatesByCountry(countryID uuid.UUID) ([]*models.State, *errors.ApiError) {
	country, err := s.countryRepo.GetCountry(countryID)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	if country == nil {
		return nil, errors.NotFoundError("country does not exist")
	}

	states, err := s.stateRepo.GetStatesByCountry(countryID)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return states, nil
}

func (s *locationService) CreateState(countryID uuid.UUID, req payloads.CreateStateRequest) (*models.State, *errors.ApiError) {
	country, err := s.countryRepo.GetCountry(countryID)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	if country == nil {
		return nil, errors.NotFoundError("country does not exist")
	}

	state, err := s.stateRepo.GetStateByName(countryID, req.Name)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	if state != nil {
		return nil, errors.BadRequestError("duplicate state name for this country")
	}

	state = &models.State{
		ID:        uuid.New(),
		Name:      req.Name,
		Code:      req.Code,
		CountryID: countryID,
	}
	if err := s.transactionManager.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		return s.stateRepo.CreateState(tx, state)
	}); err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return state, nil
}

func (s *locationService) CreateCity(countryID, stateID uuid.UUID, req payloads.CreateCityRequest) (*models.City, *errors.ApiError) {
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

	city, err := s.cityRepo.GetCity(payloads.GetCityFilter{StateID: &stateID, Name: &req.Name})
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	if city != nil {
		return nil, errors.BadRequestError("duplicate city name for this state")
	}

	city = &models.City{
		ID:      uuid.New(),
		Name:    req.Name,
		StateID: stateID,
	}
	if err := s.transactionManager.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		return s.cityRepo.CreateCity(tx, city)
	}); err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return city, nil
}
