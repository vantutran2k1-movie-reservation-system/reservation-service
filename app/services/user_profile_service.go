package services

import (
	"time"

	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"gorm.io/gorm"
)

var CreateUserProfile = func(db *gorm.DB, userID uuid.UUID, firstName, lastName, phoneNumber, dateOfBirth string) (*models.UserProfile, *errors.ApiError) {
	err := db.Where(&models.UserProfile{UserID: userID}).First(&models.UserProfile{}).Error

	if err == nil {
		return nil, errors.BadRequestError("Duplicate profile for user %s", userID)
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
	if err := db.Create(&p).Error; err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return &p, nil
}
