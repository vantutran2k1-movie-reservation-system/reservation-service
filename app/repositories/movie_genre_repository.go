package repositories

import (
	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"gorm.io/gorm"
)

type MovieGenreRepository interface {
	UpdateGenresOfMovie(tx *gorm.DB, movieID uuid.UUID, genreIDs []uuid.UUID) error
	DeleteByMovieId(tx *gorm.DB, movieID uuid.UUID) error
	DeleteByGenreId(tx *gorm.DB, genreId uuid.UUID) error
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

func (r *movieGenreRepository) DeleteByMovieId(tx *gorm.DB, movieID uuid.UUID) error {
	return tx.Delete(&models.MovieGenre{}, "movie_id = ?", movieID).Error
}

func (r *movieGenreRepository) DeleteByGenreId(tx *gorm.DB, genreId uuid.UUID) error {
	return tx.Delete(&models.MovieGenre{}, "genre_id = ?", genreId).Error
}
