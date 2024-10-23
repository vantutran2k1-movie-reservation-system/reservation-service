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
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"go.uber.org/mock/gomock"
)

func TestUserController_GetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockUserService(ctrl)
	controller := UserController{
		UserService: service,
	}

	gin.SetMode(gin.TestMode)

	session := utils.GenerateUserSession()
	user := utils.GenerateUser()

	t.Run("successful user retrieval", func(t *testing.T) {
		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set(constants.UserSession, session)
			c.Next()
		})
		router.GET("/user", controller.GetUser)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/user", nil)

		service.EXPECT().GetUser(session.UserID).Return(user, nil).Times(1)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), user.Email)
		assert.Contains(t, w.Body.String(), user.ID.String())
	})

	t.Run("session not found", func(t *testing.T) {
		router := gin.Default()
		router.GET("/user", controller.GetUser)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/user", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("error retrieving user", func(t *testing.T) {
		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set(constants.UserSession, session)
			c.Next()
		})
		router.GET("/user", controller.GetUser)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/user", nil)

		service.EXPECT().GetUser(session.UserID).Return(nil, errors.BadRequestError("user not found"))

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "user not found")
	})
}

func TestUserController_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockUserService(ctrl)
	controller := UserController{
		UserService: service,
	}

	gin.SetMode(gin.TestMode)

	payload := utils.GenerateCreateUserRequest()

	t.Run("successful user creation", func(t *testing.T) {
		router := gin.Default()
		router.POST("/user", controller.CreateUser)

		service.EXPECT().CreateUser(payload).Return(&models.User{Email: payload.Email}, nil)

		reqBody := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, payload.Email, payload.Password)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/user", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), payload.Email)
	})

	t.Run("validation error", func(t *testing.T) {
		router := gin.Default()
		router.POST("/user", controller.CreateUser)

		reqBody := `{"email": "invalid-email", "password": "short"}`

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/user", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "errors")
	})

	t.Run("service error", func(t *testing.T) {
		router := gin.Default()
		router.POST("/user", controller.CreateUser)

		service.EXPECT().CreateUser(payload).Return(nil, errors.InternalServerError("Service error"))

		reqBody := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, payload.Email, payload.Password)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/user", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "error")
	})
}

func TestUserController_LoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_services.NewMockUserService(ctrl)
	userController := UserController{
		UserService: mockUserService,
	}

	gin.SetMode(gin.TestMode)

	token := utils.GenerateLoginToken()
	payload := utils.GenerateLoginUserRequest()

	t.Run("successful login", func(t *testing.T) {
		router := gin.Default()
		router.POST("/login", userController.LoginUser)

		mockUserService.EXPECT().LoginUser(payload).Return(token, nil)

		reqBody := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, payload.Email, payload.Password)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "token")
		assert.Contains(t, w.Body.String(), token.TokenValue)
	})

	t.Run("validation error", func(t *testing.T) {
		router := gin.Default()
		router.POST("/login", userController.LoginUser)

		reqBody := `{"email": "invalid-email", "password": ""}`

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "errors")
	})

	t.Run("service error", func(t *testing.T) {
		router := gin.Default()
		router.POST("/login", userController.LoginUser)

		mockUserService.EXPECT().LoginUser(payload).Return(nil, errors.UnauthorizedError("Invalid credentials"))

		reqBody := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, payload.Email, payload.Password)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid credentials")
	})
}

func TestUserController_UpdateUserPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockUserService(ctrl)
	controller := UserController{
		UserService: service,
	}

	gin.SetMode(gin.TestMode)

	session := utils.GenerateUserSession()
	payload := utils.GenerateUpdatePasswordRequest()

	t.Run("successful password update", func(t *testing.T) {
		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set(constants.UserSession, session)
			c.Next()
		})
		router.PUT("/users/password", controller.UpdateUserPassword)

		service.EXPECT().UpdateUserPassword(session.UserID, payload).Return(nil)

		reqBody := fmt.Sprintf(`{"password": "%s"}`, payload.Password)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/users/password", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Password is updated successfully")
	})

	t.Run("validation error", func(t *testing.T) {
		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set(constants.UserSession, session)
			c.Next()
		})
		router.PUT("/users/password", controller.UpdateUserPassword)

		invalidReqBody := `{"password": ""}`

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/users/password", bytes.NewBufferString(invalidReqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "errors")
	})

	t.Run("session retrieval error", func(t *testing.T) {
		router := gin.Default()
		router.PUT("/users/password", controller.UpdateUserPassword)

		reqBody := fmt.Sprintf(`{"password": "%s"}`, payload.Password)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/users/password", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set(constants.UserSession, session)
			c.Next()
		})
		router.PUT("/users/password", controller.UpdateUserPassword)

		service.EXPECT().UpdateUserPassword(session.UserID, payload).Return(errors.InternalServerError("Failed to update password"))

		reqBody := fmt.Sprintf(`{"password": "%s"}`, payload.Password)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/users/password", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to update password")
	})
}

func TestUserController_CreatePasswordResetToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockUserService(ctrl)
	controller := UserController{
		UserService: service,
	}

	gin.SetMode(gin.TestMode)

	token := utils.GeneratePasswordResetToken()
	payload := utils.GenerateCreatePasswordResetTokenRequest()

	t.Run("success", func(t *testing.T) {
		router := gin.Default()
		router.POST("/users/password-reset-token", controller.CreatePasswordResetToken)

		service.EXPECT().CreatePasswordResetToken(payload).Return(token, nil).Times(1)

		reqBody := fmt.Sprintf(`{"email": "%s"}`, payload.Email)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/users/password-reset-token", bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), token.TokenValue)
	})

	t.Run("validation error", func(t *testing.T) {
		router := gin.Default()
		router.POST("/users/password-reset-token", controller.CreatePasswordResetToken)

		reqBody := fmt.Sprintf(`{"email": "%s"}`, "invalid email")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/users/password-reset-token", bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Should be a valid email address")
	})

	t.Run("service error", func(t *testing.T) {
		router := gin.Default()
		router.POST("/users/password-reset-token", controller.CreatePasswordResetToken)

		service.EXPECT().CreatePasswordResetToken(payload).Return(nil, errors.InternalServerError("service error")).Times(1)

		reqBody := fmt.Sprintf(`{"email": "%s"}`, payload.Email)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/users/password-reset-token", bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "service error")
	})
}
