package repositories

import (
	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"gorm.io/gorm"
	"time"
)

type PasswordResetTokenRepository interface {
	GetActivePasswordResetToken(tokenValue string) (*models.PasswordResetToken, error)
	GetUserActivePasswordResetTokens(userID uuid.UUID) ([]*models.PasswordResetToken, error)
	CreateToken(tx *gorm.DB, token *models.PasswordResetToken) error
	UseToken(tx *gorm.DB, token *models.PasswordResetToken) error
	RevokeTokens(tx *gorm.DB, tokens []*models.PasswordResetToken) error
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

func (r *passwordResetTokenRepository) GetUserActivePasswordResetTokens(userID uuid.UUID) ([]*models.PasswordResetToken, error) {
	var tokens []*models.PasswordResetToken
	if err := r.db.Where("user_id = ? AND is_used = ? AND expires_at > ?", userID, false, time.Now().UTC()).Find(&tokens).Error; err != nil {
		return nil, err
	}

	return tokens, nil
}

func (r *passwordResetTokenRepository) CreateToken(tx *gorm.DB, token *models.PasswordResetToken) error {
	return tx.Create(token).Error
}

func (r *passwordResetTokenRepository) UseToken(tx *gorm.DB, token *models.PasswordResetToken) error {
	return tx.Model(token).Updates(map[string]any{"is_used": true}).Error
}

func (r *passwordResetTokenRepository) RevokeTokens(tx *gorm.DB, tokens []*models.PasswordResetToken) error {
	tokenIDs := make([]uuid.UUID, len(tokens))
	for i, token := range tokens {
		tokenIDs[i] = token.ID
	}

	return tx.Model(&models.PasswordResetToken{}).Where("id IN (?)", tokenIDs).Updates(map[string]interface{}{"expires_at": time.Now().UTC()}).Error
}
