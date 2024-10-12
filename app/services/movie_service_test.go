package services

import (
	"errors"
	"net/http"
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
		repo.EXPECT().GetMovie(movie.ID).Return(nil, nil).Times(1)

		result, err := service.GetMovie(movie.ID)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "movie not found", err.Error())
	})

	t.Run("error getting movie", func(t *testing.T) {
		repo.EXPECT().GetMovie(movie.ID).Return(nil, errors.New("error getting movie")).Times(1)

		result, err := service.GetMovie(movie.ID)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting movie", err.Error())
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
		repo.EXPECT().GetMovies(limit, offset).Return(nil, errors.New("error getting movies")).Times(1)

		result, meta, err := service.GetMovies(limit, offset)

		assert.Nil(t, result)
		assert.Nil(t, meta)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting movies", err.Error())
	})

	t.Run("error counting movies", func(t *testing.T) {
		repo.EXPECT().GetMovies(limit, offset).Return(movies, nil).Times(1)
		repo.EXPECT().GetNumbersOfMovie().Return(0, errors.New("error counting movies")).Times(1)

		result, meta, err := service.GetMovies(limit, offset)

		assert.Nil(t, result)
		assert.Nil(t, meta)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error counting movies", err.Error())
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

	t.Run("error creating movie", func(t *testing.T) {
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		repo.EXPECT().CreateMovie(gomock.Any(), gomock.Any()).Return(errors.New("error creating movie")).Times(1)

		result, err := service.CreateMovie(movie.Title, movie.Description, movie.ReleaseDate, movie.DurationMinutes, movie.Language, movie.Rating, movie.CreatedBy)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error creating movie", err.Error())
	})
}

func TestMovieService_UpdateMovie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	repo := mock_repositories.NewMockMovieRepository(ctrl)
	service := NewMovieService(nil, transaction, repo)

	movie := utils.GenerateRandomMovie()

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().GetMovie(movie.ID).Return(movie, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			}).Times(1)
		repo.EXPECT().UpdateMovie(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		result, err := service.UpdateMovie(movie.ID, movie.LastUpdatedBy, movie.Title, movie.Description, movie.ReleaseDate, movie.DurationMinutes, movie.Language, movie.Rating)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, movie.ID, result.ID)
		assert.Equal(t, movie.Title, result.Title)
		assert.Equal(t, movie.Description, result.Description)
		assert.Equal(t, movie.ReleaseDate, result.ReleaseDate)
		assert.Equal(t, movie.DurationMinutes, result.DurationMinutes)
		assert.Equal(t, movie.Language, result.Language)
		assert.Equal(t, movie.Rating, result.Rating)
		assert.Equal(t, movie.CreatedBy, result.CreatedBy)
	})

	t.Run("movie not found", func(t *testing.T) {
		repo.EXPECT().GetMovie(movie.ID).Return(nil, nil).Times(1)

		result, err := service.UpdateMovie(movie.ID, movie.LastUpdatedBy, movie.Title, movie.Description, movie.ReleaseDate, movie.DurationMinutes, movie.Language, movie.Rating)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "movie not found", err.Error())
	})

	t.Run("error updating movie", func(t *testing.T) {
		repo.EXPECT().GetMovie(movie.ID).Return(movie, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			}).Times(1)
		repo.EXPECT().UpdateMovie(gomock.Any(), gomock.Any()).Return(errors.New("error updating movie")).Times(1)

		result, err := service.UpdateMovie(movie.ID, movie.LastUpdatedBy, movie.Title, movie.Description, movie.ReleaseDate, movie.DurationMinutes, movie.Language, movie.Rating)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "error updating movie", err.Error())
	})
}
