package repositories

import (
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/payloads"
	"gorm.io/gorm"
)

type CityRepository interface {
	GetCity(filter payloads.GetCityFilter) (*models.City, error)
	GetCities(filter payloads.GetCitiesFilter) ([]*models.City, error)
	CreateCity(tx *gorm.DB, city *models.City) error
}

func NewCityRepository(db *gorm.DB) CityRepository {
	return &cityRepository{db: db}
}

type cityRepository struct {
	db *gorm.DB
}

func (r *cityRepository) GetCity(filter payloads.GetCityFilter) (*models.City, error) {
	var city models.City
	if err := r.getCityFilterQuery(r.db, filter).First(&city).Error; err != nil {
		if errors.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return &city, nil
}

func (r *cityRepository) GetCities(filter payloads.GetCitiesFilter) ([]*models.City, error) {
	var cities []*models.City
	if err := r.getCitiesFilterQuery(r.db, filter).Find(&cities).Error; err != nil {
		return nil, err
	}

	return cities, nil
}

func (r *cityRepository) CreateCity(tx *gorm.DB, city *models.City) error {
	return tx.Create(city).Error
}

func (r *cityRepository) getCityFilterQuery(query *gorm.DB, filter payloads.GetCityFilter) *gorm.DB {
	if filter.ID != nil {
		query = query.Where("id = ?", filter.ID)
	}

	if filter.StateID != nil {
		query = query.Where("state_id = ?", filter.StateID)
	}

	if filter.Name != nil {
		query = query.Where("name = ?", filter.Name)
	}

	return query
}

func (r *cityRepository) getCitiesFilterQuery(query *gorm.DB, filter payloads.GetCitiesFilter) *gorm.DB {
	query = query.Where("state_id = ?", filter.StateID)

	return query
}
