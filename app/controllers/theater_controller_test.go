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
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/payloads"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestTheaterController_GetTheater(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockTheaterService(ctrl)
	controller := TheaterController{
		TheaterService: service,
	}

	theater := utils.GenerateTheater()

	router := gin.Default()
	router.GET("/theaters/:theaterId", controller.GetTheater)

	t.Run("success", func(t *testing.T) {
		service.EXPECT().GetTheater(theater.ID, true).Return(theater, nil).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/theaters/%s?%s=%v", theater.ID, constants.IncludeTheaterLocation, true), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), theater.Name)
	})

	t.Run("service error", func(t *testing.T) {
		service.EXPECT().GetTheater(theater.ID, true).Return(nil, errors.InternalServerError("service error")).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/theaters/%s?%s=%v", theater.ID, constants.IncludeTheaterLocation, true), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "service error")

	})
}

func TestTheaterController_GetNearbyTheaters(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockTheaterService(ctrl)
	controller := TheaterController{
		TheaterService: service,
	}

	theaters := utils.GenerateTheaters(3)
	distance := 10.0

	router := gin.Default()
	router.GET("/theaters/nearby", controller.GetNearbyTheaters)

	t.Run("success", func(t *testing.T) {
		service.EXPECT().GetNearbyTheaters(distance).Return(theaters, nil).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/theaters/nearby?%s=%v.0", constants.MaxDistance, distance), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		for _, theater := range theaters {
			assert.Contains(t, w.Body.String(), theater.ID.String())
			assert.Contains(t, w.Body.String(), theater.Name)
		}
	})

	t.Run("service error", func(t *testing.T) {
		service.EXPECT().GetNearbyTheaters(distance).Return(nil, errors.InternalServerError("service error")).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/theaters/nearby?%s=%v.0", constants.MaxDistance, distance), nil)
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

	theater := utils.GenerateTheater()
	payload := payloads.CreateTheaterRequest{Name: theater.Name}

	router := gin.Default()
	router.POST("/theaters", controller.CreateTheater)

	t.Run("success", func(t *testing.T) {
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
		reqBody := fmt.Sprintf(`{"name": "%s"}`, "A")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/theaters", bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Should be greater than or equal to 2")
	})

	t.Run("service error", func(t *testing.T) {
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

	theater := utils.GenerateTheater()
	location := utils.GenerateTheaterLocation()
	payload := payloads.CreateTheaterLocationRequest{
		CityID:     uuid.New(),
		Address:    "address",
		PostalCode: "700000",
		Latitude:   50.0,
		Longitude:  100.0,
	}

	router := gin.Default()
	router.POST("/theaters/:theaterId/locations", controller.CreateTheaterLocation)

	t.Run("success", func(t *testing.T) {
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
		reqBody := fmt.Sprintf(`{"city_id": "%s", "address": "%s", "postal_code": "%s", "latitude": %v, "longitude": %v}`, payload.CityID, "A", payload.PostalCode, payload.Latitude, payload.Longitude)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/theaters/%s/locations", theater.ID), bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Should be greater than or equal to 2")
	})

	t.Run("service error", func(t *testing.T) {
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

func TestTheaterController_CreateSeat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockTheaterService(ctrl)
	controller := TheaterController{
		TheaterService: service,
	}

	theater := utils.GenerateTheater()
	seat := utils.GenerateSeat()
	payload := payloads.CreateSeatPayload{
		Row:    seat.Row,
		Number: seat.Number,
		Type:   seat.Type,
	}

	router := gin.Default()
	router.POST("/theaters/:theaterId/seats", controller.CreateSeat)

	t.Run("success", func(t *testing.T) {
		service.EXPECT().CreateSeat(theater.ID, payload).Return(seat, nil).Times(1)

		reqBody := fmt.Sprintf(`{"row": "%s", "number": %d, "type": "%s"}`, seat.Row, seat.Number, seat.Type)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/theaters/%s/seats", theater.ID), bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), seat.Row)
		assert.Contains(t, w.Body.String(), fmt.Sprintf(strconv.Itoa(seat.Number)))
		assert.Contains(t, w.Body.String(), seat.Type)
	})

	t.Run("validation error", func(t *testing.T) {
		reqBody := fmt.Sprintf(`{"row": "%s", "number": %d, "type": "%s"}`, seat.Row, seat.Number, "INVALID TYPE")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/theaters/%s/seats", theater.ID), bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), fmt.Sprintf("Should be one of %s, %s", constants.Regular, constants.Vip))
	})

	t.Run("service error", func(t *testing.T) {
		service.EXPECT().CreateSeat(theater.ID, payload).Return(nil, errors.InternalServerError("service error")).Times(1)

		reqBody := fmt.Sprintf(`{"row": "%s", "number": %d, "type": "%s"}`, seat.Row, seat.Number, seat.Type)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/theaters/%s/seats", theater.ID), bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "service error")
	})
}

func TestTheaterController_UpdateTheaterLocation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockTheaterService(ctrl)
	controller := TheaterController{
		TheaterService: service,
	}

	theater := utils.GenerateTheater()
	location := utils.GenerateTheaterLocation()
	payload := payloads.UpdateTheaterLocationRequest{
		CityID:     uuid.New(),
		Address:    "address",
		PostalCode: "700000",
		Latitude:   50.0,
		Longitude:  100.0,
	}

	router := gin.Default()
	router.PUT("/theaters/:theaterId/locations", controller.UpdateTheaterLocation)

	t.Run("success", func(t *testing.T) {
		service.EXPECT().UpdateTheaterLocation(theater.ID, payload).Return(location, nil).Times(1)

		reqBody := fmt.Sprintf(`{"city_id": "%s", "address": "%s", "postal_code": "%s", "latitude": %v, "longitude": %v}`, payload.CityID, payload.Address, payload.PostalCode, payload.Latitude, payload.Longitude)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/theaters/%s/locations", theater.ID), bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), location.Address)
		assert.Contains(t, w.Body.String(), location.PostalCode)
	})

	t.Run("validation error", func(t *testing.T) {
		reqBody := fmt.Sprintf(`{"city_id": "%s", "address": "%s", "postal_code": "%s", "latitude": %v, "longitude": %v}`, payload.CityID, "A", payload.PostalCode, payload.Latitude, payload.Longitude)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/theaters/%s/locations", theater.ID), bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Should be greater than or equal to 2")
	})

	t.Run("service error", func(t *testing.T) {
		service.EXPECT().UpdateTheaterLocation(theater.ID, payload).Return(nil, errors.InternalServerError("service error")).Times(1)

		reqBody := fmt.Sprintf(`{"city_id": "%s", "address": "%s", "postal_code": "%s", "latitude": %v, "longitude": %v}`, payload.CityID, payload.Address, payload.PostalCode, payload.Latitude, payload.Longitude)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/theaters/%s/locations", theater.ID), bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "service error")
	})
}

func TestTheaterController_GetTheaters(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockTheaterService(ctrl)
	controller := TheaterController{
		TheaterService: service,
	}

	theaters := utils.GenerateTheaters(3)
	meta := &models.ResponseMeta{
		Limit:   10,
		Offset:  0,
		Total:   20,
		NextUrl: utils.GetPointerOf("/theaters"),
		PrevUrl: utils.GetPointerOf("/theaters"),
	}
	includeLocation := true
	limit := 3
	offset := 1

	router := gin.Default()
	router.GET("/theaters", controller.GetTheaters)

	t.Run("success", func(t *testing.T) {
		service.EXPECT().GetTheaters(limit, offset, includeLocation).Return(theaters, meta, nil).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/theaters?%s=%d&%s=%d&%s=%v", constants.Limit, limit, constants.Offset, offset, constants.IncludeTheaterLocation, includeLocation), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), fmt.Sprint(meta.Limit))
		assert.Contains(t, w.Body.String(), fmt.Sprint(meta.Offset))
		assert.Contains(t, w.Body.String(), fmt.Sprint(meta.Total))
		assert.Contains(t, w.Body.String(), *meta.NextUrl)
		assert.Contains(t, w.Body.String(), *meta.PrevUrl)
		for _, theater := range theaters {
			assert.Contains(t, w.Body.String(), theater.ID.String())
		}
	})

	t.Run("service error", func(t *testing.T) {
		service.EXPECT().GetTheaters(limit, offset, includeLocation).Return(nil, nil, errors.InternalServerError("service error")).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/theaters?%s=%d&%s=%d&%s=%v", constants.Limit, limit, constants.Offset, offset, constants.IncludeTheaterLocation, includeLocation), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "service error")
	})
}
