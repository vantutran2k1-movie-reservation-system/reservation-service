package repositories

import (
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/filters"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"gorm.io/gorm"
)

type TheaterRepository interface {
	GetTheater(filter filters.TheaterFilter, includeLocation bool) (*models.Theater, error)
	CreateTheater(tx *gorm.DB, theater *models.Theater) error
}

func NewTheaterRepository(db *gorm.DB) TheaterRepository {
	return &theaterRepository{
		db: db,
	}
}

type theaterRepository struct {
	db *gorm.DB
}

func (r *theaterRepository) GetTheater(filter filters.TheaterFilter, includeLocation bool) (*models.Theater, error) {
	query := filter.GetFilterQuery(r.db)
	if includeLocation {
		query = query.Preload("Location")
	}

	var theater models.Theater
	if err := query.First(&theater).Error; err != nil {
		if errors.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return &theater, nil
}

func (r *theaterRepository) CreateTheater(tx *gorm.DB, theater *models.Theater) error {
	return tx.Create(theater).Error
}
