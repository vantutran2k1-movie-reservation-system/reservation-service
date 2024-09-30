package services

// import (
// 	"errors"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_repositories"
// 	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
// 	"go.uber.org/mock/gomock"
// 	"gorm.io/gorm"
// )

// func TestUserService_GetUser(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	userRepo := mo.NewMockUserRepository(ctrl)
// 	userService := NewUserService(nil, nil, nil, userRepo, nil, nil)

// 	user := utils.GenerateRandomUser()
// 	userID := user.ID

// 	t.Run("success", func(t *testing.T) {
// 		userRepo.EXPECT().GetUser(userID).Return(user, nil).Times(1)

// 		result, err := userService.GetUser(userID)

// 		assert.Nil(t, err)
// 		assert.NotNil(t, user)
// 		assert.Equal(t, user, result)
// 	})

// 	t.Run("not found", func(t *testing.T) {
// 		userRepo.EXPECT().GetUser(userID).Return(nil, gorm.ErrRecordNotFound).Times(1)

// 		result, err := userService.GetUser(userID)

// 		assert.Nil(t, result)
// 		assert.NotNil(t, err)
// 		assert.Equal(t, "User does not exist", err.Message)
// 	})

// 	t.Run("internal error", func(t *testing.T) {
// 		userRepo.EXPECT().GetUser(userID).Return(nil, errors.New("some error")).Times(1)

// 		result, err := userService.GetUser(userID)

// 		assert.Nil(t, result)
// 		assert.NotNil(t, err)
// 		assert.Equal(t, "some error", err.Message)
// 	})
// }

// func TestUserService_CreateUser(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	userRepo := mocks.NewMockUserRepository(ctrl)
// 	authMock := mocks.Ne
// 	userService := NewUserService(nil, nil, nil, userRepo, nil, nil)

// 	user := utils.GenerateRandomUser()
// 	password := "password"

// 	t.Run("success", func(t *testing.T) {
// 		userRepo.EXPECT().FindUserByEmail(user.Email).Return(nil, gorm.ErrRecordNotFound).Times(1)
// 		auth.GenerateHashedPassword = func(pw string) (string, error) { return user.PasswordHash, nil }
// 		userRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil).Times(1)

// 		result, err := userService.CreateUser(user.Email, password)

// 		assert.Nil(t, err)
// 		assert.NotNil(t, result)
// 		assert.Equal(t, user.Email, result.Email)
// 	})

// 	t.Run("email already exists", func(t *testing.T) {
// 		userRepo.EXPECT().FindUserByEmail(user.Email).Return(user, nil).Times(1)

// 		result, err := userService.CreateUser(user.Email, password)

// 		assert.Nil(t, result)
// 		assert.NotNil(t, err)
// 		assert.Equal(t, fmt.Sprintf("Email %s already exists", user.Email), err.Message)
// 	})

// t.Run("internal error finding user", func(t *testing.T) {
// 	userRepo.EXPECT().FindUserByEmail(email).Return(nil, errors.New("database error")).Times(1)

// 	result, err := userService.CreateUser(email, password)

// 	assert.Nil(t, result)
// 	assert.NotNil(t, err)
// 	assert.Equal(t, "database error", err.Message)
// })

// t.Run("error generating password hash", func(t *testing.T) {
// 	userRepo.EXPECT().FindUserByEmail(email).Return(nil, gorm.ErrRecordNotFound).Times(1)
// 	auth.GenerateHashedPassword = func(pw string) (string, error) { return "", errors.New("hash error") }

// 	result, err := userService.CreateUser(email, password)

// 	assert.Nil(t, result)
// 	assert.NotNil(t, err)
// 	assert.Equal(t, "hash error", err.Message)
// })

// t.Run("error creating user in transaction", func(t *testing.T) {
// 	userRepo.EXPECT().FindUserByEmail(email).Return(nil, gorm.ErrRecordNotFound).Times(1)
// 	auth.GenerateHashedPassword = func(pw string) (string, error) { return hashedPassword, nil }
// 	userRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(errors.New("create error")).Times(1)

// 	result, err := userService.CreateUser(email, password)

// 	assert.Nil(t, result)
// 	assert.NotNil(t, err)
// 	assert.Equal(t, "create error", err.Message)
// })
// }
