package controllers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/constants"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_services"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"go.uber.org/mock/gomock"
)

func TestMovieController_GetMovie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockMovieService(ctrl)
	controller := MovieController{
		MovieService: service,
	}

	gin.SetMode(gin.TestMode)

	movie := utils.GenerateRandomMovie()

	t.Run("success", func(t *testing.T) {
		router := gin.Default()
		router.GET("/movies/:id", controller.GetMovie)

		service.EXPECT().GetMovie(movie.ID).Return(movie, nil).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/movies/%s", movie.ID), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), movie.Title)
		assert.Contains(t, w.Body.String(), *movie.Description)
		assert.Contains(t, w.Body.String(), movie.ReleaseDate)
		assert.Contains(t, w.Body.String(), fmt.Sprint(movie.DurationMinutes))
		assert.Contains(t, w.Body.String(), *movie.Language)
		assert.Contains(t, w.Body.String(), fmt.Sprint(*movie.Rating))
		assert.Contains(t, w.Body.String(), movie.CreatedBy.String())
	})

	t.Run("service error", func(t *testing.T) {
		router := gin.Default()
		router.GET("/movies/:id", controller.GetMovie)

		service.EXPECT().GetMovie(movie.ID).Return(nil, errors.InternalServerError("service error")).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/movies/%s", movie.ID), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "error")
		assert.Contains(t, w.Body.String(), "service error")
	})
}

func TestMovieController_GetMovies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockMovieService(ctrl)
	controller := MovieController{
		MovieService: service,
	}

	gin.SetMode(gin.TestMode)

	movies := make([]*models.Movie, 20)
	for i := 0; i < len(movies); i++ {
		movies[i] = utils.GenerateRandomMovie()
	}
	meta := utils.GenerateRandomResponseMeta()

	t.Run("success", func(t *testing.T) {
		router := gin.Default()
		router.GET("/movies", controller.GetMovies)

		service.EXPECT().GetMovies(10, 0).Return(movies, meta, nil).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/movies?limit=10&offet=0", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), fmt.Sprint(meta.Limit))
		assert.Contains(t, w.Body.String(), fmt.Sprint(meta.Offset))
		assert.Contains(t, w.Body.String(), fmt.Sprint(meta.Total))
		assert.Contains(t, w.Body.String(), *meta.NextUrl)
		assert.Contains(t, w.Body.String(), *meta.PrevUrl)

		for _, m := range movies {
			assert.Contains(t, w.Body.String(), m.ID.String())
		}
	})

	t.Run("default limit and offset when receiving invalid values", func(t *testing.T) {
		router := gin.Default()
		router.GET("/movies", controller.GetMovies)

		service.EXPECT().GetMovies(10, 0).Return(movies, meta, nil).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/movies?limit=a&offet=b", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), fmt.Sprint(meta.Limit))
		assert.Contains(t, w.Body.String(), fmt.Sprint(meta.Offset))
		assert.Contains(t, w.Body.String(), fmt.Sprint(meta.Total))
		assert.Contains(t, w.Body.String(), *meta.NextUrl)
		assert.Contains(t, w.Body.String(), *meta.PrevUrl)

		for _, m := range movies {
			assert.Contains(t, w.Body.String(), m.ID.String())
		}
	})

	t.Run("service error", func(t *testing.T) {
		router := gin.Default()
		router.GET("/movies", controller.GetMovies)

		service.EXPECT().GetMovies(10, 0).Return(nil, nil, errors.InternalServerError("service error")).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/movies?limit=10&offet=0", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "error")
		assert.Contains(t, w.Body.String(), "service error")
	})
}

func TestMovieController_CreateMovie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockMovieService(ctrl)
	controller := MovieController{
		MovieService: service,
	}

	gin.SetMode(gin.TestMode)

	session := utils.GenerateRandomUserSession()
	movie := utils.GenerateRandomMovie()
	payload := utils.GenerateRandomCreateMovieRequest()

	errors.RegisterCustomValidators()

	t.Run("success", func(t *testing.T) {
		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set(constants.USER_SESSION, session)
			c.Next()
		})
		router.POST("/movies", controller.CreateMovie)

		service.EXPECT().CreateMovie(payload.Title, payload.Description, payload.ReleaseDate, payload.DurationMinutes, payload.Language, payload.Rating, session.UserID).
			Return(movie, nil)

		reqBody := fmt.Sprintf(`{"title": "%s", "description": "%s", "release_date": "%s", "duration_minutes": %d, "language": "%s", "rating": %g}`,
			payload.Title, *payload.Description, payload.ReleaseDate, payload.DurationMinutes, *payload.Language, *payload.Rating)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/movies", bytes.NewBufferString(reqBody))
		req.Header.Set(constants.CONTENT_TYPE, constants.APPLICATION_JSON)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), movie.Title)
		assert.Contains(t, w.Body.String(), *movie.Description)
		assert.Contains(t, w.Body.String(), movie.ReleaseDate)
		assert.Contains(t, w.Body.String(), fmt.Sprint(movie.DurationMinutes))
		assert.Contains(t, w.Body.String(), *movie.Language)
		assert.Contains(t, w.Body.String(), fmt.Sprint(*movie.Rating))
		assert.Contains(t, w.Body.String(), movie.CreatedBy.String())
	})

	t.Run("validation error", func(t *testing.T) {
		router := gin.Default()
		router.POST("/movies", controller.CreateMovie)

		reqBody := fmt.Sprintf(`{"title": "%s", "description": "%s", "release_date": "%s", "duration_minutes": %d, "language": "%s", "rating": %g}`,
			payload.Title, *payload.Description, payload.ReleaseDate, payload.DurationMinutes, *payload.Language, 6.0)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/movies", bytes.NewBufferString(reqBody))
		req.Header.Set(constants.CONTENT_TYPE, constants.APPLICATION_JSON)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "errors")
		assert.Contains(t, w.Body.String(), "Should be less than or equal to 5")
	})

	t.Run("session error", func(t *testing.T) {
		router := gin.Default()
		router.POST("/movies", controller.CreateMovie)

		reqBody := fmt.Sprintf(`{"title": "%s", "description": "%s", "release_date": "%s", "duration_minutes": %d, "language": "%s", "rating": %g}`,
			payload.Title, *payload.Description, payload.ReleaseDate, payload.DurationMinutes, *payload.Language, *movie.Rating)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/movies", bytes.NewBufferString(reqBody))
		req.Header.Set(constants.CONTENT_TYPE, constants.APPLICATION_JSON)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "error")
		assert.Contains(t, w.Body.String(), "Can not get user id from request")
	})

	t.Run("service error", func(t *testing.T) {
		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set(constants.USER_SESSION, session)
			c.Next()
		})
		router.POST("/movies", controller.CreateMovie)

		service.EXPECT().CreateMovie(payload.Title, payload.Description, payload.ReleaseDate, payload.DurationMinutes, payload.Language, payload.Rating, session.UserID).
			Return(nil, errors.InternalServerError("Service error"))

		reqBody := fmt.Sprintf(`{"title": "%s", "description": "%s", "release_date": "%s", "duration_minutes": %d, "language": "%s", "rating": %g}`,
			payload.Title, *payload.Description, payload.ReleaseDate, payload.DurationMinutes, *payload.Language, *payload.Rating)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/movies", bytes.NewBufferString(reqBody))
		req.Header.Set(constants.CONTENT_TYPE, constants.APPLICATION_JSON)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Service error")
	})
}

func TestMovieController_UpdateMovie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockMovieService(ctrl)
	controller := MovieController{
		MovieService: service,
	}

	gin.SetMode(gin.TestMode)

	session := utils.GenerateRandomUserSession()
	movie := utils.GenerateRandomMovie()
	payload := utils.GenerateRandomCreateMovieRequest()

	errors.RegisterCustomValidators()

	t.Run("success", func(t *testing.T) {
		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set(constants.USER_SESSION, session)
			c.Next()
		})
		router.PUT("/movies/:id", controller.UpdateMovie)

		service.EXPECT().UpdateMovie(movie.ID, session.UserID, payload.Title, payload.Description, payload.ReleaseDate, payload.DurationMinutes, payload.Language, payload.Rating).
			Return(movie, nil)

		reqBody := fmt.Sprintf(`{"title": "%s", "description": "%s", "release_date": "%s", "duration_minutes": %d, "language": "%s", "rating": %g}`,
			payload.Title, *payload.Description, payload.ReleaseDate, payload.DurationMinutes, *payload.Language, *payload.Rating)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/movies/%s", movie.ID), bytes.NewBufferString(reqBody))
		req.Header.Set(constants.CONTENT_TYPE, constants.APPLICATION_JSON)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), movie.Title)
		assert.Contains(t, w.Body.String(), *movie.Description)
		assert.Contains(t, w.Body.String(), movie.ReleaseDate)
		assert.Contains(t, w.Body.String(), fmt.Sprint(movie.DurationMinutes))
		assert.Contains(t, w.Body.String(), *movie.Language)
		assert.Contains(t, w.Body.String(), fmt.Sprint(*movie.Rating))
		assert.Contains(t, w.Body.String(), movie.CreatedBy.String())
	})

	t.Run("validation error", func(t *testing.T) {
		router := gin.Default()
		router.PUT("/movies/:id", controller.UpdateMovie)

		reqBody := fmt.Sprintf(`{"title": "%s", "description": "%s", "release_date": "%s", "duration_minutes": %d, "language": "%s", "rating": %g}`,
			payload.Title, *payload.Description, payload.ReleaseDate, payload.DurationMinutes, *payload.Language, 6.0)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/movies/%s", movie.ID), bytes.NewBufferString(reqBody))
		req.Header.Set(constants.CONTENT_TYPE, constants.APPLICATION_JSON)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "errors")
		assert.Contains(t, w.Body.String(), "Should be less than or equal to 5")
	})

	t.Run("service error", func(t *testing.T) {
		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set(constants.USER_SESSION, session)
			c.Next()
		})
		router.PUT("/movies/:id", controller.UpdateMovie)

		service.EXPECT().UpdateMovie(movie.ID, session.UserID, payload.Title, payload.Description, payload.ReleaseDate, payload.DurationMinutes, payload.Language, payload.Rating).
			Return(nil, errors.InternalServerError("Service error"))

		reqBody := fmt.Sprintf(`{"title": "%s", "description": "%s", "release_date": "%s", "duration_minutes": %d, "language": "%s", "rating": %g}`,
			payload.Title, *payload.Description, payload.ReleaseDate, payload.DurationMinutes, *payload.Language, *payload.Rating)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/movies/%s", movie.ID), bytes.NewBufferString(reqBody))
		req.Header.Set(constants.CONTENT_TYPE, constants.APPLICATION_JSON)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Service error")
	})
}
