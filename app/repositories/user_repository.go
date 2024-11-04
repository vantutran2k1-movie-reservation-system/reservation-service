package repositories

import (
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/filters"
	"time"

	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUser(filter filters.UserFilter) (*models.User, error)
	CreateUser(tx *gorm.DB, user *models.User) error
	UpdatePassword(tx *gorm.DB, user *models.User, password string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUser(filter filters.UserFilter) (*models.User, error) {
	var u models.User
	if err := filter.GetFilterQuery(r.db).First(&u).Error; err != nil {
		if errors.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return &u, nil
}

func (r *userRepository) CreateUser(tx *gorm.DB, user *models.User) error {
	return tx.Create(user).Error
}

func (r *userRepository) UpdatePassword(tx *gorm.DB, user *models.User, password string) (*models.User, error) {
	if err := tx.Model(user).Updates(map[string]any{"password_hash": password, "updated_at": time.Now().UTC()}).Error; err != nil {
		return nil, err
	}

	return user, nil
}
