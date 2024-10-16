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
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStateController_GetStatesByCountry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockStateService(ctrl)
	controller := StateController{
		StateService: service,
	}

	gin.SetMode(gin.TestMode)

	countryID := uuid.New()

	t.Run("success", func(t *testing.T) {
		router := gin.Default()
		router.GET("/countries/:countryId/states", controller.GetStatesByCountry)

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
		router := gin.Default()
		router.GET("/countries/:countryId/states", controller.GetStatesByCountry)

		service.EXPECT().GetStatesByCountry(countryID).Return(nil, errors.InternalServerError("service error")).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/countries/%s/states", countryID), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "service error")
	})
}

func TestStateController_CreateState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockStateService(ctrl)
	controller := StateController{
		StateService: service,
	}

	gin.SetMode(gin.TestMode)

	state := utils.GenerateState()
	payload := utils.GenerateCreateStateRequest()

	t.Run("success", func(t *testing.T) {
		router := gin.Default()
		router.POST("/countries/:countryId/states", controller.CreateState)

		service.EXPECT().CreateState(state.CountryID, payload.Name, payload.Code).Return(state, nil).Times(1)

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
		router := gin.Default()
		router.POST("/countries/:countryId/states", controller.CreateState)

		reqBody := fmt.Sprintf(`{"name": "%s", "code": "%s"}`, payload.Name, *payload.Code)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/countries/%s/states", "invalid id"), bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid country id")
	})

	t.Run("validation error", func(t *testing.T) {
		router := gin.Default()
		router.POST("/countries/:countryId/states", controller.CreateState)

		reqBody := fmt.Sprintf(`{"name": "%s", "code": "%s"}`, payload.Name, "A")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/countries/%s/states", state.CountryID), bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Should be greater than or equal to 2")
	})

	t.Run("service error", func(t *testing.T) {
		router := gin.Default()
		router.POST("/countries/:countryId/states", controller.CreateState)

		service.EXPECT().CreateState(state.CountryID, payload.Name, payload.Code).Return(nil, errors.InternalServerError("service error")).Times(1)

		reqBody := fmt.Sprintf(`{"name": "%s", "code": "%s"}`, payload.Name, *payload.Code)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/countries/%s/states", state.CountryID), bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "service error")
	})
}
