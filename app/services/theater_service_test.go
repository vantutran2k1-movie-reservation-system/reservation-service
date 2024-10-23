package services

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_repositories"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_transaction"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"net/http"
	"testing"
)

func TestTheaterService_CreateTheater(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	repo := mock_repositories.NewMockTheaterRepository(ctrl)
	service := NewTheaterService(nil, transaction, repo)

	theater := utils.GenerateTheater()
	req := utils.GenerateCreateTheaterRequest()

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().GetTheaterByName(req.Name).Return(nil, nil).Times(1)
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
		repo.EXPECT().GetTheaterByName(req.Name).Return(theater, nil).Times(1)

		result, err := service.CreateTheater(req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.StatusCode)
		assert.Equal(t, "duplicate theater name", err.Error())
	})

	t.Run("error getting theater", func(t *testing.T) {
		repo.EXPECT().GetTheaterByName(req.Name).Return(nil, errors.New("error getting theater")).Times(1)

		result, err := service.CreateTheater(req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting theater", err.Error())
	})

	t.Run("error creating theater", func(t *testing.T) {
		repo.EXPECT().GetTheaterByName(req.Name).Return(nil, nil).Times(1)
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
