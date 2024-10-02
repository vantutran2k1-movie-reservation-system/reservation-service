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

func TestUserProfileController_GetProfileByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockUserProfileService(ctrl)
	controller := UserProfileController{
		UserProfileService: service,
	}

	gin.SetMode(gin.TestMode)

	session := utils.GenerateSampleUserSession()
	profile := utils.GenerateSampleUserProfile()

	t.Run("successful profile retrieval", func(t *testing.T) {
		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set(constants.USER_SESSION, session)
			c.Next()
		})
		router.GET("/profiles", controller.GetProfileByUserID)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/profiles", nil)

		service.EXPECT().GetProfileByUserID(session.UserID).Return(profile, nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), profile.FirstName)
		assert.Contains(t, w.Body.String(), profile.LastName)
	})

	t.Run("session retrieval error", func(t *testing.T) {
		router := gin.Default()
		router.GET("/profiles", controller.GetProfileByUserID)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/profiles", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set(constants.USER_SESSION, session)
			c.Next()
		})
		router.GET("/profiles", controller.GetProfileByUserID)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/profiles", nil)

		service.EXPECT().GetProfileByUserID(session.UserID).Return(nil, errors.InternalServerError("Failed to retrieve profile"))

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to retrieve profile")
	})
}

func TestUserProfileController_CreateUserProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockUserProfileService(ctrl)
	controller := UserProfileController{
		UserProfileService: service,
	}

	gin.SetMode(gin.TestMode)

	session := utils.GenerateSampleUserSession()
	profile := utils.GenerateSampleUserProfile()
	payload := utils.GenerateSampleCreateUserProfileRequest()

	t.Run("successful profile creation", func(t *testing.T) {
		errors.RegisterCustomValidators()

		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set(constants.USER_SESSION, session)
			c.Next()
		})
		router.POST("/profiles", controller.CreateUserProfile)

		service.EXPECT().
			CreateUserProfile(session.UserID, payload.FirstName, payload.LastName, payload.PhoneNumber, gomock.Any()).
			Return(profile, nil)

		reqBody := fmt.Sprintf(`{"first_name": "%s", "last_name": "%s", "phone_number": "%s","date_of_birth": "%s"}`,
			payload.FirstName, payload.LastName, *payload.PhoneNumber, *payload.DateOfBirth)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/profiles", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), profile.FirstName)
		assert.Contains(t, w.Body.String(), profile.LastName)
	})

	t.Run("validation error", func(t *testing.T) {
		errors.RegisterCustomValidators()

		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set(constants.USER_SESSION, session)
			c.Next()
		})
		router.POST("/profiles", controller.CreateUserProfile)

		reqBody := fmt.Sprintf(`{"first_name": "", "last_name": "%s", "phone_number": "%s","date_of_birth": "2024-13-01"}`,
			payload.LastName, *payload.PhoneNumber)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/profiles", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "This field is required")
		assert.Contains(t, w.Body.String(), "Should be a valid date with format YYYY-MM-DD")
	})

	t.Run("session retrieval error", func(t *testing.T) {
		errors.RegisterCustomValidators()

		router := gin.Default()
		router.POST("/profiles", controller.CreateUserProfile)

		reqBody := fmt.Sprintf(`{"first_name": "%s", "last_name": "%s", "phone_number": "%s","date_of_birth": "%s"}`,
			payload.FirstName, payload.LastName, *payload.PhoneNumber, *payload.DateOfBirth)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/profiles", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		errors.RegisterCustomValidators()

		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set(constants.USER_SESSION, session)
			c.Next()
		})
		router.POST("/profiles", controller.CreateUserProfile)

		reqBody := fmt.Sprintf(`{"first_name": "%s", "last_name": "%s", "phone_number": "%s","date_of_birth": "%s"}`,
			payload.FirstName, payload.LastName, *payload.PhoneNumber, *payload.DateOfBirth)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/profiles", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		service.EXPECT().
			CreateUserProfile(session.UserID, payload.FirstName, payload.LastName, payload.PhoneNumber, payload.DateOfBirth).
			Return(nil, errors.InternalServerError("Failed to create profile"))

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to create profile")
	})
}

func TestUserProfileController_UpdateUserProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockUserProfileService(ctrl)
	controller := UserProfileController{
		UserProfileService: service,
	}

	gin.SetMode(gin.TestMode)

	session := utils.GenerateSampleUserSession()
	profile := utils.GenerateSampleUserProfile()
	payload := utils.GenerateSampleUpdateUserProfileRequest()

	t.Run("successful profile update", func(t *testing.T) {
		errors.RegisterCustomValidators()

		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set(constants.USER_SESSION, session)
			c.Next()
		})
		router.PUT("/profiles", controller.UpdateUserProfile)

		service.EXPECT().
			UpdateUserProfile(session.UserID, payload.FirstName, payload.LastName, payload.PhoneNumber, payload.DateOfBirth).
			Return(profile, nil)

		reqBody := fmt.Sprintf(`{"first_name": "%s", "last_name": "%s", "phone_number": "%s","date_of_birth": "%s"}`,
			payload.FirstName, payload.LastName, *payload.PhoneNumber, *payload.DateOfBirth)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/profiles", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), profile.FirstName)
		assert.Contains(t, w.Body.String(), profile.LastName)
	})

	t.Run("validation error", func(t *testing.T) {
		errors.RegisterCustomValidators()

		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set(constants.USER_SESSION, session)
			c.Next()
		})
		router.PUT("/profiles", controller.CreateUserProfile)

		reqBody := fmt.Sprintf(`{"first_name": "", "last_name": "%s", "phone_number": "%s","date_of_birth": "2024-13-01"}`,
			payload.LastName, *payload.PhoneNumber)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/profiles", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "This field is required")
		assert.Contains(t, w.Body.String(), "Should be a valid date with format YYYY-MM-DD")
	})

	t.Run("session retrieval error", func(t *testing.T) {
		errors.RegisterCustomValidators()

		router := gin.Default()
		router.PUT("/profiles", controller.CreateUserProfile)

		reqBody := fmt.Sprintf(`{"first_name": "%s", "last_name": "%s", "phone_number": "%s","date_of_birth": "%s"}`,
			payload.FirstName, payload.LastName, *payload.PhoneNumber, *payload.DateOfBirth)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/profiles", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		errors.RegisterCustomValidators()

		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set(constants.USER_SESSION, session)
			c.Next()
		})
		router.PUT("/profiles", controller.UpdateUserProfile)

		service.EXPECT().
			UpdateUserProfile(session.UserID, payload.FirstName, payload.LastName, payload.PhoneNumber, payload.DateOfBirth).
			Return(nil, errors.InternalServerError("Failed to update profile"))

		reqBody := fmt.Sprintf(`{"first_name": "%s", "last_name": "%s", "phone_number": "%s","date_of_birth": "%s"}`,
			payload.FirstName, payload.LastName, *payload.PhoneNumber, *payload.DateOfBirth)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/profiles", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to update profile")
	})
}