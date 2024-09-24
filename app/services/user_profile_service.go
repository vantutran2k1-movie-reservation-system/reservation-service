package services

import (
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
}

type userProfileService struct {
	db              *gorm.DB
	userProfileRepo repositories.UserProfileRepository
}

func NewUserProfileService(db *gorm.DB, userProfileRepo repositories.UserProfileRepository) UserProfileService {
	return &userProfileService{db: db, userProfileRepo: userProfileRepo}
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
	if err := transaction.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
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
	if err := transaction.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		return s.userProfileRepo.UpdateUserProfile(tx, p)
	}); err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return p, nil
}
