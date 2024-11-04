package controllers

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/constants"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_services"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTheaterController_GetTheater(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockTheaterService(ctrl)
	controller := TheaterController{
		TheaterService: service,
	}

	gin.SetMode(gin.TestMode)

	theater := utils.GenerateTheater()

	t.Run("success", func(t *testing.T) {
		router := gin.Default()
		router.GET("/theaters/:theaterId", controller.GetTheater)

		service.EXPECT().GetTheater(theater.ID, true).Return(theater, nil).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/theaters/%s?%s=%v", theater.ID, constants.IncludeTheaterLocation, true), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), theater.Name)
	})

	t.Run("service error", func(t *testing.T) {
		router := gin.Default()
		router.GET("/theaters/:theaterId", controller.GetTheater)

		service.EXPECT().GetTheater(theater.ID, true).Return(nil, errors.InternalServerError("service error")).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/theaters/%s?%s=%v", theater.ID, constants.IncludeTheaterLocation, true), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "service error")

	})
}

func TestTheaterController_CreateTheater(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockTheaterService(ctrl)
	controller := TheaterController{
		TheaterService: service,
	}

	gin.SetMode(gin.TestMode)

	theater := utils.GenerateTheater()
	payload := utils.GenerateCreateTheaterRequest()

	t.Run("success", func(t *testing.T) {
		router := gin.Default()
		router.POST("/theaters", controller.CreateTheater)

		service.EXPECT().CreateTheater(payload).Return(theater, nil).Times(1)

		reqBody := fmt.Sprintf(`{"name": "%s"}`, payload.Name)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/theaters", bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), theater.Name)
	})

	t.Run("validation error", func(t *testing.T) {
		router := gin.Default()
		router.POST("/theaters", controller.CreateTheater)

		reqBody := fmt.Sprintf(`{"name": "%s"}`, "A")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/theaters", bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Should be greater than or equal to 2")
	})

	t.Run("service error", func(t *testing.T) {
		router := gin.Default()
		router.POST("/theaters", controller.CreateTheater)

		service.EXPECT().CreateTheater(payload).Return(nil, errors.InternalServerError("service error")).Times(1)

		reqBody := fmt.Sprintf(`{"name": "%s"}`, payload.Name)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/theaters", bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "service error")
	})
}

func TestTheaterController_CreateTheaterLocation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockTheaterService(ctrl)
	controller := TheaterController{
		TheaterService: service,
	}

	gin.SetMode(gin.TestMode)

	theater := utils.GenerateTheater()
	location := utils.GenerateTheaterLocation()
	payload := utils.GenerateCreateTheaterLocationRequest()

	t.Run("success", func(t *testing.T) {
		router := gin.Default()
		router.POST("/theaters/:theaterId/locations", controller.CreateTheaterLocation)

		service.EXPECT().CreateTheaterLocation(theater.ID, payload).Return(location, nil).Times(1)

		reqBody := fmt.Sprintf(`{"city_id": "%s", "address": "%s", "postal_code": "%s", "latitude": %v, "longitude": %v}`, payload.CityID, payload.Address, payload.PostalCode, payload.Latitude, payload.Longitude)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/theaters/%s/locations", theater.ID), bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), location.Address)
		assert.Contains(t, w.Body.String(), location.PostalCode)
	})

	t.Run("validation error", func(t *testing.T) {
		router := gin.Default()
		router.POST("/theaters/:theaterId/locations", controller.CreateTheaterLocation)

		reqBody := fmt.Sprintf(`{"city_id": "%s", "address": "%s", "postal_code": "%s", "latitude": %v, "longitude": %v}`, payload.CityID, "A", payload.PostalCode, payload.Latitude, payload.Longitude)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/theaters/%s/locations", theater.ID), bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Should be greater than or equal to 2")
	})

	t.Run("service error", func(t *testing.T) {
		router := gin.Default()
		router.POST("/theaters/:theaterId/locations", controller.CreateTheaterLocation)

		service.EXPECT().CreateTheaterLocation(theater.ID, payload).Return(nil, errors.InternalServerError("service error")).Times(1)

		reqBody := fmt.Sprintf(`{"city_id": "%s", "address": "%s", "postal_code": "%s", "latitude": %v, "longitude": %v}`, payload.CityID, payload.Address, payload.PostalCode, payload.Latitude, payload.Longitude)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/theaters/%s/locations", theater.ID), bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "service error")
	})
}
