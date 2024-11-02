package repositories

import (
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/filters"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"gorm.io/gorm"
)

type CountryRepository interface {
	GetCountry(filter filters.CountryFilter) (*models.Country, error)
	GetCountries(filter filters.CountryFilter) ([]*models.Country, error)
	CreateCountry(tx *gorm.DB, country *models.Country) error
}

func NewCountryRepository(db *gorm.DB) CountryRepository {
	return &countryRepository{
		db: db,
	}
}

type countryRepository struct {
	db *gorm.DB
}

func (r *countryRepository) GetCountry(filter filters.CountryFilter) (*models.Country, error) {
	var country models.Country
	if err := filter.GetFilterQuery(r.db).First(&country).Error; err != nil {
		if errors.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return &country, nil
}

func (r *countryRepository) GetCountries(filter filters.CountryFilter) ([]*models.Country, error) {
	var countries []*models.Country
	if err := filter.GetFilterQuery(r.db).Find(&countries).Error; err != nil {
		return nil, err
	}

	return countries, nil
}

func (r *countryRepository) CreateCountry(tx *gorm.DB, country *models.Country) error {
	return tx.Create(country).Error
}
