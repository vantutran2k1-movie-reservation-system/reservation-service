package services

import (
	"errors"
	"fmt"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"net/http"
	"os"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_auth"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_repositories"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_transaction"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestUserService_GetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repositories.NewMockUserRepository(ctrl)
	service := NewUserService(nil, nil, nil, nil, repo, nil, nil, nil)

	user := utils.GenerateRandomUser()

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().GetUser(user.ID).Return(user, nil).Times(1)

		result, err := service.GetUser(user.ID)

		assert.NotNil(t, user)
		assert.Nil(t, err)
		assert.Equal(t, user, result)
	})

	t.Run("user not found", func(t *testing.T) {
		repo.EXPECT().GetUser(user.ID).Return(nil, nil).Times(1)

		result, err := service.GetUser(user.ID)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "user does not exist", err.Error())
	})

	t.Run("error getting user", func(t *testing.T) {
		repo.EXPECT().GetUser(user.ID).Return(nil, errors.New("repo error")).Times(1)

		result, err := service.GetUser(user.ID)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "repo error", err.Error())
	})
}

func TestUserService_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	auth := mock_auth.NewMockAuthenticator(ctrl)
	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	repo := mock_repositories.NewMockUserRepository(ctrl)
	service := NewUserService(nil, nil, auth, transaction, repo, nil, nil, nil)

	user := utils.GenerateRandomUser()
	password := utils.GenerateRandomPassword()

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().GetUserByEmail(user.Email).Return(nil, nil).Times(1)
		auth.EXPECT().GenerateHashedPassword(password).Return(user.PasswordHash, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		repo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		result, err := service.CreateUser(user.Email, password)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, user.Email, result.Email)
		assert.Equal(t, user.PasswordHash, result.PasswordHash)
	})

	t.Run("duplicate email", func(t *testing.T) {
		repo.EXPECT().GetUserByEmail(user.Email).Return(user, nil).Times(1)

		result, err := service.CreateUser(user.Email, password)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.StatusCode)
		assert.Equal(t, fmt.Sprintf("email %s already exists", user.Email), err.Error())
	})

	t.Run("error getting user", func(t *testing.T) {
		repo.EXPECT().GetUserByEmail(user.Email).Return(nil, errors.New("error getting user")).Times(1)

		result, err := service.CreateUser(user.Email, password)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting user", err.Error())
	})

	t.Run("error generating password hash", func(t *testing.T) {
		repo.EXPECT().GetUserByEmail(user.Email).Return(nil, nil).Times(1)
		auth.EXPECT().GenerateHashedPassword(password).Return("", errors.New("error generating password hash"))

		result, err := service.CreateUser(user.Email, password)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error generating password hash", err.Error())
	})

	t.Run("error creating user", func(t *testing.T) {
		repo.EXPECT().GetUserByEmail(user.Email).Return(nil, nil).Times(1)
		auth.EXPECT().GenerateHashedPassword(password).Return(user.PasswordHash, nil)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		repo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(errors.New("error creating user")).Times(1)

		result, err := service.CreateUser(user.Email, password)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error creating user", err.Error())
	})
}

func TestUserService_LoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	auth := mock_auth.NewMockAuthenticator(ctrl)
	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	userRepo := mock_repositories.NewMockUserRepository(ctrl)
	loginTokenRepo := mock_repositories.NewMockLoginTokenRepository(ctrl)
	userSessionRepo := mock_repositories.NewMockUserSessionRepository(ctrl)

	service := NewUserService(nil, nil, auth, transaction, userRepo, loginTokenRepo, userSessionRepo, nil)

	user := utils.GenerateRandomUser()
	token := utils.GenerateRandomLoginToken()
	password := utils.GenerateRandomPassword()

	t.Run("success", func(t *testing.T) {
		userRepo.EXPECT().GetUserByEmail(user.Email).Return(user, nil).Times(1)
		auth.EXPECT().DoPasswordsMatch(user.PasswordHash, password).Return(true).Times(1)
		auth.EXPECT().GenerateLoginToken().Return(token.TokenValue).Times(1)
		loginTokenRepo.EXPECT().GetActiveLoginToken(token.TokenValue).Return(nil, nil).Times(1)
		os.Setenv("LOGIN_TOKEN_EXPIRES_AFTER_MINUTES", "60")
		defer os.Unsetenv("LOGIN_TOKEN_EXPIRES_AFTER_MINUTES")
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		loginTokenRepo.EXPECT().CreateLoginToken(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		transaction.EXPECT().ExecuteInRedisTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(rdb *redis.Client, fn func(tx *redis.Tx) error) error {
				return fn(nil)
			},
		).Times(1)
		userSessionRepo.EXPECT().GetUserSessionID(token.TokenValue).Return(token.TokenValue).Times(1)
		userSessionRepo.EXPECT().CreateUserSession(token.TokenValue, gomock.Any(), gomock.Any()).Return(nil).Times(1)

		result, err := service.LoginUser(user.Email, password)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, user.ID, result.UserID)
		assert.Equal(t, token.TokenValue, result.TokenValue)
	})

	t.Run("invalid email", func(t *testing.T) {
		userRepo.EXPECT().GetUserByEmail(user.Email).Return(nil, nil).Times(1)

		result, err := service.LoginUser(user.Email, password)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusUnauthorized, err.StatusCode)
		assert.Equal(t, fmt.Sprintf("invalid email %s", user.Email), err.Error())
	})

	t.Run("invalid password", func(t *testing.T) {
		userRepo.EXPECT().GetUserByEmail(user.Email).Return(user, nil).Times(1)
		auth.EXPECT().DoPasswordsMatch(user.PasswordHash, password).Return(false).Times(1)

		result, err := service.LoginUser(user.Email, password)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusUnauthorized, err.StatusCode)
		assert.Equal(t, "invalid password", err.Error())
	})

	t.Run("error getting token", func(t *testing.T) {
		userRepo.EXPECT().GetUserByEmail(user.Email).Return(user, nil).Times(1)
		auth.EXPECT().DoPasswordsMatch(user.PasswordHash, password).Return(true).Times(1)
		auth.EXPECT().GenerateLoginToken().Return(token.TokenValue).Times(1)
		loginTokenRepo.EXPECT().GetActiveLoginToken(token.TokenValue).Return(nil, errors.New("error getting token")).Times(1)

		result, err := service.LoginUser(user.Email, password)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting token", err.Error())
	})

	t.Run("token already exists", func(t *testing.T) {
		userRepo.EXPECT().GetUserByEmail(user.Email).Return(user, nil).Times(1)
		auth.EXPECT().DoPasswordsMatch(user.PasswordHash, password).Return(true).Times(1)
		auth.EXPECT().GenerateLoginToken().Return(token.TokenValue).Times(1)
		loginTokenRepo.EXPECT().GetActiveLoginToken(token.TokenValue).Return(token, nil).Times(1)

		result, err := service.LoginUser(user.Email, password)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "token value already exists", err.Error())
	})

	t.Run("error getting expiry time", func(t *testing.T) {
		userRepo.EXPECT().GetUserByEmail(user.Email).Return(user, nil).Times(1)
		auth.EXPECT().DoPasswordsMatch(user.PasswordHash, password).Return(true).Times(1)
		auth.EXPECT().GenerateLoginToken().Return(token.TokenValue).Times(1)
		loginTokenRepo.EXPECT().GetActiveLoginToken(token.TokenValue).Return(nil, nil).Times(1)

		result, err := service.LoginUser(user.Email, password)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "invalid token expiry time: strconv.Atoi: parsing \"\": invalid syntax", err.Error())
	})

	t.Run("error creating token", func(t *testing.T) {
		userRepo.EXPECT().GetUserByEmail(user.Email).Return(user, nil).Times(1)
		auth.EXPECT().DoPasswordsMatch(user.PasswordHash, password).Return(true).Times(1)
		auth.EXPECT().GenerateLoginToken().Return(token.TokenValue).Times(1)
		loginTokenRepo.EXPECT().GetActiveLoginToken(token.TokenValue).Return(nil, nil).Times(1)
		os.Setenv("LOGIN_TOKEN_EXPIRES_AFTER_MINUTES", "60")
		defer os.Unsetenv("LOGIN_TOKEN_EXPIRES_AFTER_MINUTES")
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		loginTokenRepo.EXPECT().CreateLoginToken(gomock.Any(), gomock.Any()).Return(errors.New("error creating token")).Times(1)

		result, err := service.LoginUser(user.Email, password)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error creating token", err.Error())
	})

	t.Run("error creating user session", func(t *testing.T) {
		userRepo.EXPECT().GetUserByEmail(user.Email).Return(user, nil).Times(1)
		auth.EXPECT().DoPasswordsMatch(user.PasswordHash, password).Return(true).Times(1)
		auth.EXPECT().GenerateLoginToken().Return(token.TokenValue).Times(1)
		loginTokenRepo.EXPECT().GetActiveLoginToken(token.TokenValue).Return(nil, nil).Times(1)
		os.Setenv("LOGIN_TOKEN_EXPIRES_AFTER_MINUTES", "60")
		defer os.Unsetenv("LOGIN_TOKEN_EXPIRES_AFTER_MINUTES")
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		loginTokenRepo.EXPECT().CreateLoginToken(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		transaction.EXPECT().ExecuteInRedisTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(rdb *redis.Client, fn func(tx *redis.Tx) error) error {
				return fn(nil)
			},
		).Times(1)
		userSessionRepo.EXPECT().GetUserSessionID(token.TokenValue).Return(token.TokenValue).Times(1)
		userSessionRepo.EXPECT().CreateUserSession(token.TokenValue, gomock.Any(), gomock.Any()).Return(errors.New("error creating session")).Times(1)

		result, err := service.LoginUser(user.Email, password)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error creating session", err.Error())
	})
}

func TestUserService_LogoutUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	userSessionRepo := mock_repositories.NewMockUserSessionRepository(ctrl)
	loginTokenRepo := mock_repositories.NewMockLoginTokenRepository(ctrl)

	service := NewUserService(nil, nil, nil, transaction, nil, loginTokenRepo, userSessionRepo, nil)

	token := utils.GenerateRandomLoginToken()

	t.Run("success", func(t *testing.T) {
		loginTokenRepo.EXPECT().GetActiveLoginToken(token.TokenValue).Return(token, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		loginTokenRepo.EXPECT().RevokeLoginToken(gomock.Any(), token).Return(nil).Times(1)
		transaction.EXPECT().ExecuteInRedisTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(rdb *redis.Client, fn func(tx *redis.Tx) error) error {
				return fn(nil)
			},
		).Times(1)
		userSessionRepo.EXPECT().GetUserSessionID(token.TokenValue).Return(token.TokenValue).Times(1)
		userSessionRepo.EXPECT().DeleteUserSession(token.TokenValue).Return(nil).Times(1)

		err := service.LogoutUser(token.TokenValue)

		assert.Nil(t, err)
	})

	t.Run("token not found", func(t *testing.T) {
		loginTokenRepo.EXPECT().GetActiveLoginToken(token.TokenValue).Return(nil, nil).Times(1)

		err := service.LogoutUser(token.TokenValue)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.StatusCode)
		assert.Equal(t, "token not found", err.Error())
	})

	t.Run("error getting token", func(t *testing.T) {
		loginTokenRepo.EXPECT().GetActiveLoginToken(token.TokenValue).Return(nil, errors.New("error getting token")).Times(1)

		err := service.LogoutUser(token.TokenValue)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting token", err.Error())
	})

	t.Run("error revoking token", func(t *testing.T) {
		loginTokenRepo.EXPECT().GetActiveLoginToken(token.TokenValue).Return(token, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		loginTokenRepo.EXPECT().RevokeLoginToken(gomock.Any(), token).Return(errors.New("error revoking token")).Times(1)

		err := service.LogoutUser(token.TokenValue)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error revoking token", err.Error())
	})

	t.Run("error deleting session", func(t *testing.T) {
		loginTokenRepo.EXPECT().GetActiveLoginToken(token.TokenValue).Return(token, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		loginTokenRepo.EXPECT().RevokeLoginToken(gomock.Any(), token).Return(nil).Times(1)
		transaction.EXPECT().ExecuteInRedisTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(rdb *redis.Client, fn func(tx *redis.Tx) error) error {
				return fn(nil)
			},
		).Times(1)
		userSessionRepo.EXPECT().GetUserSessionID(token.TokenValue).Return(token.TokenValue).Times(1)
		userSessionRepo.EXPECT().DeleteUserSession(token.TokenValue).Return(errors.New("error deleting session")).Times(1)

		err := service.LogoutUser(token.TokenValue)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error deleting session", err.Error())
	})
}

func TestUserService_UpdateUserPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	auth := mock_auth.NewMockAuthenticator(ctrl)
	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	userRepo := mock_repositories.NewMockUserRepository(ctrl)
	loginTokenRepo := mock_repositories.NewMockLoginTokenRepository(ctrl)
	userSessionRepo := mock_repositories.NewMockUserSessionRepository(ctrl)

	service := NewUserService(nil, nil, auth, transaction, userRepo, loginTokenRepo, userSessionRepo, nil)

	user := utils.GenerateRandomUser()
	updatedUser := utils.GenerateRandomUser()
	updatedUser.ID = user.ID
	password := utils.GenerateRandomPassword()

	t.Run("success", func(t *testing.T) {
		userRepo.EXPECT().GetUser(user.ID).Return(user, nil).Times(1)
		auth.EXPECT().DoPasswordsMatch(user.PasswordHash, password).Return(false).Times(1)
		auth.EXPECT().GenerateHashedPassword(password).Return(updatedUser.PasswordHash, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		userRepo.EXPECT().UpdatePassword(gomock.Any(), user, updatedUser.PasswordHash).Return(updatedUser, nil).Times(1)
		loginTokenRepo.EXPECT().RevokeUserLoginTokens(gomock.Any(), user.ID).Return(nil).Times(1)
		transaction.EXPECT().ExecuteInRedisTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(rdb *redis.Client, fn func(tx *redis.Tx) error) error {
				return fn(nil)
			},
		).Times(1)
		userSessionRepo.EXPECT().DeleteUserSessions(user.ID).Return(nil).Times(1)

		err := service.UpdateUserPassword(user.ID, password)

		assert.Nil(t, err)
	})

	t.Run("user not found", func(t *testing.T) {
		userRepo.EXPECT().GetUser(user.ID).Return(nil, nil).Times(1)

		err := service.UpdateUserPassword(user.ID, password)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "user does not exist", err.Error())
	})

	t.Run("error getting user", func(t *testing.T) {
		userRepo.EXPECT().GetUser(user.ID).Return(nil, errors.New("error getting user")).Times(1)

		err := service.UpdateUserPassword(user.ID, password)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting user", err.Error())
	})

	t.Run("error password same as current value", func(t *testing.T) {
		userRepo.EXPECT().GetUser(user.ID).Return(user, nil).Times(1)
		auth.EXPECT().DoPasswordsMatch(user.PasswordHash, password).Return(true).Times(1)

		err := service.UpdateUserPassword(user.ID, password)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.StatusCode)
		assert.Equal(t, "new password can not be the same as current value", err.Error())
	})

	t.Run("error generating hashed password", func(t *testing.T) {
		userRepo.EXPECT().GetUser(user.ID).Return(user, nil).Times(1)
		auth.EXPECT().DoPasswordsMatch(user.PasswordHash, password).Return(false).Times(1)
		auth.EXPECT().GenerateHashedPassword(password).Return("", errors.New("hashing error")).Times(1)

		err := service.UpdateUserPassword(user.ID, password)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "hashing error", err.Error())
	})

	t.Run("error updating password", func(t *testing.T) {
		userRepo.EXPECT().GetUser(user.ID).Return(user, nil).Times(1)
		auth.EXPECT().DoPasswordsMatch(user.PasswordHash, password).Return(false).Times(1)
		auth.EXPECT().GenerateHashedPassword(password).Return(updatedUser.PasswordHash, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		userRepo.EXPECT().UpdatePassword(gomock.Any(), user, updatedUser.PasswordHash).Return(nil, errors.New("error updating password")).Times(1)

		err := service.UpdateUserPassword(user.ID, password)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error updating password", err.Error())
	})

	t.Run("error revoking tokens", func(t *testing.T) {
		userRepo.EXPECT().GetUser(user.ID).Return(user, nil).Times(1)
		auth.EXPECT().DoPasswordsMatch(user.PasswordHash, password).Return(false).Times(1)
		auth.EXPECT().GenerateHashedPassword(password).Return(updatedUser.PasswordHash, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		userRepo.EXPECT().UpdatePassword(gomock.Any(), user, updatedUser.PasswordHash).Return(updatedUser, nil).Times(1)
		loginTokenRepo.EXPECT().RevokeUserLoginTokens(gomock.Any(), user.ID).Return(errors.New("error revoking tokens")).Times(1)

		err := service.UpdateUserPassword(user.ID, password)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error revoking tokens", err.Error())
	})

	t.Run("error deleting sessions", func(t *testing.T) {
		userRepo.EXPECT().GetUser(user.ID).Return(user, nil).Times(1)
		auth.EXPECT().DoPasswordsMatch(user.PasswordHash, password).Return(false).Times(1)
		auth.EXPECT().GenerateHashedPassword(password).Return(updatedUser.PasswordHash, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		userRepo.EXPECT().UpdatePassword(gomock.Any(), user, updatedUser.PasswordHash).Return(updatedUser, nil).Times(1)
		loginTokenRepo.EXPECT().RevokeUserLoginTokens(gomock.Any(), user.ID).Return(nil).Times(1)
		transaction.EXPECT().ExecuteInRedisTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(rdb *redis.Client, fn func(tx *redis.Tx) error) error {
				return fn(nil)
			},
		).Times(1)
		userSessionRepo.EXPECT().DeleteUserSessions(user.ID).Return(errors.New("error deleting sessions")).Times(1)

		err := service.UpdateUserPassword(user.ID, password)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error deleting sessions", err.Error())
	})
}

func TestUserService_CreatePasswordResetToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	auth := mock_auth.NewMockAuthenticator(ctrl)
	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	userRepo := mock_repositories.NewMockUserRepository(ctrl)
	tokenRepo := mock_repositories.NewMockPasswordResetTokenRepository(ctrl)

	service := NewUserService(nil, nil, auth, transaction, userRepo, nil, nil, tokenRepo)

	user := utils.GenerateRandomUser()
	token := utils.GenerateRandomPasswordResetToken()

	t.Run("success", func(t *testing.T) {
		userRepo.EXPECT().GetUserByEmail(user.Email).Return(user, nil).Times(1)
		auth.EXPECT().GeneratePasswordResetToken().Return(token.TokenValue).Times(1)
		tokenRepo.EXPECT().GetActivePasswordResetToken(token.TokenValue).Return(nil, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		os.Setenv("PASSWORD_RESET_TOKEN_EXPIRES_AFTER_MINUTES", "60")
		defer os.Unsetenv("PASSWORD_RESET_TOKEN_EXPIRES_AFTER_MINUTES")
		tokenRepo.EXPECT().CreateToken(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		result, err := service.CreatePasswordResetToken(user.Email)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, user.ID, result.UserID)
		assert.Equal(t, token.TokenValue, result.TokenValue)
		assert.Equal(t, false, result.IsUsed)
	})

	t.Run("error getting user", func(t *testing.T) {
		userRepo.EXPECT().GetUserByEmail(user.Email).Return(nil, errors.New("error getting user")).Times(1)

		result, err := service.CreatePasswordResetToken(user.Email)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting user", err.Error())
	})

	t.Run("error email not found", func(t *testing.T) {
		userRepo.EXPECT().GetUserByEmail(user.Email).Return(nil, nil).Times(1)

		result, err := service.CreatePasswordResetToken(user.Email)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, fmt.Sprintf("email %s does not exist", user.Email), err.Error())
	})

	t.Run("error getting active tokens", func(t *testing.T) {
		userRepo.EXPECT().GetUserByEmail(user.Email).Return(user, nil).Times(1)
		auth.EXPECT().GeneratePasswordResetToken().Return(token.TokenValue).Times(1)
		tokenRepo.EXPECT().GetActivePasswordResetToken(token.TokenValue).Return(nil, errors.New("error getting tokens")).Times(1)

		result, err := service.CreatePasswordResetToken(user.Email)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting tokens", err.Error())
	})

	t.Run("duplicate token", func(t *testing.T) {
		userRepo.EXPECT().GetUserByEmail(user.Email).Return(user, nil).Times(1)
		auth.EXPECT().GeneratePasswordResetToken().Return(token.TokenValue).Times(1)
		tokenRepo.EXPECT().GetActivePasswordResetToken(token.TokenValue).Return(token, nil).Times(1)

		result, err := service.CreatePasswordResetToken(user.Email)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "token value already exists", err.Error())
	})

	t.Run("error getting expire time", func(t *testing.T) {
		userRepo.EXPECT().GetUserByEmail(user.Email).Return(user, nil).Times(1)
		auth.EXPECT().GeneratePasswordResetToken().Return(token.TokenValue).Times(1)
		tokenRepo.EXPECT().GetActivePasswordResetToken(token.TokenValue).Return(nil, nil).Times(1)

		result, err := service.CreatePasswordResetToken(user.Email)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Contains(t, err.Error(), "invalid token expiry time")
	})

	t.Run("error creating token", func(t *testing.T) {
		userRepo.EXPECT().GetUserByEmail(user.Email).Return(user, nil).Times(1)
		auth.EXPECT().GeneratePasswordResetToken().Return(token.TokenValue).Times(1)
		tokenRepo.EXPECT().GetActivePasswordResetToken(token.TokenValue).Return(nil, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		os.Setenv("PASSWORD_RESET_TOKEN_EXPIRES_AFTER_MINUTES", "60")
		defer os.Unsetenv("PASSWORD_RESET_TOKEN_EXPIRES_AFTER_MINUTES")
		tokenRepo.EXPECT().CreateToken(gomock.Any(), gomock.Any()).Return(errors.New("error creating token")).Times(1)

		result, err := service.CreatePasswordResetToken(user.Email)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error creating token", err.Error())
	})
}

func TestUserService_ResetUserPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	auth := mock_auth.NewMockAuthenticator(ctrl)
	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	userRepo := mock_repositories.NewMockUserRepository(ctrl)
	sessionRepo := mock_repositories.NewMockUserSessionRepository(ctrl)
	loginTokenRepo := mock_repositories.NewMockLoginTokenRepository(ctrl)
	resetTokenRepo := mock_repositories.NewMockPasswordResetTokenRepository(ctrl)

	service := NewUserService(nil, nil, auth, transaction, userRepo, loginTokenRepo, sessionRepo, resetTokenRepo)

	resetToken := utils.GenerateRandomPasswordResetToken()
	user := utils.GenerateRandomUser()
	password := utils.GenerateRandomPassword()
	hashedPassword := utils.GenerateRandomHashedPassword()
	allResetTokens := make([]*models.PasswordResetToken, 3)
	for i := 0; i < len(allResetTokens)-1; i++ {
		allResetTokens[i] = utils.GenerateRandomPasswordResetToken()
	}
	allResetTokens[len(allResetTokens)-1] = resetToken

	t.Run("success", func(t *testing.T) {
		resetTokenRepo.EXPECT().GetActivePasswordResetToken(resetToken.TokenValue).Return(resetToken, nil).Times(1)
		userRepo.EXPECT().GetUser(resetToken.UserID).Return(user, nil).Times(1)
		auth.EXPECT().GenerateHashedPassword(password).Return(hashedPassword, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		userRepo.EXPECT().UpdatePassword(gomock.Any(), user, hashedPassword).Return(user, nil).Times(1)
		loginTokenRepo.EXPECT().RevokeUserLoginTokens(gomock.Any(), user.ID).Return(nil).Times(1)
		resetTokenRepo.EXPECT().UseToken(gomock.Any(), resetToken).Return(nil).Times(1)
		resetTokenRepo.EXPECT().GetUserActivePasswordResetTokens(resetToken.UserID).Return(allResetTokens, nil).Times(1)
		resetTokenRepo.EXPECT().RevokeTokens(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		transaction.EXPECT().ExecuteInRedisTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(rdb *redis.Client, fn func(tx *redis.Tx) error) error {
				return fn(nil)
			},
		).Times(1)
		sessionRepo.EXPECT().DeleteUserSessions(user.ID).Return(nil).Times(1)

		err := service.ResetUserPassword(resetToken.TokenValue, password)

		assert.Nil(t, err)
	})

	t.Run("reset token not found", func(t *testing.T) {
		resetTokenRepo.EXPECT().GetActivePasswordResetToken(resetToken.TokenValue).Return(nil, nil).Times(1)

		err := service.ResetUserPassword(resetToken.TokenValue, password)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusUnauthorized, err.StatusCode)
		assert.Equal(t, "invalid or expired token", err.Error())
	})

	t.Run("error getting active token", func(t *testing.T) {
		resetTokenRepo.EXPECT().GetActivePasswordResetToken(resetToken.TokenValue).Return(nil, errors.New("error getting token")).Times(1)

		err := service.ResetUserPassword(resetToken.TokenValue, password)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting token", err.Error())
	})

	t.Run("user not found", func(t *testing.T) {
		resetTokenRepo.EXPECT().GetActivePasswordResetToken(resetToken.TokenValue).Return(resetToken, nil).Times(1)
		userRepo.EXPECT().GetUser(resetToken.UserID).Return(nil, nil).Times(1)

		err := service.ResetUserPassword(resetToken.TokenValue, password)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "user not found", err.Error())
	})

	t.Run("error getting user", func(t *testing.T) {
		resetTokenRepo.EXPECT().GetActivePasswordResetToken(resetToken.TokenValue).Return(resetToken, nil).Times(1)
		userRepo.EXPECT().GetUser(resetToken.UserID).Return(nil, errors.New("error getting user")).Times(1)

		err := service.ResetUserPassword(resetToken.TokenValue, password)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting user", err.Error())
	})

	t.Run("error generating password", func(t *testing.T) {
		resetTokenRepo.EXPECT().GetActivePasswordResetToken(resetToken.TokenValue).Return(resetToken, nil).Times(1)
		userRepo.EXPECT().GetUser(resetToken.UserID).Return(user, nil).Times(1)
		auth.EXPECT().GenerateHashedPassword(password).Return("", errors.New("error generating password")).Times(1)

		err := service.ResetUserPassword(resetToken.TokenValue, password)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error generating password", err.Error())
	})

	t.Run("error updating password", func(t *testing.T) {
		resetTokenRepo.EXPECT().GetActivePasswordResetToken(resetToken.TokenValue).Return(resetToken, nil).Times(1)
		userRepo.EXPECT().GetUser(resetToken.UserID).Return(user, nil).Times(1)
		auth.EXPECT().GenerateHashedPassword(password).Return(hashedPassword, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		userRepo.EXPECT().UpdatePassword(gomock.Any(), user, hashedPassword).Return(nil, errors.New("error updating password")).Times(1)

		err := service.ResetUserPassword(resetToken.TokenValue, password)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error updating password", err.Error())
	})

	t.Run("error revoking user login tokens", func(t *testing.T) {
		resetTokenRepo.EXPECT().GetActivePasswordResetToken(resetToken.TokenValue).Return(resetToken, nil).Times(1)
		userRepo.EXPECT().GetUser(resetToken.UserID).Return(user, nil).Times(1)
		auth.EXPECT().GenerateHashedPassword(password).Return(hashedPassword, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		userRepo.EXPECT().UpdatePassword(gomock.Any(), user, hashedPassword).Return(user, nil).Times(1)
		loginTokenRepo.EXPECT().RevokeUserLoginTokens(gomock.Any(), user.ID).Return(errors.New("error revoking tokens")).Times(1)

		err := service.ResetUserPassword(resetToken.TokenValue, password)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error revoking tokens", err.Error())
	})

	t.Run("error using reset token", func(t *testing.T) {
		resetTokenRepo.EXPECT().GetActivePasswordResetToken(resetToken.TokenValue).Return(resetToken, nil).Times(1)
		userRepo.EXPECT().GetUser(resetToken.UserID).Return(user, nil).Times(1)
		auth.EXPECT().GenerateHashedPassword(password).Return(hashedPassword, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		userRepo.EXPECT().UpdatePassword(gomock.Any(), user, hashedPassword).Return(user, nil).Times(1)
		loginTokenRepo.EXPECT().RevokeUserLoginTokens(gomock.Any(), user.ID).Return(nil).Times(1)
		resetTokenRepo.EXPECT().UseToken(gomock.Any(), resetToken).Return(errors.New("error using token")).Times(1)

		err := service.ResetUserPassword(resetToken.TokenValue, password)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error using token", err.Error())
	})

	t.Run("error getting password reset tokens", func(t *testing.T) {
		resetTokenRepo.EXPECT().GetActivePasswordResetToken(resetToken.TokenValue).Return(resetToken, nil).Times(1)
		userRepo.EXPECT().GetUser(resetToken.UserID).Return(user, nil).Times(1)
		auth.EXPECT().GenerateHashedPassword(password).Return(hashedPassword, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		userRepo.EXPECT().UpdatePassword(gomock.Any(), user, hashedPassword).Return(user, nil).Times(1)
		loginTokenRepo.EXPECT().RevokeUserLoginTokens(gomock.Any(), user.ID).Return(nil).Times(1)
		resetTokenRepo.EXPECT().UseToken(gomock.Any(), resetToken).Return(nil).Times(1)
		resetTokenRepo.EXPECT().GetUserActivePasswordResetTokens(resetToken.UserID).Return(nil, errors.New("error getting tokens")).Times(1)

		err := service.ResetUserPassword(resetToken.TokenValue, password)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting tokens", err.Error())
	})

	t.Run("error revoking reset tokens", func(t *testing.T) {
		resetTokenRepo.EXPECT().GetActivePasswordResetToken(resetToken.TokenValue).Return(resetToken, nil).Times(1)
		userRepo.EXPECT().GetUser(resetToken.UserID).Return(user, nil).Times(1)
		auth.EXPECT().GenerateHashedPassword(password).Return(hashedPassword, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		userRepo.EXPECT().UpdatePassword(gomock.Any(), user, hashedPassword).Return(user, nil).Times(1)
		loginTokenRepo.EXPECT().RevokeUserLoginTokens(gomock.Any(), user.ID).Return(nil).Times(1)
		resetTokenRepo.EXPECT().UseToken(gomock.Any(), resetToken).Return(nil).Times(1)
		resetTokenRepo.EXPECT().GetUserActivePasswordResetTokens(resetToken.UserID).Return(allResetTokens, nil).Times(1)
		resetTokenRepo.EXPECT().RevokeTokens(gomock.Any(), gomock.Any()).Return(errors.New("error revoking tokens")).Times(1)

		err := service.ResetUserPassword(resetToken.TokenValue, password)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error revoking tokens", err.Error())
	})

	t.Run("error deleting user sessions", func(t *testing.T) {
		resetTokenRepo.EXPECT().GetActivePasswordResetToken(resetToken.TokenValue).Return(resetToken, nil).Times(1)
		userRepo.EXPECT().GetUser(resetToken.UserID).Return(user, nil).Times(1)
		auth.EXPECT().GenerateHashedPassword(password).Return(hashedPassword, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		userRepo.EXPECT().UpdatePassword(gomock.Any(), user, hashedPassword).Return(user, nil).Times(1)
		loginTokenRepo.EXPECT().RevokeUserLoginTokens(gomock.Any(), user.ID).Return(nil).Times(1)
		resetTokenRepo.EXPECT().UseToken(gomock.Any(), resetToken).Return(nil).Times(1)
		resetTokenRepo.EXPECT().GetUserActivePasswordResetTokens(resetToken.UserID).Return(allResetTokens, nil).Times(1)
		resetTokenRepo.EXPECT().RevokeTokens(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		transaction.EXPECT().ExecuteInRedisTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(rdb *redis.Client, fn func(tx *redis.Tx) error) error {
				return fn(nil)
			},
		).Times(1)
		sessionRepo.EXPECT().DeleteUserSessions(user.ID).Return(errors.New("error deleting tokens")).Times(1)

		err := service.ResetUserPassword(resetToken.TokenValue, password)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error deleting tokens", err.Error())
	})
}
