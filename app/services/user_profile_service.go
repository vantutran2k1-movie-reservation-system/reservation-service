package services

import (
	"fmt"
	"mime/multipart"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/repositories"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/transaction"
	"gorm.io/gorm"
)

type UserProfileService interface {
	GetProfileByUserID(userID uuid.UUID) (*models.UserProfile, *errors.ApiError)
	CreateUserProfile(userID uuid.UUID, firstName, lastName string, phoneNumber, dateOfBirth *string) (*models.UserProfile, *errors.ApiError)
	UpdateUserProfile(userID uuid.UUID, firstName, lastName string, phoneNumber, dateOfBirth *string) (*models.UserProfile, *errors.ApiError)
	UpdateProfilePicture(userID uuid.UUID, file *multipart.FileHeader) *errors.ApiError
	DeleteProfilePicture(userID uuid.UUID) *errors.ApiError
}

type userProfileService struct {
	db                 *gorm.DB
	transactionManager transaction.TransactionManager
	userProfileRepo    repositories.UserProfileRepository
	profilePictureRepo repositories.ProfilePictureRepository
}

func NewUserProfileService(
	db *gorm.DB,
	transactionManager transaction.TransactionManager,
	userProfileRepo repositories.UserProfileRepository,
	profilePictureRepo repositories.ProfilePictureRepository,
) UserProfileService {
	return &userProfileService{
		db:                 db,
		transactionManager: transactionManager,
		userProfileRepo:    userProfileRepo,
		profilePictureRepo: profilePictureRepo,
	}
}

func (s *userProfileService) GetProfileByUserID(userID uuid.UUID) (*models.UserProfile, *errors.ApiError) {
	p, err := s.userProfileRepo.GetProfileByUserID(userID)
	if err != nil {
		if errors.IsRecordNotFoundError(err) {
			return nil, errors.BadRequestError("User profile does not exist")
		}

		return nil, errors.InternalServerError(err.Error())
	}

	return p, nil
}

func (s *userProfileService) CreateUserProfile(userID uuid.UUID, firstName, lastName string, phoneNumber, dateOfBirth *string) (*models.UserProfile, *errors.ApiError) {
	_, err := s.userProfileRepo.GetProfileByUserID(userID)
	if err == nil {
		return nil, errors.BadRequestError("Duplicate profile for current user")
	}
	if !errors.IsRecordNotFoundError(err) {
		return nil, errors.InternalServerError(err.Error())
	}

	p := models.UserProfile{
		ID:          uuid.New(),
		UserID:      userID,
		FirstName:   firstName,
		LastName:    lastName,
		PhoneNumber: phoneNumber,
		DateOfBirth: dateOfBirth,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	if err := s.transactionManager.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		return s.userProfileRepo.CreateUserProfile(tx, &p)
	}); err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return &p, nil
}

func (s *userProfileService) UpdateUserProfile(userID uuid.UUID, firstName, lastName string, phoneNumber, dateOfBirth *string) (*models.UserProfile, *errors.ApiError) {
	p, err := s.userProfileRepo.GetProfileByUserID(userID)
	if err != nil {
		if errors.IsRecordNotFoundError(err) {
			return nil, errors.BadRequestError("User profile does not exist")
		}

		return nil, errors.InternalServerError(err.Error())
	}

	p.FirstName = firstName
	p.LastName = lastName
	p.PhoneNumber = phoneNumber
	p.DateOfBirth = dateOfBirth
	p.UpdatedAt = time.Now().UTC()
	if err := s.transactionManager.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		return s.userProfileRepo.UpdateUserProfile(tx, p)
	}); err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return p, nil
}

func (s *userProfileService) UpdateProfilePicture(userID uuid.UUID, file *multipart.FileHeader) *errors.ApiError {
	bucketName := os.Getenv("MINIO_PROFILE_PICTURE_BUCKET_NAME")

	bucketExists := s.profilePictureRepo.BucketExists(bucketName)
	if !bucketExists {
		if err := s.profilePictureRepo.CreateBucket(bucketName); err != nil {
			return errors.InternalServerError(err.Error())
		}
	}

	objectName := fmt.Sprintf("%s/%d", userID, time.Now().Unix())
	if err := s.profilePictureRepo.CreateProfilePicture(file, bucketName, objectName); err != nil {
		return errors.InternalServerError(err.Error())
	}

	if err := s.transactionManager.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		return s.userProfileRepo.UpdateProfilePicture(tx, userID, objectName)
	}); err != nil {
		return errors.InternalServerError(err.Error())
	}

	return nil
}

func (s *userProfileService) DeleteProfilePicture(userID uuid.UUID) *errors.ApiError {
	if err := s.transactionManager.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		return s.userProfileRepo.DeleteProfilePicture(tx, userID)
	}); err != nil {
		return errors.InternalServerError(err.Error())
	}

	return nil
}
