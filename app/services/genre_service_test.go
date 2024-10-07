package services

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_repositories"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_transaction"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestGenreService_CreateGenre(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	repo := mock_repositories.NewMockGenreRepository(ctrl)
	service := NewGenreService(nil, transaction, repo)

	genre := utils.GenerateRandomGenre()

	t.Run("success", func(t *testing.T) {
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)

		repo.EXPECT().GetGenreByName(genre.Name).Return(nil, gorm.ErrRecordNotFound).Times(1)
		repo.EXPECT().CreateGenre(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		result, err := service.CreateGenre(genre.Name)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, genre.Name, result.Name)
	})

	t.Run("duplicate genre name", func(t *testing.T) {
		repo.EXPECT().GetGenreByName(genre.Name).Return(genre, nil).Times(1)

		result, err := service.CreateGenre(genre.Name)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "Duplicate genre name", err.Error())
	})

	t.Run("error getting genre", func(t *testing.T) {
		repo.EXPECT().GetGenreByName(genre.Name).Return(nil, errors.New("db error")).Times(1)

		result, err := service.CreateGenre(genre.Name)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Error())
	})

	t.Run("error creating genre", func(t *testing.T) {
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)

		repo.EXPECT().GetGenreByName(genre.Name).Return(nil, gorm.ErrRecordNotFound).Times(1)
		repo.EXPECT().CreateGenre(gomock.Any(), gomock.Any()).Return(errors.New("creating error")).Times(1)

		result, err := service.CreateGenre(genre.Name)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "creating error", err.Error())
	})
}
