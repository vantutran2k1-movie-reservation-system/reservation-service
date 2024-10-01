package services

import (
	"errors"
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

	mockUserProfileRepo := mock_repositories.NewMockUserProfileRepository(ctrl)

	userProfileService := NewUserProfileService(nil, nil, nil, mockUserProfileRepo)

	userProfile := utils.GenerateRandomUserProfile()

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

	userProfileService := NewUserProfileService(nil, nil, transactionMock, userProfileRepo)

	profile := utils.GenerateRandomUserProfile()

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

	userProfileService := NewUserProfileService(nil, nil, transactionMock, userProfileRepo)

	firstName := "First"
	lastName := "Last"
	phoneNumber := "1234567890"
	dateOfBirth := "1990-01-01"

	profile := utils.GenerateRandomUserProfile()

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

// func TestUserProfileService_UpdateProfilePicture(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	bucketName := "test-bucket"
// 	os.Setenv("MINIO_PROFILE_PICTURE_BUCKET_NAME", bucketName)
// 	defer os.Unsetenv("MINIO_PROFILE_PICTURE_BUCKET_NAME")

// 	minioMock := mocks.NewMockMinioClient(ctrl)
// 	transactionMock := mock_transaction.NewMockTransactionManager(ctrl)
// 	userProfileRepo := mock_repositories.NewMockUserProfileRepository(ctrl)

// 	userProfileService := NewUserProfileService(nil, nil, transactionMock, userProfileRepo)

// 	objectName := fmt.Sprintf("%s/%d", userID, time.Now().Unix())
// 	file := utils.GenerateRandomFileHeader()

// 	srcFileContent := "fake image content"
// 	srcFile := ioutil.NopCloser(strings.NewReader(srcFileContent))

// 	t.Run("success", func(t *testing.T) {
// 		minioMock.EXPECT().PutObject(gomock.Any(), bucketName, objectName, gomock.Any(), file.Size, gomock.Any()).
// 			Return(minio.UploadInfo{}, nil).Times(1)

// 		// 2. ExecuteInTransaction should be called and succeed.
// 		transactionMock.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
// 			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
// 				// Simulate the transaction function being called.
// 				return fn(nil)
// 			}).Times(1)

// 		// 3. UpdateProfilePicture should be called within the transaction.
// 		userProfileRepo.EXPECT().UpdateProfilePicture(gomock.Any(), userID, objectName).Return(nil).Times(1)

// 		// Mock the file.Open call to return our fake srcFile.
// 		mockFile := mocks.NewMockFile(ctrl)
// 		mockFile.EXPECT().Open().Return(srcFile, nil).Times(1)
// 		defer srcFile.Close() // Close the file after the test.

// 		// Call the service method.
// 		err := userProfileService.UpdateProfilePicture(userID, file)

// 		// Assertions.
// 		assert.Nil(t, err)
// 	})
// }
