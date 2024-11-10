package repositories

import (
	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/filters"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"gorm.io/gorm"
)

type GenreRepository interface {
	GetGenre(filter filters.GenreFilter) (*models.Genre, error)
	GetGenres(filter filters.GenreFilter) ([]*models.Genre, error)
	GetGenreIDs(filter filters.GenreFilter) ([]uuid.UUID, error)
	CreateGenre(tx *gorm.DB, genre *models.Genre) error
	UpdateGenre(tx *gorm.DB, genre *models.Genre) error
	DeleteGenre(tx *gorm.DB, genre *models.Genre) error
}

func NewGenreRepository(db *gorm.DB) GenreRepository {
	return &genreRepository{db: db}
}

type genreRepository struct {
	db *gorm.DB
}

func (r *genreRepository) GetGenre(filter filters.GenreFilter) (*models.Genre, error) {
	var g models.Genre
	if err := filter.GetFilterQuery(r.db).First(&g).Error; err != nil {
		if errors.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return &g, nil
}

func (r *genreRepository) GetGenres(filter filters.GenreFilter) ([]*models.Genre, error) {
	var genres []*models.Genre
	if err := filter.GetFilterQuery(r.db).Find(&genres).Error; err != nil {
		return nil, err
	}

	return genres, nil
}

func (r *genreRepository) GetGenreIDs(filter filters.GenreFilter) ([]uuid.UUID, error) {
	var genres []*models.Genre
	if err := filter.GetFilterQuery(r.db).Select("id").Find(&genres).Error; err != nil {
		return nil, err
	}

	ids := make([]uuid.UUID, len(genres))
	for i, g := range genres {
		ids[i] = g.ID
	}

	return ids, nil
}

func (r *genreRepository) CreateGenre(tx *gorm.DB, genre *models.Genre) error {
	return tx.Create(genre).Error
}

func (r *genreRepository) UpdateGenre(tx *gorm.DB, genre *models.Genre) error {
	return tx.Save(genre).Error
}

func (r *genreRepository) DeleteGenre(tx *gorm.DB, genre *models.Genre) error {
	return tx.Delete(genre).Error
}
