package repositories

import (
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/filters"
	"time"

	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"gorm.io/gorm"
)

type UserProfileRepository interface {
	GetProfile(filter filters.UserProfileFilter) (*models.UserProfile, error)
	CreateUserProfile(tx *gorm.DB, profile *models.UserProfile) error
	UpdateUserProfile(tx *gorm.DB, profile *models.UserProfile) error
	CreateOrUpdateUserProfile(tx *gorm.DB, profile *models.UserProfile) error
	UpdateProfilePicture(tx *gorm.DB, profile *models.UserProfile, url *string) (*models.UserProfile, error)
}

type userProfileRepository struct {
	db *gorm.DB
}

func NewUserProfileRepository(db *gorm.DB) UserProfileRepository {
	return &userProfileRepository{db: db}
}

func (r *userProfileRepository) GetProfile(filter filters.UserProfileFilter) (*models.UserProfile, error) {
	var p models.UserProfile
	if err := filter.GetFilterQuery(r.db).First(&p).Error; err != nil {
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

func (r *userProfileRepository) CreateOrUpdateUserProfile(tx *gorm.DB, profile *models.UserProfile) error {
	filter := filters.UserProfileFilter{
		Filter: &filters.SingleFilter{},
		UserID: &filters.Condition{Operator: filters.OpEqual, Value: profile.UserID},
	}
	p, err := r.GetProfile(filter)
	if err != nil {
		return err
	}
	if p != nil {
		return tx.Model(p).Updates(map[string]any{
			"first_name":          profile.FirstName,
			"last_name":           profile.LastName,
			"phone_number":        profile.PhoneNumber,
			"date_of_birth":       profile.DateOfBirth,
			"profile_picture_url": profile.ProfilePictureUrl,
			"bio":                 profile.Bio,
			"updated_at":          time.Now().UTC(),
		}).Error
	}

	return tx.Create(profile).Error
}

func (r *userProfileRepository) UpdateProfilePicture(tx *gorm.DB, profile *models.UserProfile, url *string) (*models.UserProfile, error) {
	if err := tx.Model(profile).Updates(map[string]any{"profile_picture_url": url, "updated_at": time.Now().UTC()}).Error; err != nil {
		return nil, err
	}

	return profile, nil
}
