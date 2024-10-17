package repositories

import (
	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"gorm.io/gorm"
)

type StateRepository interface {
	GetStateByName(countryID uuid.UUID, name string) (*models.State, error)
	GetStatesByCountry(countryID uuid.UUID) ([]*models.State, error)
	CreateState(tx *gorm.DB, state *models.State) error
}

func NewStateRepository(db *gorm.DB) StateRepository {
	return &stateRepository{db}
}

type stateRepository struct {
	db *gorm.DB
}

func (r *stateRepository) GetStateByName(countryID uuid.UUID, name string) (*models.State, error) {
	var state models.State
	if err := r.db.Where("country_id = ? AND name = ?", countryID, name).First(&state).Error; err != nil {
		if errors.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return &state, nil
}

func (r *stateRepository) GetStatesByCountry(countryID uuid.UUID) ([]*models.State, error) {
	var states []*models.State
	if err := r.db.Where("country_id = ?", countryID).Find(&states).Error; err != nil {
		return nil, err
	}

	return states, nil
}

func (r *stateRepository) CreateState(tx *gorm.DB, state *models.State) error {
	return tx.Create(state).Error
}
