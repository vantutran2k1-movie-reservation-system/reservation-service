package services

import (
	"time"

	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"gorm.io/gorm"
)

var CreateUser = func(db *gorm.DB, email string, password string) (*models.User, *errors.ApiError) {
	err := db.Where(&models.User{Email: email}).First(&models.User{}).Error

	if err == nil {
		return nil, errors.BadRequestError("Email %s already exists", email)
	}

	if !errors.IsRecordNotFoundError(err) {
		return nil, errors.InternalServerError(err.Error())
	}

	u := models.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: password,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	if err := db.Create(&u).Error; err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return &u, nil
}
