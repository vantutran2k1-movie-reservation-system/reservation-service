package repositories

import (
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/filters"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"gorm.io/gorm"
)

type TheaterRepository interface {
	GetTheater(filter filters.TheaterFilter, includeLocation bool) (*models.Theater, error)
	GetTheaters(filter filters.TheaterFilter, includeLocation bool) ([]*models.Theater, error)
	GetNumbersOfTheater(filter filters.TheaterFilter) (int, error)
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

func (r *theaterRepository) GetTheaters(filter filters.TheaterFilter, includeLocation bool) ([]*models.Theater, error) {
	query := filter.GetFilterQuery(r.db)
	if includeLocation {
		query = query.Preload("Location")
	}

	var theaters []*models.Theater
	if err := query.Find(&theaters).Error; err != nil {
		return nil, err
	}

	return theaters, nil
}

func (r *theaterRepository) GetNumbersOfTheater(filter filters.TheaterFilter) (int, error) {
	var count int64
	if err := filter.GetFilterQuery(r.db).Model(&models.Theater{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

func (r *theaterRepository) CreateTheater(tx *gorm.DB, theater *models.Theater) error {
	return tx.Create(theater).Error
}
