package repositories

import (
	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"gorm.io/gorm"
)

type CountryRepository interface {
	GetCountry(id uuid.UUID) (*models.Country, error)
	GetCountryByName(name string) (*models.Country, error)
	GetCountryByCode(code string) (*models.Country, error)
	GetCountries() ([]*models.Country, error)
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

func (r *countryRepository) GetCountry(id uuid.UUID) (*models.Country, error) {
	var country models.Country
	if err := r.db.Where("id = ?", id).First(&country).Error; err != nil {
		if errors.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return &country, nil
}

func (r *countryRepository) GetCountryByName(name string) (*models.Country, error) {
	var country models.Country
	if err := r.db.Where("name = ?", name).First(&country).Error; err != nil {
		if errors.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return &country, nil
}

func (r *countryRepository) GetCountryByCode(code string) (*models.Country, error) {
	var country models.Country
	if err := r.db.Where("code = ?", code).First(&country).Error; err != nil {
		if errors.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return &country, nil
}

func (r *countryRepository) GetCountries() ([]*models.Country, error) {
	var countries []*models.Country
	if err := r.db.Find(&countries).Error; err != nil {
		return nil, err
	}

	return countries, nil
}

func (r *countryRepository) CreateCountry(tx *gorm.DB, country *models.Country) error {
	return tx.Create(country).Error
}
