package repositories

import (
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"gorm.io/gorm"
)

type MovieRepository interface {
	CreateMovie(tx *gorm.DB, movie *models.Movie) error
}

type movieRepository struct {
	db *gorm.DB
}

func NewMovieRepository(db *gorm.DB) MovieRepository {
	return &movieRepository{db: db}
}

func (r *movieRepository) CreateMovie(tx *gorm.DB, movie *models.Movie) error {
	return tx.Save(movie).Error
}
