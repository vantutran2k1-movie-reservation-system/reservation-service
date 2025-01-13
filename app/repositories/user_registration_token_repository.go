package repositories

import (
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/filters"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"gorm.io/gorm"
)

type UserRegistrationTokenRepository interface {
	GetToken(filter filters.UserRegistrationTokenFilter) (*models.UserRegistrationToken, error)
	CreateToken(tx *gorm.DB, token *models.UserRegistrationToken) error
	UseToken(tx *gorm.DB, token *models.UserRegistrationToken) error
}

func NewUserRegistrationTokenRepository(db *gorm.DB) UserRegistrationTokenRepository {
	return &userRegistrationTokenRepository{db}
}

type userRegistrationTokenRepository struct {
	db *gorm.DB
}

func (r *userRegistrationTokenRepository) GetToken(filter filters.UserRegistrationTokenFilter) (*models.UserRegistrationToken, error) {
	var token models.UserRegistrationToken
	if err := filter.GetFilterQuery(r.db).First(&token).Error; err != nil {
		if errors.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return &token, nil
}

func (r *userRegistrationTokenRepository) CreateToken(tx *gorm.DB, token *models.UserRegistrationToken) error {
	return tx.Create(token).Error
}

func (r *userRegistrationTokenRepository) UseToken(tx *gorm.DB, token *models.UserRegistrationToken) error {
	return tx.Model(token).Updates(map[string]any{"is_used": true}).Error
}
