package services

import (
	"time"

	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/sessions"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
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

	hashedPassword, err := utils.GenerateHashedPassword(password)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	u := models.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	if err := db.Create(&u).Error; err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return &u, nil
}

var LoginUser = func(db *gorm.DB, email string, password string) (string, *errors.ApiError) {
	var u models.User
	if err := db.Where(&models.User{Email: email}).First(&u).Error; err != nil {
		if errors.IsRecordNotFoundError(err) {
			return "", errors.BadRequestError("Invalid email %s", email)
		}

		return "", errors.InternalServerError(err.Error())
	}

	if err := utils.CompareHashAndPassword(u.PasswordHash, password); err != nil {
		return "", errors.BadRequestError("Invalid password")
	}

	token, err := utils.GenerateToken(u.ID)
	if err != nil {
		return "", errors.InternalServerError(err.Error())
	}

	if err := sessions.CreateSession(token, u.ID); err != nil {
		return "", errors.InternalServerError(err.Error())
	}
	return token, nil
}
