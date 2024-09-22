package repositories

import (
	"time"

	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"gorm.io/gorm"
)

type LoginTokenRepository interface {
	GetActiveLoginToken(tokenValue string) (*models.LoginToken, error)
	CreateLoginToken(tx *gorm.DB, loginToken *models.LoginToken) error
	RevokeLoginToken(tx *gorm.DB, loginToken *models.LoginToken) error
	RevokeUserLoginTokens(tx *gorm.DB, userID uuid.UUID) error
}

type loginTokenRepository struct {
	db *gorm.DB
}

func NewLoginTokenRepository(db *gorm.DB) LoginTokenRepository {
	return &loginTokenRepository{db: db}
}

func (r *loginTokenRepository) GetActiveLoginToken(tokenValue string) (*models.LoginToken, error) {
	var t models.LoginToken
	err := r.db.Where("token_value = ? AND expires_at > ?", tokenValue, time.Now().UTC()).First(&t).Error
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *loginTokenRepository) CreateLoginToken(tx *gorm.DB, loginToken *models.LoginToken) error {
	return tx.Create(loginToken).Error
}

func (r *loginTokenRepository) RevokeLoginToken(tx *gorm.DB, loginToken *models.LoginToken) error {
	loginToken.ExpiresAt = time.Now().UTC()
	return tx.Save(&loginToken).Error
}

func (r *loginTokenRepository) RevokeUserLoginTokens(tx *gorm.DB, userID uuid.UUID) error {
	var tokens []*models.LoginToken
	if err := r.db.Where("user_id = ? AND expires_at > ?", userID, time.Now().UTC()).Find(&tokens).Error; err != nil {
		return err
	}

	for _, t := range tokens {
		t.ExpiresAt = time.Now().UTC()
	}

	return tx.Save(&tokens).Error
}
