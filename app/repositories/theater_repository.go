package repositories

import (
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/payloads"
	"gorm.io/gorm"
)

type TheaterRepository interface {
	GetTheater(filter payloads.GetTheaterFilter) (*models.Theater, error)
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

func (r *theaterRepository) GetTheater(filter payloads.GetTheaterFilter) (*models.Theater, error) {
	var theater models.Theater
	if err := r.getTheaterFilterQuery(r.db, filter).First(&theater).Error; err != nil {
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

func (r *theaterRepository) getTheaterFilterQuery(query *gorm.DB, filter payloads.GetTheaterFilter) *gorm.DB {
	if filter.ID != nil {
		query = query.Where("id = ?", filter.ID)
	}

	if filter.Name != nil {
		query = query.Where("name = ?", *filter.Name)
	}

	if filter.IncludeLocation != nil && *filter.IncludeLocation {
		query = query.Preload("Location")
	}

	return query
}
