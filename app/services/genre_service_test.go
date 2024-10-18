package services

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_repositories"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_transaction"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestGenreService_GetGenre(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repositories.NewMockGenreRepository(ctrl)
	service := NewGenreService(nil, nil, repo)

	genre := utils.GenerateGenre()

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().GetGenre(genre.ID).Return(genre, nil).Times(1)

		result, err := service.GetGenre(genre.ID)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, genre, result)
	})

	t.Run("genre not found", func(t *testing.T) {
		repo.EXPECT().GetGenre(genre.ID).Return(nil, nil).Times(1)

		result, err := service.GetGenre(genre.ID)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "genre not found", err.Error())
	})

	t.Run("error getting genre", func(t *testing.T) {
		repo.EXPECT().GetGenre(genre.ID).Return(nil, errors.New("error getting genre")).Times(1)

		result, err := service.GetGenre(genre.ID)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting genre", err.Error())
	})
}

func TestGenreService_GetGenres(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repositories.NewMockGenreRepository(ctrl)
	service := NewGenreService(nil, nil, repo)

	t.Run("success", func(t *testing.T) {
		genres := utils.GenerateGenres(3)

		repo.EXPECT().GetGenres().Return(genres, nil).Times(1)

		result, err := service.GetGenres()

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, genres, result)
	})

	t.Run("error getting genres", func(t *testing.T) {
		repo.EXPECT().GetGenres().Return(nil, errors.New("error getting genres")).Times(1)

		result, err := service.GetGenres()

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting genres", err.Error())
	})
}

func TestGenreService_CreateGenre(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	repo := mock_repositories.NewMockGenreRepository(ctrl)
	service := NewGenreService(nil, transaction, repo)

	genre := utils.GenerateGenre()

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().GetGenreByName(genre.Name).Return(nil, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
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
		assert.Equal(t, http.StatusBadRequest, err.StatusCode)
		assert.Equal(t, "duplicate genre name", err.Error())
	})

	t.Run("error getting genre", func(t *testing.T) {
		repo.EXPECT().GetGenreByName(genre.Name).Return(nil, errors.New("error getting genre")).Times(1)

		result, err := service.CreateGenre(genre.Name)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting genre", err.Error())
	})

	t.Run("error creating genre", func(t *testing.T) {
		repo.EXPECT().GetGenreByName(genre.Name).Return(nil, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		repo.EXPECT().CreateGenre(gomock.Any(), gomock.Any()).Return(errors.New("error creating genre")).Times(1)

		result, err := service.CreateGenre(genre.Name)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error creating genre", err.Error())
	})
}

func TestGenreService_UpdateGenre(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	repo := mock_repositories.NewMockGenreRepository(ctrl)
	service := NewGenreService(nil, transaction, repo)

	genre := utils.GenerateGenre()

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().GetGenre(genre.ID).Return(genre, nil).Times(1)
		repo.EXPECT().GetGenreByName(genre.Name).Return(nil, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		repo.EXPECT().UpdateGenre(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		result, err := service.UpdateGenre(genre.ID, genre.Name)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, genre.ID, result.ID)
		assert.Equal(t, genre.Name, result.Name)
	})

	t.Run("genre not found", func(t *testing.T) {
		repo.EXPECT().GetGenre(genre.ID).Return(nil, nil).Times(1)

		result, err := service.UpdateGenre(genre.ID, genre.Name)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "genre not found", err.Error())
	})

	t.Run("error getting genre", func(t *testing.T) {
		repo.EXPECT().GetGenre(genre.ID).Return(nil, errors.New("error getting genre")).Times(1)

		result, err := service.UpdateGenre(genre.ID, genre.Name)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting genre", err.Error())
	})

	t.Run("duplicate genre name", func(t *testing.T) {
		repo.EXPECT().GetGenre(genre.ID).Return(genre, nil).Times(1)
		repo.EXPECT().GetGenreByName(genre.Name).Return(genre, nil).Times(1)

		result, err := service.UpdateGenre(genre.ID, genre.Name)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.StatusCode)
		assert.Equal(t, "duplicate genre name", err.Error())
	})

	t.Run("error getting genre by name", func(t *testing.T) {
		repo.EXPECT().GetGenre(genre.ID).Return(genre, nil).Times(1)
		repo.EXPECT().GetGenreByName(genre.Name).Return(nil, errors.New("error getting genre")).Times(1)

		result, err := service.UpdateGenre(genre.ID, genre.Name)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting genre", err.Error())
	})

	t.Run("error updating genre", func(t *testing.T) {
		repo.EXPECT().GetGenre(genre.ID).Return(genre, nil).Times(1)
		repo.EXPECT().GetGenreByName(genre.Name).Return(nil, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		repo.EXPECT().UpdateGenre(gomock.Any(), gomock.Any()).Return(errors.New("error updating genre")).Times(1)

		result, err := service.UpdateGenre(genre.ID, genre.Name)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error updating genre", err.Error())
	})
}
