package controllers

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
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

func TestShowController_GetActiveShows(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockShowService(ctrl)
	controller := ShowController{
		ShowService: service,
	}

	router := gin.Default()
	router.GET("/shows/active", controller.GetActiveShows)

	limit := 3
	offset := 0
	shows := utils.GenerateShows(3)

	t.Run("success", func(t *testing.T) {
		service.EXPECT().GetShows(constants.Active, limit, offset).Return(shows, nil).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/shows/active?%s=%d&%s=%d", constants.Limit, limit, constants.Offset, offset), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		for _, show := range shows {
			assert.Contains(t, w.Body.String(), show.Id.String())
		}
	})

	t.Run("service error", func(t *testing.T) {
		service.EXPECT().GetShows(constants.Active, limit, offset).Return(nil, errors.InternalServerError("service error")).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/shows/active?%s=%d&%s=%d", constants.Limit, limit, constants.Offset, offset), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "service error")
	})
}

func TestShowController_CreateShow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockShowService(ctrl)
	controller := ShowController{
		ShowService: service,
	}

	router := gin.Default()
	router.POST("/shows", controller.CreateShow)

	show := utils.GenerateShow()
	payload := payloads.CreateShowRequest{
		MovieId:   *show.MovieId,
		TheaterId: *show.TheaterId,
		StartTime: show.StartTime,
		EndTime:   show.EndTime,
		Status:    show.Status,
	}

	t.Run("success", func(t *testing.T) {
		service.EXPECT().CreateShow(gomock.Any()).Return(show, nil).Times(1)

		reqBody := fmt.Sprintf(
			`{"movie_id": "%s", "theater_id": "%s", "start_time": "%s", "end_time": "%s", "status": "%s"}`,
			payload.MovieId, payload.TheaterId, payload.StartTime.Format(constants.DateTimeFormat), payload.EndTime.Format(constants.DateTimeFormat), payload.Status,
		)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/shows", bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), payload.MovieId.String())
	})

	t.Run("validation error", func(t *testing.T) {
		reqBody := fmt.Sprintf(
			`{"movie_id": "%s", "theater_id": "%s", "start_time": "%s", "end_time": "%s", "status": "%s"}`,
			payload.MovieId, payload.TheaterId, payload.StartTime.Format(constants.DateTimeFormat), payload.EndTime.Format(constants.DateTimeFormat), "Invalid status",
		)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/shows", bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Should be one of ACTIVE, CANCELLED, COMPLETED, EXPIRED, SCHEDULED, ON-HOLD")
	})

	t.Run("service error", func(t *testing.T) {
		service.EXPECT().CreateShow(gomock.Any()).Return(nil, errors.InternalServerError("service error")).Times(1)

		reqBody := fmt.Sprintf(
			`{"movie_id": "%s", "theater_id": "%s", "start_time": "%s", "end_time": "%s", "status": "%s"}`,
			payload.MovieId, payload.TheaterId, payload.StartTime.Format(constants.DateTimeFormat), payload.EndTime.Format(constants.DateTimeFormat), payload.Status,
		)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/shows", bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "service error")
	})
}
