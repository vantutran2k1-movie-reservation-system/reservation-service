package repositories

import (
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/filters"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/payloads"
	"gorm.io/gorm"
)

type TheaterRepository interface {
	GetTheater(filter filters.TheaterFilter, includeLocation bool) (*models.Theater, error)
	GetTheaters(filter filters.TheaterFilter, includeLocation bool) ([]*models.Theater, error)
	GetNearbyTheatersWithLocations(lat, lon, distance float64) ([]*payloads.GetTheaterWithLocationResult, error)
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

func (r *theaterRepository) GetNearbyTheatersWithLocations(lat, lon, distance float64) ([]*payloads.GetTheaterWithLocationResult, error) {
	var theaterLocations []*payloads.GetTheaterWithLocationResult
	query := `
		WITH data AS (
			SELECT
				t.id, t.name,
				tl.id AS location_id, tl.city_id, tl.address, tl.postal_code, tl.latitude, tl.longitude,
				(6371 * acos(cos(radians(?)) * cos(radians(tl.latitude)) * cos(radians(tl.longitude) - radians(?)) + sin(radians(?)) * sin(radians(tl.latitude)))) AS distance
			FROM theaters t
			JOIN theater_locations tl ON t.id = tl.theater_id
		)
		SELECT *
		FROM data
		WHERE distance < ?
		ORDER BY distance
	`
	if err := r.db.Raw(query, lat, lon, lat, distance).Scan(&theaterLocations).Error; err != nil {
		return nil, err
	}

	return theaterLocations, nil
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
