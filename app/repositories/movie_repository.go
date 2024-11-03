package repositories

import (
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/filters"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"gorm.io/gorm"
)

type MovieRepository interface {
	GetMovie(filter filters.MovieFilter, includeGenres bool) (*models.Movie, error)
	GetMovies(filter filters.MovieFilter) ([]*models.Movie, error)
	GetNumbersOfMovie(filter filters.MovieFilter) (int, error)
	CreateMovie(tx *gorm.DB, movie *models.Movie) error
	UpdateMovie(tx *gorm.DB, movie *models.Movie) error
}

type movieRepository struct {
	db *gorm.DB
}

func NewMovieRepository(db *gorm.DB) MovieRepository {
	return &movieRepository{db: db}
}

func (r *movieRepository) GetMovie(filter filters.MovieFilter, includeGenres bool) (*models.Movie, error) {
	query := filter.GetFilterQuery(r.db)
	if includeGenres {
		query = query.Preload("Genres")
	}

	var m models.Movie
	if err := query.First(&m).Error; err != nil {
		if errors.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return &m, nil
}

func (r *movieRepository) GetMovies(filter filters.MovieFilter) ([]*models.Movie, error) {
	var movies []*models.Movie
	if err := filter.GetFilterQuery(r.db).Find(&movies).Error; err != nil {
		return nil, err
	}

	return movies, nil
}

func (r *movieRepository) GetNumbersOfMovie(filter filters.MovieFilter) (int, error) {
	var count int64
	if err := filter.GetFilterQuery(r.db).Model(&models.Movie{}).Count(&count).Error; err != nil {
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
