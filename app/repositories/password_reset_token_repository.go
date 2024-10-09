package repositories

import (
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"gorm.io/gorm"
	"time"
)

type PasswordResetTokenRepository interface {
	GetActivePasswordResetToken(tokenValue string) (*models.PasswordResetToken, error)
	CreateToken(tx *gorm.DB, token *models.PasswordResetToken) error
}

func NewPasswordResetTokenRepository(db *gorm.DB) PasswordResetTokenRepository {
	return &passwordResetTokenRepository{db: db}
}

type passwordResetTokenRepository struct {
	db *gorm.DB
}

func (r *passwordResetTokenRepository) GetActivePasswordResetToken(tokenValue string) (*models.PasswordResetToken, error) {
	var token models.PasswordResetToken
	if err := r.db.Where("token_value = ? AND is_used = ? AND expires_at > ?", tokenValue, false, time.Now().UTC()).First(&token).Error; err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *passwordResetTokenRepository) CreateToken(tx *gorm.DB, token *models.PasswordResetToken) error {
	return tx.Create(token).Error
}
