package services

import (
	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/constants"
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
	GetShow(id uuid.UUID, userEmail *string) (*models.Show, *errors.ApiError)
	GetShows(status constants.ShowStatus, limit, offset int) ([]*models.Show, *errors.ApiError)
	CreateShow(req payloads.CreateShowRequest) (*models.Show, *errors.ApiError)
	ScheduleUpdateShowStatus() error
}

func NewShowService(
	db *gorm.DB,
	transactionManager transaction.TransactionManager,
	showRepo repositories.ShowRepository,
	movieRepo repositories.MovieRepository,
	theaterRepo repositories.TheaterRepository,
	featureFlagRepo repositories.FeatureFlagRepository,
) ShowService {
	return &showService{
		db:                 db,
		transactionManager: transactionManager,
		showRepo:           showRepo,
		movieRepo:          movieRepo,
		theaterRepo:        theaterRepo,
		featureFlagRepo:    featureFlagRepo,
	}
}

type showService struct {
	db                 *gorm.DB
	transactionManager transaction.TransactionManager
	showRepo           repositories.ShowRepository
	movieRepo          repositories.MovieRepository
	theaterRepo        repositories.TheaterRepository
	featureFlagRepo    repositories.FeatureFlagRepository
}

func (s *showService) GetShow(id uuid.UUID, userEmail *string) (*models.Show, *errors.ApiError) {
	show, err := s.showRepo.GetShow(filters.ShowFilter{
		Filter: &filters.SingleFilter{},
		Id:     &filters.Condition{Operator: filters.OpEqual, Value: id.String()},
	})
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	if show == nil || !(show.Status == constants.Active || s.adminUser(userEmail)) {
		return nil, errors.NotFoundError("show not found")
	}

	return show, nil
}

func (s *showService) GetShows(status constants.ShowStatus, limit, offset int) ([]*models.Show, *errors.ApiError) {
	shows, err := s.showRepo.GetShows(filters.ShowFilter{
		Filter: &filters.MultiFilter{Limit: &limit, Offset: &offset},
		Status: &filters.Condition{Operator: filters.OpEqual, Value: status},
	})
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return shows, nil
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

func (s *showService) ScheduleUpdateShowStatus() error {
	if err := s.transactionManager.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		if err := s.showRepo.ScheduleActivateShows(tx, time.Hour*72); err != nil {
			return err
		}

		return s.showRepo.ScheduleCompleteShows(tx)
	}); err != nil {
		return err
	}

	return nil
}

func (s *showService) adminUser(email *string) bool {
	return email != nil && s.featureFlagRepo.HasFlagEnabled(*email, constants.CanModifyShows)
}
