package repositories

import (
	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/filters"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"gorm.io/gorm"
	"time"
)

type PasswordResetTokenRepository interface {
	GetToken(filter filters.PasswordResetTokenFilter) (*models.PasswordResetToken, error)
	GetTokens(filter filters.PasswordResetTokenFilter) ([]*models.PasswordResetToken, error)
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

func (r *passwordResetTokenRepository) GetToken(filter filters.PasswordResetTokenFilter) (*models.PasswordResetToken, error) {
	var token models.PasswordResetToken
	if err := filter.GetFilterQuery(r.db).First(&token).Error; err != nil {
		if errors.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return &token, nil
}

func (r *passwordResetTokenRepository) GetTokens(filter filters.PasswordResetTokenFilter) ([]*models.PasswordResetToken, error) {
	var tokens []*models.PasswordResetToken
	if err := filter.GetFilterQuery(r.db).Find(&tokens).Error; err != nil {
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

	filter := filters.PasswordResetTokenFilter{
		Filter: &filters.MultiFilter{},
		ID:     &filters.Condition{Operator: filters.OpIn, Value: tokenIDs},
	}

	return filter.GetFilterQuery(tx).Model(&models.PasswordResetToken{}).Updates(map[string]any{"expires_at": time.Now().UTC()}).Error
}
