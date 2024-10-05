package services

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_repositories"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_transaction"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestMovieService_GetMovie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repositories.NewMockMovieRepository(ctrl)
	service := NewMovieService(nil, nil, repo)

	movie := utils.GenerateRandomMovie()

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().GetMovie(movie.ID).Return(movie, nil).Times(1)

		result, err := service.GetMovie(movie.ID)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, movie, result)
	})

	t.Run("movie not found", func(t *testing.T) {
		repo.EXPECT().GetMovie(movie.ID).Return(nil, gorm.ErrRecordNotFound).Times(1)

		result, err := service.GetMovie(movie.ID)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "Movie not found", err.Error())
	})

	t.Run("db error", func(t *testing.T) {
		repo.EXPECT().GetMovie(movie.ID).Return(nil, errors.New("db error")).Times(1)

		result, err := service.GetMovie(movie.ID)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Error())
	})
}

func TestMovieService_GetMovies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repositories.NewMockMovieRepository(ctrl)
	service := NewMovieService(nil, nil, repo)

	movies := make([]*models.Movie, 20)
	for i := 0; i < len(movies); i++ {
		movies[i] = utils.GenerateRandomMovie()
	}

	limit := 10
	offset := 0

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().GetMovies(limit, offset).Return(movies, nil).Times(1)
		repo.EXPECT().GetNumbersOfMovie().Return(len(movies), nil).Times(1)

		result, meta, err := service.GetMovies(limit, offset)

		assert.NotNil(t, result)
		assert.NotNil(t, meta)
		assert.Nil(t, err)

		for i, m := range result {
			assert.Equal(t, movies[i], m)
		}

		nextUrl := "/movies?limit=10&offset=10"
		expectedMeta := models.ResponseMeta{
			Limit:   limit,
			Offset:  offset,
			Total:   len(movies),
			NextUrl: &nextUrl,
		}
		assert.Equal(t, &expectedMeta, meta)
	})

	t.Run("error getting movies", func(t *testing.T) {
		repo.EXPECT().GetMovies(limit, offset).Return(nil, errors.New("db error")).Times(1)

		result, meta, err := service.GetMovies(limit, offset)

		assert.Nil(t, result)
		assert.Nil(t, meta)
		assert.NotNil(t, err)

		assert.Equal(t, "db error", err.Error())
	})

	t.Run("error counting movies", func(t *testing.T) {
		repo.EXPECT().GetMovies(limit, offset).Return(movies, nil).Times(1)
		repo.EXPECT().GetNumbersOfMovie().Return(0, errors.New("db error")).Times(1)

		result, meta, err := service.GetMovies(limit, offset)

		assert.Nil(t, result)
		assert.Nil(t, meta)
		assert.NotNil(t, err)

		assert.Equal(t, "db error", err.Error())
	})
}

func TestMovieService_CreateMovie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	repo := mock_repositories.NewMockMovieRepository(ctrl)
	service := NewMovieService(nil, transaction, repo)

	movie := utils.GenerateRandomMovie()

	t.Run("success", func(t *testing.T) {
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)

		repo.EXPECT().CreateMovie(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		result, err := service.CreateMovie(movie.Title, movie.Description, movie.ReleaseDate, movie.DurationMinutes, movie.Language, movie.Rating, movie.CreatedBy)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, movie.Title, result.Title)
		assert.Equal(t, movie.Description, result.Description)
		assert.Equal(t, movie.ReleaseDate, result.ReleaseDate)
		assert.Equal(t, movie.DurationMinutes, result.DurationMinutes)
		assert.Equal(t, movie.Language, result.Language)
		assert.Equal(t, movie.Rating, result.Rating)
		assert.Equal(t, movie.CreatedBy, result.CreatedBy)
	})

	t.Run("db error", func(t *testing.T) {
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)

		repo.EXPECT().CreateMovie(gomock.Any(), gomock.Any()).Return(errors.New("db error")).Times(1)

		result, err := service.CreateMovie(movie.Title, movie.Description, movie.ReleaseDate, movie.DurationMinutes, movie.Language, movie.Rating, movie.CreatedBy)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Message)
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

		result, err := service.UpdateMovie(currentMovie.ID, updatedMovie.LastUpdatedBy, updatedMovie.Title, updatedMovie.Description, updatedMovie.ReleaseDate, updatedMovie.DurationMinutes, updatedMovie.Language, updatedMovie.Rating)

		assert.NotNil(t, result)
		assert.Nil(t, err)
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

	t.Run("movie not found", func(t *testing.T) {
		repo.EXPECT().GetMovie(currentMovie.ID).Return(nil, gorm.ErrRecordNotFound).Times(1)

		result, err := service.UpdateMovie(updatedMovie.ID, updatedMovie.LastUpdatedBy, updatedMovie.Title, updatedMovie.Description, updatedMovie.ReleaseDate, updatedMovie.DurationMinutes, updatedMovie.Language, updatedMovie.Rating)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "Movie not found", err.Error())
	})

	t.Run("db error", func(t *testing.T) {
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			}).Times(1)

		repo.EXPECT().GetMovie(currentMovie.ID).Return(currentMovie, nil).Times(1)
		repo.EXPECT().UpdateMovie(gomock.Any(), gomock.Any()).Return(errors.New("db error")).Times(1)

		result, err := service.UpdateMovie(updatedMovie.ID, updatedMovie.LastUpdatedBy, updatedMovie.Title, updatedMovie.Description, updatedMovie.ReleaseDate, updatedMovie.DurationMinutes, updatedMovie.Language, updatedMovie.Rating)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Error())
	})
}
