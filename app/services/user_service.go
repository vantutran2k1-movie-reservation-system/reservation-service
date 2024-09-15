package services

import (
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/auth"
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

	hashedPassword, err := auth.GenerateHashedPassword(password)
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

var LoginUser = func(db *gorm.DB, rdb *redis.Client, email string, password string) (string, *errors.ApiError) {
	var u models.User
	if err := db.Where(&models.User{Email: email}).First(&u).Error; err != nil {
		if errors.IsRecordNotFoundError(err) {
			return "", errors.BadRequestError("Invalid email %s", email)
		}

		return "", errors.InternalServerError(err.Error())
	}

	if err := auth.CompareHashAndPassword(u.PasswordHash, password); err != nil {
		return "", errors.BadRequestError("Invalid password")
	}

	token, err := auth.GenerateJwtToken(u.ID)
	if err != nil {
		return "", errors.InternalServerError(err.Error())
	}

	if err := CreateSession(rdb, token, u.ID); err != nil {
		return "", errors.InternalServerError(err.Error())
	}

	if err := CreateLoginToken(db, token, u.ID); err != nil {
		return "", errors.InternalServerError(err.Error())
	}

	return token.TokenValue, nil
}

var LogoutUser = func(db *gorm.DB, rdb *redis.Client, tokenValue string) *errors.ApiError {
	if err := DeleteSession(rdb, tokenValue); err != nil {
		return errors.InternalServerError(err.Error())
	}

	if err := RevokeLoginToken(db, tokenValue); err != nil {
		return err
	}

	return nil
}

var UpdatePassword = func(db *gorm.DB, rdb *redis.Client, userID uuid.UUID, newPassword string) *errors.ApiError {
	var u models.User
	if err := db.Where(&models.User{ID: userID}).First(&u).Error; err != nil {
		return errors.InternalServerError(err.Error())
	}

	if err := auth.CompareHashAndPassword(u.PasswordHash, newPassword); err == nil {
		return errors.BadRequestError("New password can not be the same as current value")
	}

	p, err := auth.GenerateHashedPassword(newPassword)
	if err != nil {
		return errors.InternalServerError(err.Error())
	}

	u.PasswordHash = p
	u.UpdatedAt = time.Now().UTC()

	if err := db.Save(&u).Error; err != nil {
		return errors.InternalServerError(err.Error())
	}

	if err := DeleteUserSessions(rdb, userID); err != nil {
		return errors.InternalServerError(err.Error())
	}

	if err := RevokeUserLoginTokens(db, userID); err != nil {
		return err
	}

	return nil
}
