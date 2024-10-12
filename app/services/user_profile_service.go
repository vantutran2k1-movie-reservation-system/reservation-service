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
	UpdateProfilePicture(userID uuid.UUID, file *multipart.FileHeader) (*models.UserProfile, *errors.ApiError)
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
		return nil, errors.InternalServerError(err.Error())

	}
	if p == nil {
		return nil, errors.NotFoundError("user profile does not exist")
	}

	return p, nil
}

func (s *userProfileService) CreateUserProfile(userID uuid.UUID, firstName, lastName string, phoneNumber, dateOfBirth *string) (*models.UserProfile, *errors.ApiError) {
	u, err := s.userProfileRepo.GetProfileByUserID(userID)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	if u != nil {
		return nil, errors.BadRequestError("duplicate profile for current user")
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
		return nil, errors.InternalServerError(err.Error())

	}
	if p == nil {
		return nil, errors.NotFoundError("user profile does not exist")
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

func (s *userProfileService) UpdateProfilePicture(userID uuid.UUID, file *multipart.FileHeader) (*models.UserProfile, *errors.ApiError) {
	p, err := s.userProfileRepo.GetProfileByUserID(userID)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())

	}
	if p == nil {
		return nil, errors.NotFoundError("user profile does not exist")
	}

	objectName := fmt.Sprintf("%s/%d", userID, time.Now().Unix())

	bucketName := os.Getenv("MINIO_PROFILE_PICTURE_BUCKET_NAME")
	if err := s.profilePictureRepo.CreateProfilePicture(file, bucketName, objectName); err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	if err := s.transactionManager.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		profile, err := s.userProfileRepo.UpdateProfilePicture(tx, p, &objectName)
		p = profile
		return err
	}); err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return p, nil
}

func (s *userProfileService) DeleteProfilePicture(userID uuid.UUID) *errors.ApiError {
	p, err := s.userProfileRepo.GetProfileByUserID(userID)
	if err != nil {
		return errors.InternalServerError(err.Error())

	}
	if p == nil {
		return errors.NotFoundError("user profile does not exist")
	}

	if err := s.transactionManager.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		_, err := s.userProfileRepo.UpdateProfilePicture(tx, p, nil)
		return err
	}); err != nil {
		return errors.InternalServerError(err.Error())
	}

	return nil
}
