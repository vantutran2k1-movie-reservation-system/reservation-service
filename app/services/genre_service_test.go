package services

import (
	"errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/filters"
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
	service := NewGenreService(nil, nil, repo, nil)

	genre := utils.GenerateGenre()
	filter := filters.GenreFilter{
		Filter: &filters.SingleFilter{Logic: filters.And},
		ID:     &filters.Condition{Operator: filters.OpEqual, Value: genre.ID},
	}

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().GetGenre(filter).Return(genre, nil).Times(1)

		result, err := service.GetGenre(genre.ID)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, genre, result)
	})

	t.Run("genre not found", func(t *testing.T) {
		repo.EXPECT().GetGenre(filter).Return(nil, nil).Times(1)

		result, err := service.GetGenre(genre.ID)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "genre not found", err.Error())
	})

	t.Run("error getting genre", func(t *testing.T) {
		repo.EXPECT().GetGenre(filter).Return(nil, errors.New("error getting genre")).Times(1)

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
	service := NewGenreService(nil, nil, repo, nil)

	filter := filters.GenreFilter{
		Filter: &filters.MultiFilter{Logic: filters.And},
	}

	t.Run("success", func(t *testing.T) {
		genres := utils.GenerateGenres(3)

		repo.EXPECT().GetGenres(filter).Return(genres, nil).Times(1)

		result, err := service.GetGenres()

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, genres, result)
	})

	t.Run("error getting genres", func(t *testing.T) {
		repo.EXPECT().GetGenres(filter).Return(nil, errors.New("error getting genres")).Times(1)

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
	service := NewGenreService(nil, transaction, repo, nil)

	genre := utils.GenerateGenre()
	req := utils.GenerateCreateGenreRequest()

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().GetGenre(gomock.Any()).Return(nil, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		repo.EXPECT().CreateGenre(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		result, err := service.CreateGenre(req)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, req.Name, result.Name)
	})

	t.Run("duplicate genre name", func(t *testing.T) {
		repo.EXPECT().GetGenre(gomock.Any()).Return(genre, nil).Times(1)

		result, err := service.CreateGenre(req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.StatusCode)
		assert.Equal(t, "duplicate genre name", err.Error())
	})

	t.Run("error getting genre", func(t *testing.T) {
		repo.EXPECT().GetGenre(gomock.Any()).Return(nil, errors.New("error getting genre")).Times(1)

		result, err := service.CreateGenre(req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting genre", err.Error())
	})

	t.Run("error creating genre", func(t *testing.T) {
		repo.EXPECT().GetGenre(gomock.Any()).Return(nil, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		repo.EXPECT().CreateGenre(gomock.Any(), gomock.Any()).Return(errors.New("error creating genre")).Times(1)

		result, err := service.CreateGenre(req)

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
	service := NewGenreService(nil, transaction, repo, nil)

	genre := utils.GenerateGenre()
	req := utils.GenerateUpdateGenreRequest()
	idFilter := filters.GenreFilter{
		Filter: &filters.SingleFilter{Logic: filters.And},
		ID:     &filters.Condition{Operator: filters.OpEqual, Value: genre.ID},
	}
	nameFilter := filters.GenreFilter{
		Filter: &filters.SingleFilter{Logic: filters.And},
		Name:   &filters.Condition{Operator: filters.OpEqual, Value: req.Name},
	}

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().GetGenre(gomock.Eq(idFilter)).Return(genre, nil).Times(1)
		repo.EXPECT().GetGenre(gomock.Eq(nameFilter)).Return(nil, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		repo.EXPECT().UpdateGenre(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		result, err := service.UpdateGenre(genre.ID, req)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, genre.ID, result.ID)
		assert.Equal(t, req.Name, result.Name)
	})

	t.Run("genre not found", func(t *testing.T) {
		repo.EXPECT().GetGenre(gomock.Eq(idFilter)).Return(nil, nil).Times(1)

		result, err := service.UpdateGenre(genre.ID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "genre not found", err.Error())
	})

	t.Run("error getting genre", func(t *testing.T) {
		repo.EXPECT().GetGenre(gomock.Eq(idFilter)).Return(nil, errors.New("error getting genre")).Times(1)

		result, err := service.UpdateGenre(genre.ID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting genre", err.Error())
	})

	t.Run("duplicate genre name", func(t *testing.T) {
		repo.EXPECT().GetGenre(gomock.Eq(idFilter)).Return(genre, nil).Times(1)
		repo.EXPECT().GetGenre(gomock.Eq(nameFilter)).Return(genre, nil).Times(1)

		result, err := service.UpdateGenre(genre.ID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.StatusCode)
		assert.Equal(t, "duplicate genre name", err.Error())
	})

	t.Run("error getting genre by name", func(t *testing.T) {
		repo.EXPECT().GetGenre(gomock.Eq(idFilter)).Return(genre, nil).Times(1)
		repo.EXPECT().GetGenre(gomock.Eq(nameFilter)).Return(nil, errors.New("error getting genre")).Times(1)

		result, err := service.UpdateGenre(genre.ID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting genre", err.Error())
	})

	t.Run("error updating genre", func(t *testing.T) {
		repo.EXPECT().GetGenre(gomock.Eq(idFilter)).Return(genre, nil).Times(1)
		repo.EXPECT().GetGenre(gomock.Eq(nameFilter)).Return(nil, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		repo.EXPECT().UpdateGenre(gomock.Any(), gomock.Any()).Return(errors.New("error updating genre")).Times(1)

		result, err := service.UpdateGenre(genre.ID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error updating genre", err.Error())
	})
}

func TestGenreService_DeleteGenre(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	genreRepo := mock_repositories.NewMockGenreRepository(ctrl)
	movieGenreRepo := mock_repositories.NewMockMovieGenreRepository(ctrl)
	service := NewGenreService(nil, transaction, genreRepo, movieGenreRepo)

	genre := utils.GenerateGenre()
	filter := filters.GenreFilter{
		Filter: &filters.SingleFilter{Logic: filters.And},
		ID:     &filters.Condition{Operator: filters.OpEqual, Value: genre.ID},
	}

	t.Run("success", func(t *testing.T) {
		genreRepo.EXPECT().GetGenre(filter).Return(genre, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		movieGenreRepo.EXPECT().DeleteByGenreId(gomock.Any(), genre.ID).Return(nil).Times(1)
		genreRepo.EXPECT().DeleteGenre(gomock.Any(), genre).Return(nil).Times(1)

		err := service.DeleteGenre(genre.ID)

		assert.Nil(t, err)
	})

	t.Run("genre not found", func(t *testing.T) {
		genreRepo.EXPECT().GetGenre(filter).Return(nil, nil).Times(1)

		err := service.DeleteGenre(genre.ID)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "genre not found", err.Error())
	})

	t.Run("error getting genre", func(t *testing.T) {
		genreRepo.EXPECT().GetGenre(filter).Return(nil, errors.New("error getting genre")).Times(1)

		err := service.DeleteGenre(genre.ID)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting genre", err.Error())
	})

	t.Run("error deleting movie genres", func(t *testing.T) {
		genreRepo.EXPECT().GetGenre(filter).Return(genre, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		movieGenreRepo.EXPECT().DeleteByGenreId(gomock.Any(), genre.ID).Return(errors.New("error deleting movie genres")).Times(1)

		err := service.DeleteGenre(genre.ID)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error deleting movie genres", err.Error())
	})

	t.Run("error deleting genre", func(t *testing.T) {
		genreRepo.EXPECT().GetGenre(filter).Return(genre, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		movieGenreRepo.EXPECT().DeleteByGenreId(gomock.Any(), genre.ID).Return(nil).Times(1)
		genreRepo.EXPECT().DeleteGenre(gomock.Any(), genre).Return(errors.New("error deleting genre")).Times(1)

		err := service.DeleteGenre(genre.ID)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error deleting genre", err.Error())
	})
}
