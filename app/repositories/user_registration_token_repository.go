package repositories

import (
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"gorm.io/gorm"
)

type UserRegistrationTokenRepository interface {
	CreateToken(tx *gorm.DB, token *models.UserRegistrationToken) error
}

func NewUserRegistrationTokenRepository(db *gorm.DB) UserRegistrationTokenRepository {
	return &userRegistrationTokenRepository{db}
}

type userRegistrationTokenRepository struct {
	db *gorm.DB
}

func (u *userRegistrationTokenRepository) CreateToken(tx *gorm.DB, token *models.UserRegistrationToken) error {
	return tx.Create(token).Error
}
