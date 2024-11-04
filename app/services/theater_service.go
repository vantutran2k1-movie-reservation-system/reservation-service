package services

import (
	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/filters"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/payloads"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/repositories"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/transaction"
	"gorm.io/gorm"
)

type TheaterService interface {
	GetTheater(id uuid.UUID, includeLocation bool) (*models.Theater, *errors.ApiError)
	CreateTheater(req payloads.CreateTheaterRequest) (*models.Theater, *errors.ApiError)
	CreateTheaterLocation(theaterID uuid.UUID, req payloads.CreateTheaterLocationRequest) (*models.TheaterLocation, *errors.ApiError)
}

func NewTheaterService(
	db *gorm.DB,
	transactionManager transaction.TransactionManager,
	theaterRepo repositories.TheaterRepository,
	theaterLocationRepo repositories.TheaterLocationRepository,
	cityRepo repositories.CityRepository,
) TheaterService {
	return &theaterService{
		db:                  db,
		transactionManager:  transactionManager,
		theaterRepo:         theaterRepo,
		theaterLocationRepo: theaterLocationRepo,
		cityRepo:            cityRepo,
	}
}

type theaterService struct {
	db                  *gorm.DB
	transactionManager  transaction.TransactionManager
	theaterRepo         repositories.TheaterRepository
	theaterLocationRepo repositories.TheaterLocationRepository
	cityRepo            repositories.CityRepository
}

func (s *theaterService) GetTheater(id uuid.UUID, includeLocation bool) (*models.Theater, *errors.ApiError) {
	filter := filters.TheaterFilter{
		Filter: &filters.SingleFilter{},
		ID:     &filters.Condition{Operator: filters.OpEqual, Value: id},
	}
	t, err := s.theaterRepo.GetTheater(filter, includeLocation)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	if t == nil {
		return nil, errors.NotFoundError("theater not found")
	}

	return t, nil
}

func (s *theaterService) CreateTheater(req payloads.CreateTheaterRequest) (*models.Theater, *errors.ApiError) {
	filter := filters.TheaterFilter{
		Filter: &filters.SingleFilter{},
		Name:   &filters.Condition{Operator: filters.OpEqual, Value: &req.Name},
	}
	t, err := s.theaterRepo.GetTheater(filter, false)
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

func (s *theaterService) CreateTheaterLocation(theaterID uuid.UUID, req payloads.CreateTheaterLocationRequest) (*models.TheaterLocation, *errors.ApiError) {
	t, err := s.theaterRepo.GetTheater(filters.TheaterFilter{
		Filter: &filters.SingleFilter{},
		ID:     &filters.Condition{Operator: filters.OpEqual, Value: theaterID},
	}, true)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	if t == nil {
		return nil, errors.NotFoundError("theater not found")
	}
	if t.Location != nil {
		return nil, errors.BadRequestError("duplicate location for this theater")
	}

	cityFilter := filters.CityFilter{
		Filter: &filters.SingleFilter{Logic: filters.And},
		ID:     &filters.Condition{Operator: filters.OpEqual, Value: req.CityID},
	}
	c, err := s.cityRepo.GetCity(cityFilter)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	if c == nil {
		return nil, errors.BadRequestError("invalid city id")
	}

	l := &models.TheaterLocation{
		ID:         uuid.New(),
		TheaterID:  theaterID,
		CityID:     req.CityID,
		Address:    req.Address,
		PostalCode: req.PostalCode,
		Latitude:   req.Latitude,
		Longitude:  req.Longitude,
	}
	if err := s.transactionManager.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		return s.theaterLocationRepo.CreateTheaterLocation(tx, l)
	}); err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return l, nil
}
