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

func TestGenreController_GetGenre(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockGenreService(ctrl)
	controller := GenreController{
		GenreService: service,
	}

	gin.SetMode(gin.TestMode)

	genre := utils.GenerateGenre()

	t.Run("success", func(t *testing.T) {
		router := gin.Default()
		router.GET("/genres/:id", controller.GetGenre)

		service.EXPECT().GetGenre(genre.ID).Return(genre, nil).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/genres/%s", genre.ID), nil)
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), genre.Name)
	})

	t.Run("invalid id", func(t *testing.T) {
		router := gin.Default()
		router.GET("/genres/:id", controller.GetGenre)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/genres/%s", "test id"), nil)
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid genre id")
	})

	t.Run("genre not found", func(t *testing.T) {
		router := gin.Default()
		router.GET("/genres/:id", controller.GetGenre)

		service.EXPECT().GetGenre(gomock.Any()).Return(nil, errors.NotFoundError("not found")).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/genres/%s", genre.ID), nil)
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "not found")
	})

	t.Run("service error", func(t *testing.T) {
		router := gin.Default()
		router.GET("/genres/:id", controller.GetGenre)

		service.EXPECT().GetGenre(genre.ID).Return(nil, errors.InternalServerError("service error")).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/genres/%s", genre.ID), nil)
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "service error")
	})
}

func TestGenreController_GetGenres(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockGenreService(ctrl)
	controller := GenreController{
		GenreService: service,
	}

	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		router := gin.Default()
		router.GET("/genres", controller.GetGenres)

		genres := utils.GenerateGenres(3)

		service.EXPECT().GetGenres().Return(genres, nil).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/genres", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		for _, genre := range genres {
			assert.Contains(t, w.Body.String(), genre.Name)
		}
	})

	t.Run("service error", func(t *testing.T) {
		router := gin.Default()
		router.GET("/genres", controller.GetGenres)

		service.EXPECT().GetGenres().Return(nil, errors.InternalServerError("service error")).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/genres", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "service error")
	})
}

func TestGenreController_CreateGenre(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockGenreService(ctrl)
	controller := GenreController{
		GenreService: service,
	}

	gin.SetMode(gin.TestMode)

	genre := utils.GenerateGenre()
	payload := utils.GenerateCreateGenreRequest()

	t.Run("success", func(t *testing.T) {
		router := gin.Default()
		router.POST("/genres", controller.CreateGenre)

		service.EXPECT().CreateGenre(payload.Name).Return(genre, nil).Times(1)

		reqBody := fmt.Sprintf(`{"name": "%s"}`, payload.Name)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/genres", bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
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
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
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
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "error")
		assert.Contains(t, w.Body.String(), "service error")
	})
}

func TestGenreController_UpdateGenre(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockGenreService(ctrl)
	controller := GenreController{
		GenreService: service,
	}

	gin.SetMode(gin.TestMode)

	genre := utils.GenerateGenre()
	payload := utils.GenerateUpdateGenreRequest()

	t.Run("success", func(t *testing.T) {
		router := gin.Default()
		router.PUT("/genres/:id", controller.UpdateGenre)

		service.EXPECT().UpdateGenre(genre.ID, payload.Name).Return(genre, nil).Times(1)

		reqBody := fmt.Sprintf(`{"name": "%s"}`, payload.Name)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/genres/%s", genre.ID), bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), genre.Name)
	})

	t.Run("error validating data", func(t *testing.T) {
		router := gin.Default()
		router.PUT("/genres/:id", controller.UpdateGenre)

		reqBody := `{}`

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/genres/%s", genre.ID), bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "This field is required")
	})

	t.Run("error updating genre", func(t *testing.T) {
		router := gin.Default()
		router.PUT("/genres/:id", controller.UpdateGenre)

		service.EXPECT().UpdateGenre(genre.ID, payload.Name).Return(nil, errors.InternalServerError("error updating genre")).Times(1)

		reqBody := fmt.Sprintf(`{"name": "%s"}`, payload.Name)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/genres/%s", genre.ID), bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "error")
		assert.Contains(t, w.Body.String(), "error updating genre")
	})
}
