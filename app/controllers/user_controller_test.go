package controllers

import (
	"bytes"
	"fmt"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/context"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/payloads"
	"net/http"
	"net/http/httptest"
	"strconv"
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

	user := utils.GenerateUser()

	router := gin.Default()
	router.GET("/users/:userId", controller.GetUser)

	t.Run("success", func(t *testing.T) {
		service.EXPECT().GetUser(user.ID, true).Return(user, nil).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/users/%s?%s=%v", user.ID, constants.IncludeUserProfile, true), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), user.Email)
		assert.Contains(t, w.Body.String(), user.ID.String())
	})

	t.Run("error getting user", func(t *testing.T) {
		service.EXPECT().GetUser(user.ID, true).Return(nil, errors.InternalServerError("error getting user")).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/users/%s?%s=%v", user.ID, constants.IncludeUserProfile, true), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "error getting user")
	})
}

func TestUserController_GetCurrentUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockUserService(ctrl)
	controller := UserController{
		UserService: service,
	}

	session := utils.GenerateUserSession()
	user := utils.GenerateUser()

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		context.SetRequestContext(c, context.RequestContext{UserSession: session})
		c.Next()
	})
	router.GET("/users/me", controller.GetCurrentUser)

	t.Run("success", func(t *testing.T) {
		service.EXPECT().GetUser(session.UserID, true).Return(user, nil).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/users/me?%s=%v", constants.IncludeUserProfile, true), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), user.Email)
		assert.Contains(t, w.Body.String(), user.ID.String())
	})

	t.Run("session not found", func(t *testing.T) {
		routerErr := gin.Default()
		routerErr.GET("/users/me", controller.GetCurrentUser)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/users/me?%s=%v", constants.IncludeUserProfile, true), nil)
		routerErr.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("error retrieving user", func(t *testing.T) {
		service.EXPECT().GetUser(session.UserID, true).Return(nil, errors.BadRequestError("user not found"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/users/me?%s=%v", constants.IncludeUserProfile, true), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "user not found")
	})
}

func TestUserController_UserExistsByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockUserService(ctrl)
	controller := UserController{
		UserService: service,
	}

	user := utils.GenerateUser()

	router := gin.Default()
	router.GET("/users/exists", controller.UserExistsByEmail)

	t.Run("success", func(t *testing.T) {
		service.EXPECT().UserExistsByEmail(user.Email).Return(true, nil).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/users/exists?%s=%s", constants.Email, user.Email), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), strconv.FormatBool(true))
	})

	t.Run("user not exists", func(t *testing.T) {
		service.EXPECT().UserExistsByEmail(user.Email).Return(false, nil).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/users/exists?%s=%s", constants.Email, user.Email), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), strconv.FormatBool(false))
	})

	t.Run("error getting user", func(t *testing.T) {
		service.EXPECT().UserExistsByEmail(user.Email).Return(false, errors.InternalServerError("error getting user")).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/users/exists?%s=%s", constants.Email, user.Email), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "error getting user")
	})
}

func TestUserController_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockUserService(ctrl)
	controller := UserController{
		UserService: service,
	}

	payload := payloads.CreateUserRequest{
		Email:    "example@example.com",
		Password: "password",
		Profile: payloads.CreateUserProfileRequest{
			FirstName:   "First",
			LastName:    "Last",
			PhoneNumber: utils.GetPointerOf("0000000000"),
			DateOfBirth: utils.GetPointerOf("1970-01-01"),
		},
	}

	router := gin.Default()
	router.POST("/user", controller.CreateUser)

	t.Run("successful user creation", func(t *testing.T) {
		service.EXPECT().CreateUser(payload).Return(&models.User{Email: payload.Email}, nil)

		reqBody := fmt.Sprintf(`{"email": "%s", "password": "%s", "profile": {"first_name": "%s", "last_name": "%s", "phone_number": "%s", "date_of_birth": "%s"}}`, payload.Email, payload.Password, payload.Profile.FirstName, payload.Profile.LastName, *payload.Profile.PhoneNumber, *payload.Profile.DateOfBirth)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/user", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), payload.Email)
	})

	t.Run("validation error", func(t *testing.T) {
		reqBody := `{"email": "invalid-email", "password": "short"}`

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/user", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "errors")
	})

	t.Run("service error", func(t *testing.T) {
		service.EXPECT().CreateUser(payload).Return(nil, errors.InternalServerError("Service error"))

		reqBody := fmt.Sprintf(`{"email": "%s", "password": "%s", "profile": {"first_name": "%s", "last_name": "%s", "phone_number": "%s", "date_of_birth": "%s"}}`, payload.Email, payload.Password, payload.Profile.FirstName, payload.Profile.LastName, *payload.Profile.PhoneNumber, *payload.Profile.DateOfBirth)
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

	token := utils.GenerateLoginToken()
	payload := payloads.LoginUserRequest{
		Email:    "example@example.com",
		Password: "test password",
	}

	router := gin.Default()
	router.POST("/login", userController.LoginUser)

	t.Run("successful login", func(t *testing.T) {
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
		reqBody := `{"email": "invalid-email", "password": ""}`

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "errors")
	})

	t.Run("service error", func(t *testing.T) {
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

func TestUserController_VerifyUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockUserService(ctrl)
	controller := UserController{
		UserService: service,
	}

	token := utils.GenerateUserRegistrationToken()

	router := gin.Default()
	router.POST("/users/verify", controller.VerifyUser)

	t.Run("success", func(t *testing.T) {
		service.EXPECT().VerifyUser(token.TokenValue).Return(nil).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/users/verify", nil)
		req.Header.Set(constants.UserVerificationToken, token.TokenValue)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		service.EXPECT().VerifyUser(token.TokenValue).Return(errors.InternalServerError("service error")).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/users/verify", nil)
		req.Header.Set(constants.UserVerificationToken, token.TokenValue)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "service error")
	})
}

func TestUserController_UpdateUserPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockUserService(ctrl)
	controller := UserController{
		UserService: service,
	}

	session := utils.GenerateUserSession()
	payload := payloads.UpdatePasswordRequest{Password: "example password"}

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		context.SetRequestContext(c, context.RequestContext{UserSession: session})
		c.Next()
	})
	router.PUT("/users/password", controller.UpdateUserPassword)

	t.Run("successful password update", func(t *testing.T) {
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
		invalidReqBody := `{"password": ""}`

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/users/password", bytes.NewBufferString(invalidReqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "errors")
	})

	t.Run("session retrieval error", func(t *testing.T) {
		routerErr := gin.Default()
		routerErr.PUT("/users/password", controller.UpdateUserPassword)

		reqBody := fmt.Sprintf(`{"password": "%s"}`, payload.Password)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/users/password", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		routerErr.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
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

	token := utils.GeneratePasswordResetToken()
	payload := payloads.CreatePasswordResetTokenRequest{
		Email: "example@example.com",
	}

	router := gin.Default()
	router.POST("/users/password-reset-token", controller.CreatePasswordResetToken)

	t.Run("success", func(t *testing.T) {
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
		reqBody := fmt.Sprintf(`{"email": "%s"}`, "invalid email")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/users/password-reset-token", bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Should be a valid email address")
	})

	t.Run("service error", func(t *testing.T) {
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

func TestUserController_ResetPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_services.NewMockUserService(ctrl)
	controller := UserController{
		UserService: service,
	}

	token := utils.GeneratePasswordResetToken()
	payload := payloads.ResetPasswordRequest{Password: "test password"}

	router := gin.Default()
	router.POST("/users/password-reset", controller.ResetPassword)

	t.Run("success", func(t *testing.T) {
		service.EXPECT().ResetUserPassword(token.TokenValue, payload).Return(nil).Times(1)

		reqBody := fmt.Sprintf(`{"password": "%s"}`, payload.Password)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/users/password-reset", bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		req.Header.Set(constants.ResetToken, token.TokenValue)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("validation error", func(t *testing.T) {
		reqBody := fmt.Sprintf(`{"password": "%s"}`, "A")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/users/password-reset", bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		req.Header.Set(constants.ResetToken, token.TokenValue)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Should be greater than or equal to 8")
	})

	t.Run("service error", func(t *testing.T) {
		service.EXPECT().ResetUserPassword(token.TokenValue, payload).Return(errors.InternalServerError("service error")).Times(1)

		reqBody := fmt.Sprintf(`{"password": "%s"}`, payload.Password)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/users/password-reset", bytes.NewBufferString(reqBody))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		req.Header.Set(constants.ResetToken, token.TokenValue)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "service error")
	})
}
