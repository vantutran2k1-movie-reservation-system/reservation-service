package repositories

import (
	"time"

	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"gorm.io/gorm"
)

type UserProfileRepository interface {
	GetProfileByUserID(userID uuid.UUID) (*models.UserProfile, error)
	CreateUserProfile(tx *gorm.DB, profile *models.UserProfile) error
	UpdateUserProfile(tx *gorm.DB, profile *models.UserProfile) error
	UpdateProfilePicture(tx *gorm.DB, userID uuid.UUID, profilePictureUrl string) error
}

type userProfileRepository struct {
	db *gorm.DB
}

func NewUserProfileRepository(db *gorm.DB) UserProfileRepository {
	return &userProfileRepository{db: db}
}

func (r *userProfileRepository) GetProfileByUserID(userID uuid.UUID) (*models.UserProfile, error) {
	var p models.UserProfile
	if err := r.db.Where(&models.UserProfile{UserID: userID}).First(&p).Error; err != nil {
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

func (r *userProfileRepository) UpdateProfilePicture(tx *gorm.DB, userID uuid.UUID, profilePictureUrl string) error {
	p, err := r.GetProfileByUserID(userID)
	if err != nil {
		return err
	}

	p.ProfilePictureUrl = &profilePictureUrl
	p.UpdatedAt = time.Now().UTC()

	return tx.Save(p).Error
}
