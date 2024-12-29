package services

import (
	"errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/constants"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/filters"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_transaction"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/payloads"
	"gorm.io/gorm"
	"mime/multipart"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_repositories"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"go.uber.org/mock/gomock"
)

func TestUserProfileService_GetProfileByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock_repositories.NewMockUserRepository(ctrl)
	service := NewUserProfileService(nil, nil, userRepo, nil, nil)

	user, profile := setupUserWithProfile()

	filter := filters.UserFilter{
		Filter: &filters.SingleFilter{},
		ID:     &filters.Condition{Operator: filters.OpEqual, Value: user.ID},
	}

	t.Run("success", func(t *testing.T) {
		userRepo.EXPECT().GetUser(filter, true).Return(user, nil).Times(1)

		result, err := service.GetProfileByUserID(profile.UserID)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, profile, result)
	})

	t.Run("user not found", func(t *testing.T) {
		userRepo.EXPECT().GetUser(filter, true).Return(nil, nil).Times(1)

		result, err := service.GetProfileByUserID(profile.UserID)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "user does not exist", err.Error())
	})

	t.Run("profile not found", func(t *testing.T) {
		userRepo.EXPECT().GetUser(filter, true).Return(&models.User{}, nil).Times(1)

		result, err := service.GetProfileByUserID(profile.UserID)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "user profile does not exist", err.Error())
	})

	t.Run("error getting user", func(t *testing.T) {
		userRepo.EXPECT().GetUser(filter, true).Return(nil, errors.New("error getting user")).Times(1)

		result, err := service.GetProfileByUserID(profile.UserID)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting user", err.Error())
	})
}

func TestUserProfileService_UpdateUserProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	userRepo := mock_repositories.NewMockUserRepository(ctrl)
	profileRepo := mock_repositories.NewMockUserProfileRepository(ctrl)
	service := NewUserProfileService(nil, transaction, userRepo, profileRepo, nil)

	user, profile := setupUserWithProfile()
	req := payloads.UpdateUserProfileRequest{
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: utils.GetPointerOf("0000000000"),
		DateOfBirth: utils.GetPointerOf("1970-01-01"),
	}
	filter := filters.UserFilter{
		Filter: &filters.SingleFilter{},
		ID:     &filters.Condition{Operator: filters.OpEqual, Value: user.ID},
	}

	t.Run("success", func(t *testing.T) {
		userRepo.EXPECT().GetUser(filter, true).Return(user, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			}).Times(1)
		profileRepo.EXPECT().UpdateUserProfile(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		result, err := service.UpdateUserProfile(profile.UserID, req)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, profile.ID, result.ID)
		assert.Equal(t, profile.UserID, result.UserID)
		assert.Equal(t, req.FirstName, result.FirstName)
		assert.Equal(t, req.LastName, result.LastName)
		assert.Equal(t, req.PhoneNumber, result.PhoneNumber)
		assert.Equal(t, req.DateOfBirth, result.DateOfBirth)
	})

	t.Run("user not found", func(t *testing.T) {
		userRepo.EXPECT().GetUser(filter, true).Return(nil, nil).Times(1)

		result, err := service.UpdateUserProfile(profile.UserID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "user does not exist", err.Error())
	})

	t.Run("error getting user", func(t *testing.T) {
		userRepo.EXPECT().GetUser(filter, true).Return(nil, errors.New("error getting user")).Times(1)

		result, err := service.UpdateUserProfile(profile.UserID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting user", err.Error())
	})

	t.Run("profile not found", func(t *testing.T) {
		userRepo.EXPECT().GetUser(filter, true).Return(&models.User{}, nil).Times(1)

		result, err := service.UpdateUserProfile(profile.UserID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "user profile does not exist", err.Error())
	})

	t.Run("error updating profile", func(t *testing.T) {
		userRepo.EXPECT().GetUser(filter, true).Return(user, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			}).Times(1)
		profileRepo.EXPECT().UpdateUserProfile(gomock.Any(), gomock.Any()).Return(errors.New("error updating profile")).Times(1)

		result, err := service.UpdateUserProfile(profile.UserID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error updating profile", err.Error())
	})
}

func TestUserProfileService_UpdateProfilePicture(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bucketName := "test-bucket"
	os.Setenv("MINIO_PROFILE_PICTURE_BUCKET_NAME", bucketName)
	defer os.Unsetenv("MINIO_PROFILE_PICTURE_BUCKET_NAME")

	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	userRepo := mock_repositories.NewMockUserRepository(ctrl)
	profileRepo := mock_repositories.NewMockUserProfileRepository(ctrl)
	profilePictureRepo := mock_repositories.NewMockProfilePictureRepository(ctrl)
	service := NewUserProfileService(nil, transaction, userRepo, profileRepo, profilePictureRepo)

	user, profile := setupUserWithProfile()
	file := &multipart.FileHeader{
		Filename: "file name",
		Size:     100,
		Header:   map[string][]string{constants.ContentType: {constants.ImagePng}},
	}
	filter := filters.UserFilter{
		Filter: &filters.SingleFilter{},
		ID:     &filters.Condition{Operator: filters.OpEqual, Value: user.ID},
	}

	t.Run("success", func(t *testing.T) {
		userRepo.EXPECT().GetUser(filter, true).Return(user, nil).Times(1)
		profilePictureRepo.EXPECT().CreateProfilePicture(file, bucketName, gomock.Any()).Return(nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		profileRepo.EXPECT().UpdateProfilePicture(gomock.Any(), profile, gomock.Any()).Return(profile, nil).Times(1)

		result, err := service.UpdateProfilePicture(profile.UserID, file)

		assert.NotNil(t, result)
		assert.Nil(t, err)
	})

	t.Run("user not found", func(t *testing.T) {
		userRepo.EXPECT().GetUser(filter, true).Return(nil, nil).Times(1)

		result, err := service.UpdateProfilePicture(profile.UserID, file)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "user does not exist", err.Error())
	})

	t.Run("error getting user", func(t *testing.T) {
		userRepo.EXPECT().GetUser(filter, true).Return(nil, errors.New("error getting user")).Times(1)

		result, err := service.UpdateProfilePicture(profile.UserID, file)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting user", err.Error())
	})

	t.Run("profile not found", func(t *testing.T) {
		userRepo.EXPECT().GetUser(filter, true).Return(&models.User{}, nil).Times(1)

		result, err := service.UpdateProfilePicture(profile.UserID, file)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "user profile does not exist", err.Error())
	})

	t.Run("error creating picture", func(t *testing.T) {
		userRepo.EXPECT().GetUser(filter, true).Return(user, nil).Times(1)
		profilePictureRepo.EXPECT().CreateProfilePicture(file, bucketName, gomock.Any()).Return(errors.New("error creating picture")).Times(1)

		result, err := service.UpdateProfilePicture(profile.UserID, file)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error creating picture", err.Error())
	})

	t.Run("error updating profile", func(t *testing.T) {
		userRepo.EXPECT().GetUser(filter, true).Return(user, nil).Times(1)
		profilePictureRepo.EXPECT().CreateProfilePicture(file, bucketName, gomock.Any()).Return(nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		profileRepo.EXPECT().UpdateProfilePicture(gomock.Any(), profile, gomock.Any()).Return(nil, errors.New("error updating profile")).Times(1)

		result, err := service.UpdateProfilePicture(profile.UserID, file)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error updating profile", err.Error())
	})
}

func TestUserProfileService_DeleteProfilePicture(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	userRepo := mock_repositories.NewMockUserRepository(ctrl)
	profileRepo := mock_repositories.NewMockUserProfileRepository(ctrl)
	service := NewUserProfileService(nil, transaction, userRepo, profileRepo, nil)

	user, profile := setupUserWithProfile()
	filter := filters.UserFilter{
		Filter: &filters.SingleFilter{},
		ID:     &filters.Condition{Operator: filters.OpEqual, Value: user.ID},
	}

	t.Run("success", func(t *testing.T) {
		userRepo.EXPECT().GetUser(filter, true).Return(user, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		profileRepo.EXPECT().UpdateProfilePicture(gomock.Any(), profile, nil).Return(profile, nil).Times(1)

		err := service.DeleteProfilePicture(profile.UserID)

		assert.Nil(t, err)
	})

	t.Run("user not found", func(t *testing.T) {
		userRepo.EXPECT().GetUser(filter, true).Return(nil, nil).Times(1)

		err := service.DeleteProfilePicture(profile.UserID)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "user does not exist", err.Error())
	})

	t.Run("error getting user", func(t *testing.T) {
		userRepo.EXPECT().GetUser(filter, true).Return(nil, errors.New("error getting user")).Times(1)

		err := service.DeleteProfilePicture(profile.UserID)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting user", err.Error())
	})

	t.Run("profile not found", func(t *testing.T) {
		userRepo.EXPECT().GetUser(filter, true).Return(&models.User{}, nil).Times(1)

		err := service.DeleteProfilePicture(profile.UserID)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "user profile does not exist", err.Error())
	})

	t.Run("error deleting profile picture", func(t *testing.T) {
		userRepo.EXPECT().GetUser(filter, true).Return(user, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		profileRepo.EXPECT().UpdateProfilePicture(gomock.Any(), profile, nil).Return(nil, errors.New("error deleting profile picture")).Times(1)

		err := service.DeleteProfilePicture(profile.UserID)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error deleting profile picture", err.Error())
	})
}

func setupUserWithProfile() (*models.User, *models.UserProfile) {
	user := utils.GenerateUser()
	profile := utils.GenerateUserProfile()

	profile.UserID = user.ID
	user.Profile = profile

	return user, profile
}
