package services

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/constants"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/filters"
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

	flagRepo := mock_repositories.NewMockFeatureFlagRepository(ctrl)
	movieRepo := mock_repositories.NewMockMovieRepository(ctrl)
	service := NewMovieService(nil, nil, movieRepo, nil, nil, flagRepo)

	email := "test@example.com"
	movie := utils.GenerateMovie()
	movie.IsActive = false
	filter := filters.MovieFilter{
		Filter:    &filters.SingleFilter{},
		ID:        &filters.Condition{Operator: filters.OpEqual, Value: movie.ID},
		IsDeleted: &filters.Condition{Operator: filters.OpEqual, Value: false},
	}

	t.Run("success", func(t *testing.T) {
		movieRepo.EXPECT().GetMovie(gomock.Eq(filter), false).Return(movie, nil).Times(1)
		flagRepo.EXPECT().HasFlagEnabled(email, constants.CanModifyMovies).Return(true).Times(1)

		result, err := service.GetMovie(movie.ID, &email, false)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, movie, result)
	})

	t.Run("movie not found", func(t *testing.T) {
		movieRepo.EXPECT().GetMovie(gomock.Eq(filter), false).Return(nil, nil).Times(1)

		result, err := service.GetMovie(movie.ID, &email, false)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "movie not found", err.Error())
	})

	t.Run("error getting movie", func(t *testing.T) {
		movieRepo.EXPECT().GetMovie(gomock.Eq(filter), false).Return(nil, errors.New("error getting movie")).Times(1)

		result, err := service.GetMovie(movie.ID, &email, false)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting movie", err.Error())
	})

	t.Run("unauthenticated user", func(t *testing.T) {
		movieRepo.EXPECT().GetMovie(gomock.Eq(filter), false).Return(movie, nil).Times(1)

		result, err := service.GetMovie(movie.ID, nil, false)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusForbidden, err.StatusCode)
		assert.Equal(t, "permission denied", err.Error())
	})

	t.Run("unauthorized user", func(t *testing.T) {
		movieRepo.EXPECT().GetMovie(gomock.Eq(filter), false).Return(movie, nil).Times(1)
		flagRepo.EXPECT().HasFlagEnabled(email, constants.CanModifyMovies).Return(false).Times(1)

		result, err := service.GetMovie(movie.ID, &email, false)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusForbidden, err.StatusCode)
		assert.Equal(t, "permission denied", err.Error())
	})
}

func TestMovieService_GetMovies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	movieRepo := mock_repositories.NewMockMovieRepository(ctrl)
	flagRepo := mock_repositories.NewMockFeatureFlagRepository(ctrl)
	service := NewMovieService(nil, nil, movieRepo, nil, nil, flagRepo)

	userEmail := "test@example.com"
	movies := utils.GenerateMovies(20)
	limit := 10
	offset := 0
	includeGenres := true
	getFilter := filters.MovieFilter{
		Filter:    &filters.MultiFilter{Limit: &limit, Offset: &offset},
		IsDeleted: &filters.Condition{Operator: filters.OpEqual, Value: false},
	}
	countFilter := filters.MovieFilter{
		Filter:    &filters.SingleFilter{},
		IsDeleted: &filters.Condition{Operator: filters.OpEqual, Value: false},
	}

	t.Run("success for normal user", func(t *testing.T) {
		flagRepo.EXPECT().HasFlagEnabled(userEmail, constants.CanModifyMovies).Return(false).Times(1)

		normalGetFilter := getFilter
		normalGetFilter.IsActive = &filters.Condition{Operator: filters.OpEqual, Value: true}
		movieRepo.EXPECT().GetMoviesWithGenres(normalGetFilter).Return(movies, nil).Times(1)

		normalCountFilter := countFilter
		normalCountFilter.IsActive = &filters.Condition{Operator: filters.OpEqual, Value: true}
		movieRepo.EXPECT().GetNumbersOfMovie(normalCountFilter).Return(len(movies), nil).Times(1)

		result, meta, err := service.GetMovies(limit, offset, &userEmail, includeGenres)

		assert.NotNil(t, result)
		assert.NotNil(t, meta)
		assert.Nil(t, err)
		assert.Equal(t, len(movies), len(result))

		for i, m := range result {
			assert.Equal(t, movies[i], m)
		}

		nextUrl := fmt.Sprintf("/movies?%s=10&%s=10&includeGenres=%v", constants.Limit, constants.Offset, includeGenres)
		expectedMeta := models.ResponseMeta{
			Limit:   limit,
			Offset:  offset,
			Total:   len(movies),
			NextUrl: &nextUrl,
		}
		assert.Equal(t, &expectedMeta, meta)
	})

	t.Run("success for admin user", func(t *testing.T) {
		flagRepo.EXPECT().HasFlagEnabled(userEmail, constants.CanModifyMovies).Return(true).Times(1)
		movieRepo.EXPECT().GetMoviesWithGenres(gomock.Eq(getFilter)).Return(movies, nil).Times(1)
		movieRepo.EXPECT().GetNumbersOfMovie(gomock.Eq(countFilter)).Return(len(movies), nil).Times(1)

		result, meta, err := service.GetMovies(limit, offset, &userEmail, includeGenres)

		assert.NotNil(t, result)
		assert.NotNil(t, meta)
		assert.Nil(t, err)
		assert.Equal(t, len(movies), len(result))

		for i, m := range result {
			assert.Equal(t, movies[i], m)
		}

		nextUrl := fmt.Sprintf("/movies?%s=10&%s=10&includeGenres=%v", constants.Limit, constants.Offset, includeGenres)
		expectedMeta := models.ResponseMeta{
			Limit:   limit,
			Offset:  offset,
			Total:   len(movies),
			NextUrl: &nextUrl,
		}
		assert.Equal(t, &expectedMeta, meta)
	})

	t.Run("error getting movies", func(t *testing.T) {
		flagRepo.EXPECT().HasFlagEnabled(userEmail, constants.CanModifyMovies).Return(true).Times(1)
		movieRepo.EXPECT().GetMoviesWithGenres(gomock.Eq(getFilter)).Return(nil, errors.New("error getting movies")).Times(1)

		result, meta, err := service.GetMovies(limit, offset, &userEmail, includeGenres)

		assert.Nil(t, result)
		assert.Nil(t, meta)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting movies", err.Error())
	})

	t.Run("error counting movies", func(t *testing.T) {
		flagRepo.EXPECT().HasFlagEnabled(userEmail, constants.CanModifyMovies).Return(true).Times(1)
		movieRepo.EXPECT().GetMoviesWithGenres(gomock.Eq(getFilter)).Return(movies, nil).Times(1)
		movieRepo.EXPECT().GetNumbersOfMovie(gomock.Eq(countFilter)).Return(0, errors.New("error counting movies")).Times(1)

		result, meta, err := service.GetMovies(limit, offset, &userEmail, includeGenres)

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
	service := NewMovieService(nil, transaction, repo, nil, nil, nil)

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
		assert.Equal(t, req.IsActive, &result.IsActive)
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
	service := NewMovieService(nil, transaction, repo, nil, nil, nil)

	movie := utils.GenerateMovie()
	req := utils.GenerateUpdateMovieRequest()
	filter := filters.MovieFilter{
		Filter:    &filters.SingleFilter{},
		ID:        &filters.Condition{Operator: filters.OpEqual, Value: movie.ID},
		IsDeleted: &filters.Condition{Operator: filters.OpEqual, Value: false},
	}

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().GetMovie(filter, false).Return(movie, nil).Times(1)
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
		assert.Equal(t, req.IsActive, &result.IsActive)
		assert.Equal(t, movie.CreatedBy, result.CreatedBy)
	})

	t.Run("movie not found", func(t *testing.T) {
		repo.EXPECT().GetMovie(filter, false).Return(nil, nil).Times(1)

		result, err := service.UpdateMovie(movie.ID, movie.LastUpdatedBy, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "movie not found", err.Error())
	})

	t.Run("error updating movie", func(t *testing.T) {
		repo.EXPECT().GetMovie(filter, false).Return(movie, nil).Times(1)
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
	service := NewMovieService(nil, transaction, movieRepo, genreRepo, movieGenreRepo, nil)

	movie := utils.GenerateMovie()
	allGenreIds := make([]uuid.UUID, 3)
	for i := 0; i < len(allGenreIds); i++ {
		allGenreIds[i] = utils.GenerateGenre().ID
	}
	updatedGenreIds := []uuid.UUID{allGenreIds[0], allGenreIds[1]}
	filter := filters.MovieFilter{
		Filter:    &filters.SingleFilter{},
		ID:        &filters.Condition{Operator: filters.OpEqual, Value: movie.ID},
		IsDeleted: &filters.Condition{Operator: filters.OpEqual, Value: false},
	}

	t.Run("success", func(t *testing.T) {
		movieRepo.EXPECT().GetMovie(filter, false).Return(movie, nil).Times(1)
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
		movieRepo.EXPECT().GetMovie(filter, false).Return(nil, nil).Times(1)

		err := service.AssignGenres(movie.ID, updatedGenreIds)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "movie not found", err.Error())
	})

	t.Run("error getting movie", func(t *testing.T) {
		movieRepo.EXPECT().GetMovie(filter, false).Return(nil, errors.New("error getting movie")).Times(1)

		err := service.AssignGenres(movie.ID, updatedGenreIds)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting movie", err.Error())
	})

	t.Run("error getting genres", func(t *testing.T) {
		movieRepo.EXPECT().GetMovie(filter, false).Return(movie, nil).Times(1)
		genreRepo.EXPECT().GetGenreIDs(gomock.Any()).Return(nil, errors.New("error getting genres")).Times(1)

		err := service.AssignGenres(movie.ID, updatedGenreIds)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting genres", err.Error())

	})

	t.Run("error updated genres not found", func(t *testing.T) {
		movieRepo.EXPECT().GetMovie(filter, false).Return(movie, nil).Times(1)
		genreRepo.EXPECT().GetGenreIDs(gomock.Any()).Return(allGenreIds, nil).Times(1)

		err := service.AssignGenres(movie.ID, []uuid.UUID{uuid.New(), uuid.New()})

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.StatusCode)
		assert.Equal(t, "invalid genre ids", err.Error())
	})

	t.Run("error updating movie", func(t *testing.T) {
		movieRepo.EXPECT().GetMovie(filter, false).Return(movie, nil).Times(1)
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
