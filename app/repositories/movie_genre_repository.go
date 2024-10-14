package repositories

import (
	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"gorm.io/gorm"
)

type MovieGenreRepository interface {
	UpdateGenresOfMovie(tx *gorm.DB, movieID uuid.UUID, genreIDs []uuid.UUID) error
}

func NewMovieGenreRepository(db *gorm.DB) MovieGenreRepository {
	return &movieGenreRepository{db}
}

type movieGenreRepository struct {
	db *gorm.DB
}

func (r *movieGenreRepository) UpdateGenresOfMovie(tx *gorm.DB, movieID uuid.UUID, genreIDs []uuid.UUID) error {
	if err := tx.Where("movie_id = ?", movieID).Delete(&models.MovieGenre{}).Error; err != nil {
		return err
	}

	if len(genreIDs) == 0 {
		return nil
	}

	newGenres := make([]models.MovieGenre, len(genreIDs))
	for i, genreID := range genreIDs {
		newGenres[i] = models.MovieGenre{
			MovieID: movieID,
			GenreID: genreID,
		}
	}

	return tx.Create(&newGenres).Error
}
