package services

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_repositories"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_transaction"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestUserProfileService_GetProfileByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserProfileRepo := mock_repositories.NewMockUserProfileRepository(ctrl)

	userProfileService := NewUserProfileService(nil, nil, mockUserProfileRepo, nil)

	userProfile := utils.GenerateSampleUserProfile()

	t.Run("success", func(t *testing.T) {
		mockUserProfileRepo.EXPECT().GetProfileByUserID(userProfile.UserID).Return(userProfile, nil).Times(1)

		result, err := userProfileService.GetProfileByUserID(userProfile.UserID)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, userProfile, result)
	})

	t.Run("user profile not found", func(t *testing.T) {
		mockUserProfileRepo.EXPECT().GetProfileByUserID(userProfile.UserID).Return(nil, gorm.ErrRecordNotFound).Times(1)

		result, err := userProfileService.GetProfileByUserID(userProfile.UserID)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "User profile does not exist", err.Message)
	})

	t.Run("internal server error", func(t *testing.T) {
		mockUserProfileRepo.EXPECT().GetProfileByUserID(userProfile.UserID).Return(nil, errors.New("db error")).Times(1)

		result, err := userProfileService.GetProfileByUserID(userProfile.UserID)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Message)
	})
}

func TestUserProfileService_CreateUserProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transactionMock := mock_transaction.NewMockTransactionManager(ctrl)

	userProfileRepo := mock_repositories.NewMockUserProfileRepository(ctrl)

	userProfileService := NewUserProfileService(nil, transactionMock, userProfileRepo, nil)

	profile := utils.GenerateSampleUserProfile()

	t.Run("success", func(t *testing.T) {
		transactionMock.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			}).Times(1)

		userProfileRepo.EXPECT().CreateUserProfile(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		userProfileRepo.EXPECT().GetProfileByUserID(profile.UserID).Return(nil, gorm.ErrRecordNotFound).Times(1)

		result, err := userProfileService.CreateUserProfile(profile.UserID, profile.FirstName, profile.LastName, profile.PhoneNumber, profile.DateOfBirth)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, profile.FirstName, result.FirstName)
		assert.Equal(t, profile.LastName, result.LastName)
		assert.Equal(t, profile.PhoneNumber, result.PhoneNumber)
		assert.Equal(t, profile.DateOfBirth, result.DateOfBirth)
	})

	t.Run("duplicate profile error", func(t *testing.T) {
		userProfileRepo.EXPECT().GetProfileByUserID(profile.UserID).Return(profile, nil).Times(1)

		result, err := userProfileService.CreateUserProfile(profile.UserID, profile.FirstName, profile.LastName, profile.PhoneNumber, profile.DateOfBirth)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "Duplicate profile for current user", err.Message)
	})

	t.Run("internal error on GetProfileByUserID", func(t *testing.T) {
		userProfileRepo.EXPECT().GetProfileByUserID(profile.UserID).Return(nil, errors.New("db error")).Times(1)

		result, err := userProfileService.CreateUserProfile(profile.UserID, profile.FirstName, profile.LastName, profile.PhoneNumber, profile.DateOfBirth)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Message)
	})

	t.Run("internal error on CreateUserProfile", func(t *testing.T) {
		transactionMock.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			}).Times(1)

		userProfileRepo.EXPECT().GetProfileByUserID(profile.UserID).Return(nil, gorm.ErrRecordNotFound).Times(1)
		userProfileRepo.EXPECT().CreateUserProfile(gomock.Any(), gomock.Any()).Return(errors.New("create error")).Times(1)

		result, err := userProfileService.CreateUserProfile(profile.UserID, profile.FirstName, profile.LastName, profile.PhoneNumber, profile.DateOfBirth)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "create error", err.Message)
	})
}

func TestUserProfileService_UpdateUserProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transactionMock := mock_transaction.NewMockTransactionManager(ctrl)

	userProfileRepo := mock_repositories.NewMockUserProfileRepository(ctrl)

	userProfileService := NewUserProfileService(nil, transactionMock, userProfileRepo, nil)

	firstName := "First"
	lastName := "Last"
	phoneNumber := "1234567890"
	dateOfBirth := "1990-01-01"

	profile := utils.GenerateSampleUserProfile()

	t.Run("success", func(t *testing.T) {
		transactionMock.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			}).Times(1)

		userProfileRepo.EXPECT().GetProfileByUserID(profile.UserID).Return(profile, nil).Times(1)
		userProfileRepo.EXPECT().UpdateUserProfile(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		result, err := userProfileService.UpdateUserProfile(profile.UserID, firstName, lastName, &phoneNumber, &dateOfBirth)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, firstName, result.FirstName)
		assert.Equal(t, lastName, result.LastName)
		assert.Equal(t, &phoneNumber, result.PhoneNumber)
		assert.Equal(t, &dateOfBirth, result.DateOfBirth)
	})

	t.Run("profile not found", func(t *testing.T) {
		userProfileRepo.EXPECT().GetProfileByUserID(profile.UserID).Return(nil, gorm.ErrRecordNotFound).Times(1)

		result, err := userProfileService.UpdateUserProfile(profile.UserID, firstName, lastName, &phoneNumber, &dateOfBirth)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "User profile does not exist", err.Message)
	})

	t.Run("internal error on GetProfileByUserID", func(t *testing.T) {
		userProfileRepo.EXPECT().GetProfileByUserID(profile.UserID).Return(nil, errors.New("db error")).Times(1)

		result, err := userProfileService.UpdateUserProfile(profile.UserID, firstName, lastName, &phoneNumber, &dateOfBirth)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Message)
	})

	t.Run("internal error on UpdateUserProfile", func(t *testing.T) {
		userProfileRepo.EXPECT().GetProfileByUserID(profile.UserID).Return(profile, nil).Times(1)
		transactionMock.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			}).Times(1)
		userProfileRepo.EXPECT().UpdateUserProfile(gomock.Any(), gomock.Any()).Return(errors.New("update error")).Times(1)

		result, err := userProfileService.UpdateUserProfile(profile.UserID, firstName, lastName, &phoneNumber, &dateOfBirth)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "update error", err.Message)
	})
}

func TestUserProfileService_UpdateProfilePicture(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bucketName := "test-bucket"
	os.Setenv("MINIO_PROFILE_PICTURE_BUCKET_NAME", bucketName)
	defer os.Unsetenv("MINIO_PROFILE_PICTURE_BUCKET_NAME")

	transactionMock := mock_transaction.NewMockTransactionManager(ctrl)
	userProfileRepo := mock_repositories.NewMockUserProfileRepository(ctrl)
	profilePictureRepo := mock_repositories.NewMockProfilePictureRepository(ctrl)

	userProfileService := NewUserProfileService(nil, transactionMock, userProfileRepo, profilePictureRepo)

	userID := uuid.New()
	objectName := fmt.Sprintf("%s/%d", userID, time.Now().Unix())
	file := utils.GenerateSampleFileHeader()

	t.Run("bucket does not exist and create succeeds", func(t *testing.T) {
		os.Setenv("MINIO_PROFILE_PICTURE_BUCKET_NAME", bucketName)
		defer os.Unsetenv("MINIO_PROFILE_PICTURE_BUCKET_NAME")

		transactionMock.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)

		profilePictureRepo.EXPECT().BucketExists(bucketName).Return(false).Times(1)
		profilePictureRepo.EXPECT().CreateBucket(bucketName).Return(nil).Times(1)
		profilePictureRepo.EXPECT().CreateProfilePicture(file, bucketName, objectName).Return(nil).Times(1)
		userProfileRepo.EXPECT().UpdateProfilePicture(gomock.Any(), userID, objectName).Return(nil).Times(1)

		err := userProfileService.UpdateProfilePicture(userID, file)
		assert.Nil(t, err)
	})

	t.Run("bucket exists", func(t *testing.T) {
		os.Setenv("MINIO_PROFILE_PICTURE_BUCKET_NAME", bucketName)
		defer os.Unsetenv("MINIO_PROFILE_PICTURE_BUCKET_NAME")

		transactionMock.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)

		profilePictureRepo.EXPECT().BucketExists(bucketName).Return(true).Times(1)
		profilePictureRepo.EXPECT().CreateProfilePicture(file, bucketName, objectName).Return(nil).Times(1)
		userProfileRepo.EXPECT().UpdateProfilePicture(gomock.Any(), userID, objectName).Return(nil).Times(1)

		err := userProfileService.UpdateProfilePicture(userID, file)
		assert.Nil(t, err)
	})

	t.Run("error creating bucket", func(t *testing.T) {
		os.Setenv("MINIO_PROFILE_PICTURE_BUCKET_NAME", bucketName)
		defer os.Unsetenv("MINIO_PROFILE_PICTURE_BUCKET_NAME")

		profilePictureRepo.EXPECT().BucketExists(bucketName).Return(false).Times(1)
		profilePictureRepo.EXPECT().CreateBucket(bucketName).Return(errors.New("create bucket error")).Times(1)

		err := userProfileService.UpdateProfilePicture(userID, file)
		assert.Error(t, err)
		assert.Equal(t, "create bucket error", err.Error())
	})

	t.Run("error uploading profile picture", func(t *testing.T) {
		os.Setenv("MINIO_PROFILE_PICTURE_BUCKET_NAME", bucketName)
		defer os.Unsetenv("MINIO_PROFILE_PICTURE_BUCKET_NAME")

		profilePictureRepo.EXPECT().BucketExists(bucketName).Return(false).Times(1)
		profilePictureRepo.EXPECT().CreateBucket(bucketName).Return(nil).Times(1)
		profilePictureRepo.EXPECT().CreateProfilePicture(file, bucketName, objectName).Return(errors.New("upload error")).Times(1)

		err := userProfileService.UpdateProfilePicture(userID, file)
		assert.Error(t, err)
		assert.Equal(t, "upload error", err.Error())
	})
}

func TestUserProfileService_DeleteProfilePicture(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transactionMock := mock_transaction.NewMockTransactionManager(ctrl)
	userProfileRepo := mock_repositories.NewMockUserProfileRepository(ctrl)
	profilePictureRepo := mock_repositories.NewMockProfilePictureRepository(ctrl)

	userID := uuid.New()
	service := NewUserProfileService(nil, transactionMock, userProfileRepo, profilePictureRepo)

	t.Run("successful deletion", func(t *testing.T) {
		transactionMock.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)

		userProfileRepo.EXPECT().DeleteProfilePicture(gomock.Any(), userID).Return(nil).Times(1)

		err := service.DeleteProfilePicture(userID)

		assert.Nil(t, err)
	})

	t.Run("error during transaction", func(t *testing.T) {
		transactionMock.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).Return(errors.New("transaction error")).Times(1)

		err := service.DeleteProfilePicture(userID)

		assert.Error(t, err)
		assert.Equal(t, "transaction error", err.Error())
	})

	t.Run("error deleting profile picture", func(t *testing.T) {
		transactionMock.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(tx *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(tx)
			},
		).Times(1)

		userProfileRepo.EXPECT().DeleteProfilePicture(gomock.Any(), userID).Return(errors.New("delete error")).Times(1)

		err := service.DeleteProfilePicture(userID)
		assert.Error(t, err)
		assert.Equal(t, "delete error", err.Error())
	})
}
