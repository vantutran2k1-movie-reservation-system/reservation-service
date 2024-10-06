package services

import (
	"errors"
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

	userProfile := utils.GenerateRandomUserProfile()

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().GetProfileByUserID(userProfile.UserID).Return(userProfile, nil).Times(1)

		result, err := service.GetProfileByUserID(userProfile.UserID)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, userProfile, result)
	})

	t.Run("user not found", func(t *testing.T) {
		repo.EXPECT().GetProfileByUserID(userProfile.UserID).Return(nil, gorm.ErrRecordNotFound).Times(1)

		result, err := service.GetProfileByUserID(userProfile.UserID)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "User profile does not exist", err.Error())
	})

	t.Run("db error", func(t *testing.T) {
		repo.EXPECT().GetProfileByUserID(userProfile.UserID).Return(nil, errors.New("db error")).Times(1)

		result, err := service.GetProfileByUserID(userProfile.UserID)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Error())
	})
}

func TestUserProfileService_CreateUserProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	repo := mock_repositories.NewMockUserProfileRepository(ctrl)
	service := NewUserProfileService(nil, transaction, repo, nil)

	profile := utils.GenerateRandomUserProfile()

	t.Run("success", func(t *testing.T) {
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			}).Times(1)

		repo.EXPECT().CreateUserProfile(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		repo.EXPECT().GetProfileByUserID(profile.UserID).Return(nil, gorm.ErrRecordNotFound).Times(1)

		result, err := service.CreateUserProfile(profile.UserID, profile.FirstName, profile.LastName, profile.PhoneNumber, profile.DateOfBirth)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, profile.FirstName, result.FirstName)
		assert.Equal(t, profile.LastName, result.LastName)
		assert.Equal(t, profile.PhoneNumber, result.PhoneNumber)
		assert.Equal(t, profile.DateOfBirth, result.DateOfBirth)
	})

	t.Run("duplicate user", func(t *testing.T) {
		repo.EXPECT().GetProfileByUserID(profile.UserID).Return(profile, nil).Times(1)

		result, err := service.CreateUserProfile(profile.UserID, profile.FirstName, profile.LastName, profile.PhoneNumber, profile.DateOfBirth)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "Duplicate profile for current user", err.Message)
	})

	t.Run("db error getting profile", func(t *testing.T) {
		repo.EXPECT().GetProfileByUserID(profile.UserID).Return(nil, errors.New("db error")).Times(1)

		result, err := service.CreateUserProfile(profile.UserID, profile.FirstName, profile.LastName, profile.PhoneNumber, profile.DateOfBirth)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Error())
	})

	t.Run("db error creating profile", func(t *testing.T) {
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			}).Times(1)

		repo.EXPECT().GetProfileByUserID(profile.UserID).Return(nil, gorm.ErrRecordNotFound).Times(1)
		repo.EXPECT().CreateUserProfile(gomock.Any(), gomock.Any()).Return(errors.New("db error")).Times(1)

		result, err := service.CreateUserProfile(profile.UserID, profile.FirstName, profile.LastName, profile.PhoneNumber, profile.DateOfBirth)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Error())
	})
}

func TestUserProfileService_UpdateUserProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	repo := mock_repositories.NewMockUserProfileRepository(ctrl)
	service := NewUserProfileService(nil, transaction, repo, nil)

	currentProfile := utils.GenerateRandomUserProfile()
	updatedProfile := utils.GenerateRandomUserProfile()
	updatedProfile.ID = currentProfile.ID
	updatedProfile.UserID = currentProfile.UserID

	t.Run("success", func(t *testing.T) {
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			}).Times(1)

		repo.EXPECT().GetProfileByUserID(currentProfile.UserID).Return(currentProfile, nil).Times(1)
		repo.EXPECT().UpdateUserProfile(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		result, err := service.UpdateUserProfile(currentProfile.UserID, updatedProfile.FirstName, updatedProfile.LastName, updatedProfile.PhoneNumber, updatedProfile.DateOfBirth)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, updatedProfile.ID, result.ID)
		assert.Equal(t, updatedProfile.UserID, result.UserID)
		assert.Equal(t, updatedProfile.FirstName, result.FirstName)
		assert.Equal(t, updatedProfile.LastName, result.LastName)
		assert.Equal(t, updatedProfile.PhoneNumber, result.PhoneNumber)
		assert.Equal(t, updatedProfile.DateOfBirth, result.DateOfBirth)

	})

	t.Run("profile not found", func(t *testing.T) {
		repo.EXPECT().GetProfileByUserID(currentProfile.UserID).Return(nil, gorm.ErrRecordNotFound).Times(1)

		result, err := service.UpdateUserProfile(currentProfile.UserID, updatedProfile.FirstName, updatedProfile.LastName, updatedProfile.PhoneNumber, updatedProfile.DateOfBirth)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "User profile does not exist", err.Message)
	})

	t.Run("db error getting user", func(t *testing.T) {
		repo.EXPECT().GetProfileByUserID(currentProfile.UserID).Return(nil, errors.New("db error")).Times(1)

		result, err := service.UpdateUserProfile(currentProfile.UserID, updatedProfile.FirstName, updatedProfile.LastName, updatedProfile.PhoneNumber, updatedProfile.DateOfBirth)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Message)
	})

	t.Run("db error updating profile", func(t *testing.T) {
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			}).Times(1)

		repo.EXPECT().GetProfileByUserID(currentProfile.UserID).Return(currentProfile, nil).Times(1)
		repo.EXPECT().UpdateUserProfile(gomock.Any(), gomock.Any()).Return(errors.New("db error")).Times(1)

		result, err := service.UpdateUserProfile(currentProfile.UserID, updatedProfile.FirstName, updatedProfile.LastName, updatedProfile.PhoneNumber, updatedProfile.DateOfBirth)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Message)
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

	profile := utils.GenerateRandomUserProfile()
	file := utils.GenerateRandomFileHeader()

	t.Run("success", func(t *testing.T) {
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)

		profileRepo.EXPECT().GetProfileByUserID(profile.UserID).Return(profile, nil).Times(1)
		profilePictureRepo.EXPECT().CreateProfilePicture(file, bucketName, gomock.Any()).Return(nil).Times(1)
		profileRepo.EXPECT().UpdateProfilePicture(gomock.Any(), profile, gomock.Any()).Return(profile, nil).Times(1)

		result, err := service.UpdateProfilePicture(profile.UserID, file)

		assert.NotNil(t, result)
		assert.Nil(t, err)
	})

	t.Run("db error updating profile", func(t *testing.T) {
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)

		profileRepo.EXPECT().GetProfileByUserID(profile.UserID).Return(profile, nil).Times(1)
		profilePictureRepo.EXPECT().CreateProfilePicture(file, bucketName, gomock.Any()).Return(nil).Times(1)
		profileRepo.EXPECT().UpdateProfilePicture(gomock.Any(), profile, gomock.Any()).Return(nil, errors.New("db error")).Times(1)

		result, err := service.UpdateProfilePicture(profile.UserID, file)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Error())
	})

	t.Run("service error uploading picture", func(t *testing.T) {
		profileRepo.EXPECT().GetProfileByUserID(profile.UserID).Return(profile, nil).Times(1)
		profilePictureRepo.EXPECT().CreateProfilePicture(file, bucketName, gomock.Any()).Return(errors.New("service error")).Times(1)

		result, err := service.UpdateProfilePicture(profile.UserID, file)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "service error", err.Error())
	})
}

func TestUserProfileService_DeleteProfilePicture(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	profileRepo := mock_repositories.NewMockUserProfileRepository(ctrl)
	profilePictureRepo := mock_repositories.NewMockProfilePictureRepository(ctrl)
	service := NewUserProfileService(nil, transaction, profileRepo, profilePictureRepo)

	profile := utils.GenerateRandomUserProfile()

	t.Run("success", func(t *testing.T) {
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)

		profileRepo.EXPECT().GetProfileByUserID(profile.UserID).Return(profile, nil).Times(1)
		profileRepo.EXPECT().UpdateProfilePicture(gomock.Any(), profile, gomock.Any()).Return(profile, nil).Times(1)

		err := service.DeleteProfilePicture(profile.UserID)

		assert.Nil(t, err)
	})

	t.Run("error deleting profile picture", func(t *testing.T) {
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)

		profileRepo.EXPECT().GetProfileByUserID(profile.UserID).Return(profile, nil).Times(1)
		profileRepo.EXPECT().UpdateProfilePicture(gomock.Any(), profile, gomock.Any()).Return(nil, errors.New("db error")).Times(1)

		err := service.DeleteProfilePicture(profile.UserID)

		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Error())
	})
}
