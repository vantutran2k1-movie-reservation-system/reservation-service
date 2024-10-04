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

func TestMovieService_CreateMovie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	movieRepo := mock_repositories.NewMockMovieRepository(ctrl)
	service := NewMovieService(nil, transaction, movieRepo)

	movie := utils.GenerateRandomMovie()

	t.Run("success", func(t *testing.T) {
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		movieRepo.EXPECT().CreateMovie(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		result, err := service.CreateMovie(movie.CreatedBy, movie.Title, movie.Description, movie.ReleaseDate, movie.DurationMinutes, movie.Language, movie.Rating)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, movie.Title, result.Title)
	})

	t.Run("error creating movie", func(t *testing.T) {
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		movieRepo.EXPECT().CreateMovie(gomock.Any(), gomock.Any()).Return(errors.New("create error")).Times(1)

		result, err := service.CreateMovie(movie.CreatedBy, movie.Title, movie.Description, movie.ReleaseDate, movie.DurationMinutes, movie.Language, movie.Rating)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "create error", err.Message)
	})
}

func TestMovieService_UpdateMovie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	repo := mock_repositories.NewMockMovieRepository(ctrl)
	service := NewMovieService(nil, transaction, repo)

	currentMovie := utils.GenerateRandomMovie()
	updatedMovie := utils.GenerateRandomMovie()
	updatedMovie.ID = currentMovie.ID
	updatedMovie.CreatedBy = currentMovie.CreatedBy

	t.Run("success", func(t *testing.T) {
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			}).Times(1)

		repo.EXPECT().GetMovie(currentMovie.ID).Return(currentMovie, nil).Times(1)
		repo.EXPECT().UpdateMovie(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		result, err := service.UpdateMovie(updatedMovie.ID, updatedMovie.LastUpdatedBy, updatedMovie.Title, updatedMovie.Description, updatedMovie.ReleaseDate, updatedMovie.DurationMinutes, updatedMovie.Language, updatedMovie.Rating)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, updatedMovie.ID, result.ID)
		assert.Equal(t, updatedMovie.Title, result.Title)
		assert.Equal(t, updatedMovie.Description, result.Description)
		assert.Equal(t, updatedMovie.ReleaseDate, result.ReleaseDate)
		assert.Equal(t, updatedMovie.DurationMinutes, result.DurationMinutes)
		assert.Equal(t, updatedMovie.Language, result.Language)
		assert.Equal(t, updatedMovie.Rating, result.Rating)
		assert.Equal(t, currentMovie.CreatedAt, result.CreatedAt)
		assert.Equal(t, currentMovie.CreatedBy, result.CreatedBy)
		assert.Equal(t, updatedMovie.LastUpdatedBy, result.LastUpdatedBy)
	})
}
