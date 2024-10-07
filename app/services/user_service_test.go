package services

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_auth"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_repositories"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_transaction"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestUserService_GetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repositories.NewMockUserRepository(ctrl)
	service := NewUserService(nil, nil, nil, nil, repo, nil, nil)

	user := utils.GenerateRandomUser()

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().GetUser(user.ID).Return(user, nil).Times(1)

		result, err := service.GetUser(user.ID)

		assert.NotNil(t, user)
		assert.Nil(t, err)
		assert.Equal(t, user, result)
	})

	t.Run("user not found", func(t *testing.T) {
		repo.EXPECT().GetUser(user.ID).Return(nil, gorm.ErrRecordNotFound).Times(1)

		result, err := service.GetUser(user.ID)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "User does not exist", err.Message)
	})

	t.Run("db error", func(t *testing.T) {
		repo.EXPECT().GetUser(user.ID).Return(nil, errors.New("db error")).Times(1)

		result, err := service.GetUser(user.ID)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Message)
	})
}

func TestUserService_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	auth := mock_auth.NewMockAuthenticator(ctrl)
	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	repo := mock_repositories.NewMockUserRepository(ctrl)
	service := NewUserService(nil, nil, auth, transaction, repo, nil, nil)

	user := utils.GenerateRandomUser()
	password := "password"

	t.Run("success", func(t *testing.T) {
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)

		repo.EXPECT().FindUserByEmail(user.Email).Return(nil, gorm.ErrRecordNotFound).Times(1)
		auth.EXPECT().GenerateHashedPassword(password).Return(user.PasswordHash, nil).Times(1)
		repo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		result, err := service.CreateUser(user.Email, password)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, user.Email, result.Email)
	})

	t.Run("duplicate email error", func(t *testing.T) {
		repo.EXPECT().FindUserByEmail(user.Email).Return(user, nil).Times(1)

		result, err := service.CreateUser(user.Email, password)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, fmt.Sprintf("Email %s already exists", user.Email), err.Error())
	})

	t.Run("error getting user", func(t *testing.T) {
		repo.EXPECT().FindUserByEmail(user.Email).Return(nil, errors.New("db error")).Times(1)

		result, err := service.CreateUser(user.Email, password)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Message)
	})

	t.Run("error generating password hash", func(t *testing.T) {
		auth.EXPECT().GenerateHashedPassword(password).Return("", errors.New("hash error"))
		repo.EXPECT().FindUserByEmail(user.Email).Return(nil, gorm.ErrRecordNotFound).Times(1)

		result, err := service.CreateUser(user.Email, password)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "hash error", err.Message)
	})

	t.Run("error creating user", func(t *testing.T) {
		repo.EXPECT().FindUserByEmail(user.Email).Return(nil, gorm.ErrRecordNotFound).Times(1)
		auth.EXPECT().GenerateHashedPassword(password).Return(user.PasswordHash, nil)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		repo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(errors.New("db error")).Times(1)

		result, err := service.CreateUser(user.Email, password)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Message)
	})
}

func TestUserService_LoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	os.Setenv("AUTH_TOKEN_EXPIRES_AFTER_MINUTES", "60")
	defer os.Unsetenv("AUTH_TOKEN_EXPIRES_AFTER_MINUTES")

	auth := mock_auth.NewMockAuthenticator(ctrl)
	transaction := mock_transaction.NewMockTransactionManager(ctrl)

	userRepo := mock_repositories.NewMockUserRepository(ctrl)
	loginTokenRepo := mock_repositories.NewMockLoginTokenRepository(ctrl)
	userSessionRepo := mock_repositories.NewMockUserSessionRepository(ctrl)

	service := NewUserService(nil, nil, auth, transaction, userRepo, loginTokenRepo, userSessionRepo)

	user := utils.GenerateRandomUser()
	token := uuid.NewString()
	password := "password"

	t.Run("success", func(t *testing.T) {
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		transaction.EXPECT().ExecuteInRedisTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(rdb *redis.Client, fn func(tx *redis.Tx) error) error {
				return fn(nil)
			},
		).Times(1)

		auth.EXPECT().GenerateToken().Return(token).Times(1)
		auth.EXPECT().DoPasswordsMatch(user.PasswordHash, password).Return(true).Times(1)
		userRepo.EXPECT().FindUserByEmail(user.Email).Return(user, nil).Times(1)
		userSessionRepo.EXPECT().GetUserSessionID(gomock.Any()).Return(token).Times(1)
		loginTokenRepo.EXPECT().GetActiveLoginToken(gomock.Any()).Return(nil, gorm.ErrRecordNotFound).Times(1)
		loginTokenRepo.EXPECT().CreateLoginToken(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		userSessionRepo.EXPECT().CreateUserSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

		result, err := service.LoginUser(user.Email, password)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, user.ID, result.UserID)
		assert.Equal(t, token, result.TokenValue)
	})

	t.Run("invalid email", func(t *testing.T) {
		userRepo.EXPECT().FindUserByEmail(user.Email).Return(nil, gorm.ErrRecordNotFound).Times(1)

		result, err := service.LoginUser(user.Email, password)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, fmt.Sprintf("Invalid email %s", user.Email), err.Message)
	})

	t.Run("invalid password", func(t *testing.T) {
		auth.EXPECT().DoPasswordsMatch(user.PasswordHash, password).Return(false).Times(1)
		userRepo.EXPECT().FindUserByEmail(user.Email).Return(user, nil).Times(1)

		result, err := service.LoginUser(user.Email, password)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "Invalid password", err.Message)
	})

	t.Run("token already exists", func(t *testing.T) {
		auth.EXPECT().DoPasswordsMatch(user.PasswordHash, password).Return(true).Times(1)
		userRepo.EXPECT().FindUserByEmail(user.Email).Return(user, nil).Times(1)
		auth.EXPECT().GenerateToken().Return(token).Times(1)
		loginTokenRepo.EXPECT().GetActiveLoginToken(token).Return(&models.LoginToken{}, nil).Times(1)

		result, err := service.LoginUser(user.Email, password)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "Token value already exists", err.Message)
	})

	t.Run("error creating login token", func(t *testing.T) {
		auth.EXPECT().DoPasswordsMatch(user.PasswordHash, password).Return(true).Times(1)
		auth.EXPECT().GenerateToken().Return(token).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)

		userRepo.EXPECT().FindUserByEmail(user.Email).Return(user, nil).Times(1)
		loginTokenRepo.EXPECT().GetActiveLoginToken(token).Return(nil, gorm.ErrRecordNotFound).Times(1)
		loginTokenRepo.EXPECT().CreateLoginToken(gomock.Any(), gomock.Any()).Return(errors.New("create token error")).Times(1)

		result, err := service.LoginUser(user.Email, password)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "create token error", err.Message)
	})

	t.Run("error creating user session", func(t *testing.T) {
		userRepo.EXPECT().FindUserByEmail(user.Email).Return(user, nil).Times(1)
		auth.EXPECT().DoPasswordsMatch(user.PasswordHash, password).Return(true).Times(1)
		auth.EXPECT().GenerateToken().Return(token).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		transaction.EXPECT().ExecuteInRedisTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(rdb *redis.Client, fn func(tx *redis.Tx) error) error {
				return fn(nil)
			},
		).Times(1)

		loginTokenRepo.EXPECT().GetActiveLoginToken(token).Return(nil, gorm.ErrRecordNotFound).Times(1)
		loginTokenRepo.EXPECT().CreateLoginToken(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		userSessionRepo.EXPECT().GetUserSessionID(gomock.Any()).Return(token).Times(1)
		userSessionRepo.EXPECT().CreateUserSession(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("session creation error")).Times(1)

		result, err := service.LoginUser(user.Email, password)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "session creation error", err.Message)
	})
}

func TestUserService_LogoutUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transactionManager := mock_transaction.NewMockTransactionManager(ctrl)
	userSessionRepo := mock_repositories.NewMockUserSessionRepository(ctrl)
	loginTokenRepo := mock_repositories.NewMockLoginTokenRepository(ctrl)

	userService := NewUserService(nil, nil, nil, transactionManager, nil, loginTokenRepo, userSessionRepo)

	token := utils.GenerateRandomLoginToken()

	t.Run("success", func(t *testing.T) {
		transactionManager.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		transactionManager.EXPECT().ExecuteInRedisTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(rdb *redis.Client, fn func(tx *redis.Tx) error) error {
				return fn(nil)
			},
		).Times(1)

		loginTokenRepo.EXPECT().GetActiveLoginToken(token.TokenValue).Return(token, nil).Times(1)
		loginTokenRepo.EXPECT().RevokeLoginToken(gomock.Any(), token).Return(nil).Times(1)
		userSessionRepo.EXPECT().GetUserSessionID(gomock.Any()).Return(token.TokenValue).Times(1)
		userSessionRepo.EXPECT().DeleteUserSession(gomock.Any()).Return(nil).Times(1)

		err := userService.LogoutUser(token.TokenValue)

		assert.Nil(t, err)
	})

	t.Run("token not found", func(t *testing.T) {
		loginTokenRepo.EXPECT().GetActiveLoginToken(token.TokenValue).Return(nil, gorm.ErrRecordNotFound).Times(1)

		err := userService.LogoutUser(token.TokenValue)

		assert.NotNil(t, err)
		assert.Equal(t, "Token not found", err.Message)
	})

	t.Run("internal server error on GetActiveLoginToken", func(t *testing.T) {
		loginTokenRepo.EXPECT().GetActiveLoginToken(token.TokenValue).Return(nil, errors.New("some error")).Times(1)

		err := userService.LogoutUser(token.TokenValue)

		assert.NotNil(t, err)
		assert.Equal(t, "some error", err.Message)
	})

	t.Run("internal server error on ExecuteInTransaction", func(t *testing.T) {
		transactionManager.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).Return(errors.New("db transaction error")).Times(1)

		loginTokenRepo.EXPECT().GetActiveLoginToken(token.TokenValue).Return(token, nil).Times(1)

		err := userService.LogoutUser(token.TokenValue)

		assert.NotNil(t, err)
		assert.Equal(t, "db transaction error", err.Message)
	})

	t.Run("internal server error on ExecuteInRedisTransaction", func(t *testing.T) {
		transactionManager.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)

		loginTokenRepo.EXPECT().GetActiveLoginToken(token.TokenValue).Return(token, nil).Times(1)
		loginTokenRepo.EXPECT().RevokeLoginToken(gomock.Any(), token).Return(nil).Times(1)

		transactionManager.EXPECT().ExecuteInRedisTransaction(gomock.Any(), gomock.Any()).Return(errors.New("redis transaction error")).Times(1)

		err := userService.LogoutUser(token.TokenValue)

		assert.NotNil(t, err)
		assert.Equal(t, "redis transaction error", err.Message)
	})
}

func TestUserService_UpdateUserPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authMock := mock_auth.NewMockAuthenticator(ctrl)
	transactionManager := mock_transaction.NewMockTransactionManager(ctrl)

	userRepo := mock_repositories.NewMockUserRepository(ctrl)
	loginTokenRepo := mock_repositories.NewMockLoginTokenRepository(ctrl)
	userSessionRepo := mock_repositories.NewMockUserSessionRepository(ctrl)

	userService := NewUserService(nil, nil, authMock, transactionManager, userRepo, loginTokenRepo, userSessionRepo)

	user := utils.GenerateRandomUser()
	password := "password"

	t.Run("success", func(t *testing.T) {
		userRepo.EXPECT().GetUser(user.ID).Return(user, nil).Times(1)

		authMock.EXPECT().DoPasswordsMatch(user.PasswordHash, password).Return(false).Times(1)
		authMock.EXPECT().GenerateHashedPassword(password).Return(user.PasswordHash, nil).Times(1)

		transactionManager.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		transactionManager.EXPECT().ExecuteInRedisTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(rdb *redis.Client, fn func(tx *redis.Tx) error) error {
				return fn(nil)
			},
		).Times(1)

		userRepo.EXPECT().UpdatePassword(gomock.Any(), user, user.PasswordHash).Return(user, nil).Times(1)
		loginTokenRepo.EXPECT().RevokeUserLoginTokens(gomock.Any(), user.ID).Return(nil).Times(1)

		userSessionRepo.EXPECT().DeleteUserSessions(user.ID).Return(nil).Times(1)

		err := userService.UpdateUserPassword(user.ID, password)

		assert.Nil(t, err)
	})

	t.Run("user not found", func(t *testing.T) {
		userRepo.EXPECT().GetUser(user.ID).Return(nil, gorm.ErrRecordNotFound).Times(1)

		err := userService.UpdateUserPassword(user.ID, password)

		assert.NotNil(t, err)
		assert.Equal(t, "User does not exist", err.Message)
	})

	t.Run("error getting user", func(t *testing.T) {

		userRepo.EXPECT().GetUser(user.ID).Return(nil, errors.New("db error")).Times(1)
		err := userService.UpdateUserPassword(user.ID, password)

		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Message)
	})

	t.Run("error password same as current value", func(t *testing.T) {
		authMock.EXPECT().DoPasswordsMatch(user.PasswordHash, password).Return(true).Times(1)

		userRepo.EXPECT().GetUser(user.ID).Return(user, nil).Times(1)

		err := userService.UpdateUserPassword(user.ID, password)

		assert.NotNil(t, err)
		assert.Equal(t, "New password can not be the same as current value", err.Message)
	})

	t.Run("error generating hashed password", func(t *testing.T) {
		authMock.EXPECT().DoPasswordsMatch(user.PasswordHash, password).Return(false).Times(1)
		authMock.EXPECT().GenerateHashedPassword(password).Return("", errors.New("hashing error")).Times(1)

		userRepo.EXPECT().GetUser(user.ID).Return(user, nil).Times(1)

		err := userService.UpdateUserPassword(user.ID, password)

		assert.NotNil(t, err)
		assert.Equal(t, "hashing error", err.Message)
	})
}
