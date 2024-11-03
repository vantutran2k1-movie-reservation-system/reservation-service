package services

import (
	"errors"
	"github.com/google/uuid"
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
	service := NewMovieService(nil, nil, repo, nil, nil)

	movie := utils.GenerateMovie()

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().GetMovie(movie.ID, false).Return(movie, nil).Times(1)

		result, err := service.GetMovie(movie.ID, false)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, movie, result)
	})

	t.Run("movie not found", func(t *testing.T) {
		repo.EXPECT().GetMovie(movie.ID, false).Return(nil, nil).Times(1)

		result, err := service.GetMovie(movie.ID, false)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "movie not found", err.Error())
	})

	t.Run("error getting movie", func(t *testing.T) {
		repo.EXPECT().GetMovie(movie.ID, false).Return(nil, errors.New("error getting movie")).Times(1)

		result, err := service.GetMovie(movie.ID, false)

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
	service := NewMovieService(nil, nil, repo, nil, nil)

	movies := utils.GenerateMovies(20)
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
	service := NewMovieService(nil, transaction, repo, nil, nil)

	movie := utils.GenerateMovie()
	req := utils.GenerateCreateMovieRequest()

	t.Run("success", func(t *testing.T) {
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		repo.EXPECT().CreateMovie(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		result, err := service.CreateMovie(req, movie.CreatedBy)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, req.Title, result.Title)
		assert.Equal(t, req.Description, result.Description)
		assert.Equal(t, req.ReleaseDate, result.ReleaseDate)
		assert.Equal(t, req.DurationMinutes, result.DurationMinutes)
		assert.Equal(t, req.Language, result.Language)
		assert.Equal(t, req.Rating, result.Rating)
		assert.Equal(t, movie.CreatedBy, result.CreatedBy)
	})

	t.Run("error creating movie", func(t *testing.T) {
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		repo.EXPECT().CreateMovie(gomock.Any(), gomock.Any()).Return(errors.New("error creating movie")).Times(1)

		result, err := service.CreateMovie(req, movie.CreatedBy)

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
	service := NewMovieService(nil, transaction, repo, nil, nil)

	movie := utils.GenerateMovie()
	req := utils.GenerateUpdateMovieRequest()

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().GetMovie(movie.ID, false).Return(movie, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			}).Times(1)
		repo.EXPECT().UpdateMovie(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		result, err := service.UpdateMovie(movie.ID, movie.LastUpdatedBy, req)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, movie.ID, result.ID)
		assert.Equal(t, req.Title, result.Title)
		assert.Equal(t, req.Description, result.Description)
		assert.Equal(t, req.ReleaseDate, result.ReleaseDate)
		assert.Equal(t, req.DurationMinutes, result.DurationMinutes)
		assert.Equal(t, req.Language, result.Language)
		assert.Equal(t, req.Rating, result.Rating)
		assert.Equal(t, movie.CreatedBy, result.CreatedBy)
	})

	t.Run("movie not found", func(t *testing.T) {
		repo.EXPECT().GetMovie(movie.ID, false).Return(nil, nil).Times(1)

		result, err := service.UpdateMovie(movie.ID, movie.LastUpdatedBy, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "movie not found", err.Error())
	})

	t.Run("error updating movie", func(t *testing.T) {
		repo.EXPECT().GetMovie(movie.ID, false).Return(movie, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			}).Times(1)
		repo.EXPECT().UpdateMovie(gomock.Any(), gomock.Any()).Return(errors.New("error updating movie")).Times(1)

		result, err := service.UpdateMovie(movie.ID, movie.LastUpdatedBy, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "error updating movie", err.Error())
	})
}

func TestMovieService_AssignGenres(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	movieRepo := mock_repositories.NewMockMovieRepository(ctrl)
	genreRepo := mock_repositories.NewMockGenreRepository(ctrl)
	movieGenreRepo := mock_repositories.NewMockMovieGenreRepository(ctrl)
	service := NewMovieService(nil, transaction, movieRepo, genreRepo, movieGenreRepo)

	movie := utils.GenerateMovie()
	allGenreIds := make([]uuid.UUID, 3)
	for i := 0; i < len(allGenreIds); i++ {
		allGenreIds[i] = utils.GenerateGenre().ID
	}
	updatedGenreIds := []uuid.UUID{allGenreIds[0], allGenreIds[1]}

	t.Run("success", func(t *testing.T) {
		movieRepo.EXPECT().GetMovie(movie.ID, false).Return(movie, nil).Times(1)
		genreRepo.EXPECT().GetGenreIDs(gomock.Any()).Return(allGenreIds, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			}).Times(1)
		movieGenreRepo.EXPECT().UpdateGenresOfMovie(gomock.Any(), movie.ID, updatedGenreIds).Return(nil).Times(1)

		err := service.AssignGenres(movie.ID, updatedGenreIds)

		assert.Nil(t, err)
	})

	t.Run("movie not found", func(t *testing.T) {
		movieRepo.EXPECT().GetMovie(movie.ID, false).Return(nil, nil).Times(1)

		err := service.AssignGenres(movie.ID, updatedGenreIds)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "movie not found", err.Error())
	})

	t.Run("error getting movie", func(t *testing.T) {
		movieRepo.EXPECT().GetMovie(movie.ID, false).Return(nil, errors.New("error getting movie")).Times(1)

		err := service.AssignGenres(movie.ID, updatedGenreIds)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting movie", err.Error())
	})

	t.Run("error getting genres", func(t *testing.T) {
		movieRepo.EXPECT().GetMovie(movie.ID, false).Return(movie, nil).Times(1)
		genreRepo.EXPECT().GetGenreIDs(gomock.Any()).Return(nil, errors.New("error getting genres")).Times(1)

		err := service.AssignGenres(movie.ID, updatedGenreIds)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting genres", err.Error())

	})

	t.Run("error updated genres not found", func(t *testing.T) {
		movieRepo.EXPECT().GetMovie(movie.ID, false).Return(movie, nil).Times(1)
		genreRepo.EXPECT().GetGenreIDs(gomock.Any()).Return(allGenreIds, nil).Times(1)

		err := service.AssignGenres(movie.ID, []uuid.UUID{uuid.New(), uuid.New()})

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.StatusCode)
		assert.Equal(t, "invalid genre ids", err.Error())
	})

	t.Run("error updating movie", func(t *testing.T) {
		movieRepo.EXPECT().GetMovie(movie.ID, false).Return(movie, nil).Times(1)
		genreRepo.EXPECT().GetGenreIDs(gomock.Any()).Return(allGenreIds, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			}).Times(1)
		movieGenreRepo.EXPECT().UpdateGenresOfMovie(gomock.Any(), movie.ID, updatedGenreIds).Return(errors.New("error updating movie")).Times(1)

		err := service.AssignGenres(movie.ID, updatedGenreIds)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error updating movie", err.Error())
	})
}
