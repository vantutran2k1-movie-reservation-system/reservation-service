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
	"time"
)

type ShowService interface {
	CreateShow(req payloads.CreateShowRequest) (*models.Show, *errors.ApiError)
}

func NewShowService(
	db *gorm.DB,
	transactionManager transaction.TransactionManager,
	showRepo repositories.ShowRepository,
	movieRepo repositories.MovieRepository,
	theaterRepo repositories.TheaterRepository,
) ShowService {
	return &showService{
		db:                 db,
		transactionManager: transactionManager,
		showRepo:           showRepo,
		movieRepo:          movieRepo,
		theaterRepo:        theaterRepo,
	}
}

type showService struct {
	db                 *gorm.DB
	transactionManager transaction.TransactionManager
	showRepo           repositories.ShowRepository
	movieRepo          repositories.MovieRepository
	theaterRepo        repositories.TheaterRepository
}

func (s *showService) CreateShow(req payloads.CreateShowRequest) (*models.Show, *errors.ApiError) {
	movie, err := s.movieRepo.GetMovie(filters.MovieFilter{
		Filter: &filters.SingleFilter{},
		ID:     &filters.Condition{Operator: filters.OpEqual, Value: req.MovieId},
	}, false)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	if movie == nil {
		return nil, errors.BadRequestError("movie not found")
	}

	theater, err := s.theaterRepo.GetTheater(filters.TheaterFilter{
		Filter: &filters.SingleFilter{},
		ID:     &filters.Condition{Operator: filters.OpEqual, Value: req.TheaterId},
	}, false)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	if theater == nil {
		return nil, errors.BadRequestError("theater not found")
	}

	valid, err := s.showRepo.IsShowInValidTimeRange(req.TheaterId, req.StartTime, req.EndTime)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	if !valid {
		return nil, errors.BadRequestError("invalid time range for this show")
	}

	currentTime := time.Now().UTC()
	show := &models.Show{
		Id:        uuid.New(),
		MovieId:   &req.MovieId,
		TheaterId: &req.TheaterId,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Status:    req.Status,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}
	if err := s.transactionManager.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		return s.showRepo.CreateShow(tx, show)
	}); err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return show, nil
}
