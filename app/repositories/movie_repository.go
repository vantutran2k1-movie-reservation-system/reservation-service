package repositories

import (
	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"gorm.io/gorm"
)

type MovieRepository interface {
	GetMovie(id uuid.UUID) (*models.Movie, error)
	GetMovies(limit, offset int) ([]*models.Movie, error)
	GetNumbersOfMovie() (int, error)
	CreateMovie(tx *gorm.DB, movie *models.Movie) error
	UpdateMovie(tx *gorm.DB, movie *models.Movie) error
}

type movieRepository struct {
	db *gorm.DB
}

func NewMovieRepository(db *gorm.DB) MovieRepository {
	return &movieRepository{db: db}
}

func (r *movieRepository) GetMovie(id uuid.UUID) (*models.Movie, error) {
	var m models.Movie
	if err := r.db.Where(&models.Movie{ID: id}).First(&m).Error; err != nil {
		return nil, err
	}

	return &m, nil
}

func (r *movieRepository) GetMovies(limit, offset int) ([]*models.Movie, error) {
	var movies []*models.Movie
	if err := r.db.Limit(limit).Offset(offset).Find(&movies).Error; err != nil {
		return nil, err
	}

	return movies, nil
}

func (r *movieRepository) GetNumbersOfMovie() (int, error) {
	var count int64
	if err := r.db.Model(&models.Movie{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

func (r *movieRepository) CreateMovie(tx *gorm.DB, movie *models.Movie) error {
	return tx.Create(movie).Error
}

func (r *movieRepository) UpdateMovie(tx *gorm.DB, movie *models.Movie) error {
	return tx.Save(movie).Error
}
