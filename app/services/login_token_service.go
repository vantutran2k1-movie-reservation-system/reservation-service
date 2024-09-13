package services

import (
	"time"

	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"gorm.io/gorm"
)

var CreateLoginToken = func(db *gorm.DB, token *utils.AuthToken, userID uuid.UUID) *errors.ApiError {
	err := db.Where(
		"token_value = ? AND expires_at > ?",
		token.TokenValue,
		time.Now().UTC(),
	).First(&models.LoginToken{}).Error

	if err == nil {
		return errors.InternalServerError("Token %s already exists", token.TokenValue)
	}

	if !errors.IsRecordNotFoundError(err) {
		return errors.InternalServerError(err.Error())
	}

	t := models.LoginToken{
		ID:         uuid.New(),
		UserID:     userID,
		TokenValue: token.TokenValue,
		CreatedAt:  token.CreatedAt,
		ExpiresAt:  token.CreatedAt.Add(token.ValidDuration),
	}
	if err := db.Create(&t).Error; err != nil {
		return errors.InternalServerError(err.Error())
	}

	return nil
}
