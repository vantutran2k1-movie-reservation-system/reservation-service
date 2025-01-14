package controllers

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/constants"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_services"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/payloads"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLocationController_GetCountries(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockLocationService(ctrl)
	controller := LocationController{
		LocationService: service,
	}

	router := gin.Default()
	router.GET("/countries", controller.GetCountries)

	t.Run("success", func(t *testing.T) {
		countries := utils.GenerateCountries(3)

		service.EXPECT().GetCountries().Return(countries, nil).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/countries", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		for _, country := range countries {
			assert.Contains(t, w.Body.String(), country.Name)
			assert.Contains(t, w.Body.String(), country.Code)
		}
	})

	t.Run("service error", func(t *testing.T) {
		service.EXPECT().GetCountries().Return(nil, errors.InternalServerError("service error")).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/countries", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "service error")
	})
}

func TestLocationController_CreateCountry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockLocationService(ctrl)
	controller := LocationController{
		LocationService: service,
	}

	router := gin.Default()
	router.POST("/countries", controller.CreateCountry)

	country := utils.GenerateCountry()
	payload := payloads.CreateCountryRequest{
		Name: country.Name,
		Code: country.Code,
	}

	t.Run("success", func(t *testing.T) {
		service.EXPECT().CreateCountry(payload).Return(country, nil).Times(1)

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
		service.EXPECT().CreateCountry(payload).Return(nil, errors.InternalServerError("service error")).Times(1)

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

func TestLocationController_GetStatesByCountry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockLocationService(ctrl)
	controller := LocationController{
		LocationService: service,
	}

	router := gin.Default()
	router.GET("/countries/:countryId/states", controller.GetStatesByCountry)

	countryID := uuid.New()

	t.Run("success", func(t *testing.T) {
		states := utils.GenerateStates(3)

		service.EXPECT().GetStatesByCountry(countryID).Return(states, nil).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/countries/%s/states", countryID), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		for _, state := range states {
			assert.Contains(t, w.Body.String(), state.ID.String())
			assert.Contains(t, w.Body.String(), state.Name)
			assert.Contains(t, w.Body.String(), *state.Code)
		}
	})

	t.Run("service error", func(t *testing.T) {
		service.EXPECT().GetStatesByCountry(countryID).Return(nil, errors.InternalServerError("service error")).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/countries/%s/states", countryID), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "service error")
	})
}

func TestLocationController_CreateState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockLocationService(ctrl)
	controller := LocationController{
		LocationService: service,
	}

	router := gin.Default()
	router.POST("/countries/:countryId/states", controller.CreateState)

	state := utils.GenerateState()
	payload := payloads.CreateStateRequest{
		Name: state.Name,
		Code: state.Code,
	}

	t.Run("success", func(t *testing.T) {
		service.EXPECT().CreateState(state.CountryID, payload).Return(state, nil).Times(1)

		reqBody := fmt.Sprintf(`{"name": "%s", "code": "%s"}`, payload.Name, *payload.Code)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/countries/%s/states", state.CountryID), bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), state.Name)
		assert.Contains(t, w.Body.String(), *state.Code)
	})

	t.Run("invalid country id", func(t *testing.T) {
		reqBody := fmt.Sprintf(`{"name": "%s", "code": "%s"}`, payload.Name, *payload.Code)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/countries/%s/states", "invalid id"), bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid country id")
	})

	t.Run("validation error", func(t *testing.T) {
		reqBody := fmt.Sprintf(`{"name": "%s", "code": "%s"}`, payload.Name, "A")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/countries/%s/states", state.CountryID), bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Should be greater than or equal to 2")
	})

	t.Run("service error", func(t *testing.T) {
		service.EXPECT().CreateState(state.CountryID, payload).Return(nil, errors.InternalServerError("service error")).Times(1)

		reqBody := fmt.Sprintf(`{"name": "%s", "code": "%s"}`, payload.Name, *payload.Code)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/countries/%s/states", state.CountryID), bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "service error")
	})
}

func TestLocationController_GetCitiesByState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockLocationService(ctrl)
	controller := LocationController{
		LocationService: service,
	}

	router := gin.Default()
	router.GET("/countries/:countryId/states/:stateId/cities", controller.GetCitiesByState)

	cities := utils.GenerateCities(3)

	t.Run("success", func(t *testing.T) {
		service.EXPECT().GetCitiesByState(gomock.Any(), gomock.Any()).Return(cities, nil).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/countries/%s/states/%s/cities", uuid.NewString(), uuid.NewString()), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		for _, city := range cities {
			assert.Contains(t, w.Body.String(), city.ID.String())
		}
	})

	t.Run("service error", func(t *testing.T) {
		service.EXPECT().GetCitiesByState(gomock.Any(), gomock.Any()).Return(nil, errors.InternalServerError("service error")).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/countries/%s/states/%s/cities", uuid.NewString(), uuid.NewString()), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "service error")
	})
}

func TestLocationController_CreateCity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockLocationService(ctrl)
	controller := LocationController{
		LocationService: service,
	}

	router := gin.Default()
	router.POST("/countries/:countryId/states/:stateId/cities", controller.CreateCity)

	state := utils.GenerateState()
	city := utils.GenerateCity()
	payload := payloads.CreateCityRequest{
		Name: city.Name,
	}

	t.Run("success", func(t *testing.T) {
		service.EXPECT().CreateCity(state.CountryID, city.StateID, payload).Return(city, nil).Times(1)

		reqBody := fmt.Sprintf(`{"name": "%s"}`, payload.Name)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/countries/%s/states/%s/cities", state.CountryID, city.StateID), bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), city.Name)
	})

	t.Run("invalid country id", func(t *testing.T) {
		reqBody := fmt.Sprintf(`{"name": "%s"}`, payload.Name)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/countries/%s/states/%s/cities", "invalid id", city.StateID), bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid country id")
	})

	t.Run("invalid state id", func(t *testing.T) {
		reqBody := fmt.Sprintf(`{"name": "%s"}`, payload.Name)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/countries/%s/states/%s/cities", state.CountryID, "invalid id"), bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid state id")
	})

	t.Run("validation error", func(t *testing.T) {
		reqBody := fmt.Sprintf(`{"name": "%s"}`, "A")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/countries/%s/states/%s/cities", state.CountryID, city.StateID), bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Should be greater than or equal to 2")
	})

	t.Run("service error", func(t *testing.T) {
		service.EXPECT().CreateCity(state.CountryID, city.StateID, payload).Return(nil, errors.InternalServerError("service error")).Times(1)

		reqBody := fmt.Sprintf(`{"name": "%s"}`, payload.Name)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/countries/%s/states/%s/cities", state.CountryID, city.StateID), bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "service error")
	})
}
