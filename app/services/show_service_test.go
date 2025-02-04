package services

import (
	"errors"
	"github.com/stretchr/testify/assert"
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
)

func TestShowService_CreateShow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	showRepo := mock_repositories.NewMockShowRepository(ctrl)
	movieRepo := mock_repositories.NewMockMovieRepository(ctrl)
	theaterRepo := mock_repositories.NewMockTheaterRepository(ctrl)
	service := NewShowService(nil, transaction, showRepo, movieRepo, theaterRepo)

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
