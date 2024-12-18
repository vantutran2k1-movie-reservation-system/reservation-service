package services

import (
	"fmt"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/filters"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/payloads"
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
	UpdateUserProfile(userID uuid.UUID, req payloads.UpdateUserProfileRequest) (*models.UserProfile, *errors.ApiError)
	UpdateProfilePicture(userID uuid.UUID, file *multipart.FileHeader) (*models.UserProfile, *errors.ApiError)
	DeleteProfilePicture(userID uuid.UUID) *errors.ApiError
}

type userProfileService struct {
	db                 *gorm.DB
	transactionManager transaction.TransactionManager
	userRepo           repositories.UserRepository
	userProfileRepo    repositories.UserProfileRepository
	profilePictureRepo repositories.ProfilePictureRepository
}

func NewUserProfileService(
	db *gorm.DB,
	transactionManager transaction.TransactionManager,
	userRepo repositories.UserRepository,
	userProfileRepo repositories.UserProfileRepository,
	profilePictureRepo repositories.ProfilePictureRepository,
) UserProfileService {
	return &userProfileService{
		db:                 db,
		transactionManager: transactionManager,
		userRepo:           userRepo,
		userProfileRepo:    userProfileRepo,
		profilePictureRepo: profilePictureRepo,
	}
}

func (s *userProfileService) GetProfileByUserID(userID uuid.UUID) (*models.UserProfile, *errors.ApiError) {
	return s.getUserProfileAndVerifyExist(userID)
}

func (s *userProfileService) UpdateUserProfile(userID uuid.UUID, req payloads.UpdateUserProfileRequest) (*models.UserProfile, *errors.ApiError) {
	p, apiErr := s.getUserProfileAndVerifyExist(userID)
	if apiErr != nil {
		return nil, apiErr
	}

	p.FirstName = req.FirstName
	p.LastName = req.LastName
	p.PhoneNumber = req.PhoneNumber
	p.DateOfBirth = req.DateOfBirth
	p.UpdatedAt = time.Now().UTC()
	if err := s.transactionManager.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		return s.userProfileRepo.UpdateUserProfile(tx, p)
	}); err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return p, nil
}

func (s *userProfileService) UpdateProfilePicture(userID uuid.UUID, file *multipart.FileHeader) (*models.UserProfile, *errors.ApiError) {
	p, apiErr := s.getUserProfileAndVerifyExist(userID)
	if apiErr != nil {
		return nil, apiErr
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
	p, apiErr := s.getUserProfileAndVerifyExist(userID)
	if apiErr != nil {
		return apiErr
	}

	if err := s.transactionManager.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		_, err := s.userProfileRepo.UpdateProfilePicture(tx, p, nil)
		return err
	}); err != nil {
		return errors.InternalServerError(err.Error())
	}

	return nil
}

func (s *userProfileService) getUser(id uuid.UUID) (*models.User, error) {
	return s.userRepo.GetUser(filters.UserFilter{
		Filter: &filters.SingleFilter{},
		ID:     &filters.Condition{Operator: filters.OpEqual, Value: id},
	}, true)
}

func (s *userProfileService) getUserProfileAndVerifyExist(userID uuid.UUID) (*models.UserProfile, *errors.ApiError) {
	u, err := s.getUser(userID)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	if u == nil {
		return nil, errors.NotFoundError("user does not exist")
	}
	if u.Profile == nil {
		return nil, errors.NotFoundError("user profile does not exist")
	}

	return u.Profile, nil
}
