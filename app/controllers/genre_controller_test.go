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

func TestGenreController_CreateGenre(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockGenreService(ctrl)
	controller := GenreController{
		GenreService: service,
	}

	gin.SetMode(gin.TestMode)

	genre := utils.GenerateRandomGenre()
	payload := utils.GenerateRandomCreateGenreRequest()

	t.Run("success", func(t *testing.T) {
		router := gin.Default()
		router.POST("/genres", controller.CreateGenre)

		service.EXPECT().CreateGenre(payload.Name).Return(genre, nil).Times(1)

		reqBody := fmt.Sprintf(`{"name": "%s"}`, payload.Name)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/genres", bytes.NewBufferString(reqBody))
		req.Header.Set(constants.CONTENT_TYPE, constants.APPLICATION_JSON)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), genre.Name)
	})

	t.Run("validation error", func(t *testing.T) {
		router := gin.Default()
		router.POST("/genres", controller.CreateGenre)

		reqBody := `{}`

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/genres", bytes.NewBufferString(reqBody))
		req.Header.Set(constants.CONTENT_TYPE, constants.APPLICATION_JSON)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "This field is required")
	})

	t.Run("service error", func(t *testing.T) {
		router := gin.Default()
		router.POST("/genres", controller.CreateGenre)

		service.EXPECT().CreateGenre(payload.Name).Return(nil, errors.InternalServerError("service error"))

		reqBody := fmt.Sprintf(`{"name": "%s"}`, payload.Name)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/genres", bytes.NewBufferString(reqBody))
		req.Header.Set(constants.CONTENT_TYPE, constants.APPLICATION_JSON)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "error")
		assert.Contains(t, w.Body.String(), "service error")
	})
}
