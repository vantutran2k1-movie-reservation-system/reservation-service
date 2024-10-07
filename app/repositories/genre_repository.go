package repositories

import (
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"gorm.io/gorm"
)

type GenreRepository interface {
	GetGenreByName(name string) (*models.Genre, error)
	CreateGenre(tx *gorm.DB, genre *models.Genre) error
}

func NewGenreRepository(db *gorm.DB) GenreRepository {
	return &genreRepository{db: db}
}

type genreRepository struct {
	db *gorm.DB
}

func (r *genreRepository) GetGenreByName(name string) (*models.Genre, error) {
	var g models.Genre
	if err := r.db.Where(&models.Genre{Name: name}).First(&g).Error; err != nil {
		return nil, err
	}

	return &g, nil
}

func (r *genreRepository) CreateGenre(tx *gorm.DB, genre *models.Genre) error {
	return tx.Create(genre).Error
}
