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
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"go.uber.org/mock/gomock"
)

func TestMovieController_GetUser(t *testing.T) {
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

	t.Run("successful movie creation", func(t *testing.T) {
		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set(constants.USER_SESSION, session)
			c.Next()
		})
		router.POST("/movies", controller.CreateMovie)

		service.EXPECT().CreateMovie(session.UserID, payload.Title, payload.Description, payload.ReleaseDate, payload.DurationMinutes, payload.Language, payload.Rating).
			Return(movie, nil)

		reqBody := fmt.Sprintf(`{"title": "%s", "description": "%s", "release_date": "%s", "duration_minutes": %d, "language": "%s", "rating": %g}`,
			payload.Title, *payload.Description, payload.ReleaseDate, payload.DurationMinutes, *payload.Language, *payload.Rating)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/movies", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), movie.Title)
	})

	t.Run("validation error", func(t *testing.T) {
		router := gin.Default()
		router.POST("/movies", controller.CreateMovie)

		reqBody := fmt.Sprintf(`{"title": "%s", "description": "%s", "release_date": "%s", "duration_minutes": %d, "language": "%s", "rating": %g}`,
			payload.Title, *payload.Description, payload.ReleaseDate, payload.DurationMinutes, *payload.Language, 6.0)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/movies", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
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
		router.POST("/movies", controller.CreateMovie)

		service.EXPECT().CreateMovie(session.UserID, payload.Title, payload.Description, payload.ReleaseDate, payload.DurationMinutes, payload.Language, payload.Rating).
			Return(nil, errors.InternalServerError("Service error"))

		reqBody := fmt.Sprintf(`{"title": "%s", "description": "%s", "release_date": "%s", "duration_minutes": %d, "language": "%s", "rating": %g}`,
			payload.Title, *payload.Description, payload.ReleaseDate, payload.DurationMinutes, *payload.Language, *payload.Rating)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/movies", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "error")
	})
}
