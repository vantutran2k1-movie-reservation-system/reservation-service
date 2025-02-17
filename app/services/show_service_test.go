package services

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/constants"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/filters"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_repositories"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_transaction"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/payloads"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"net/http"
	"testing"
	"time"
)

func TestShowService_GetShow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	flagRepo := mock_repositories.NewMockFeatureFlagRepository(ctrl)
	showRepo := mock_repositories.NewMockShowRepository(ctrl)
	service := NewShowService(nil, nil, showRepo, nil, nil, flagRepo)

	show := utils.GenerateShow()
	show.Status = constants.Completed
	email := utils.GetPointerOf("example@gmail.com")
	filter := filters.ShowFilter{
		Filter: &filters.SingleFilter{},
		Id:     &filters.Condition{Operator: filters.OpEqual, Value: show.Id.String()},
	}

	t.Run("success", func(t *testing.T) {
		showRepo.EXPECT().GetShow(filter).Return(show, nil).Times(1)
		flagRepo.EXPECT().HasFlagEnabled(*email, constants.CanModifyShows).Return(true).Times(1)

		result, err := service.GetShow(show.Id, email)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, show, result)
	})

	t.Run("error getting show", func(t *testing.T) {
		showRepo.EXPECT().GetShow(filter).Return(nil, errors.New("error getting show")).Times(1)

		result, err := service.GetShow(show.Id, email)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting show", err.Error())
	})

	t.Run("user not have permission", func(t *testing.T) {
		showRepo.EXPECT().GetShow(filter).Return(show, nil).Times(1)
		flagRepo.EXPECT().HasFlagEnabled(*email, constants.CanModifyShows).Return(false).Times(1)

		result, err := service.GetShow(show.Id, email)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "show not found", err.Error())
	})
}

func TestShowService_GetShows(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repositories.NewMockShowRepository(ctrl)
	service := NewShowService(nil, nil, repo, nil, nil, nil)

	shows := utils.GenerateShows(3)
	limit := 3
	offset := 0
	filter := filters.ShowFilter{
		Filter: &filters.MultiFilter{Limit: &limit, Offset: &offset},
		Status: &filters.Condition{Operator: filters.OpEqual, Value: constants.Active},
	}

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().GetShows(filter).Return(shows, nil).Times(1)

		result, err := service.GetShows(constants.Active, limit, offset)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, shows, result)
	})

	t.Run("error getting shows", func(t *testing.T) {
		repo.EXPECT().GetShows(filter).Return(nil, errors.New("error getting shows")).Times(1)

		result, err := service.GetShows(constants.Active, limit, offset)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.EqualError(t, err, "error getting shows")
	})
}

func TestShowService_CreateShow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	showRepo := mock_repositories.NewMockShowRepository(ctrl)
	movieRepo := mock_repositories.NewMockMovieRepository(ctrl)
	theaterRepo := mock_repositories.NewMockTheaterRepository(ctrl)
	service := NewShowService(nil, transaction, showRepo, movieRepo, theaterRepo, nil)

	show := utils.GenerateShow()
	req := payloads.CreateShowRequest{
		MovieId:   *show.MovieId,
		TheaterId: *show.TheaterId,
		StartTime: show.StartTime,
		EndTime:   show.EndTime,
		Status:    show.Status,
	}
	movieFilter := filters.MovieFilter{
		Filter: &filters.SingleFilter{},
		ID:     &filters.Condition{Operator: filters.OpEqual, Value: req.MovieId},
	}
	theaterFilter := filters.TheaterFilter{
		Filter: &filters.SingleFilter{},
		ID:     &filters.Condition{Operator: filters.OpEqual, Value: req.TheaterId},
	}

	t.Run("success", func(t *testing.T) {
		movieRepo.EXPECT().GetMovie(movieFilter, false).Return(&models.Movie{}, nil).Times(1)
		theaterRepo.EXPECT().GetTheater(theaterFilter, false).Return(&models.Theater{}, nil).Times(1)
		showRepo.EXPECT().IsShowInValidTimeRange(req.TheaterId, req.StartTime, req.EndTime).Return(true, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		showRepo.EXPECT().CreateShow(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		result, err := service.CreateShow(req)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, show.MovieId, result.MovieId)
		assert.Equal(t, show.TheaterId, result.TheaterId)
		assert.Equal(t, show.StartTime, result.StartTime)
		assert.Equal(t, show.EndTime, result.EndTime)
		assert.Equal(t, show.Status, result.Status)
	})

	t.Run("movie not found", func(t *testing.T) {
		movieRepo.EXPECT().GetMovie(movieFilter, false).Return(nil, nil).Times(1)

		result, err := service.CreateShow(req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.StatusCode)
		assert.EqualError(t, err, "movie not found")
	})

	t.Run("error getting movie", func(t *testing.T) {
		movieRepo.EXPECT().GetMovie(movieFilter, false).Return(nil, errors.New("error getting movie")).Times(1)

		result, err := service.CreateShow(req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.EqualError(t, err, "error getting movie")
	})

	t.Run("theater not found", func(t *testing.T) {
		movieRepo.EXPECT().GetMovie(movieFilter, false).Return(&models.Movie{}, nil).Times(1)
		theaterRepo.EXPECT().GetTheater(theaterFilter, false).Return(nil, nil).Times(1)

		result, err := service.CreateShow(req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.StatusCode)
		assert.EqualError(t, err, "theater not found")
	})

	t.Run("error getting theater", func(t *testing.T) {
		movieRepo.EXPECT().GetMovie(movieFilter, false).Return(&models.Movie{}, nil).Times(1)
		theaterRepo.EXPECT().GetTheater(theaterFilter, false).Return(nil, errors.New("error getting theater")).Times(1)

		result, err := service.CreateShow(req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.EqualError(t, err, "error getting theater")
	})

	t.Run("not valid time", func(t *testing.T) {
		movieRepo.EXPECT().GetMovie(movieFilter, false).Return(&models.Movie{}, nil).Times(1)
		theaterRepo.EXPECT().GetTheater(theaterFilter, false).Return(&models.Theater{}, nil).Times(1)
		showRepo.EXPECT().IsShowInValidTimeRange(req.TheaterId, req.StartTime, req.EndTime).Return(false, nil).Times(1)

		result, err := service.CreateShow(req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.StatusCode)
		assert.EqualError(t, err, "invalid time range for this show")
	})

	t.Run("error checking time range", func(t *testing.T) {
		movieRepo.EXPECT().GetMovie(movieFilter, false).Return(&models.Movie{}, nil).Times(1)
		theaterRepo.EXPECT().GetTheater(theaterFilter, false).Return(&models.Theater{}, nil).Times(1)
		showRepo.EXPECT().IsShowInValidTimeRange(req.TheaterId, req.StartTime, req.EndTime).Return(false, errors.New("error checking time range")).Times(1)

		result, err := service.CreateShow(req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.EqualError(t, err, "error checking time range")
	})

	t.Run("error creating show", func(t *testing.T) {
		movieRepo.EXPECT().GetMovie(movieFilter, false).Return(&models.Movie{}, nil).Times(1)
		theaterRepo.EXPECT().GetTheater(theaterFilter, false).Return(&models.Theater{}, nil).Times(1)
		showRepo.EXPECT().IsShowInValidTimeRange(req.TheaterId, req.StartTime, req.EndTime).Return(true, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		showRepo.EXPECT().CreateShow(gomock.Any(), gomock.Any()).Return(errors.New("error creating show")).Times(1)

		result, err := service.CreateShow(req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.EqualError(t, err, "error creating show")
	})
}

func TestShowService_ScheduleUpdateShowStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	repo := mock_repositories.NewMockShowRepository(ctrl)
	service := NewShowService(nil, transaction, repo, nil, nil, nil)

	t.Run("success", func(t *testing.T) {
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		repo.EXPECT().ScheduleActivateShows(gomock.Any(), time.Hour*72).Return(nil).Times(1)
		repo.EXPECT().ScheduleCompleteShows(gomock.Any()).Return(nil).Times(1)

		err := service.ScheduleUpdateShowStatus()

		assert.Nil(t, err)
	})

	t.Run("error activating shows", func(t *testing.T) {
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		repo.EXPECT().ScheduleActivateShows(gomock.Any(), time.Hour*72).Return(errors.New("error activating shows")).Times(1)

		err := service.ScheduleUpdateShowStatus()

		assert.NotNil(t, err)
		assert.EqualError(t, err, "error activating shows")
	})

	t.Run("error completing shows", func(t *testing.T) {
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		repo.EXPECT().ScheduleActivateShows(gomock.Any(), time.Hour*72).Return(nil).Times(1)
		repo.EXPECT().ScheduleCompleteShows(gomock.Any()).Return(errors.New("error completing shows")).Times(1)

		err := service.ScheduleUpdateShowStatus()

		assert.NotNil(t, err)
		assert.EqualError(t, err, "error completing shows")
	})
}
