package services

import (
	"errors"
	"net/http"
	"os"
	"testing"

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

	repo := mock_repositories.NewMockUserProfileRepository(ctrl)
	service := NewUserProfileService(nil, nil, repo, nil)

	userProfile := utils.GenerateUserProfile()

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().GetProfileByUserID(userProfile.UserID).Return(userProfile, nil).Times(1)

		result, err := service.GetProfileByUserID(userProfile.UserID)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, userProfile, result)
	})

	t.Run("user not found", func(t *testing.T) {
		repo.EXPECT().GetProfileByUserID(userProfile.UserID).Return(nil, nil).Times(1)

		result, err := service.GetProfileByUserID(userProfile.UserID)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "user profile does not exist", err.Error())
	})

	t.Run("error getting user", func(t *testing.T) {
		repo.EXPECT().GetProfileByUserID(userProfile.UserID).Return(nil, errors.New("error getting user")).Times(1)

		result, err := service.GetProfileByUserID(userProfile.UserID)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting user", err.Error())
	})
}

func TestUserProfileService_CreateUserProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	repo := mock_repositories.NewMockUserProfileRepository(ctrl)
	service := NewUserProfileService(nil, transaction, repo, nil)

	profile := utils.GenerateUserProfile()
	req := utils.GenerateCreateUserProfileRequest()

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().GetProfileByUserID(profile.UserID).Return(nil, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			}).Times(1)
		repo.EXPECT().CreateUserProfile(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		result, err := service.CreateUserProfile(profile.UserID, req)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, req.FirstName, result.FirstName)
		assert.Equal(t, req.LastName, result.LastName)
		assert.Equal(t, req.PhoneNumber, result.PhoneNumber)
		assert.Equal(t, req.DateOfBirth, result.DateOfBirth)
	})

	t.Run("duplicate user", func(t *testing.T) {
		repo.EXPECT().GetProfileByUserID(profile.UserID).Return(profile, nil).Times(1)

		result, err := service.CreateUserProfile(profile.UserID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.StatusCode)
		assert.Equal(t, "duplicate profile for current user", err.Error())
	})

	t.Run("error getting profile", func(t *testing.T) {
		repo.EXPECT().GetProfileByUserID(profile.UserID).Return(nil, errors.New("error getting profile")).Times(1)

		result, err := service.CreateUserProfile(profile.UserID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting profile", err.Error())
	})

	t.Run("error creating profile", func(t *testing.T) {
		repo.EXPECT().GetProfileByUserID(profile.UserID).Return(nil, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			}).Times(1)
		repo.EXPECT().CreateUserProfile(gomock.Any(), gomock.Any()).Return(errors.New("error creating profile")).Times(1)

		result, err := service.CreateUserProfile(profile.UserID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error creating profile", err.Error())
	})
}

func TestUserProfileService_UpdateUserProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	repo := mock_repositories.NewMockUserProfileRepository(ctrl)
	service := NewUserProfileService(nil, transaction, repo, nil)

	profile := utils.GenerateUserProfile()
	req := utils.GenerateUpdateUserProfileRequest()

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().GetProfileByUserID(profile.UserID).Return(profile, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			}).Times(1)
		repo.EXPECT().UpdateUserProfile(gomock.Any(), gomock.Any()).Return(nil).Times(1)

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

	t.Run("profile not found", func(t *testing.T) {
		repo.EXPECT().GetProfileByUserID(profile.UserID).Return(nil, nil).Times(1)

		result, err := service.UpdateUserProfile(profile.UserID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "user profile does not exist", err.Error())
	})

	t.Run("error getting user", func(t *testing.T) {
		repo.EXPECT().GetProfileByUserID(profile.UserID).Return(nil, errors.New("error getting user")).Times(1)

		result, err := service.UpdateUserProfile(profile.UserID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting user", err.Error())
	})

	t.Run("error updating profile", func(t *testing.T) {
		repo.EXPECT().GetProfileByUserID(profile.UserID).Return(profile, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			}).Times(1)
		repo.EXPECT().UpdateUserProfile(gomock.Any(), gomock.Any()).Return(errors.New("error updating profile")).Times(1)

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
	profileRepo := mock_repositories.NewMockUserProfileRepository(ctrl)
	profilePictureRepo := mock_repositories.NewMockProfilePictureRepository(ctrl)
	service := NewUserProfileService(nil, transaction, profileRepo, profilePictureRepo)

	profile := utils.GenerateUserProfile()
	file := utils.GenerateFileHeader()

	t.Run("success", func(t *testing.T) {
		profileRepo.EXPECT().GetProfileByUserID(profile.UserID).Return(profile, nil).Times(1)
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

	t.Run("profile not found", func(t *testing.T) {
		profileRepo.EXPECT().GetProfileByUserID(profile.UserID).Return(nil, nil).Times(1)

		result, err := service.UpdateProfilePicture(profile.UserID, file)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "user profile does not exist", err.Error())
	})

	t.Run("error getting profile", func(t *testing.T) {
		profileRepo.EXPECT().GetProfileByUserID(profile.UserID).Return(nil, errors.New("error getting profile")).Times(1)

		result, err := service.UpdateProfilePicture(profile.UserID, file)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting profile", err.Error())
	})

	t.Run("error creating picture", func(t *testing.T) {
		profileRepo.EXPECT().GetProfileByUserID(profile.UserID).Return(profile, nil).Times(1)
		profilePictureRepo.EXPECT().CreateProfilePicture(file, bucketName, gomock.Any()).Return(errors.New("error creating picture")).Times(1)

		result, err := service.UpdateProfilePicture(profile.UserID, file)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error creating picture", err.Error())
	})

	t.Run("error updating profile", func(t *testing.T) {
		profileRepo.EXPECT().GetProfileByUserID(profile.UserID).Return(profile, nil).Times(1)
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
	repo := mock_repositories.NewMockUserProfileRepository(ctrl)
	service := NewUserProfileService(nil, transaction, repo, nil)

	profile := utils.GenerateUserProfile()

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().GetProfileByUserID(profile.UserID).Return(profile, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		repo.EXPECT().UpdateProfilePicture(gomock.Any(), profile, nil).Return(profile, nil).Times(1)

		err := service.DeleteProfilePicture(profile.UserID)

		assert.Nil(t, err)
	})

	t.Run("profile not found", func(t *testing.T) {
		repo.EXPECT().GetProfileByUserID(profile.UserID).Return(nil, nil).Times(1)

		err := service.DeleteProfilePicture(profile.UserID)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "user profile does not exist", err.Error())
	})

	t.Run("error getting profile", func(t *testing.T) {
		repo.EXPECT().GetProfileByUserID(profile.UserID).Return(nil, errors.New("error getting profile")).Times(1)

		err := service.DeleteProfilePicture(profile.UserID)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting profile", err.Error())
	})

	t.Run("error deleting profile picture", func(t *testing.T) {
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)

		repo.EXPECT().GetProfileByUserID(profile.UserID).Return(profile, nil).Times(1)
		repo.EXPECT().UpdateProfilePicture(gomock.Any(), profile, nil).Return(nil, errors.New("error deleting profile picture")).Times(1)

		err := service.DeleteProfilePicture(profile.UserID)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error deleting profile picture", err.Error())
	})
}
