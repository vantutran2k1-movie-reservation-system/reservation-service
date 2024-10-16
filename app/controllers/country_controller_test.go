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

func TestCountryController_CreateCountry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockCountryService(ctrl)
	controller := CountryController{
		CountryService: service,
	}

	gin.SetMode(gin.TestMode)

	country := utils.GenerateRandomCountry()
	payload := utils.GenerateRandomCreateCountryRequest()

	t.Run("success", func(t *testing.T) {
		router := gin.Default()
		router.POST("/countries", controller.CreateCountry)

		service.EXPECT().CreateCountry(payload.Name, payload.Code).Return(country, nil).Times(1)

		reqBody := fmt.Sprintf(`{"name": "%s", "code": "%s"}`, payload.Name, payload.Code)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/countries", bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), country.Name)
		assert.Contains(t, w.Body.String(), country.Code)
	})

	t.Run("validation error", func(t *testing.T) {
		router := gin.Default()
		router.POST("/countries", controller.CreateCountry)

		reqBody := fmt.Sprintf(`{"name": "%s", "code": "%s"}`, payload.Name, "INVALID_CODE")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/countries", bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "errors")
		assert.Contains(t, w.Body.String(), "Should be a valid length of 2")
	})

	t.Run("service error", func(t *testing.T) {
		router := gin.Default()
		router.POST("/countries", controller.CreateCountry)

		service.EXPECT().CreateCountry(payload.Name, payload.Code).Return(nil, errors.InternalServerError("service error")).Times(1)

		reqBody := fmt.Sprintf(`{"name": "%s", "code": "%s"}`, payload.Name, payload.Code)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/countries", bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "error")
		assert.Contains(t, w.Body.String(), "service error")
	})
}
