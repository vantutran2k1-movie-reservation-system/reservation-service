package services

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/filters"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_repositories"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_transaction"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"net/http"
	"testing"
)

func TestTheaterService_GetTheater(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repositories.NewMockTheaterRepository(ctrl)
	service := NewTheaterService(nil, nil, repo, nil, nil)

	theater := utils.GenerateTheater()
	filter := filters.TheaterFilter{
		Filter: &filters.SingleFilter{},
		ID:     &filters.Condition{Operator: filters.OpEqual, Value: theater.ID},
	}

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().GetTheater(filter, false).Return(theater, nil).Times(1)

		result, err := service.GetTheater(theater.ID, false)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, theater, result)
	})

	t.Run("theater not found", func(t *testing.T) {
		repo.EXPECT().GetTheater(filter, false).Return(nil, nil).Times(1)

		result, err := service.GetTheater(theater.ID, false)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "theater not found", err.Error())
	})

	t.Run("error getting theater", func(t *testing.T) {
		repo.EXPECT().GetTheater(filter, false).Return(nil, errors.New("error getting theater")).Times(1)

		result, err := service.GetTheater(theater.ID, false)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting theater", err.Error())
	})
}

func TestTheaterService_CreateTheater(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	repo := mock_repositories.NewMockTheaterRepository(ctrl)
	service := NewTheaterService(nil, transaction, repo, nil, nil)

	theater := utils.GenerateTheater()
	req := utils.GenerateCreateTheaterRequest()
	filter := filters.TheaterFilter{
		Filter: &filters.SingleFilter{},
		Name:   &filters.Condition{Operator: filters.OpEqual, Value: &req.Name},
	}

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().GetTheater(filter, false).Return(nil, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		repo.EXPECT().CreateTheater(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		result, err := service.CreateTheater(req)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, req.Name, result.Name)
	})

	t.Run("duplicate theater name", func(t *testing.T) {
		repo.EXPECT().GetTheater(filter, false).Return(theater, nil).Times(1)

		result, err := service.CreateTheater(req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.StatusCode)
		assert.Equal(t, "duplicate theater name", err.Error())
	})

	t.Run("error getting theater", func(t *testing.T) {
		repo.EXPECT().GetTheater(filter, false).Return(nil, errors.New("error getting theater")).Times(1)

		result, err := service.CreateTheater(req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting theater", err.Error())
	})

	t.Run("error creating theater", func(t *testing.T) {
		repo.EXPECT().GetTheater(filter, false).Return(nil, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		repo.EXPECT().CreateTheater(gomock.Any(), gomock.Any()).Return(errors.New("error creating theater")).Times(1)

		result, err := service.CreateTheater(req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error creating theater", err.Error())
	})
}

func TestTheaterService_CreateTheaterLocation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	theaterRepo := mock_repositories.NewMockTheaterRepository(ctrl)
	theaterLocationRepo := mock_repositories.NewMockTheaterLocationRepository(ctrl)
	cityRepo := mock_repositories.NewMockCityRepository(ctrl)
	service := NewTheaterService(nil, transaction, theaterRepo, theaterLocationRepo, cityRepo)

	theater := utils.GenerateTheater()
	city := utils.GenerateCity()
	req := utils.GenerateCreateTheaterLocationRequest()
	cityFilter := filters.CityFilter{
		Filter: &filters.SingleFilter{Logic: filters.And},
		ID:     &filters.Condition{Operator: filters.OpEqual, Value: req.CityID},
	}
	theaterFilter := filters.TheaterFilter{
		Filter: &filters.SingleFilter{},
		ID:     &filters.Condition{Operator: filters.OpEqual, Value: theater.ID},
	}

	t.Run("success", func(t *testing.T) {
		theaterRepo.EXPECT().GetTheater(theaterFilter, true).Return(theater, nil).Times(1)
		cityRepo.EXPECT().GetCity(cityFilter).Return(city, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		theaterLocationRepo.EXPECT().CreateTheaterLocation(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		result, err := service.CreateTheaterLocation(theater.ID, req)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, req.CityID, result.CityID)
		assert.Equal(t, req.Address, result.Address)
		assert.Equal(t, req.PostalCode, result.PostalCode)
		assert.Equal(t, req.Latitude, result.Latitude)
		assert.Equal(t, req.Longitude, result.Longitude)
	})

	t.Run("theater not found", func(t *testing.T) {
		theaterRepo.EXPECT().GetTheater(theaterFilter, true).Return(nil, nil).Times(1)

		result, err := service.CreateTheaterLocation(theater.ID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "theater not found", err.Error())
	})

	t.Run("error getting theater", func(t *testing.T) {
		theaterRepo.EXPECT().GetTheater(theaterFilter, true).Return(nil, errors.New("error getting theater")).Times(1)

		result, err := service.CreateTheaterLocation(theater.ID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting theater", err.Error())
	})

	t.Run("duplicate theater location", func(t *testing.T) {
		th := utils.GenerateTheater()
		th.Location = utils.GenerateTheaterLocation()
		theaterRepo.EXPECT().GetTheater(theaterFilter, true).Return(th, nil).Times(1)

		result, err := service.CreateTheaterLocation(theater.ID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.StatusCode)
		assert.Equal(t, "duplicate location for this theater", err.Error())
	})

	t.Run("city not found", func(t *testing.T) {
		theaterRepo.EXPECT().GetTheater(theaterFilter, true).Return(theater, nil).Times(1)
		cityRepo.EXPECT().GetCity(cityFilter).Return(nil, nil).Times(1)

		result, err := service.CreateTheaterLocation(theater.ID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.StatusCode)
		assert.Equal(t, "invalid city id", err.Error())
	})

	t.Run("error getting city", func(t *testing.T) {
		theaterRepo.EXPECT().GetTheater(theaterFilter, true).Return(theater, nil).Times(1)
		cityRepo.EXPECT().GetCity(cityFilter).Return(nil, errors.New("error getting city")).Times(1)

		result, err := service.CreateTheaterLocation(theater.ID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting city", err.Error())
	})

	t.Run("error creating location", func(t *testing.T) {
		theaterRepo.EXPECT().GetTheater(theaterFilter, true).Return(theater, nil).Times(1)
		cityRepo.EXPECT().GetCity(cityFilter).Return(city, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		theaterLocationRepo.EXPECT().CreateTheaterLocation(gomock.Any(), gomock.Any()).Return(errors.New("error creating location")).Times(1)

		result, err := service.CreateTheaterLocation(theater.ID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error creating location", err.Error())
	})
}
