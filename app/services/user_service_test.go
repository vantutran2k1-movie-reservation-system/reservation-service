package services

import (
	"errors"
	"fmt"
	"testing"

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

	userRepo := mock_repositories.NewMockUserRepository(ctrl)
	userService := NewUserService(nil, nil, nil, nil, nil, userRepo, nil, nil)

	user := utils.GenerateRandomUser()

	t.Run("success", func(t *testing.T) {
		userRepo.EXPECT().GetUser(user.ID).Return(user, nil).Times(1)

		result, err := userService.GetUser(user.ID)

		assert.Nil(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, user, result)
	})

	t.Run("not found", func(t *testing.T) {
		userRepo.EXPECT().GetUser(user.ID).Return(nil, gorm.ErrRecordNotFound).Times(1)

		result, err := userService.GetUser(user.ID)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "User does not exist", err.Message)
	})

	t.Run("internal error", func(t *testing.T) {
		userRepo.EXPECT().GetUser(user.ID).Return(nil, errors.New("some error")).Times(1)

		result, err := userService.GetUser(user.ID)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "some error", err.Message)
	})
}

func TestUserService_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock_repositories.NewMockUserRepository(ctrl)
	authMock := mock_auth.NewMockAuthenticator(ctrl)
	transactionMock := mock_transaction.NewMockTransactionManager(ctrl)
	userService := NewUserService(nil, nil, authMock, nil, transactionMock, userRepo, nil, nil)

	user := utils.GenerateRandomUser()
	password := "password"

	t.Run("success", func(t *testing.T) {
		userRepo.EXPECT().FindUserByEmail(user.Email).Return(nil, gorm.ErrRecordNotFound).Times(1)
		authMock.EXPECT().GenerateHashedPassword(password).Return(user.PasswordHash, nil)
		transactionMock.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		userRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		result, err := userService.CreateUser(user.Email, password)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, user.Email, result.Email)
	})

	t.Run("email already exists", func(t *testing.T) {
		userRepo.EXPECT().FindUserByEmail(user.Email).Return(user, nil).Times(1)

		result, err := userService.CreateUser(user.Email, password)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, fmt.Sprintf("Email %s already exists", user.Email), err.Message)
	})

	t.Run("internal error finding user", func(t *testing.T) {
		userRepo.EXPECT().FindUserByEmail(user.Email).Return(nil, errors.New("database error")).Times(1)

		result, err := userService.CreateUser(user.Email, password)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "database error", err.Message)
	})

	t.Run("error generating password hash", func(t *testing.T) {
		userRepo.EXPECT().FindUserByEmail(user.Email).Return(nil, gorm.ErrRecordNotFound).Times(1)
		authMock.EXPECT().GenerateHashedPassword(password).Return("", errors.New("hash error"))

		result, err := userService.CreateUser(user.Email, password)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "hash error", err.Message)
	})

	t.Run("error creating user in transaction", func(t *testing.T) {
		userRepo.EXPECT().FindUserByEmail(user.Email).Return(nil, gorm.ErrRecordNotFound).Times(1)
		authMock.EXPECT().GenerateHashedPassword(password).Return(user.PasswordHash, nil)
		transactionMock.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		userRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(errors.New("create error")).Times(1)

		result, err := userService.CreateUser(user.Email, password)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "create error", err.Message)
	})
}

// func TestUserService_LoginUser(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	os.Setenv("AUTH_TOKEN_EXPIRES_AFTER_MINUTES", "60")
// 	defer os.Unsetenv("AUTH_TOKEN_EXPIRES_AFTER_MINUTES")

// 	userRepo := mock_repositories.NewMockUserRepository(ctrl)
// 	loginTokenRepo := mock_repositories.NewMockLoginTokenRepository(ctrl)
// 	userSessionRepo := mock_repositories.NewMockUserSessionRepository(ctrl)
// 	authMock := mock_auth.NewMockAuthenticator(ctrl)
// 	tokenGeneratorMock := mock_auth.NewMockTokenGenerator(ctrl)

// 	userService := NewUserService(nil, nil, authMock, tokenGeneratorMock, userRepo, loginTokenRepo, userSessionRepo)

// 	user := utils.GenerateRandomUser()
// 	token := utils.GenerateRandomAuthToken()

// 	password := "password"

// 	t.Run("successful login", func(t *testing.T) {
// 		userRepo.EXPECT().FindUserByEmail(user.Email).Return(user, nil).Times(1)
// 		authMock.EXPECT().IsPasswordsMatch(user.PasswordHash, password).Return(true).Times(1)
// 		tokenGeneratorMock.EXPECT().GenerateToken().Return(token, nil).Times(1)
// 		loginTokenRepo.EXPECT().GetActiveLoginToken(token.TokenValue).Return(nil, gorm.ErrRecordNotFound).Times(1)
// 		loginTokenRepo.EXPECT().CreateLoginToken(gomock.Any(), gomock.Any()).Return(nil).Times(1)
// 		userSessionRepo.EXPECT().CreateUserSession(gomock.Any(), token.ValidDuration, gomock.Any()).Return(nil).Times(1)

// 		result, err := userService.LoginUser(user.Email, password)

// 		assert.Nil(t, err)
// 		assert.NotNil(t, result)
// 		assert.Equal(t, user.ID, result.UserID)
// 		assert.Equal(t, token.TokenValue, result.TokenValue)
// 	})

// 	t.Run("invalid email", func(t *testing.T) {
// 		userRepo.EXPECT().FindUserByEmail(email).Return(nil, gorm.ErrRecordNotFound).Times(1)

// 		result, err := userService.LoginUser(email, password)

// 		assert.Nil(t, result)
// 		assert.NotNil(t, err)
// 		assert.Equal(t, "Invalid email test@example.com", err.Message)
// 	})

// 	t.Run("invalid password", func(t *testing.T) {
// 		userRepo.EXPECT().FindUserByEmail(email).Return(user, nil).Times(1)
// 		authMock.EXPECT().IsPasswordsMatch(hashedPassword, invalidPassword).Return(false).Times(1)

// 		result, err := userService.LoginUser(email, invalidPassword)

// 		assert.Nil(t, result)
// 		assert.NotNil(t, err)
// 		assert.Equal(t, "Invalid password", err.Message)
// 	})

// 	t.Run("token generation failure", func(t *testing.T) {
// 		userRepo.EXPECT().FindUserByEmail(email).Return(user, nil).Times(1)
// 		authMock.EXPECT().IsPasswordsMatch(hashedPassword, password).Return(true).Times(1)
// 		authMock.EXPECT().GenerateToken().Return(nil, errors.New("token generation error")).Times(1)

// 		result, err := userService.LoginUser(email, password)

// 		assert.Nil(t, result)
// 		assert.NotNil(t, err)
// 		assert.Equal(t, "token generation error", err.Message)
// 	})

// 	t.Run("token already exists", func(t *testing.T) {
// 		userRepo.EXPECT().FindUserByEmail(email).Return(user, nil).Times(1)
// 		authMock.EXPECT().IsPasswordsMatch(hashedPassword, password).Return(true).Times(1)
// 		authMock.EXPECT().GenerateToken().Return(token, nil).Times(1)
// 		loginTokenRepo.EXPECT().GetActiveLoginToken(tokenValue).Return(&models.LoginToken{}, nil).Times(1)

// 		result, err := userService.LoginUser(email, password)

// 		assert.Nil(t, result)
// 		assert.NotNil(t, err)
// 		assert.Equal(t, "Token value already exists", err.Message)
// 	})

// 	t.Run("error creating login token", func(t *testing.T) {
// 		userRepo.EXPECT().FindUserByEmail(email).Return(user, nil).Times(1)
// 		authMock.EXPECT().IsPasswordsMatch(hashedPassword, password).Return(true).Times(1)
// 		authMock.EXPECT().GenerateToken().Return(token, nil).Times(1)
// 		loginTokenRepo.EXPECT().GetActiveLoginToken(tokenValue).Return(nil, gorm.ErrRecordNotFound).Times(1)
// 		loginTokenRepo.EXPECT().CreateLoginToken(gomock.Any(), gomock.Any()).Return(errors.New("create token error")).Times(1)

// 		result, err := userService.LoginUser(email, password)

// 		assert.Nil(t, result)
// 		assert.NotNil(t, err)
// 		assert.Equal(t, "create token error", err.Message)
// 	})

// 	t.Run("error creating user session", func(t *testing.T) {
// 		userRepo.EXPECT().FindUserByEmail(email).Return(user, nil).Times(1)
// 		authMock.EXPECT().IsPasswordsMatch(hashedPassword, password).Return(true).Times(1)
// 		authMock.EXPECT().GenerateToken().Return(token, nil).Times(1)
// 		loginTokenRepo.EXPECT().GetActiveLoginToken(tokenValue).Return(nil, gorm.ErrRecordNotFound).Times(1)
// 		loginTokenRepo.EXPECT().CreateLoginToken(gomock.Any(), gomock.Any()).Return(nil).Times(1)
// 		userSessionRepo.EXPECT().CreateUserSession(gomock.Any(), validDuration, gomock.Any()).Return(errors.New("session creation error")).Times(1)

// 		result, err := userService.LoginUser(email, password)

// 		assert.Nil(t, result)
// 		assert.NotNil(t, err)
// 		assert.Equal(t, "session creation error", err.Message)
// 	})
// }
