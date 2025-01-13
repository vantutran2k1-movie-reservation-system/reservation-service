package repositories

import (
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/filters"
	"time"

	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUser(filter filters.UserFilter, includeProfile bool) (*models.User, error)
	UserExists(filter filters.UserFilter) (bool, error)
	CreateUser(tx *gorm.DB, user *models.User) error
	CreateOrUpdateUser(tx *gorm.DB, user *models.User) error
	UpdatePassword(tx *gorm.DB, user *models.User, password string) (*models.User, error)
	VerifyUser(tx *gorm.DB, user *models.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUser(filter filters.UserFilter, includeProfile bool) (*models.User, error) {
	query := filter.GetFilterQuery(r.db)
	if includeProfile {
		query = query.Preload("Profile")
	}

	var u models.User
	if err := query.First(&u).Error; err != nil {
		if errors.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return &u, nil
}

func (r *userRepository) UserExists(filter filters.UserFilter) (bool, error) {
	query := filter.GetFilterQuery(r.db)
	var u models.User
	if err := query.Select("id").First(&u).Error; err != nil {
		if errors.IsRecordNotFoundError(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (r *userRepository) CreateUser(tx *gorm.DB, user *models.User) error {
	return tx.Create(user).Error
}

func (r *userRepository) CreateOrUpdateUser(tx *gorm.DB, user *models.User) error {
	filter := filters.UserFilter{
		Filter: &filters.SingleFilter{},
		Email:  &filters.Condition{Operator: filters.OpEqual, Value: user.Email},
	}
	u, err := r.GetUser(filter, false)
	if err != nil {
		return err
	}
	if u != nil {
		return tx.Model(u).Updates(map[string]any{
			"password_hash": user.PasswordHash,
			"is_active":     user.IsActive,
			"is_verified":   user.IsVerified,
			"updated_at":    time.Now().UTC(),
		}).Error
	}

	return tx.Create(user).Error
}

func (r *userRepository) UpdatePassword(tx *gorm.DB, user *models.User, password string) (*models.User, error) {
	if err := tx.Model(user).Updates(map[string]any{"password_hash": password, "updated_at": time.Now().UTC()}).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) VerifyUser(tx *gorm.DB, user *models.User) error {
	if err := tx.Model(user).Updates(map[string]any{"is_verified": true, "updated_at": time.Now().UTC()}).Error; err != nil {
		return err
	}

	return nil
}
