package repositories

import (
	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"gorm.io/gorm"
)

type GenreRepository interface {
	GetGenre(id uuid.UUID) (*models.Genre, error)
	GetGenreByName(name string) (*models.Genre, error)
	GetGenres() ([]*models.Genre, error)
	CreateGenre(tx *gorm.DB, genre *models.Genre) error
}

func NewGenreRepository(db *gorm.DB) GenreRepository {
	return &genreRepository{db: db}
}

type genreRepository struct {
	db *gorm.DB
}

func (r *genreRepository) GetGenre(id uuid.UUID) (*models.Genre, error) {
	var g models.Genre
	if err := r.db.Where(&models.Genre{ID: id}).First(&g).Error; err != nil {
		return nil, err
	}

	return &g, nil
}

func (r *genreRepository) GetGenreByName(name string) (*models.Genre, error) {
	var g models.Genre
	if err := r.db.Where(&models.Genre{Name: name}).First(&g).Error; err != nil {
		return nil, err
	}

	return &g, nil
}

func (r *genreRepository) GetGenres() ([]*models.Genre, error) {
	var genres []*models.Genre
	if err := r.db.Find(&genres).Error; err != nil {
		return nil, err
	}

	return genres, nil
}

func (r *genreRepository) CreateGenre(tx *gorm.DB, genre *models.Genre) error {
	return tx.Create(genre).Error
}
