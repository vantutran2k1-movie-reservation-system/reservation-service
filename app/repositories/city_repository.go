package repositories

import (
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/filters"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"gorm.io/gorm"
)

type CityRepository interface {
	GetCity(filter filters.CityFilter) (*models.City, error)
	GetCities(filter filters.CityFilter) ([]*models.City, error)
	CreateCity(tx *gorm.DB, city *models.City) error
}

func NewCityRepository(db *gorm.DB) CityRepository {
	return &cityRepository{db: db}
}

type cityRepository struct {
	db *gorm.DB
}

func (r *cityRepository) GetCity(filter filters.CityFilter) (*models.City, error) {
	var city models.City
	if err := filter.GetFilterQuery(r.db).First(&city).Error; err != nil {
		if errors.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return &city, nil
}

func (r *cityRepository) GetCities(filter filters.CityFilter) ([]*models.City, error) {
	var cities []*models.City
	if err := filter.GetFilterQuery(r.db).Find(&cities).Error; err != nil {
		return nil, err
	}

	return cities, nil
}

func (r *cityRepository) CreateCity(tx *gorm.DB, city *models.City) error {
	return tx.Create(city).Error
}
