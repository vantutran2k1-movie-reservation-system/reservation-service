package repositories

import (
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"time"

	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"gorm.io/gorm"
)

type UserProfileRepository interface {
	GetProfileByUserID(userID uuid.UUID) (*models.UserProfile, error)
	CreateUserProfile(tx *gorm.DB, profile *models.UserProfile) error
	UpdateUserProfile(tx *gorm.DB, profile *models.UserProfile) error
	UpdateProfilePicture(tx *gorm.DB, profile *models.UserProfile, url *string) (*models.UserProfile, error)
}

type userProfileRepository struct {
	db *gorm.DB
}

func NewUserProfileRepository(db *gorm.DB) UserProfileRepository {
	return &userProfileRepository{db: db}
}

func (r *userProfileRepository) GetProfileByUserID(userID uuid.UUID) (*models.UserProfile, error) {
	var p models.UserProfile
	if err := r.db.Where("user_id = ?", userID).First(&p).Error; err != nil {
		if errors.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return &p, nil
}

func (r *userProfileRepository) CreateUserProfile(tx *gorm.DB, profile *models.UserProfile) error {
	return tx.Create(profile).Error
}

func (r *userProfileRepository) UpdateUserProfile(tx *gorm.DB, profile *models.UserProfile) error {
	return tx.Save(profile).Error
}

func (r *userProfileRepository) UpdateProfilePicture(tx *gorm.DB, profile *models.UserProfile, url *string) (*models.UserProfile, error) {
	if err := tx.Model(profile).Updates(map[string]any{"profile_picture_url": url, "updated_at": time.Now().UTC()}).Error; err != nil {
		return nil, err
	}

	return profile, nil
}
