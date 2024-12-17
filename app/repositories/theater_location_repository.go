package repositories

import (
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/filters"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"gorm.io/gorm"
)

type TheaterLocationRepository interface {
	GetLocation(filter filters.TheaterLocationFilter) (*models.TheaterLocation, error)
	CreateTheaterLocation(tx *gorm.DB, location *models.TheaterLocation) error
	UpdateTheaterLocation(tx *gorm.DB, location *models.TheaterLocation) error
}

func NewTheaterLocationRepository(db *gorm.DB) TheaterLocationRepository {
	return &theaterLocationRepository{
		db: db,
	}
}

type theaterLocationRepository struct {
	db *gorm.DB
}

func (r *theaterLocationRepository) GetLocation(filter filters.TheaterLocationFilter) (*models.TheaterLocation, error) {
	var location models.TheaterLocation
	if err := filter.GetFilterQuery(r.db).First(&location).Error; err != nil {
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

func (r *theaterLocationRepository) UpdateTheaterLocation(tx *gorm.DB, location *models.TheaterLocation) error {
	return tx.Save(location).Error
}
