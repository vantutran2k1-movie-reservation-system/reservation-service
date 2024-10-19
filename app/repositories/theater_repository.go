package repositories

import (
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"gorm.io/gorm"
)

type TheaterRepository interface {
	GetTheaterByName(name string) (*models.Theater, error)
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

func (r *theaterRepository) GetTheaterByName(name string) (*models.Theater, error) {
	var theater models.Theater
	if err := r.db.Where("name = ?", name).First(&theater).Error; err != nil {
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
