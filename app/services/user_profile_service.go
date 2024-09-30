package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
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
	minioClient        *minio.Client
	transactionManager transaction.TransactionManager
	userProfileRepo    repositories.UserProfileRepository
}

func NewUserProfileService(
	db *gorm.DB,
	minioClient *minio.Client,
	transactionManager transaction.TransactionManager,
	userProfileRepo repositories.UserProfileRepository,
) UserProfileService {
	return &userProfileService{
		db:                 db,
		minioClient:        minioClient,
		transactionManager: transactionManager,
		userProfileRepo:    userProfileRepo,
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
	srcFile, err := file.Open()
	if err != nil {
		return errors.InternalServerError(err.Error())
	}
	defer srcFile.Close()

	ctx := context.Background()
	bucketName := os.Getenv("MINIO_PROFILE_PICTURE_BUCKET_NAME")
	objectName := fmt.Sprintf("%s/%d", userID, time.Now().Unix())
	contentType := file.Header.Get("Content-Type")

	err = s.createMinioBucketIfNotExists(ctx, bucketName)
	if err != nil {
		return errors.InternalServerError(err.Error())
	}

	_, err = s.minioClient.PutObject(ctx, bucketName, objectName, srcFile, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
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

func (s *userProfileService) createMinioBucketIfNotExists(ctx context.Context, bucketName string) error {
	exists, err := s.minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	return s.minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
}
