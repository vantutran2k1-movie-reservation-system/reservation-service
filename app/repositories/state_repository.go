package repositories

import (
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/filters"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"gorm.io/gorm"
)

type StateRepository interface {
	GetState(filter filters.StateFilter) (*models.State, error)
	GetStates(filter filters.StateFilter) ([]*models.State, error)
	CreateState(tx *gorm.DB, state *models.State) error
}

func NewStateRepository(db *gorm.DB) StateRepository {
	return &stateRepository{db}
}

type stateRepository struct {
	db *gorm.DB
}

func (r *stateRepository) GetState(filter filters.StateFilter) (*models.State, error) {
	var state models.State
	if err := filter.GetFilterQuery(r.db).First(&state).Error; err != nil {
		if errors.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return &state, nil
}

func (r *stateRepository) GetStates(filter filters.StateFilter) ([]*models.State, error) {
	var states []*models.State
	if err := filter.GetFilterQuery(r.db).Find(&states).Error; err != nil {
		return nil, err
	}

	return states, nil
}

func (r *stateRepository) CreateState(tx *gorm.DB, state *models.State) error {
	return tx.Create(state).Error
}
