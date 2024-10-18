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

func TestCityController_CreateCity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockCityService(ctrl)
	controller := CityController{
		CityService: service,
	}

	gin.SetMode(gin.TestMode)

	state := utils.GenerateState()
	city := utils.GenerateCity()
	payload := utils.GenerateCreateCityRequest()

	t.Run("success", func(t *testing.T) {
		router := gin.Default()
		router.POST("/countries/:countryId/states/:stateId/cities", controller.CreateCity)

		service.EXPECT().CreateCity(state.CountryID, city.StateID, payload.Name).Return(city, nil).Times(1)

		reqBody := fmt.Sprintf(`{"name": "%s"}`, payload.Name)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/countries/%s/states/%s/cities", state.CountryID, city.StateID), bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), city.Name)
	})

	t.Run("invalid country id", func(t *testing.T) {
		router := gin.Default()
		router.POST("/countries/:countryId/states/:stateId/cities", controller.CreateCity)

		reqBody := fmt.Sprintf(`{"name": "%s"}`, payload.Name)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/countries/%s/states/%s/cities", "invalid id", city.StateID), bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid country id")
	})

	t.Run("invalid state id", func(t *testing.T) {
		router := gin.Default()
		router.POST("/countries/:countryId/states/:stateId/cities", controller.CreateCity)

		reqBody := fmt.Sprintf(`{"name": "%s"}`, payload.Name)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/countries/%s/states/%s/cities", state.CountryID, "invalid id"), bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid state id")
	})

	t.Run("validation error", func(t *testing.T) {
		router := gin.Default()
		router.POST("/countries/:countryId/states/:stateId/cities", controller.CreateCity)

		reqBody := fmt.Sprintf(`{"name": "%s"}`, "A")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/countries/%s/states/%s/cities", state.CountryID, city.StateID), bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Should be greater than or equal to 2")
	})

	t.Run("service error", func(t *testing.T) {
		router := gin.Default()
		router.POST("/countries/:countryId/states/:stateId/cities", controller.CreateCity)

		service.EXPECT().CreateCity(state.CountryID, city.StateID, payload.Name).Return(nil, errors.InternalServerError("service error")).Times(1)

		reqBody := fmt.Sprintf(`{"name": "%s"}`, payload.Name)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/countries/%s/states/%s/cities", state.CountryID, city.StateID), bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "service error")
	})
}
