package repositories

import (
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/filters"
	"time"

	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"gorm.io/gorm"
)

type LoginTokenRepository interface {
	GetLoginToken(filter filters.LoginTokenFilter) (*models.LoginToken, error)
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

func (r *loginTokenRepository) GetLoginToken(filter filters.LoginTokenFilter) (*models.LoginToken, error) {
	var t models.LoginToken
	if err := filter.GetFilterQuery(r.db).First(&t).Error; err != nil {
		if errors.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return &t, nil
}

func (r *loginTokenRepository) CreateLoginToken(tx *gorm.DB, loginToken *models.LoginToken) error {
	return tx.Create(loginToken).Error
}

func (r *loginTokenRepository) RevokeLoginToken(tx *gorm.DB, loginToken *models.LoginToken) error {
	return tx.Model(loginToken).Updates(map[string]any{"expires_at": time.Now().UTC()}).Error
}

func (r *loginTokenRepository) RevokeUserLoginTokens(tx *gorm.DB, userID uuid.UUID) error {
	return tx.Model(&models.LoginToken{}).Where("user_id = ? AND expires_at > ?", userID, time.Now().UTC()).Updates(map[string]any{"expires_at": time.Now().UTC()}).Error
}
