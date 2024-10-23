package repositories

import (
	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"gorm.io/gorm"
)

type TheaterLocationRepository interface {
	GetLocationByTheaterID(theaterId uuid.UUID) (*models.TheaterLocation, error)
	CreateTheaterLocation(tx *gorm.DB, location *models.TheaterLocation) error
}

func NewTheaterLocationRepository(db *gorm.DB) TheaterLocationRepository {
	return &theaterLocationRepository{
		db: db,
	}
}

type theaterLocationRepository struct {
	db *gorm.DB
}

func (r *theaterLocationRepository) GetLocationByTheaterID(theaterId uuid.UUID) (*models.TheaterLocation, error) {
	var location models.TheaterLocation
	if err := r.db.Where("theater_id = ?", theaterId).First(&location).Error; err != nil {
		if errors.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return &location, nil
}

func (r *theaterLocationRepository) CreateTheaterLocation(tx *gorm.DB, location *models.TheaterLocation) error {
	return tx.Create(location).Error
}
